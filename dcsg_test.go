package main

import "testing"

func Test_getProjectDirectory(t *testing.T) {

	fileName := "/var/www/some-project/docker-compose.yml"
	result, _ := getProjectDirectory(fileName)

	expected := "/var/www/some-project"
	if result != expected {
		t.Fail()
		t.Logf("getProjectDirectory(%s) returned %q instead of %q", fileName, result, expected)
	}
}

func Test_getProjectName(t *testing.T) {

	fileName := "/var/www/some-project/docker-compose.yml"
	result, _ := getProjectName(fileName)

	expected := "someproject"
	if result != expected {
		t.Fail()
		t.Logf("getProjectName(%s) returned %q instead of %q", fileName, result, expected)
	}

}
