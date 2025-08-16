# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.4.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.3.1...v1.4.0) (2025-08-16)


### Features

* remove weekend schedule restriction from renovate config ([d4ce4a0](https://github.com/d0ugal/mqtt-exporter/commit/d4ce4a0cdccf765aeb06b0f8c5479aa90ab22244))
* upgrade to Go 1.25 ([730d116](https://github.com/d0ugal/mqtt-exporter/commit/730d116afea5c218aed910bdaeac5c5c7af3064f))


### Bug Fixes

* revert golangci-lint config to version 2 for compatibility ([c814455](https://github.com/d0ugal/mqtt-exporter/commit/c814455cc64cc5861e68d86379b924109b09cc64))
* update golangci-lint config for Go 1.25 compatibility ([fca876d](https://github.com/d0ugal/mqtt-exporter/commit/fca876dc55a0e05322b81f9db7a5951b8abe9e06))

## [1.3.1](https://github.com/d0ugal/mqtt-exporter/compare/v1.3.0...v1.3.1) (2025-08-14)


### Bug Fixes

* ensure correct version reporting in release builds ([a02760e](https://github.com/d0ugal/mqtt-exporter/commit/a02760eb5760c4a35b2317668fc5a35ae2f352c5))

## [1.3.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.2.0...v1.3.0) (2025-08-14)


### Features

* add version info metric and subtle version display in h1 header ([be5fa68](https://github.com/d0ugal/mqtt-exporter/commit/be5fa6832fb7d9521a6d9d36083348d20697522d))
* add version to title, separate version info, and add copyright footer with GitHub links ([1d5551d](https://github.com/d0ugal/mqtt-exporter/commit/1d5551d9ce6e2add115da7af6ab4767d8442ac91))


### Bug Fixes

* update Dockerfile to inject version information during build ([decdd29](https://github.com/d0ugal/mqtt-exporter/commit/decdd29bc52d8dbd147f9a958e6c37d94a626432))

## [1.2.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.1.2...v1.2.0) (2025-08-13)


### Features

* add dynamic version information to web UI and CLI ([043b046](https://github.com/d0ugal/mqtt-exporter/commit/043b0462736ad7f58a1056680e72cd8b7ecd87fd))

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
