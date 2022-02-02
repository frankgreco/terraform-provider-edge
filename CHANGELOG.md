# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.4] - 2022-02-01
### Added
- Support for all protocols that EdgeOS supports.

## [0.1.3] - 2022-02-01
### Added
- Partial `edge_port_group` update functionality. The description can now be updated.
### Changed
- Fixed a critical bug that caused all updates to panic.
- Updated `edge-sdk-go` dependency which fixed a bug which prevented the optional field `edge_firewall_ruleset.description` from not being set.

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
