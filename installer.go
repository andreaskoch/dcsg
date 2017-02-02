package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type installer interface {
	Install(projectDirectory, dockerComposeFileName, projectName string) error
}

type systemdInstaller struct {
	systemdDirectory string
	commandExecutor  Executor
}

func (installer *systemdInstaller) Install(projectDirectory, dockerComposeFileName, projectName string) error {

	serviceName := getServiceName(projectName)
	serviceViewModel := serviceDefinition{
		ProjectName:       projectName,
		ProjectDirectory:  projectDirectory,
		DockerComposeFile: dockerComposeFileName,
	}

	if err := createSystemdService(serviceViewModel, installer.systemdDirectory); err != nil {
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

func createSystemdService(service serviceDefinition, targetDirectory string) error {
	serviceFilePath := filepath.Join(targetDirectory, getServiceName(service.ProjectName))
	file, err := os.OpenFile(serviceFilePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0664)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to open the systemd service file: %q", serviceFilePath))
	}

	defer file.Close()

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
After=network.target

[Service]
WorkingDirectory={{ .ProjectDirectory }}
ExecStartPre=/usr/bin/env docker-compose -p "{{ .ProjectName }}" -f "{{ .DockerComposeFile }}" pull
ExecStart=/usr/bin/env docker-compose -p "{{ .ProjectName }}" -f "{{ .DockerComposeFile }}" up
ExecStop=/usr/bin/env docker-compose -p "{{ .ProjectName }}" -f "{{ .DockerComposeFile }}" stop
ExecStopPost=/usr/bin/env docker-compose -p "{{ .ProjectName }}" -f "{{ .DockerComposeFile }}" down

[Install]
WantedBy=network-online.target
`

type serviceDefinition struct {
	ProjectName       string
	ProjectDirectory  string
	DockerComposeFile string
}
