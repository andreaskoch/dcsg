# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

### Changed
- Set the working directory (see: https://www.freedesktop.org/software/systemd/man/systemd.exec.html) for the docker-compose systemd services. Otherwise docker-compose cannot locate `env_file`s (see: https://docs.docker.com/compose/compose-file/#/envfile)

## [v0.1.0-alpha] - 2016-12-18

The prototype

### Added
- A prototype for the dcsg - A systemd service generator for Docker Compose projects
