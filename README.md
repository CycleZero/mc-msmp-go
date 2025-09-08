
# Minecraft Server Management Protocol Client (mc-msmp-go)

[![Go Reference](https://pkg.go.dev/badge/github.com/CycleZero/mc-msmp-go.svg)](https://pkg.go.dev/github.com/CycleZero/mc-msmp-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/CycleZero/mc-msmp-go)](https://goreportcard.com/report/github.com/CycleZero/mc-msmp-go)

一个基于Go语言的Minecraft服务端管理协议(MSMP)客户端实现，用于通过WebSocket与Minecraft Java版服务端进行通信和管理。

## 功能特性

- 基于WebSocket的通信
- 实现Minecraft服务端管理协议的所有主要功能
- 支持自动重连
- 线程安全的设计
- 易于使用的API接口

## 支持的管理功能

### 白名单管理
- `allowlist` - 获取白名单列表
- `allowlist/set` - 设置白名单
- `allowlist/add` - 添加白名单玩家
- `allowlist/remove` - 移除白名单玩家
- `allowlist/clear` - 清空白名单

### 封禁玩家列表管理
- `bans` - 获取封禁玩家列表
- `bans/set` - 设置封禁玩家列表
- `bans/add` - 添加封禁玩家
- `bans/remove` - 移除封禁玩家
- `bans/clear` - 清空封禁玩家列表

### 封禁IP列表管理
- `ip_bans` - 获取封禁IP列表
- `ip_bans/set` - 设置封禁IP列表
- `ip_bans/add` - 添加封禁IP
- `ip_bans/remove` - 移除封禁IP
- `ip_bans/clear` - 清空封禁IP列表

### 玩家管理
- `players` - 获取在线玩家列表
- `players/kick` - 踢出玩家

### 管理员列表管理
- `operators` - 获取管理员列表
- `operators/set` - 设置管理员列表
- `operators/add` - 添加管理员
- `operators/remove` - 移除管理员
- `operators/clear` - 清空管理员列表

### 服务端状态管理
- `server/status` - 获取服务端状态
- `server/save` - 保存服务端数据
- `server/stop` - 停止服务端
- `server/system_message` - 发送系统消息

### 服务端设置管理
- `serversettings/*` - 获取服务端设置
- `serversettings/*/set` - 设置服务端设置

### 游戏规则管理
- `gamerules` - 获取游戏规则
- `gamerules/update` - 更新游戏规则

## 安装

```
bash
go get github.com/CycleZero/mc-msmp-go
```
## 使用示例

```
go
package main

import "github.com/CycleZero/mc-msmp-go/client"

func main() {
// 创建客户端实例
url := "ws://localhost:25585"
cli := client.NewMsmpClient(url)

    // 连接到服务端
    err := cli.Connect()
    if err != nil {
        panic(err)
    }
    defer cli.Disconnect()
    
    // 使用白名单功能
    cli.AllowlistSet([]subdto.PlayerDto{
        {Id: "player-uuid", Name: "player-name"},
    })
    
    // 获取在线玩家
    cli.Players()
}
```
## 配置要求

要使用服务端管理协议，需要在Minecraft服务端配置中进行以下设置：

1. 设置 `management-server-enabled=true` 开启管理协议
2. 配置 `management-server-host` 和 `management-server-port`（默认端口为25585）
3. 注意：管理协议未实现认证，将端口暴露在公网上具有极高风险

## 注意事项

- 本项目基于Minecraft Java版1.21.9+的服务端管理协议实现
- 确保服务端已启用管理协议功能
- 不要在生产环境中将管理协议端口暴露在公网上

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request来改进这个项目。
