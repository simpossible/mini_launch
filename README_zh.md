# mini_launch

[English](README.md)

一个轻量的跨平台服务守护进程管理工具，支持 macOS 和 Linux。

利用操作系统原生的服务管理机制来管理你的本地开发服务：
- **macOS**: `launchd`（通过 `launchctl`）
- **Linux**: `systemd` 用户服务（通过 `systemctl --user`）

## 为什么写这个工具？

手写 plist 或者 systemd 的 service 文件太麻烦了。`mini_launch` 帮你生成配置，一条命令搞定，不用再去啃 XML 语法。

同时也给团队定了个简单约定：

- **统一的服务目录** — 所有服务放在 `$HOME/servers/` 下，谁都能找到。
- **统一的日志位置** — 日志固定输出到服务目录下的 `std.log`，不用到处翻。
- **环境变量自动注入** — 从 `~/.zshrc` 或 `~/.bashrc` 读取 `export` 语句，守护进程和你的终端用同一套环境变量。

## 快速开始

```bash
# 1. 在 $HOME/servers/ 下创建目录，放入一个可执行文件
mkdir -p $HOME/servers/myapp
cp /path/to/myapp $HOME/servers/myapp/

# 2. 进入目录
cd $HOME/servers/myapp

# 3. 初始化并启动，搞定
mini_launch initial
mini_launch start
```

> **注意：** 每个服务目录下必须只有一个可执行文件，`mini_launch` 会自动识别。

## 工作原理

服务统一放在 `$HOME/servers/` 下，每个服务一个子目录（支持嵌套），目录里放一个可执行文件：

```
$HOME/servers/
├── myapp/
│   └── myapp            # 可执行文件
├── myserver/
│   └── myserver         # 可执行文件
└── web/
    └── api/
        └── api-server   # 可执行文件 → 服务名: web_api
```

- **服务名** = 相对于 `$HOME/servers/` 的路径，`/` 替换为 `_`
- **日志输出** → 服务目录下的 `std.log`
- **环境变量** → 自动从 shell 配置文件中提取

## 安装

### Homebrew（macOS）

```bash
brew tap simpossible/tap
brew install mini_launch
```

### 一键安装（Linux & macOS）

```bash
curl -fsSL https://raw.githubusercontent.com/simpossible/mini_launch/main/install.sh | bash
```

### DEB / RPM（Linux）

从 [最新 Release](https://github.com/simpossible/mini_launch/releases/latest) 下载对应的包：

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

### 源码编译

```bash
git clone https://github.com/simpossible/mini_launch.git
cd mini_launch
go build -o mini_launch .
```

## 使用方法

所有命令都可以在服务目录下直接运行（不需要指定服务名），也可以在任意目录通过服务名操作。

```bash
cd $HOME/servers/myapp

mini_launch initial    # 生成守护进程配置
mini_launch start      # 启动服务
mini_launch status     # 查看状态
mini_launch restart    # 重启
mini_launch stop       # 停止
mini_launch remove     # 停止并删除配置
```

也可以在任意目录通过服务名操作：

```bash
mini_launch start myapp
mini_launch status myapp
mini_launch stop myapp
```

### 查看所有服务

```bash
mini_launch list
mini_launch status      # 查看所有服务状态
```

## 命令一览

| 命令 | 说明 |
|------|------|
| `initial` | 扫描当前目录，生成守护进程配置 |
| `start [service]` | 启动服务守护进程 |
| `stop [service]` | 停止服务守护进程 |
| `restart [service]` | 重启服务守护进程 |
| `status [service]` | 查看服务状态（无参数则显示全部） |
| `list` | 列出所有已发现的服务 |
| `remove [service]` | 停止并移除服务配置 |

## 平台细节

### macOS (launchd)

- 配置位置：`<服务目录>/com.minilaunch.<服务名>.plist`
- 自动重启（`KeepAlive: true`）且加载即启动（`RunAtLoad: true`）
- 通过 `launchctl load/unload` 管理

### Linux (systemd --user)

- 配置位置：`~/.config/systemd/user/mini-launch-<服务名>.service`
- 崩溃后延迟 5 秒自动重启
- 通过 `systemctl --user start/stop/restart` 管理
- 使用 `loginctl enable-linger` 可确保注销后服务继续运行

## 环境变量

`mini_launch` 会自动从 shell 配置文件（zsh 用 `~/.zshrc`，bash 用 `~/.bashrc`）中提取 `export VAR=value` 形式的环境变量，注入到守护进程中。引用了其他变量的（如 `$PATH`）会被跳过。

## License

MIT
