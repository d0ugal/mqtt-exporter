# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.2](https://github.com/d0ugal/mqtt-exporter/compare/v1.1.1...v1.1.2) (2025-08-13)


### Bug Fixes

* **docker:** update Alpine base image to 3.22.1 for better security and reproducibility ([c2c9279](https://github.com/d0ugal/mqtt-exporter/commit/c2c92795893d3ce27a9f57533e74772bd9503e6a))

## [1.1.1](https://github.com/d0ugal/mqtt-exporter/compare/v1.1.0...v1.1.1) (2025-08-11)


### Bug Fixes

* automatically detect environment variables for configuration ([e777617](https://github.com/d0ugal/mqtt-exporter/commit/e777617be18ac50e165ae6fe780ec03ddd5e8d16))

## [1.1.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.0.0...v1.1.0) (2025-08-11)


### Features

* add environment variable configuration support ([0e66112](https://github.com/d0ugal/mqtt-exporter/commit/0e66112d553f60c6e45dec665a65d93e7b452934))

## 1.0.0 (2025-08-11)


### Features

* add release-please for automated releases ([cecdcb8](https://github.com/d0ugal/mqtt-exporter/commit/cecdcb80ff736ca33020931ea63d0d8a1166dae9))
* initial MQTT exporter implementation ([e17a08c](https://github.com/d0ugal/mqtt-exporter/commit/e17a08c6343dc6d08713a28c9c3f5c829b4bc8ea))


### Bug Fixes

* resolve CI issues and improve Docker setup ([86cda88](https://github.com/d0ugal/mqtt-exporter/commit/86cda88c7982680fc51f8fe499ce756592016514))

## [Unreleased]

### Added
- Initial MQTT exporter implementation
- Prometheus metrics for MQTT message monitoring
- Configuration management with YAML support
- Structured logging with JSON/text format
- HTTP server with metrics and health endpoints
- Docker support with multi-stage builds
- CI/CD pipeline with GitHub Actions
- Comprehensive test coverage
- Development workflow with linting and formatting
- Automated dependency updates with Renovate

### Changed
- N/A

### Deprecated
- N/A

### Removed
- N/A

### Fixed
- N/A

### Security
- N/A
