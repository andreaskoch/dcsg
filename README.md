# dcsg: A systemd service generator for docker-compose

dcsg is a command-line utility for Linux that generates systemd services for Docker Compose projects.

If you have one or more docker compose projects running on your server you might want **create a systemd service** for each of them.

And **dcsg** is here to help you with just that. Quickly create systemd services from docker-compose files.

## Usage

`dcsg <action> [<path-to-a-docker-compose-file>] [<docker-compose-project-name]>`

### Install

Register a Docker Compose project as a systemd service:

```bash
dcsg install docker-compose.yml
```

### Uninstall

Uninstall the systemd service for the given Docker Compose project:

```bash
dcsg uninstall docker-compose.yml
```

### Help

Show the available actions and arguments:

```bash
dcsg help
```

The help for a specific action:

```bash
dcsg install --help
```

## What does dcsg do?

**dcsg** doesn't do much.

For the `install` action **dcsg** creates a systemd service (see: [installer.go](installer.go)):

1. Create a systemd service definition in `/etc/systemd/system/<project-name>.service`
2. Execute systemctl `daemon-reload`, `enable` and `start`

The name of the service created will be the project name of your docker-compose project.

For the `uninstall` action **dcsg** remove the systemd service it created earlier (see: [uninstaller.go](uninstaller.go)):

1. Execute a systemctl `stop` and `disable` for the docker-compose service
2. Delete the systemd service definition from `/etc/systemd/system/<project-name>.service`
3. Execute sytemctl `daemon-reload`

## Download

You can download pre-built binaries for Linux (64bit, ARM 5, ARM 6 and ARM 7) from [github.com » andreaskoch » dcsg » releases](/releases/latest):

- [Download for dcsg (Linux 64bit)](https://github.com/andreaskoch/dcsg/releases/download/v0.1.1-alpha/dcsg_linux_amd64)
- [Download for dcsg (Linux ARM5)](https://github.com/andreaskoch/dcsg/releases/download/v0.1.1-alpha/dcsg_linux_arm5)
- [Download for dcsg (Linux ARM6)](https://github.com/andreaskoch/dcsg/releases/download/v0.1.1-alpha/dcsg_linux_arm6)
- [Download for dcsg (Linux ARM7)](https://github.com/andreaskoch/dcsg/releases/download/v0.1.1-alpha/dcsg_linux_arm7)

```bash
curl -L https://github.com/andreaskoch/dcsg/releases/download/v0.1.1-alpha/dcsg_linux_amd64 > /usr/local/bin/dcsg
chmod +x /usr/local/bin/dcsg
```


## Build

If you have go installed you can use `go get` to fetch and build **dcsg**:

```bash
go get github.com/andreaskoch/dcsg
```

To **cross-compile dcsg** for the different Linux architectures (AMD64, ARM5, ARM6 and ARM7) you can use the `crosscompile` action of the make script:

```bash
go get github.com/andreaskoch/dcsg
cd $GOPATH/github.com/andreaskoch/dcsg
make crosscompile
```

## Licensing

»dcsg« is licensed under the **Apache License, Version 2.0**. See [LICENSE](LICENSE) for the full license text.
