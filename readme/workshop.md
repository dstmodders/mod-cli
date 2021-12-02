# workshop

## Overview

Package workshop has been designed to prepare a mod directory or archive for
Steam Workshop. It allows including only the essential files based on ignore
paths. In the future, it may also include the features to automatically publish
your mod.

By default, it ignores the following globs:

- `.*`
- `codecov.yml`
- `config.ld`
- `lcov.info`
- `luacov.*`
- `Makefile`
- `spec/`

You may override the defaults by using the configuration file. See
[Configuration][] to learn more.

- [Usage](#usage)
- [Configuration][]
- [Examples](#examples)

## Usage

```txt
$ mod workshop -h
usage: mod workshop [<flags>] [<path>]

Steam Workshop tools.

Flags:
  -h, --help              Show context-sensitive help (also try --help-long and --help-man).
  -c, --config=".modcli"  Path to configurations file.
  -v, --version           Show application version.
  -n, --name="workshop"   Name of destination directory/archive.
  -z, --zip               Create a ZIP archive instead.

Args:
  [<path>]  Path to mod directory.
```

## Configuration

```yml
workshop:
  ignore:
    - "*.md"
    - "*.zip"
    - ".*"
    - "Makefile"
    - "codecov.yml"
    - "config.ld"
    - "docs/"
    - "lcov.info"
    - "luacov.*"
    - "modicon.png"
    - "preview.png"
    - "readme/"
    - "spec/"
```

## Examples

For these examples, we will use [Dev Tools][] mod meaning that all the commands
below will be executed within its mod directory:

```shell
cd /your/path/to/mod-dev-tools/
```

### Default

```txt
$ mod workshop
[INFO]

Name: Dev Tools (dev)
Version: 0.8.0-alpha

[PATHS]

Source: /home/victor/Teams/dstmodders/mod-dev-tools
Destination: /home/victor/Teams/dstmodders/mod-dev-tools/workshop

[FILES | TOTAL: 95]

LICENSE
modicon.tex
modicon.xml
modinfo.lua
modmain.lua
scripts/devtools/api.lua
scripts/devtools/config.lua
scripts/devtools/console.lua
scripts/devtools/constants.lua
scripts/devtools/data/data.lua
scripts/devtools/data/dumpeddata.lua
scripts/devtools/data/frontenddata.lua
scripts/devtools/data/recipedata.lua
scripts/devtools/data/selecteddata.lua
scripts/devtools/data/selectedtagsdata.lua
scripts/devtools/data/worlddata.lua
scripts/devtools/data/worldstatedata.lua
scripts/devtools/debug/debug.lua
scripts/devtools/debug/debugevents.lua
scripts/devtools/debug/debugglobals.lua
scripts/devtools/debug/debugplayercontroller.lua
scripts/devtools/labels.lua
scripts/devtools/menu/menu.lua
scripts/devtools/menu/option/actionoption.lua
scripts/devtools/menu/option/checkboxoption.lua
scripts/devtools/menu/option/choicesoption.lua
scripts/devtools/menu/option/divideroption.lua
scripts/devtools/menu/option/fontoption.lua
scripts/devtools/menu/option/numericoption.lua
scripts/devtools/menu/option/option.lua
scripts/devtools/menu/option/submenuoption.lua
scripts/devtools/menu/option/togglecheckboxoption.lua
scripts/devtools/menu/submenu.lua
scripts/devtools/menu/textmenu.lua
scripts/devtools/sdk/LICENSE
scripts/devtools/sdk/sdk/config.lua
scripts/devtools/sdk/sdk/console.lua
scripts/devtools/sdk/sdk/constant.lua
scripts/devtools/sdk/sdk/debug.lua
scripts/devtools/sdk/sdk/debugupvalue.lua
scripts/devtools/sdk/sdk/dump.lua
scripts/devtools/sdk/sdk/entity.lua
scripts/devtools/sdk/sdk/frontend.lua
scripts/devtools/sdk/sdk/input.lua
scripts/devtools/sdk/sdk/method.lua
scripts/devtools/sdk/sdk/minimap.lua
scripts/devtools/sdk/sdk/modmain.lua
scripts/devtools/sdk/sdk/persistentdata.lua
scripts/devtools/sdk/sdk/player/attribute.lua
scripts/devtools/sdk/sdk/player/craft.lua
scripts/devtools/sdk/sdk/player/inventory.lua
scripts/devtools/sdk/sdk/player/minimap.lua
scripts/devtools/sdk/sdk/player/vision.lua
scripts/devtools/sdk/sdk/player.lua
scripts/devtools/sdk/sdk/remote/player.lua
scripts/devtools/sdk/sdk/remote/world.lua
scripts/devtools/sdk/sdk/remote.lua
scripts/devtools/sdk/sdk/rpc.lua
scripts/devtools/sdk/sdk/sdk.lua
scripts/devtools/sdk/sdk/temporarydata.lua
scripts/devtools/sdk/sdk/test.lua
scripts/devtools/sdk/sdk/thread.lua
scripts/devtools/sdk/sdk/time.lua
scripts/devtools/sdk/sdk/utils/chain.lua
scripts/devtools/sdk/sdk/utils/table.lua
scripts/devtools/sdk/sdk/utils/value.lua
scripts/devtools/sdk/sdk/utils.lua
scripts/devtools/sdk/sdk/vision.lua
scripts/devtools/sdk/sdk/world/savedata.lua
scripts/devtools/sdk/sdk/world/season.lua
scripts/devtools/sdk/sdk/world/weather.lua
scripts/devtools/sdk/sdk/world.lua
scripts/devtools/submenus/characterrecipessubmenu.lua
scripts/devtools/submenus/debug.lua
scripts/devtools/submenus/devtools.lua
scripts/devtools/submenus/dumpsubmenu.lua
scripts/devtools/submenus/labels.lua
scripts/devtools/submenus/language.lua
scripts/devtools/submenus/map.lua
scripts/devtools/submenus/option/toggle.lua
scripts/devtools/submenus/playerbarssubmenu.lua
scripts/devtools/submenus/playervision.lua
scripts/devtools/submenus/seasoncontrol.lua
scripts/devtools/submenus/selectsubmenu.lua
scripts/devtools/submenus/teleportsubmenu.lua
scripts/devtools/submenus/timecontrol.lua
scripts/devtools/submenus/weathercontrol.lua
scripts/devtools/tools/playercraftingtools.lua
scripts/devtools/tools/playerinventorytools.lua
scripts/devtools/tools/playertools.lua
scripts/devtools/tools/playervisiontools.lua
scripts/devtools/tools/tools.lua
scripts/devtools/tools/worldtools.lua
scripts/devtools.lua
scripts/screens/devtoolsscreen.lua
---
Done
```

[configuration]: #configuration
[dev tools]: https://github.com/dstmodders/mod-dev-tools
