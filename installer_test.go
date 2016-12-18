package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func Test_createSystemdService_NoErrorIsReturned(t *testing.T) {
	targetDirectory := os.TempDir()

	serviceViewModel := serviceDefinition{
		ProjectName:       "exampleproject",
		ProjectDirectory:  "/var/www/example-project",
		DockerComposeFile: "docker-compose.yml",
	}

	result := createSystemdService(serviceViewModel, targetDirectory)

	if result != nil {
		t.Fail()
		t.Logf("createSystemdService returned %s", result)
	}
}

func Test_createSystemdService_ServiceFileIsWritten(t *testing.T) {
	targetDirectory := os.TempDir()

	serviceViewModel := serviceDefinition{
		ProjectName:       "exampleproject",
		ProjectDirectory:  "/var/www/example-project",
		DockerComposeFile: "docker-compose.yml",
	}

	result := createSystemdService(serviceViewModel, targetDirectory)

	if result != nil {
		t.Fail()
		t.Logf("createSystemdService returned %s", result)
	}

	targetFilePath := filepath.Join(targetDirectory, fmt.Sprintf("%s.service", serviceViewModel.ProjectName))
	_, err := os.Stat(targetFilePath)
	if err != nil {
		t.Fail()
		t.Logf("createSystemdService failed to create the service file %s: %s", targetFilePath, err)
	}
}
