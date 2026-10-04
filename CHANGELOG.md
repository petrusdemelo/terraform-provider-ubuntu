# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project uses
[Semantic Versioning](https://semver.org/). Until 1.0.0, a minor release may
contain breaking changes; they are listed under **Breaking changes**.

## Unreleased

### Added

- Provider configuration: `ssh` block (key or password authentication, pinned
  host key) and `default_target` block. Connection values that are unknown at
  plan time are rejected.
- `ubuntu_os_release` data source.
