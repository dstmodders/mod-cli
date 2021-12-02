# info

## Overview

Command `info` has been designed to interpretate `modinfo.lua` and give access
to the most values through CLI. This allows automating different tasks without
bothering interpretating or parsing `modinfo.lua` on your own.

It also includes some other helpful features.

- [Usage](#usage)
- [Examples](#examples)

  - [Default](#default)
  - [Configuration](#configuration-options)
  - [Configuration (Markdown)](#configuration-options-markdown)

## Usage

```txt
$ mod info -h
usage: mod info [<flags>] [<path>]

Mod info tools.

Flags:
  -h, --help                    Show context-sensitive help (also try --help-long and --help-man).
  -c, --config=".modcli"        Path to configurations file.
  -v, --version                 Show application version.
      --compatibility           Show compatibility fields.
      --configuration           Show configuration options with their default values.
  -m, --configuration-markdown  Show configuration options with their default values as a Markdown table.
  -d, --description             Show description.
  -f, --first-line              Show first lines for values.
  -g, --general                 Show general fields.
  -n, --names                   Show variable names or options data instead of their descriptions.
  -o, --other                   Show other fields.

Args:
  [<path>]  Path to modinfo.lua.
```

## Examples

For these examples, we will use [Dev Tools][] mod meaning that all the commands
below will be executed within its mod directory:

```shell
cd your/path/to/mod-dev-tools/
```

- [Default](#default)
- [Configuration](#configuration-options)
- [Configuration (Markdown)](#configuration-options-markdown)

### Default

```txt
$ mod info
[GENERAL]

Title: Dev Tools (dev)
Author: Depressed DST Modders
Version: 0.8.0-alpha
API Version: 10

[DESCRIPTION]

Version: 0.8.0-alpha

An extendable mod, that simplifies the most common tasks for both developers and testers as an alternative to debugkeys.

[COMPATIBILITY]

Don't Starve Compatible: no
Don't Starve Together Compatible: yes
Reign Of Giants Compatible: no

[CONFIGURATION]

Toggle tools key: Right Bracket
Switch data key: X
Select key: Tab
Movement prediction key: Disabled
Pause key: P
God mode key: G
Teleport key: T
Select entity key: Z
Increase time scale key: Page Up
Decrease time scale key: Page Down
Default time scale key: Home
Reset combination: Ctrl + R
Default god mode: Enabled
Default free crafting mode: Enabled
Default labels font: Stint Ultra...
Default labels font size: 18
Default selected labels: Enabled
Default username labels: Enabled
Default username labels mode: Default
Default forced HUD visibility: Enabled
Default forced unfading: Enabled
Disable mod warning: Enabled
Debug: Disabled

[OTHER]

Icon: modicon.tex
Icon Atlas: modicon.xml
Forum Thread: -
Priority: 1.0222050664
Folder Title: mod-dev-tools
```

### Configuration options

```txt
$ mod info --configuration
Toggle tools key: Right Bracket
Switch data key: X
Select key: Tab
Movement prediction key: Disabled
Pause key: P
God mode key: G
Teleport key: T
Select entity key: Z
Increase time scale key: Page Up
Decrease time scale key: Page Down
Default time scale key: Home
Reset combination: Ctrl + R
Default god mode: Enabled
Default free crafting mode: Enabled
Default labels font: Stint Ultra...
Default labels font size: 18
Default selected labels: Enabled
Default username labels: Enabled
Default username labels mode: Default
Default forced HUD visibility: Enabled
Default forced unfading: Enabled
Disable mod warning: Enabled
Debug: Disabled
```

### Configuration options (Markdown)

```txt
$ mod info -fm
| Configuration                     | Default          | Description                                                             |
| --------------------------------- | ---------------- | ----------------------------------------------------------------------- |
| **Toggle tools key**              | _Right Bracket_  | Key used for toggling the tools                                         |
| **Switch data key**               | _X_              | Key used for switching data sidebar                                     |
| **Select key**                    | _Tab_            | Key used for selecting between menu and data sidebar                    |
| **Movement prediction key**       | _Disabled_       | Key used for toggling the movement prediction                           |
| **Pause key**                     | _P_              | Key used for pausing the game                                           |
| **God mode key**                  | _G_              | Key used for toggling god mode                                          |
| **Teleport key**                  | _T_              | Key used for (fake) teleporting on mouse position                       |
| **Select entity key**             | _Z_              | Key used for selecting an entity under mouse                            |
| **Increase time scale key**       | _Page Up_        | Key used to speed up the time scale.                                    |
| **Decrease time scale key**       | _Page Down_      | Key used to slow down the time scale.                                   |
| **Default time scale key**        | _Home_           | Key used to restore the default time scale                              |
| **Reset combination**             | _Ctrl + R_       | Key combination used for reloading all mods.                            |
| **Default god mode**              | _Enabled_        | When enabled, enables god mode by default.                              |
| **Default free crafting mode**    | _Enabled_        | When enabled, enables crafting mode by default.                         |
| **Default labels font**           | _Stint Ultra..._ | Which labels font should be used by default?                            |
| **Default labels font size**      | _18_             | Which labels font size should be used by default?                       |
| **Default selected labels**       | _Enabled_        | When enabled, show selected labels by default.                          |
| **Default username labels**       | _Enabled_        | When enabled, shows username labels by default.                         |
| **Default username labels mode**  | _Default_        | Which username labels mode should be used by default?                   |
| **Default forced HUD visibility** | _Enabled_        | When enabled, forces HUD visibility when "playerhuddirty" event occurs. |
| **Default forced unfading**       | _Enabled_        | When enabled, forces unfading when "playerfadedirty" event occurs.      |
| **Disable mod warning**           | _Enabled_        | When enabled, disables the mod warning when starting the game           |
| **Debug**                         | _Disabled_       | When enabled, displays debug data in the console.                       |
```

[dev tools]: https://github.com/dstmodders/mod-dev-tools
