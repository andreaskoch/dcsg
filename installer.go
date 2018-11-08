package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type installer interface {
	Install(projectDirectory string, dockerComposeFileName string, dockerComposeExtensionFileNames []string, projectName string) error
}

type systemdInstaller struct {
	systemdDirectory string
	commandExecutor  Executor
	dryRun           bool
	doPull           bool
}

func (installer *systemdInstaller) Install(projectDirectory string, dockerComposeFileName string, dockerComposeExtensionFileNames []string, projectName string) error {

	serviceName := getServiceName(projectName)
	serviceViewModel := serviceDefinition{
		ProjectName:                 projectName,
		ProjectDirectory:            projectDirectory,
		DockerComposeFile:           dockerComposeFileName,
		DockerComposeExtensionFiles: dockerComposeExtensionFileNames,
		DoPull:                      installer.doPull,
	}

	if err := installer.createSystemdService(serviceViewModel); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create a systemd service for project %q (Directory: %q, Docker Compose File: %q)", projectName, projectDirectory, dockerComposeFileName))
	}

	reloadError := installer.commandExecutor.Run("systemctl", "daemon-reload")
	if reloadError != nil {
		return errors.Wrap(reloadError, "Failed to reload the systemd configuration")
	}

	enableError := installer.commandExecutor.Run("systemctl", "enable", serviceName)
	if enableError != nil {
		return errors.Wrap(enableError, fmt.Sprintf("Failed to enable %q", serviceName))
	}

	startError := installer.commandExecutor.Run("systemctl", "start", serviceName)
	if startError != nil {
		return errors.Wrap(startError, fmt.Sprintf("Failed to start %q", serviceName))
	}

	return nil
}

func (installer *systemdInstaller) createSystemdService(service serviceDefinition) error {
	serviceFilePath := filepath.Join(installer.systemdDirectory, getServiceName(service.ProjectName))

	var file *os.File
	if installer.dryRun {
		log.Println("Installing this service file at:", serviceFilePath)
		file = os.Stdout
	} else {
		var err error

		file, err = os.OpenFile(serviceFilePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0664)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to open the systemd service file: %q", serviceFilePath))
		}
		defer file.Close()
	}

	serviceTemplate, err := template.New("systemdservice").Parse(serviceTemplate)
	if err != nil {
		return errors.Wrap(err, "Failed to parse systemd service template")
	}

	serviceTemplate.Execute(file, service)

	return nil
}

func getServiceName(projectName string) string {
	return fmt.Sprintf("%s.service", projectName)
}

const serviceTemplate = `[Unit]
Description={{ .ProjectName }} Service
After=network.service docker.service
Requires=docker.service

[Service]
Restart=always
RestartSec=10
TimeoutSec=300
WorkingDirectory={{ .ProjectDirectory }}
{{- if .DoPull }}
ExecStartPre=/usr/bin/env docker-compose -p "{{ .ProjectName }}" -f "{{ .DockerComposeFile }}" {{- range .DockerComposeExtensionFiles }} -f "{{ . -}}" {{ end }} pull
{{- end }}
ExecStart=/usr/bin/env docker-compose -p "{{ .ProjectName }}" -f "{{ .DockerComposeFile }}" {{- range .DockerComposeExtensionFiles }} -f "{{ . -}}" {{ end }} up
ExecStop=/usr/bin/env docker-compose -p "{{ .ProjectName }}" -f "{{ .DockerComposeFile }}" {{- range .DockerComposeExtensionFiles }} -f "{{ . -}}" {{ end }} stop
ExecStopPost=/usr/bin/env docker-compose -p "{{ .ProjectName }}" -f "{{ .DockerComposeFile }}" {{- range .DockerComposeExtensionFiles }} -f "{{ . -}}" {{ end }} down

[Install]
WantedBy=docker.service
`

type serviceDefinition struct {
	ProjectName                 string
	ProjectDirectory            string
	DockerComposeFile           string
	DockerComposeExtensionFiles []string
	DoPull                      bool
}
