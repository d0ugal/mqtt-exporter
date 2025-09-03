# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.10.1](https://github.com/d0ugal/mqtt-exporter/compare/v1.10.0...v1.10.1) (2025-09-03)


### Bug Fixes

* add build-args to pass version information in CI workflow ([6bb2566](https://github.com/d0ugal/mqtt-exporter/commit/6bb2566f3ce291ccd665fbc62d42a9d48a7b66b6))

## [1.10.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.9.1...v1.10.0) (2025-08-26)


### Features

* **docker:** use an unprivileged user during runtime ([baef927](https://github.com/d0ugal/mqtt-exporter/commit/baef927f59fafcd2e287e5d7dd9c475a8feedb15))

## [1.9.1](https://github.com/d0ugal/mqtt-exporter/compare/v1.9.0...v1.9.1) (2025-08-20)


### Bug Fixes

* remove redundant Service Information section from UI ([34cebf9](https://github.com/d0ugal/mqtt-exporter/commit/34cebf9e9567b8382de8a92a4a4982c7d03af6b8))

## [1.9.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.8.0...v1.9.0) (2025-08-20)


### Features

* optimize linting performance with caching ([1e253d0](https://github.com/d0ugal/mqtt-exporter/commit/1e253d0b5da816716edf10c15f59b305049211cb))


### Bug Fixes

* run Docker containers as current user to prevent permission issues ([bce4368](https://github.com/d0ugal/mqtt-exporter/commit/bce436845fd24f6576aabd7f093e41972ac209af))

## [1.8.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.7.0...v1.8.0) (2025-08-20)


### Features

* implement external HTML template with improved UI design ([d85743e](https://github.com/d0ugal/mqtt-exporter/commit/d85743ee39961878a8d78f6edb69b3c99dec955a))

## [1.7.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.6.0...v1.7.0) (2025-08-20)


### Features

* **api:** add pretty JSON formatting for metrics info endpoint ([fd6601c](https://github.com/d0ugal/mqtt-exporter/commit/fd6601caa30ef541b07c45d36f5d3e8fd6f225d9))
* **ui:** improve layout with grid endpoints and reorder sections ([824077e](https://github.com/d0ugal/mqtt-exporter/commit/824077e953d0a68b6252d5331ac84e8e71dcc504))
* **ui:** remove custom metrics-info endpoint and simplify UI ([906d2a2](https://github.com/d0ugal/mqtt-exporter/commit/906d2a23d2079bb49d6c0d570d956f55f173a4bb))

## [1.6.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.5.0...v1.6.0) (2025-08-19)


### Features

* **server:** add dynamic metrics information with collapsible interface ([f7be6aa](https://github.com/d0ugal/mqtt-exporter/commit/f7be6aa12c3666b0e13bc34e11c5593e5cd184e4))


### Bug Fixes

* **lint:** pre-allocate slices to resolve golangci-lint prealloc warnings ([5f7beb6](https://github.com/d0ugal/mqtt-exporter/commit/5f7beb6c31beafe0408e5fa6741d612e4046dd5e))

## [1.5.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.4.0...v1.5.0) (2025-08-17)


### Features

* add robust reconnection logic with exponential backoff ([7ee6220](https://github.com/d0ugal/mqtt-exporter/commit/7ee622007762b4c959b67cec25a8e68404fc1e70))

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
