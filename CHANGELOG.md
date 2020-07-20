# Changelog

All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased](https://github.com/atomist/k8svent/compare/v0.17.0...HEAD)

## [v0.17.0](https://github.com/atomist/k8svent/compare/v0.16.0...v0.17.0) - 2020-07-20

### Changed

-   Check for pre-releases more frequently. [09de3b5](https://github.com/atomist/k8svent/commit/09de3b55b3c84fe6cc2574d2409bd955c476e55f)
-   Use image digest to check for new version. [b066373](https://github.com/atomist/k8svent/commit/b066373bccc9446ab5d1bc14a73d48effc42a4ca)

### Removed

-   Remove deprecated specs. [5cb3e98](https://github.com/atomist/k8svent/commit/5cb3e9867bdc6a35902c7f10550f87061ae82a1c)

## [v0.16.0](https://github.com/atomist/k8svent/compare/v0.15.1...v0.16.0) - 2020-07-15

### Added

-   Send deleted pods. [99f2d81](https://github.com/atomist/k8svent/commit/99f2d813387c845e240975403d08b3317544842e)

## [v0.15.1](https://github.com/atomist/k8svent/compare/v0.15.0...v0.15.1) - 2020-06-18

### Changed

-   Send all git files to build context. [463ca42](https://github.com/atomist/k8svent/commit/463ca42d3014a0cceadbd4561b8878aa4925818b)
-   Move kube specs up one directory. [ad3f7aa](https://github.com/atomist/k8svent/commit/ad3f7aa49e3620edb1f8e946c2e4f36aa8bd93d9)

### Fixed

-   Parse version tags tolerantly. [32c1101](https://github.com/atomist/k8svent/commit/32c11015c603ef90378914f82811ee1489c19991)

## [v0.15.0](https://github.com/atomist/k8svent/compare/v0.14.0...v0.15.0) - 2020-06-18

### Added

-   Investigate auto-update feature. [#11](https://github.com/atomist/k8svent/issues/11)

### Changed

-   Multistage Docker build with dynamic version. [bc53b57](https://github.com/atomist/k8svent/commit/bc53b5764f287a0edc5d3166a2d63926efcb7c49)

### Removed

-   Remove environment from payload. [#12](https://github.com/atomist/k8svent/issues/12)

## [v0.14.0](https://github.com/atomist/k8svent/compare/v0.13.1...v0.14.0) - 2020-06-12

### Changed

-   Missing state changes. [#10](https://github.com/atomist/k8svent/issues/10)

## [v0.13.1](https://github.com/atomist/k8svent/compare/v0.13.0...v0.13.1) - 2020-06-04

### Fixed

-   Fix signature header. [729ae2b](https://github.com/atomist/k8svent/commit/729ae2bc125917f3c5793d63b2d554095dca3cab)

## [v0.13.0](https://github.com/atomist/k8svent/compare/v0.12.0...v0.13.0) - 2020-06-02

### Added

-   Add support for signed payloads. [#9](https://github.com/atomist/k8svent/issues/9)

## [0.12.0](https://github.com/atomist/k8svent/compare/0.11.0...0.12.0) - 2020-03-23

### Changed

-   Name to k8svent
-   Use Go modules
-   Update Docker base image

## [0.11.0](https://github.com/atomist/k8svent/compare/0.10.0...0.11.0) - 2019-03-12

### Changed

-   k8svent now runs as non-root user in Docker container

## [0.10.0](https://github.com/atomist/k8svent/compare/0.9.0...0.10.0) - 2018-08-06

### Changed

-   Make logs more structured

## [0.9.0](https://github.com/atomist/k8svent/compare/0.8.0...0.9.0) - 2018-07-03

### Added

-   Parse response from posting webhook and include correlation in log

## [0.8.0](https://github.com/atomist/k8svent/compare/0.7.0...0.8.0) - 2018-06-05

### Added

-   Ability to run in a single namespace

## [0.7.0](https://github.com/atomist/k8svent/compare/0.6.0...0.7.0) - 2018-04-13

### Added

-   More logging

### Changed

-   Tightened up types

## [0.6.0](https://github.com/atomist/k8svent/compare/0.5.0...0.6.0) - 2018-03-06

### Changed

-   Stop sending pod deleted events

## [0.5.1](https://github.com/atomist/k8svent/compare/0.5.0...0.5.1) - 2018-03-02

### Fixed

-   Crash due to assignment to nil annotation map

## [0.5.0](https://github.com/atomist/k8svent/compare/0.4.0...0.5.0) - 2018-03-01

### Added

-   Cache for k8svent pod annotations

## [0.4.0](https://github.com/atomist/k8svent/compare/0.3.0...0.4.0) - 2018-03-01

### Changed

-   No longer require a global webhook URL

## [0.3.0](https://github.com/atomist/k8svent/compare/0.2.0...0.3.0) - 2018-02-28

### Added

-   Resources and instructions for deploying to Kubernetes clusters
-   Support per-pod environment via annotation [#4](https://github.com/atomist/k8svent/issues/4)

## [0.2.0](https://github.com/atomist/k8svent/compare/0.1.0...0.2.0) - 2018-02-08

### Added

-   You can provide pod-specific webhook URLS in the

## [0.1.0](https://github.com/atomist/k8svent/tree/0.1.0) - 2018-01-04

### Added

-   Everything
