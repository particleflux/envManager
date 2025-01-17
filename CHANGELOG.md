# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

The entries for 1.0.0 and 1.1.0 are reconstructed from the git history, so there might have been more additions not
covered in the changelog.

## [Unreleased]


## [1.4.0-stefan] - 2022-11-15


## [1.4.0] - 2022-09-05
### Added
- `unload -a` as shorthand for `unload --all`
- Working directory-aware loading of config files
- `--local` switch for `config add mapping`
- Expand `.` in the path of directory mappings
### Changed
- Directory mappings now load dependencies
### Deprecated
- The flag `--config` was deprecated due to confusing behavior with directory-aware configuration

## [1.3.0] - 2022-02-01
### Added
- `current` command to show currently loaded profiles
- `--version` flag to display the version of envManager
- `unload --all` to unload all currently loaded profiles of this shell

## [1.2.0] - 2022-01-26
### Added
- This changelog
- [pass] Configuration option "prefix"

## [1.1.0] - 2021-07-31
### Added
- Directory mappings

## [1.0.0] - 2021-09-07
Initial version


[Unreleased]: https://github.com/DBX12/envManager/compare/v1.4.0-stefan...HEAD
[1.4.0-stefan]: https://github.com/DBX12/envManager/compare/v1.4.0...v1.4.0-stefan
[1.4.0]: https://github.com/DBX12/envManager/compare/v1.3.0...v1.4.0
[1.3.0]: https://github.com/DBX12/envManager/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/DBX12/envManager/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/DBX12/envManager/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/DBX12/envManager/releases/tag/v1.0.0