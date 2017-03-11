# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

### Changed
- Change the application version to v0.3.0

## [v0.3.0] - 2017-03-11

Timeouts and restart policy

### Added
- Add a startup timeout (see: dcsg#4)
- Add restart-policy (see: dcsg#4)

## [v0.2.0]

systemd service improvements (thanks to @hermsi1337)

### Added
- Add an animation illustrating the usage of dcsg

### Changed
- Pull the latest image before starting (see: dcsg/pull/2)
- Stop the containers before removing them (see: dcsg/pull/2)

## [v0.1.1-alpha] - 2016-12-18

Set the working directory for the systemd services

### Changed
- Set the working directory (see: https://www.freedesktop.org/software/systemd/man/systemd.exec.html) for the docker-compose systemd services. Otherwise docker-compose cannot locate `env_file`s (see: https://docs.docker.com/compose/compose-file/#/envfile)

## [v0.1.0-alpha] - 2016-12-18

The prototype

### Added
- A prototype for the dcsg - A systemd service generator for Docker Compose projects
