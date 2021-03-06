# mod-cli

## Overview

CLI modding tools for [Don't Starve Together][] to automate the most common
tasks throughout the mods' development and improve the existing workflow. The
goal is to create a single tool for all development needs and the most common
tasks.

## Commands

We will use [Dev Tools][] as an example for all commands below:

```shell
cd /your/path/to/mod-dev-tools/
```

- [changelog](./readme/changelog.md)
- [format](./readme/format.md)
- [info](./readme/info.md)
- [lint](./readme/lint.md)
- [test](./readme/test.md)
- [workshop](./readme/workshop.md)

## Contributing

### Manually building from source

1. Install Go tools (1.17+): [https://go.dev/doc/install](https://go.dev/doc/install)
2. Clone the repository: `git clone https://github.com/dstmodders/mod-cli.git`
3. Build and install: `make install`

You can also use `go run .` to run it directly during the development.

## License

Released under the [MIT License](https://opensource.org/licenses/MIT).

[dev tools]: https://github.com/dstmodders/mod-dev-tools
[don't starve together]: https://www.klei.com/games/dont-starve-together
