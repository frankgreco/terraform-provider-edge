# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.6.0] - 2022-08-22
### Added
- Apple M1 support.

## [0.5.0] - 2022-08-22
### Added
- Acceptance tests for `edge_firewall_port_group`.
### Changed
- Added a new validator to all of the `description` attributes. While it can be `null`, if it is set, it must have a length of at least 1.

## [0.1.6] - 2022-02-07
### Added
- Support for optional field `edge_firewall_ruleset.default_logging`.
- Support for optional field `edge_firewall_ruleset.rules[*].log`.
### Changed
- Updated `edge-sdk-go` dependency which fixed a bug that prevented updates from working for certain resources.

## [0.1.5] - 2022-02-05
### Changed
- Updated `edge-sdk-go` dependency which fixed a bug that prevented null values for fields in the state field of ruleset if state was not null.

## [0.1.4] - 2022-02-01
### Added
- Support for all protocols that EdgeOS supports.

## [0.1.3] - 2022-02-01
### Added
- Partial `edge_port_group` update functionality. The description can now be updated.
### Changed
- Fixed a critical bug that caused all updates to panic.
- Updated `edge-sdk-go` dependency which fixed a bug that prevented the optional field `edge_firewall_ruleset.description` from not being set.

## [0.1.2] - 2022-01-15
### Changed
- Updated `edge-sdk-go` dependency which fixed a critical bug in `edge_firewall_ruleset` creation.
- Update port range bug in guided firewall example.

## [0.1.1] - 2022-01-09
### Added
- Initial stable release.
- Terraform provider.
- Terraform resource `edge_firewall_address_group` to manage EdgeOS address groups.
- Terraform resource `edge_firewall_port_group` to manage EdgeOS port groups.
- Terraform resource `edge_firewall_ruleset` to manage EdgeOS rulesets.
- Terraform resource `edge_firewall_ruleset_attachment` to attach rulesets to interfaces.
- Terraform data source `edge_interface_ethernet` to get details about EdgeOS ethernet interfaces.
- Terraform documentation.
