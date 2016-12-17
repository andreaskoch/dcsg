# dcsg - Systemd service generator for docker-compose

dcsg generates systemd services for your Docker Compose projects

## Usage

dcsg `<action>` -f `<path-to-a-docker-compose-file>`

Register a Docker Compose project as a systemd service:

```bas
dcsg install -f docker-compose.yml
```

Uninstall the systemd service for the given Docker Compose project:

```bash
dscg uninstall -f docker-compose.yml
```
