# Changelog

All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased](https://github.com/atomist/k8svent/compare/v0.15.0...HEAD)

## [0.15.0](https://github.com/atomist/k8svent/compare/v0.14.0...v0.15.0) - 2020-06-18

### Added

-   Investigate auto-update feature. [#11](https://github.com/atomist/k8svent/issues/11)

### Changed

-   Multistage Docker build with dynamic version. [bc53b57](https://github.com/atomist/k8svent/commit/bc53b5764f287a0edc5d3166a2d63926efcb7c49)

### Removed

-   Remove environment from payload. [#12](https://github.com/atomist/k8svent/issues/12)

## [0.14.0](https://github.com/atomist/k8svent/compare/v0.13.1...v0.14.0) - 2020-06-12

### Changed

-   Missing state changes. [#10](https://github.com/atomist/k8svent/issues/10)

## [0.13.1](https://github.com/atomist/k8svent/compare/v0.13.0...v0.13.1) - 2020-06-04

### Fixed

-   Fix signature header. [729ae2b](https://github.com/atomist/k8svent/commit/729ae2bc125917f3c5793d63b2d554095dca3cab)

## [0.13.0](https://github.com/atomist/k8svent/compare/v0.12.0...v0.13.0) - 2020-06-02

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
