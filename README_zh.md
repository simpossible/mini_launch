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

## 工作原理

服务统一放在 `$HOME/servers/` 下，每个服务一个子目录（支持嵌套），目录里放一个可执行文件：

```
$HOME/servers/
├── mysql/
│   └── mysqld          # 可执行文件
├── redis/
│   └── redis-server    # 可执行文件
└── web/
    └── api/
        └── myapp       # 可执行文件 → 服务名: web_api
```

- **服务名** = 相对于 `$HOME/servers/` 的路径，`/` 替换为 `_`
- **日志输出** → 服务目录下的 `std.log`
- **环境变量** → 自动从 shell 配置文件中提取

## 安装

### Homebrew（推荐）

```bash
brew tap simpossible/tap
brew install mini_launch
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

### 初始化服务

在服务目录下执行：

```bash
cd $HOME/servers/mysql
mini_launch initial
```

输出：
```
Service name: mysql
Directory:    /Users/you/servers/mysql
Executable:   /Users/you/servers/mysql/mysqld
Log file:     /Users/you/servers/mysql/std.log
Generated: /Users/you/servers/mysql/com.minilaunch.mysql.plist

Service 'mysql' initialized successfully.
Use 'mini_launch start mysql' to start the service.
```

### 启动服务

```bash
# 指定服务名
mini_launch start mysql

# 或者在服务目录下直接运行
cd $HOME/servers/mysql
mini_launch start
```

### 停止服务

```bash
mini_launch stop mysql
```

### 重启服务

```bash
mini_launch restart mysql
```

### 查看状态

```bash
# 查看所有服务
mini_launch status

# 查看指定服务
mini_launch status mysql
```

### 列出所有服务

```bash
mini_launch list
```

### 移除服务

停止服务并删除守护进程配置：

```bash
mini_launch remove mysql
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
