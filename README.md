# mini_launch

[中文](README_zh.md)

A simple cross-platform service daemon manager for macOS and Linux.

`mini_launch` manages long-running services using native platform tools:
- **macOS**: `launchd` via `launchctl`
- **Linux**: `systemd` user services via `systemctl --user`

## Why mini_launch?

Writing plist or systemd service files by hand is tedious. `mini_launch` generates the config for you — one command and you're done.

It also establishes a simple convention for teams:

- **Unified service directory** — all services live under `$HOME/servers/`, easy to find.
- **Unified log location** — logs go to `std.log` in the service directory, no more hunting around.
- **Auto environment injection** — reads `export` statements from `~/.zshrc` or `~/.bashrc`, so your daemons use the same env vars as your terminal.

## Quick Start

```bash
# 1. Create a directory under $HOME/servers/ and put an executable in it
mkdir -p $HOME/servers/myapp
cp /path/to/myapp $HOME/servers/myapp/

# 2. cd into the directory
cd $HOME/servers/myapp

# 3. Initialize and start — that's it
mini_launch initial
mini_launch start
```

> **Note:** Each service directory must contain exactly one executable file. `mini_launch` will auto-detect it.

## How It Works

Services are organized under `$HOME/servers/`. Each service resides in its own subdirectory (nested directories are supported) with a single executable file:

```
$HOME/servers/
├── myapp/
│   └── myapp            # executable
├── myserver/
│   └── myserver         # executable
└── web/
    └── api/
        └── api-server   # executable → service name: web_api
```

- **Service name** = relative path from `$HOME/servers/` with `/` replaced by `_`
- **Log output** → `std.log` in the service directory
- **Environment variables** are automatically extracted from your shell config (`~/.zshrc` or `~/.bashrc`)

## Install

### Homebrew (macOS)

```bash
brew tap simpossible/tap
brew install mini_launch
```

### One-click install (Linux & macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/simpossible/mini_launch/main/install.sh | bash
```

### DEB / RPM (Linux)

Download the `.deb` or `.rpm` package from the [latest release](https://github.com/simpossible/mini_launch/releases/latest):

```bash
# Debian / Ubuntu
sudo dpkg -i mini_launch_*_linux_x86_64.deb

# RHEL / Fedora / CentOS
sudo rpm -i mini_launch_*_linux_x86_64.rpm
```

### Go install

```bash
go install github.com/simpossible/mini_launch@latest
```

### Build from source

```bash
git clone https://github.com/simpossible/mini_launch.git
cd mini_launch
go build -o mini_launch .
```

## Usage

All commands can be run from within the service directory (no need to specify a service name), or from anywhere by passing the service name as an argument.

```bash
cd $HOME/servers/myapp

mini_launch initial    # generate daemon config
mini_launch start      # start the service
mini_launch status     # check status
mini_launch restart    # restart
mini_launch stop       # stop
mini_launch remove     # stop and remove config
```

You can also operate on services by name from any directory:

```bash
mini_launch start myapp
mini_launch status myapp
mini_launch stop myapp
```

### List all services

```bash
mini_launch list
mini_launch status      # status of all services
```

## Commands

| Command | Description |
|---------|-------------|
| `initial` | Scan current directory and generate daemon config |
| `start [service]` | Start a service daemon |
| `stop [service]` | Stop a service daemon |
| `restart [service]` | Restart a service daemon |
| `status [service]` | Show service status (all if no argument) |
| `list` | List all discovered services |
| `remove [service]` | Stop and remove service configuration |

## Platform Details

### macOS (launchd)

- Config location: `<service_directory>/com.minilaunch.<service>.plist`
- Services auto-restart (`KeepAlive: true`) and start on load (`RunAtLoad: true`)
- Managed via `launchctl load/unload`

### Linux (systemd --user)

- Config location: `~/.config/systemd/user/mini-launch-<service>.service`
- Services auto-restart with 5-second delay
- Managed via `systemctl --user start/stop/restart`
- Use `loginctl enable-linger` to ensure services survive logout

## Environment Variables

`mini_launch` automatically extracts `export VAR=value` entries from your shell config file (`~/.zshrc` for zsh, `~/.bashrc` for bash) and injects them into the daemon environment. Variables that reference other variables (e.g., `$PATH`) are skipped.

## License

MIT
