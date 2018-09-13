package main

import (
	"fmt"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const applicationName = "dcsg"
const applicationVersion = "v0.4.0"

var (
	app       = kingpin.New(applicationName, fmt.Sprintf("%s creates systemd services for Docker Compose projects (%s)", applicationName, applicationVersion))
	appDryRun = app.Flag("dry-run", "Print details of what would be done but do not install anything").Short('n').Bool()

	// install
	installCommand           = app.Command("install", "Register a systemd service for the given docker-compose file")
	installDockerComposeFile = installCommand.Arg("docker-compose-file", "A docker-compose file").Default("docker-compose.yml").String()
	installProjectName       = installCommand.Arg("project-name", "The project name of the docker-compose project").String()

	// uninstall
	uninstallCommand           = app.Command("uninstall", "Uninstall the systemd service for the given docker-compose file")
	uninstallDockerComposeFile = uninstallCommand.Arg("docker-compose-file", "A docker-compose file").Default("docker-compose.yml").String()
	uinstallProjectName        = uninstallCommand.Arg("project-name", "The project name of the docker-compose project").String()
)

func init() {
	app.Version(applicationVersion)
	app.Author("Andreas Koch <andy@ak7.io>")
}

func main() {
	handleCommandlineArgument(os.Args[1:])
}

func handleCommandlineArgument(arguments []string) {

	switch kingpin.MustParse(app.Parse(arguments)) {

	case installCommand.FullCommand():
		service, err := newDscg(*installDockerComposeFile, *installProjectName, *appDryRun)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		if err := service.Install(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

	case uninstallCommand.FullCommand():
		service, err := newDscg(*uninstallDockerComposeFile, *uinstallProjectName, *appDryRun)
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
