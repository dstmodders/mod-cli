# changelog

## Overview

Command `changelog` has been designed to parse `CHANGELOG.md` and give access to
the information about existing releases. However, in the future, it may support
manipulating and generating changelogs as well if that kind of behaviour will be
needed in our CLI tools.

- [Usage](#usage)
- [Examples](#examples)

## Usage

```txt
$ mod changelog -h
usage: mod changelog [<flags>] [<path>]

Changelog tools.

Flags:
  -h, --help              Show context-sensitive help (also try --help-long and --help-man).
  -c, --config=".modcli"  Path to configuration file.
  -v, --version           Show application version.
      --count             Show total number of releases.
  -f, --first             Show first release.
  -l, --latest            Show latest release.
      --list              Show list of releases without changes.
      --list-versions     Show list of versions.

Args:
  [<path>]  Path.
```

## Examples

For these examples, we will use [Dev Tools][] mod meaning that all the commands
below will be executed within its mod directory:

```shell
cd your/path/to/mod-dev-tools/
```

### Default

```txt
$ mod changelog
[UNRELEASED]

ADDED

- Add "Hide Ground Overlay" player vision suboption
- Add support for args in the toggle checkbox option
- Add support for pause and time scale keys outside of gameplay

CHANGED

- Change author
- Migrate to the new mod SDK
- Rename and restructure some classes

REMOVED

- Remove "Hide changelog" configuration
- Remove changelog from modinfo

[0.7.0 | 2020-10-06 | https://github.com/dstmodders/mod-dev-tools/compare/v0.6.0...v0.7.0]

ADDED

- Add new "Dumped" data sidebar
- Add new "World State" data sidebar

CHANGED

- Change "Dump" submenu
- Refactor modinfo

FIXED

- Fix issue with menu update when selecting player

[0.6.0 | 2020-09-29 | https://github.com/dstmodders/mod-dev-tools/compare/v0.5.0...v0.6.0]

ADDED

- Add new "Selected Entity Tags" data sidebar
- Add support for "data_sidebar" in submenu data tables
- Add support for showing the number of data sidebars

CHANGED

- Refactor data sidebars

FIXED

- Fix issue with data sidebar scrolling position while switching

[0.5.0 | 2020-09-23 | https://github.com/dstmodders/mod-dev-tools/compare/v0.4.1...v0.5.0]

ADDED

- Add "Select key" configuration
- Add support for selecting either menu or data sidebar
- Add support for the mouse scroll in data sidebar

CHANGED

- Improve data sidebar update when selecting entity

FIXED

- Fix issue with "Switch data key" configuration

[0.4.1 | 2020-09-21 | https://github.com/dstmodders/mod-dev-tools/compare/v0.4.0...v0.4.1]

FIXED

- Fix issue with crashing related to data loading

[0.4.0 | 2020-09-21 | https://github.com/dstmodders/mod-dev-tools/compare/v0.3.0...v0.4.0]

ADDED

- Add "Locale Text Scale" in the front-end data sidebar
- Add "Toggle Locale Text Scale" in "Dev Tools" submenu
- Add new "Dev Tools" submenu
- Add new font option

CHANGED

- Refactor data sidebars

[0.3.0 | 2020-09-20 | https://github.com/dstmodders/mod-dev-tools/compare/v0.2.0...v0.3.0]

ADDED

- Add "Switch data key" configuration
- Add support for the sidebar data switching

CHANGED

- Change some configurations
- Improve overlay sizing and centring

FIXED

- Fix issue with d_gettags console command

[0.2.0 | 2020-09-11 | https://github.com/dstmodders/mod-dev-tools/compare/v0.1.0...v0.2.0]

ADDED

- Add "Hide changelog" configuration
- Add API support

CHANGED

- Enable "Player Vision" submenu on non-admin servers
- Increase mod loading priority

FIXED

- Fix issue with crashing when disabling a recipe tab

[0.1.0 | 2020-09-05]

First release.
```

[dev tools]: https://github.com/dstmodders/mod-dev-tools
