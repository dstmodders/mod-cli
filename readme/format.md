# format

## Overview

Command `format` has been designed to run code formatting tools. It acts as a
wrapper to the original CLI tools ([Prettier][] and [StyLua][]) and can run them
through [Docker][] when not locally available on the system.

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

### Default

```txt
$ mod format
[PRETTIER]

No issues found

[STYLUA]

No issues found
```

[docker]: https://www.docker.com/
[prettier]: https://prettier.io/
[stylua]: https://github.com/JohnnyMorganz/StyLua
