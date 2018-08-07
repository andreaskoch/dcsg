package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type uninstaller interface {
	Uninstall(projectDirectory, dockerComposeFileName, projectName string) error
}

type systemdUninstaller struct {
	systemdDirectory string
	commandExecutor  Executor
	dryRun           bool
}

func (uninstaller *systemdUninstaller) Uninstall(projectDirectory, dockerComposeFileName, projectName string) error {

	serviceName := getServiceName(projectName)
	serviceViewModel := serviceDefinition{
		ProjectName:       projectName,
		ProjectDirectory:  projectDirectory,
		DockerComposeFile: dockerComposeFileName,
	}

	stopError := uninstaller.commandExecutor.Run("systemctl", "stop", serviceName)
	if stopError != nil {
		return errors.Wrap(stopError, fmt.Sprintf("Failed to stop %q", serviceName))
	}

	disableError := uninstaller.commandExecutor.Run("systemctl", "disable", serviceName)
	if disableError != nil {
		return errors.Wrap(disableError, fmt.Sprintf("Failed to disable %q", serviceName))
	}

	removeError := uninstaller.removeSystemdService(serviceViewModel)
	if removeError != nil {
		return errors.Wrap(removeError, fmt.Sprintf("Failed to remove the systemd service %q", serviceViewModel.ProjectName))
	}

	reloadError := uninstaller.commandExecutor.Run("systemctl", "daemon-reload")
	if reloadError != nil {
		return errors.Wrap(reloadError, "Failed to reload the systemd configuration")
	}

	return nil
}

func (uninstaller *systemdUninstaller) removeSystemdService(service serviceDefinition) error {
	serviceFilePath := filepath.Join(uninstaller.systemdDirectory, getServiceName(service.ProjectName))

	if uninstaller.dryRun {
		log.Println("Would remove file:", serviceFilePath)
	} else {
		removeError := os.Remove(serviceFilePath)
		if removeError != nil {
			return errors.Wrap(removeError, fmt.Sprintf("Failed to remove %q", serviceFilePath))
		}
	}

	return nil
}
