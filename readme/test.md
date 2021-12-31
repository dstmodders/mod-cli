# test

## Overview

Command `test` has been designed to run Busted tests. It acts as a wrapper to
the original [Busted][] CLI tool and can run it through [Docker][] when not
locally available on the system.

- [Usage](#usage)
- [Examples](#examples)

## Usage

```txt
$ mod test -h
usage: mod test

Testing tools: Busted.

Flags:
  -h, --help              Show context-sensitive help (also try --help-long and --help-man).
  -c, --config=".modcli"  Path to configuration file.
  -v, --version           Show application version.
```

## Examples

### Default

```txt
$ mod test
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
262 successes / 0 failures / 0 errors / 0 pending : 2.350863 seconds
```

[busted]: https://olivinelabs.com/busted/
[docker]: https://www.docker.com/
