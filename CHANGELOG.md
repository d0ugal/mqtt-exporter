# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.23.2](https://github.com/d0ugal/mqtt-exporter/compare/v1.23.1...v1.23.2) (2025-10-26)


### Bug Fixes

* add internal version package and update version handling ([71c0f3a](https://github.com/d0ugal/mqtt-exporter/commit/71c0f3aa22420cf69c972cb568bdf55a0a02a00b))
* update module github.com/d0ugal/promexporter to v1.5.0 ([74e5d88](https://github.com/d0ugal/mqtt-exporter/commit/74e5d88795e1c2410a4925b247b3d8eace0103b9))
* update module github.com/prometheus/procfs to v0.19.1 ([f4c2df9](https://github.com/d0ugal/mqtt-exporter/commit/f4c2df952aa865c513965917005bc792104af853))
* use SkipVersionInfo to prevent duplicate version metrics ([a8a33a3](https://github.com/d0ugal/mqtt-exporter/commit/a8a33a311c3b7b3b3d13f970504f464999104bff))
* use WithVersionInfo to pass version info to promexporter ([9bd97fe](https://github.com/d0ugal/mqtt-exporter/commit/9bd97fe7cd9dfc5681e5512b00a42fa163f964d2))

## [1.23.1](https://github.com/d0ugal/mqtt-exporter/compare/v1.23.0...v1.23.1) (2025-10-25)


### Bug Fixes

* update module github.com/d0ugal/promexporter to v1.4.1 ([77c2525](https://github.com/d0ugal/mqtt-exporter/commit/77c25252a807fc89240ea35998804351ba8947d7))
* update module github.com/prometheus/procfs to v0.19.0 ([4bd1d17](https://github.com/d0ugal/mqtt-exporter/commit/4bd1d17c34a4c3751287eda00b47b487cd634df8))

## [1.23.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.22.0...v1.23.0) (2025-10-25)


### Features

* update promexporter to v1.4.0 ([7382e95](https://github.com/d0ugal/mqtt-exporter/commit/7382e95c974b9c34b6a0622f017ca64ce3b300e5))


### Bug Fixes

* update module github.com/d0ugal/promexporter to v1.1.0 ([9746748](https://github.com/d0ugal/mqtt-exporter/commit/97467488ec0d72aa30b1c939d7f0465c41020946))
* update module github.com/d0ugal/promexporter to v1.3.1 ([96c1aa6](https://github.com/d0ugal/mqtt-exporter/commit/96c1aa69a7dc95714109a0b3d0b1f25bd491e0e3))
* update module github.com/prometheus/procfs to v0.18.0 ([79129bd](https://github.com/d0ugal/mqtt-exporter/commit/79129bd1024599df00e9a610a84c9600d1eb7ec7))

## [1.22.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.21.0...v1.22.0) (2025-10-19)


### Features

* migrate mqtt-exporter to promexporter library ([61b99a5](https://github.com/d0ugal/mqtt-exporter/commit/61b99a5dc6706daed09c979cfa76fcc51f0c0b39))
* support duration strings for MQTT keep alive configuration ([6a0811d](https://github.com/d0ugal/mqtt-exporter/commit/6a0811daf1cc8207d87a3912f11e2f887a01b80e))
* update to promexporter v1.0.0 ([766e720](https://github.com/d0ugal/mqtt-exporter/commit/766e7208f8ebbbbb52814f1d4fc6739676aa90da))


### Bug Fixes

* add missing parse functions and update test calls ([b171357](https://github.com/d0ugal/mqtt-exporter/commit/b171357a926638e73cc46b094fe52151a9884408))
* correct Duration struct literals to use named fields ([6d63e32](https://github.com/d0ugal/mqtt-exporter/commit/6d63e3204f5f3696e56a44691ef3b3175541e13b))
* remove problematic config tests to unblock CI ([5c371f8](https://github.com/d0ugal/mqtt-exporter/commit/5c371f86ca53fcb37e7d79385651ef6fb7d86e7e))
* resolve linting issues - add empty line between embedded fields, fix godoc, and format imports ([be774f6](https://github.com/d0ugal/mqtt-exporter/commit/be774f6da870c8ac8472efef3c5aa806dc9ffbf1))
* update all remaining config test cases to use promexporter structure ([35491dc](https://github.com/d0ugal/mqtt-exporter/commit/35491dcdae8e94271164523f4fc3f511b8c3a8f6))
* update config tests to use new promexporter config structure ([03789de](https://github.com/d0ugal/mqtt-exporter/commit/03789de4e098154fb388a327149329056a8e319a))
* update go.mod after rebase ([05150ab](https://github.com/d0ugal/mqtt-exporter/commit/05150ab07486eaff3a0512162708bed73b57188b))
* update go.sum for promexporter v1.0.0 ([5a73f18](https://github.com/d0ugal/mqtt-exporter/commit/5a73f18944e726086db0e93f9be5a4618f1292e1))
* update main.go and remaining config tests for promexporter v1 ([2d2cc39](https://github.com/d0ugal/mqtt-exporter/commit/2d2cc39f590ad3c8efab362ad7f0344990eebf78))
* update module github.com/d0ugal/promexporter to v1.0.1 ([e4fe173](https://github.com/d0ugal/mqtt-exporter/commit/e4fe1738edbaefa71693eb3817d627caabd74a30))
* update remaining config test cases to use promexporter structure ([cdff132](https://github.com/d0ugal/mqtt-exporter/commit/cdff132bbcd8edd72696acab88f80a382d764918))

## [1.21.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.20.0...v1.21.0) (2025-10-14)


### Features

* improve MQTT connection robustness and auto-reconnection ([cf4547f](https://github.com/d0ugal/mqtt-exporter/commit/cf4547f301acd893c36b55043bf4de5236be532b))

## [1.20.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.19.0...v1.20.0) (2025-10-14)


### Features

* set Gin to release mode unless debug logging is enabled ([73e9595](https://github.com/d0ugal/mqtt-exporter/commit/73e95952934f30c79feb04e280656444ea16e0a2))

## [1.19.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.18.1...v1.19.0) (2025-10-14)


### Features

* add dynamic version information to web UI and CLI ([043b046](https://github.com/d0ugal/mqtt-exporter/commit/043b0462736ad7f58a1056680e72cd8b7ecd87fd))
* add environment variable configuration support ([0e66112](https://github.com/d0ugal/mqtt-exporter/commit/0e66112d553f60c6e45dec665a65d93e7b452934))
* add release-please for automated releases ([cecdcb8](https://github.com/d0ugal/mqtt-exporter/commit/cecdcb80ff736ca33020931ea63d0d8a1166dae9))
* add robust reconnection logic with exponential backoff ([7ee6220](https://github.com/d0ugal/mqtt-exporter/commit/7ee622007762b4c959b67cec25a8e68404fc1e70))
* add version info metric and subtle version display in h1 header ([be5fa68](https://github.com/d0ugal/mqtt-exporter/commit/be5fa6832fb7d9521a6d9d36083348d20697522d))
* add version to title, separate version info, and add copyright footer with GitHub links ([1d5551d](https://github.com/d0ugal/mqtt-exporter/commit/1d5551d9ce6e2add115da7af6ab4767d8442ac91))
* **api:** add pretty JSON formatting for metrics info endpoint ([fd6601c](https://github.com/d0ugal/mqtt-exporter/commit/fd6601caa30ef541b07c45d36f5d3e8fd6f225d9))
* **deps:** migrate to YAML v3 ([9d0e040](https://github.com/d0ugal/mqtt-exporter/commit/9d0e040a4d34e4f4b3efed1b77f386e79e6631db))
* **docker:** use an unprivileged user during runtime ([baef927](https://github.com/d0ugal/mqtt-exporter/commit/baef927f59fafcd2e287e5d7dd9c475a8feedb15))
* enable global automerge in Renovate config ([830d96b](https://github.com/d0ugal/mqtt-exporter/commit/830d96becb97e12681305a4d44ac98348ad359e9))
* implement external HTML template with improved UI design ([d85743e](https://github.com/d0ugal/mqtt-exporter/commit/d85743ee39961878a8d78f6edb69b3c99dec955a))
* initial MQTT exporter implementation ([e17a08c](https://github.com/d0ugal/mqtt-exporter/commit/e17a08c6343dc6d08713a28c9c3f5c829b4bc8ea))
* optimize linting performance with caching ([1e253d0](https://github.com/d0ugal/mqtt-exporter/commit/1e253d0b5da816716edf10c15f59b305049211cb))
* pin versions in documentation and examples ([8d5c837](https://github.com/d0ugal/mqtt-exporter/commit/8d5c83761058071db450356ac5b652b13134665f))
* remove weekend schedule restriction from renovate config ([d4ce4a0](https://github.com/d0ugal/mqtt-exporter/commit/d4ce4a0cdccf765aeb06b0f8c5479aa90ab22244))
* **renovate:** add docs commit message format for documentation updates ([b386e2f](https://github.com/d0ugal/mqtt-exporter/commit/b386e2f365818882119dee7d11aca7812e2991a8))
* **renovate:** use feat: commit messages for dependency updates ([5cb62a4](https://github.com/d0ugal/mqtt-exporter/commit/5cb62a40c766d289ff270d68c3ce7fe2a14dc926))
* replace latest docker tags with versioned variables for Renovate compatibility ([438fc96](https://github.com/d0ugal/mqtt-exporter/commit/438fc969572bf872947f4144b509e98ef6935b29))
* **server:** add dynamic metrics information with collapsible interface ([f7be6aa](https://github.com/d0ugal/mqtt-exporter/commit/f7be6aa12c3666b0e13bc34e11c5593e5cd184e4))
* **ui:** improve layout with grid endpoints and reorder sections ([824077e](https://github.com/d0ugal/mqtt-exporter/commit/824077e953d0a68b6252d5331ac84e8e71dcc504))
* **ui:** remove custom metrics-info endpoint and simplify UI ([906d2a2](https://github.com/d0ugal/mqtt-exporter/commit/906d2a23d2079bb49d6c0d570d956f55f173a4bb))
* update dependencies to v0.22.0 ([e51a563](https://github.com/d0ugal/mqtt-exporter/commit/e51a5637e6a4cef2c0df3638cd342a30ad3bb56c))
* update dependencies to v0.45.0 ([229c027](https://github.com/d0ugal/mqtt-exporter/commit/229c027eb43b7399d07a1beaa9191be62852a3b2))
* update dependencies to v0.67.1 ([1488eaa](https://github.com/d0ugal/mqtt-exporter/commit/1488eaabc7827bac83d58ec000d4ac978c92f705))
* update dependencies to v1.25.2 ([643895d](https://github.com/d0ugal/mqtt-exporter/commit/643895d4e8d175eaa700505c9cc21c03d23224e0))
* update dev build versioning to use semver-compatible pre-release tags ([767caba](https://github.com/d0ugal/mqtt-exporter/commit/767caba97f0ffa9c4ca2bd9f860bda23ec7fde3a))
* update module golang.org/x/crypto to v0.43.0 ([84dc5f1](https://github.com/d0ugal/mqtt-exporter/commit/84dc5f10014d901dde8bdc034c5be3fdd9cc12eb))
* update module golang.org/x/mod to v0.29.0 ([f4f9e23](https://github.com/d0ugal/mqtt-exporter/commit/f4f9e235c595882c5f3abc728f18360e343b8981))
* upgrade to Go 1.25 ([730d116](https://github.com/d0ugal/mqtt-exporter/commit/730d116afea5c218aed910bdaeac5c5c7af3064f))


### Bug Fixes

* add build-args to pass version information in CI workflow ([6bb2566](https://github.com/d0ugal/mqtt-exporter/commit/6bb2566f3ce291ccd665fbc62d42a9d48a7b66b6))
* automatically detect environment variables for configuration ([e777617](https://github.com/d0ugal/mqtt-exporter/commit/e777617be18ac50e165ae6fe780ec03ddd5e8d16))
* **ci:** add v prefix to dev tags for consistent versioning ([9d5d6f4](https://github.com/d0ugal/mqtt-exporter/commit/9d5d6f495faf3e1010c472744bdc87afb9f971ea))
* **deps:** update module github.com/eclipse/paho.mqtt.golang to v1.5.1 ([3c28ec1](https://github.com/d0ugal/mqtt-exporter/commit/3c28ec1975c19842c1fd225cd19267f4d33dd79d))
* **deps:** update module github.com/eclipse/paho.mqtt.golang to v1.5.1 ([8e31137](https://github.com/d0ugal/mqtt-exporter/commit/8e311377e991a0de7b2521883de12d3dc737590a))
* **deps:** update module github.com/gin-gonic/gin to v1.11.0 ([0c77854](https://github.com/d0ugal/mqtt-exporter/commit/0c778548093730c24b32fcee4350e5df1f09d794))
* **deps:** update module github.com/gin-gonic/gin to v1.11.0 ([69a971b](https://github.com/d0ugal/mqtt-exporter/commit/69a971badd813a13a7aafee4dcfc200a3d60ea30))
* **deps:** update module github.com/prometheus/client_golang to v1.23.1 ([bfb593d](https://github.com/d0ugal/mqtt-exporter/commit/bfb593d79df8cb5dfb2784e55a737275a794083c))
* **deps:** update module github.com/prometheus/client_golang to v1.23.1 ([2bc239b](https://github.com/d0ugal/mqtt-exporter/commit/2bc239ba2fff7fc79011823fe35fed2906981ee6))
* **deps:** update module github.com/prometheus/client_golang to v1.23.2 ([4cbb81c](https://github.com/d0ugal/mqtt-exporter/commit/4cbb81cf9af5a366f0975b8036399b8ad1d35b58))
* **deps:** update module github.com/prometheus/client_golang to v1.23.2 ([1d35111](https://github.com/d0ugal/mqtt-exporter/commit/1d3511183178a513c01cbb917251a2285f1be1af))
* **docker:** update Alpine base image to 3.22.1 for better security and reproducibility ([c2c9279](https://github.com/d0ugal/mqtt-exporter/commit/c2c92795893d3ce27a9f57533e74772bd9503e6a))
* enable indirect dependency updates in renovate config ([7822700](https://github.com/d0ugal/mqtt-exporter/commit/78227003c806ef783808aff2801f3355c51862a5))
* ensure correct version reporting in release builds ([a02760e](https://github.com/d0ugal/mqtt-exporter/commit/a02760eb5760c4a35b2317668fc5a35ae2f352c5))
* **lint:** pre-allocate slices to resolve golangci-lint prealloc warnings ([5f7beb6](https://github.com/d0ugal/mqtt-exporter/commit/5f7beb6c31beafe0408e5fa6741d612e4046dd5e))
* **lint:** resolve gosec configuration contradiction ([f200619](https://github.com/d0ugal/mqtt-exporter/commit/f200619974adc5007928b814ac5687581b7c405a))
* **lint:** resolve gosec G112 issue ([03fd06b](https://github.com/d0ugal/mqtt-exporter/commit/03fd06be678e8b76471a2ca2e169625db6cc2cf3))
* remove redundant Service Information section from UI ([34cebf9](https://github.com/d0ugal/mqtt-exporter/commit/34cebf9e9567b8382de8a92a4a4982c7d03af6b8))
* resolve CI issues and improve Docker setup ([86cda88](https://github.com/d0ugal/mqtt-exporter/commit/86cda88c7982680fc51f8fe499ce756592016514))
* revert golangci-lint config to version 2 for compatibility ([c814455](https://github.com/d0ugal/mqtt-exporter/commit/c814455cc64cc5861e68d86379b924109b09cc64))
* run Docker containers as current user to prevent permission issues ([bce4368](https://github.com/d0ugal/mqtt-exporter/commit/bce436845fd24f6576aabd7f093e41972ac209af))
* update dependency go to v1.25.3 ([81c4c86](https://github.com/d0ugal/mqtt-exporter/commit/81c4c863a7eb26851eda7ae6e080e8d61ac0c036))
* update Dockerfile to inject version information during build ([decdd29](https://github.com/d0ugal/mqtt-exporter/commit/decdd29bc52d8dbd147f9a958e6c37d94a626432))
* update golangci-lint config for Go 1.25 compatibility ([fca876d](https://github.com/d0ugal/mqtt-exporter/commit/fca876dc55a0e05322b81f9db7a5951b8abe9e06))
* update gomod commitMessagePrefix from feat to fix ([f9200a1](https://github.com/d0ugal/mqtt-exporter/commit/f9200a1054512e1f1ca8ae7980a4511a86599348))
* update module golang.org/x/tools to v0.38.0 ([a077ffb](https://github.com/d0ugal/mqtt-exporter/commit/a077ffbfd3c41a52338bd992412b617f91ce23f9))
* use actual release version as base for dev tags instead of hardcoded 0.0.0 ([7385db3](https://github.com/d0ugal/mqtt-exporter/commit/7385db38df6aa92e541af9661bdcb10ae2b7fde4))
* use fetch-depth: 0 instead of fetch-tags for full git history ([1250ad1](https://github.com/d0ugal/mqtt-exporter/commit/1250ad180632333b23d44ddef1f00eb365c7f42b))
* use fetch-tags instead of fetch-depth for GitHub Actions ([66648fe](https://github.com/d0ugal/mqtt-exporter/commit/66648fe8fd38b5a46d729139ef6d4137dcfbba54))


### Reverts

* remove unnecessary renovate config changes ([a4970cd](https://github.com/d0ugal/mqtt-exporter/commit/a4970cd179c671fbc93a4791ac2031d0377a3104))
* simplify release-please config ([80e033f](https://github.com/d0ugal/mqtt-exporter/commit/80e033f350865bb76cca848dac747b0dd21fed55))

## [1.18.1](https://github.com/d0ugal/mqtt-exporter/compare/v1.18.0...v1.18.1) (2025-10-14)


### Bug Fixes

* update module golang.org/x/tools to v0.38.0 ([a077ffb](https://github.com/d0ugal/mqtt-exporter/commit/a077ffbfd3c41a52338bd992412b617f91ce23f9))

## [1.18.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.17.0...v1.18.0) (2025-10-08)


### Features

* update dependencies to v0.22.0 ([e51a563](https://github.com/d0ugal/mqtt-exporter/commit/e51a5637e6a4cef2c0df3638cd342a30ad3bb56c))
* update module golang.org/x/crypto to v0.43.0 ([84dc5f1](https://github.com/d0ugal/mqtt-exporter/commit/84dc5f10014d901dde8bdc034c5be3fdd9cc12eb))
* update module golang.org/x/mod to v0.29.0 ([f4f9e23](https://github.com/d0ugal/mqtt-exporter/commit/f4f9e235c595882c5f3abc728f18360e343b8981))


### Bug Fixes

* update gomod commitMessagePrefix from feat to fix ([f9200a1](https://github.com/d0ugal/mqtt-exporter/commit/f9200a1054512e1f1ca8ae7980a4511a86599348))

## [1.17.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.16.0...v1.17.0) (2025-10-08)


### Features

* update dependencies to v1.25.2 ([643895d](https://github.com/d0ugal/mqtt-exporter/commit/643895d4e8d175eaa700505c9cc21c03d23224e0))

## [1.16.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.15.0...v1.16.0) (2025-10-07)


### Features

* **renovate:** use feat: commit messages for dependency updates ([5cb62a4](https://github.com/d0ugal/mqtt-exporter/commit/5cb62a40c766d289ff270d68c3ce7fe2a14dc926))
* update dependencies to v0.67.1 ([1488eaa](https://github.com/d0ugal/mqtt-exporter/commit/1488eaabc7827bac83d58ec000d4ac978c92f705))

## [1.15.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.14.0...v1.15.0) (2025-10-03)


### Features

* pin versions in documentation and examples ([8d5c837](https://github.com/d0ugal/mqtt-exporter/commit/8d5c83761058071db450356ac5b652b13134665f))
* **renovate:** add docs commit message format for documentation updates ([b386e2f](https://github.com/d0ugal/mqtt-exporter/commit/b386e2f365818882119dee7d11aca7812e2991a8))

## [1.14.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.13.6...v1.14.0) (2025-10-02)


### Features

* **deps:** migrate to YAML v3 ([9d0e040](https://github.com/d0ugal/mqtt-exporter/commit/9d0e040a4d34e4f4b3efed1b77f386e79e6631db))

## [1.13.6](https://github.com/d0ugal/mqtt-exporter/compare/v1.13.5...v1.13.6) (2025-10-02)


### Reverts

* remove unnecessary renovate config changes ([a4970cd](https://github.com/d0ugal/mqtt-exporter/commit/a4970cd179c671fbc93a4791ac2031d0377a3104))
* simplify release-please config ([80e033f](https://github.com/d0ugal/mqtt-exporter/commit/80e033f350865bb76cca848dac747b0dd21fed55))

## [1.13.5](https://github.com/d0ugal/mqtt-exporter/compare/v1.13.4...v1.13.5) (2025-10-02)


### Bug Fixes

* enable indirect dependency updates in renovate config ([7822700](https://github.com/d0ugal/mqtt-exporter/commit/78227003c806ef783808aff2801f3355c51862a5))

## [1.13.4](https://github.com/d0ugal/mqtt-exporter/compare/v1.13.3...v1.13.4) (2025-09-22)


### Bug Fixes

* **lint:** resolve gosec G112 issue ([03fd06b](https://github.com/d0ugal/mqtt-exporter/commit/03fd06be678e8b76471a2ca2e169625db6cc2cf3))

## [1.13.3](https://github.com/d0ugal/mqtt-exporter/compare/v1.13.2...v1.13.3) (2025-09-20)


### Bug Fixes

* **lint:** resolve gosec configuration contradiction ([f200619](https://github.com/d0ugal/mqtt-exporter/commit/f200619974adc5007928b814ac5687581b7c405a))

## [1.13.2](https://github.com/d0ugal/mqtt-exporter/compare/v1.13.1...v1.13.2) (2025-09-20)


### Bug Fixes

* **deps:** update module github.com/gin-gonic/gin to v1.11.0 ([0c77854](https://github.com/d0ugal/mqtt-exporter/commit/0c778548093730c24b32fcee4350e5df1f09d794))
* **deps:** update module github.com/gin-gonic/gin to v1.11.0 ([69a971b](https://github.com/d0ugal/mqtt-exporter/commit/69a971badd813a13a7aafee4dcfc200a3d60ea30))

## [1.13.1](https://github.com/d0ugal/mqtt-exporter/compare/v1.13.0...v1.13.1) (2025-09-16)


### Bug Fixes

* **deps:** update module github.com/eclipse/paho.mqtt.golang to v1.5.1 ([3c28ec1](https://github.com/d0ugal/mqtt-exporter/commit/3c28ec1975c19842c1fd225cd19267f4d33dd79d))
* **deps:** update module github.com/eclipse/paho.mqtt.golang to v1.5.1 ([8e31137](https://github.com/d0ugal/mqtt-exporter/commit/8e311377e991a0de7b2521883de12d3dc737590a))

## [1.13.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.12.2...v1.13.0) (2025-09-12)


### Features

* replace latest docker tags with versioned variables for Renovate compatibility ([438fc96](https://github.com/d0ugal/mqtt-exporter/commit/438fc969572bf872947f4144b509e98ef6935b29))

## [1.12.2](https://github.com/d0ugal/mqtt-exporter/compare/v1.12.1...v1.12.2) (2025-09-05)


### Bug Fixes

* **deps:** update module github.com/prometheus/client_golang to v1.23.2 ([4cbb81c](https://github.com/d0ugal/mqtt-exporter/commit/4cbb81cf9af5a366f0975b8036399b8ad1d35b58))
* **deps:** update module github.com/prometheus/client_golang to v1.23.2 ([1d35111](https://github.com/d0ugal/mqtt-exporter/commit/1d3511183178a513c01cbb917251a2285f1be1af))

## [1.12.1](https://github.com/d0ugal/mqtt-exporter/compare/v1.12.0...v1.12.1) (2025-09-04)


### Bug Fixes

* **deps:** update module github.com/prometheus/client_golang to v1.23.1 ([bfb593d](https://github.com/d0ugal/mqtt-exporter/commit/bfb593d79df8cb5dfb2784e55a737275a794083c))
* **deps:** update module github.com/prometheus/client_golang to v1.23.1 ([2bc239b](https://github.com/d0ugal/mqtt-exporter/commit/2bc239ba2fff7fc79011823fe35fed2906981ee6))

## [1.12.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.11.0...v1.12.0) (2025-09-04)


### Features

* update dev build versioning to use semver-compatible pre-release tags ([767caba](https://github.com/d0ugal/mqtt-exporter/commit/767caba97f0ffa9c4ca2bd9f860bda23ec7fde3a))


### Bug Fixes

* **ci:** add v prefix to dev tags for consistent versioning ([9d5d6f4](https://github.com/d0ugal/mqtt-exporter/commit/9d5d6f495faf3e1010c472744bdc87afb9f971ea))
* use actual release version as base for dev tags instead of hardcoded 0.0.0 ([7385db3](https://github.com/d0ugal/mqtt-exporter/commit/7385db38df6aa92e541af9661bdcb10ae2b7fde4))
* use fetch-depth: 0 instead of fetch-tags for full git history ([1250ad1](https://github.com/d0ugal/mqtt-exporter/commit/1250ad180632333b23d44ddef1f00eb365c7f42b))
* use fetch-tags instead of fetch-depth for GitHub Actions ([66648fe](https://github.com/d0ugal/mqtt-exporter/commit/66648fe8fd38b5a46d729139ef6d4137dcfbba54))

## [1.11.0](https://github.com/d0ugal/mqtt-exporter/compare/v1.10.1...v1.11.0) (2025-09-04)


### Features

* enable global automerge in Renovate config ([830d96b](https://github.com/d0ugal/mqtt-exporter/commit/830d96becb97e12681305a4d44ac98348ad359e9))

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
