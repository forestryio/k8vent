# Change Log

All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased]

[Unreleased]: https://github.com/atomist/k8vent/compare/0.6.0...HEAD

## [0.6.0][] - 2018-03-06

[0.6.0]: https://github.com/atomist/k8vent/compare/0.5.0...0.6.0

Undelete release

### Changed

-   Stop sending pod deleted events

## [0.5.1][] - 2018-03-02

[0.5.1]: https://github.com/atomist/k8vent/compare/0.5.0...0.5.1

Nil release

### Fixed

-   Crash due to assignment to nil annotation map

## [0.5.0][] - 2018-03-01

[0.5.0]: https://github.com/atomist/k8vent/compare/0.4.0...0.5.0

Cache release

### Added

-   Cache for k8vent pod annotations

## [0.4.0][] - 2018-03-01

[0.4.0]: https://github.com/atomist/k8vent/compare/0.3.0...0.4.0

Empty release

### Changed

-   No longer require a global webhook URL

## [0.3.0][] - 2018-02-28

[0.3.0]: https://github.com/atomist/k8vent/compare/0.2.0...0.3.0

Environment release

### Added

-   Resources and instructions for deploying to Kubernetes clusters
    using RBAC [#3][3]
-   Support per-pod environment via annotation [#4][4]

[3]: https://github.com/atomist/k8vent/issues/3
[4]: https://github.com/atomist/k8vent/issues/4

## [0.2.0][] - 2018-02-08

[0.2.0]: https://github.com/atomist/k8vent/compare/0.1.0...0.2.0

Multi-tenant release

### Added

-   You can provide pod-specific webhook URLS in the
    "atomist.com/k8vent" pod annotation [#2][2]

[2]: https://github.com/atomist/k8vent/issues/2

## [0.1.0][] - 2018-01-04

[0.1.0]: https://github.com/atomist/k8vent/tree/0.1.0

Initial release

### Added

-  Everything
