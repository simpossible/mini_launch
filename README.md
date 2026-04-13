# mini_launch

A simple cross-platform service daemon manager for macOS and Linux.

`mini_launch` manages long-running services using native platform tools:
- **macOS**: `launchd` via `launchctl`
- **Linux**: `systemd` user services via `systemctl --user`

## How It Works

Services are organized under `$HOME/servers/`. Each service resides in its own subdirectory (nested directories are supported) with a single executable file:

```
$HOME/servers/
├── mysql/
│   └── mysqld          # executable
├── redis/
│   └── redis-server    # executable
└── web/
    └── api/
        └── myapp       # executable → service name: web_api
```

- **Service name** = relative path from `$HOME/servers/` with `/` replaced by `_`
- **Log output** → `std.log` in the service directory
- **Environment variables** are automatically extracted from your shell config (`~/.zshrc` or `~/.bashrc`)

## Install

```bash
go install github.com/simpossible/mini_launch@latest
```

Or build from source:

```bash
git clone https://github.com/simpossible/mini_launch.git
cd mini_launch
go build -o mini_launch .
```

## Usage

### Initialize a service

From the service directory:

```bash
cd $HOME/servers/mysql
mini_launch initial
```

Output:
```
Service name: mysql
Directory:    /Users/you/servers/mysql
Executable:   /Users/you/servers/mysql/mysqld
Log file:     /Users/you/servers/mysql/std.log
Generated: /Users/you/Library/LaunchAgents/com.minilaunch.mysql.plist

Service 'mysql' initialized successfully.
Use 'mini_launch start mysql' to start the service.
```

### Start a service

```bash
# By service name
mini_launch start mysql

# Or from the service directory
cd $HOME/servers/mysql
mini_launch start
```

### Stop a service

```bash
mini_launch stop mysql
```

### Restart a service

```bash
mini_launch restart mysql
```

### Check status

```bash
# All services
mini_launch status

# Specific service
mini_launch status mysql
```

### List services

```bash
mini_launch list
```

### Remove a service

Stops the service and removes the daemon configuration:

```bash
mini_launch remove mysql
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

- Config location: `~/Library/LaunchAgents/com.minilaunch.<service>.plist`
- Services auto-restart (`KeepAlive: true`) and start on login (`RunAtLoad: true`)
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
