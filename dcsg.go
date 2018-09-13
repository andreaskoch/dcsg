package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func newDscg(dockerComposeFile, projectName string, dryRun bool) (*dcsg, error) {

	cleanedFilePath, err := filepath.Abs(dockerComposeFile)
	if err != nil {
		return nil, err
	}

	if !fileExists(cleanedFilePath) {
		return nil, fmt.Errorf("The Docker Compose file %q does not exist", cleanedFilePath)
	}

	if projectName == "" {
		projectNameFromDirectory, projectNameError := getProjectName(dockerComposeFile)
		if projectNameError != nil {
			return nil, projectNameError
		}

		projectName = projectNameFromDirectory
	}

	projectDirectory, err := getProjectDirectory(dockerComposeFile)
	if err != nil {
		return nil, err
	}

	dockerComposeFileName := filepath.Base(dockerComposeFile)

	systemdDirectory := "/etc/systemd/system"
	commandExecutor := newExecutor(os.Stdin, os.Stdout, os.Stderr, "", dryRun)

	return &dcsg{
		projectDirectory:      projectDirectory,
		dockerComposeFileName: dockerComposeFileName,
		projectName:           projectName,

		installer:   &systemdInstaller{systemdDirectory, commandExecutor, dryRun},
		uninstaller: &systemdUninstaller{systemdDirectory, commandExecutor, dryRun},
	}, nil
}

type dcsg struct {
	projectDirectory      string
	dockerComposeFileName string
	projectName           string

	installer   installer
	uninstaller uninstaller
}

func (service dcsg) Install() error {
	err := service.installer.Install(service.projectDirectory, service.dockerComposeFileName, service.projectName)
	if err != nil {
		return errors.Wrap(err, "Installation failed")
	}

	return nil
}

func (service dcsg) Uninstall() error {
	err := service.uninstaller.Uninstall(service.projectDirectory, service.dockerComposeFileName, service.projectName)
	if err != nil {
		return errors.Wrap(err, "Uninstall failed")
	}

	return nil
}

func getProjectDirectory(dockerComposeFile string) (string, error) {
	cleanedFilePath, err := filepath.Abs(dockerComposeFile)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Failed to determine the project directory from %q", dockerComposeFile))
	}

	return filepath.Dir(cleanedFilePath), nil
}

var invalidProjectNameCharacters = regexp.MustCompile("[^a-z0-9]")

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	return true
}

func getProjectName(dockerComposeFile string) (string, error) {
	directoryPath, err := getProjectDirectory(dockerComposeFile)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Failed to determine the project directory from %q", dockerComposeFile))
	}

	directoryName := filepath.Base(directoryPath)

	projectName := strings.ToLower(directoryName)
	projectName = strings.TrimSpace(projectName)
	projectName = invalidProjectNameCharacters.ReplaceAllString(projectName, "")

	return projectName, nil
}
