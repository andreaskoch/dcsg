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
	installCommand                     = app.Command("install", "Register a systemd service for the given docker-compose file")
	installDockerComposeFile           = installCommand.Flag("docker-compose-file", "A docker-compose file").Default("docker-compose.yml").ExistingFile()
	installDockerComposeExtensionFiles = installCommand.Flag("extend-docker-compose", "A docker-compose file").Default("docker-compose.yml").ExistingFiles()
	installProjectName                 = installCommand.Flag("project-name", "The project name of the docker-compose project").String()
	installDontPull                    = installCommand.Flag("no-pull", "The project name of the docker-compose project").Default("false").Bool()

	// uninstall
	uninstallCommand           = app.Command("uninstall", "Uninstall the systemd service for the given docker-compose file")
	uninstallDockerComposeFile = uninstallCommand.Flag("docker-compose-file", "A docker-compose file").Default("docker-compose.yml").ExistingFile()
	uninstallProjectName       = uninstallCommand.Flag("project-name", "The project name of the docker-compose project").String()
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
		service, err := newDscg(*installDockerComposeFile, *installDockerComposeExtensionFiles, *installProjectName, *appDryRun, !(*installDontPull))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		if err := service.Install(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

	case uninstallCommand.FullCommand():
		service, err := newDscg(*uninstallDockerComposeFile, *installDockerComposeExtensionFiles, *uninstallProjectName, *appDryRun, !(*installDontPull))
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
