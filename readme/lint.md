# lint

## Overview

Command `lint` has been designed to run code linting tools. It acts as a wrapper
to the original [Luacheck][] CLI tool and can run it through [Docker][] when not
locally available on the system.

- [Usage](#usage)
- [Examples](#examples)

## Usage

```txt
$ mod lint -h
usage: mod lint [<flags>]

Code linting tools: Luacheck.

Flags:
  -h, --help              Show context-sensitive help (also try --help-long and --help-man).
  -c, --config=".modcli"  Path to configuration file.
  -v, --version           Show application version.
  -d, --docker            Run through Docker.
  -l, --luacheck          Run Luacheck.
  -o, --original          Show original output instead.
```

## Examples

### Default

```txt
$ mod lint
[LUACHECK]

No issues found
```

[docker]: https://www.docker.com/
[luacheck]: https://github.com/mpeterv/luacheck
