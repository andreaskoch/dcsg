package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func getTestInstaller() *systemdInstaller {
	targetDirectory := os.TempDir()
	commandExecutor := newExecutor(os.Stdin, os.Stdout, os.Stderr, "", false)
	return &systemdInstaller{targetDirectory, commandExecutor, false}
}

func Test_createSystemdService_NoErrorIsReturned(t *testing.T) {
	installer := getTestInstaller()

	serviceViewModel := serviceDefinition{
		ProjectName:       "exampleproject",
		ProjectDirectory:  "/var/www/example-project",
		DockerComposeFile: "docker-compose.yml",
	}

	result := installer.createSystemdService(serviceViewModel)

	if result != nil {
		t.Fail()
		t.Logf("createSystemdService returned %s", result)
	}
}

func Test_createSystemdService_ServiceFileIsWritten(t *testing.T) {
	installer := getTestInstaller()

	serviceViewModel := serviceDefinition{
		ProjectName:       "exampleproject",
		ProjectDirectory:  "/var/www/example-project",
		DockerComposeFile: "docker-compose.yml",
	}

	result := installer.createSystemdService(serviceViewModel)

	if result != nil {
		t.Fail()
		t.Logf("createSystemdService returned %s", result)
	}

	targetFilePath := filepath.Join(installer.systemdDirectory, fmt.Sprintf("%s.service", serviceViewModel.ProjectName))
	_, err := os.Stat(targetFilePath)
	if err != nil {
		t.Fail()
		t.Logf("createSystemdService failed to create the service file %s: %s", targetFilePath, err)
	}
}
