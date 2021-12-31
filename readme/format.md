# format

## Overview

Command `format` has been designed to run code formatting tools. It acts as a
wrapper to the original CLI tools (Prettier and StyLua) and can run them through
Docker when not locally available on the system.

- [Usage](#usage)
- [Examples](#examples)

## Usage

```txt
$ mod format -h
usage: mod format [<flags>]

Code formatting tools: Prettier and StyLua.

Flags:
  -h, --help              Show context-sensitive help (also try --help-long and --help-man).
  -c, --config=".modcli"  Path to configuration file.
  -v, --version           Show application version.
  -d, --docker            Run through Docker.
  -f, --fix               Fix issues automatically. Beware!
  -p, --prettier          Run Prettier.
  -s, --stylua            Run StyLua.
```

## Examples

For these examples, we will use [Dev Tools][] mod meaning that all the commands
below will be executed within its mod directory:

```shell
cd your/path/to/mod-dev-tools/
```

- [Default](#default)

### Default

```txt
$ mod format
[PRETTIER]

No issues found

[STYLUA]

No issues found
```

[dev tools]: https://github.com/dstmodders/mod-dev-tools
