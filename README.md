# 🧊 dwmblocks-go
Modular status bar for dwm written in Go. 

## Features
- **No dependencies** (except for Go itself)
- **Autonomous** operation of blocks (each block has own interval for updating)
- Update blocks by **signal**
- **Asynchronous** signal processing
- Сonfiguration **without recompilation**

## Install 
``` shell
git clone https://github.com/amaretur/dwmblocks-go
make install
```

## Run
To start dwmblocks-go, run the command: 
```
dwmblocks-go &
```
For **autorun** add this command to the script which launch dwm.
If you are using an X server to run dwm, use `.xinitrc` to autorun dwmblocks-go:

``` shell
...

exec dwm &

dwmblocks-go &

...
```

## Configuration
Configuration means setting the blocks, the time of their renewal, and the signals by they can be updated.
The default config is located in `~/.config/dwmblocks/config.json`:
``` js
{
	"blocks": [
		{
			"command": "date +'Date: %d.%M.%Y'",
			"interval": 60,
			"signal": 0
		},
		{
			"command": "date +'Time: %H:%M'",
			"interval": 2,
			"signal": 0
		}
	],
	"separator": " | "
}
```

### Where:
- **command** - any command (include scripts), which can be executed from /bin/bash
- **interval** - how many time does it take between executing blocks (in seconds) 0 means "not update"
- **signal** - signal for updating one of the blocks SID must be between 1 and 30, 0 means "not signal"

* **separator** - separator will be added between all bar status blocks 

By default, dwm searches for a config along this path, but if you want to specify your config located in a different path, use the "-c" flag at startup (check `dwmblocks-go -help`):
``` shell
# ... 
dwmblocks-go -c path/to/you/config &
# ...
```
If you need to restore the default config, use the command `dwmblock-go -d path/to/you/config`. If you do not specify a path to generate the config, the config will be created by the default path (`~/.config/dwmblocks/config.json`) (check `dwmblocks-go -help`).

## Signals
To update any block, you need to send a signal using the `kill -SID PID` command, where `SID` is the **signal specified in the config, increased by `34`**:
To get the PID of dwmblocks-go use `pgrep -a dwmblocks | awk '{ printf $1 }'`.

So, the complete command to send a signal might look like this: `kill -SID $(pgrep -a dwmblocks | awk '{ printf $1 }')`.

You can use this from any program that can accept bash scripts. Example for updating a block on press Alt+Shift:

Edit `config.json`:
``` js
{
	"blocks":[
		// ...
		{
			"command": "some command",
			"interval": 0,
			"signal": 1
		},
		// ...
	],
	// ...
}
```
Edit `config.h`
``` c
// ...
static const Key keys[] = {
	// ...
	{ Mod1Mask, XK_Shift_L, spawn, SHCMD("kill -35 $(pgrep -a dwmblocks | awk '{ printf $1 }'") },
	// ...
}
// ...
```

## Uninstall
To uninstall, run:
``` shell
make uninstall
```

## Credits

This work would not have been possible without [Luke's build of dwmblocks](https://github.com/LukeSmithxyz/dwmblocks) and [0-x-f](https://github.com/0-x-f/).
