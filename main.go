package main

import (
	"fmt"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const applicationName = "dcsg"
const applicationVersion = "v0.1.0-alpha"

func main() {
	handleCommandlineArgument(os.Args[1:])
}

func handleCommandlineArgument(arguments []string) {

	app := kingpin.New(applicationName, fmt.Sprintf("%s creates systemd services for Docker Compose projects", applicationName))
	app.Version(applicationVersion)
	app.Writer(os.Stderr)
	app.Terminate(func(int) {
		return
	})

	// install
	installCommand := app.Command("install", "Register a systemd service for the given docker-compose file")
	installDockerComposeFile := installCommand.Arg("f", "A docker-compose file").Default("docker-compose.yml").String()
	installProjectName := installCommand.Arg("p", "A project name for your docker-compose project").String()

	// uninstall
	uninstallCommand := app.Command("uninstall", "Uninstall the systemd service for the given docker-compose file")
	uninstallDockerComposeFile := uninstallCommand.Arg("f", "A docker-compose file").Default("docker-compose.yml").String()
	uinstallProjectName := uninstallCommand.Arg("p", "A project name for your docker-compose project").String()

	command, err := app.Parse(arguments)
	if err != nil {
		app.Fatalf("%s", err.Error())
		return
	}

	switch command {

	case installCommand.FullCommand():
		service, err := newDscg(*installDockerComposeFile, *installProjectName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		if err := service.Install(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

	case uninstallCommand.FullCommand():
		service, err := newDscg(*uninstallDockerComposeFile, *uinstallProjectName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		if err := service.Uninstall(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

	}

}
