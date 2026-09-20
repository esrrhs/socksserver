# socksserver

[<img src="https://img.shields.io/github/license/esrrhs/socksserver">](https://github.com/esrrhs/socksserver)
[<img src="https://img.shields.io/github/languages/top/esrrhs/socksserver">](https://github.com/esrrhs/socksserver)
[![Go Report Card](https://goreportcard.com/badge/github.com/esrrhs/socksserver)](https://goreportcard.com/report/github.com/esrrhs/socksserver)
[<img src="https://img.shields.io/github/v/release/esrrhs/socksserver">](https://github.com/esrrhs/socksserver/releases)
[<img src="https://img.shields.io/github/downloads/esrrhs/socksserver/total">](https://github.com/esrrhs/socksserver/releases)
[<img src="https://img.shields.io/docker/pulls/esrrhs/socksserver">](https://hub.docker.com/repository/docker/esrrhs/socksserver)
[<img src="https://img.shields.io/github/actions/workflow/status/esrrhs/socksserver/go.yml?branch=master">](https://github.com/esrrhs/socksserver/actions)

轻量、高性能的 SOCKS5 服务器（支持认证与优雅停止）。

## 特性

- **轻量高效**：基于 Go 标准网络库与零拷贝传输机制
- **认证支持**：支持无认证模式及基于用户名/密码的认证
- **安全健壮**：连接超时保护，连接异常自动回收，避免资源泄露
- **平滑停机**：捕获系统信号（`SIGINT`/`SIGTERM`），支持优雅退出
- **多平台跨平台打包**：一键编译打包主流系统与架构（Linux/macOS/Windows/ARM等）
- **极简容器镜像**：基于 Alpine 多阶段构建，镜像大小仅约 25MB

## 快速使用

### 本地运行

无需认证：
```bash
./socksserver -l :4455
```

启用用户名与密码认证：
```bash
./socksserver -l :1080 -u myuser -p mypassword
```

### Docker 运行

```bash
docker run --name socksserver -d --privileged --network host --restart=always esrrhs/socksserver ./socksserver -u yourusername -p yourpassword -l :1080
```

也可以直接传递参数：
```bash
docker run --name socksserver -d --net=host --restart=always esrrhs/socksserver -u yourusername -p yourpassword -l :1080
```

## 参数说明

| 参数 | 默认值 | 说明 |
| :--- | :--- | :--- |
| `-l` | (必填) | 监听地址与端口，例如 `:1080` 或 `0.0.0.0:1080` |
| `-u` | 空 | 认证用户名（不填则无需认证） |
| `-p` | 空 | 认证密码 |
| `-loglevel` | `info` | 日志级别 (`debug`, `info`, `warn`, `error`) |
| `-nolog` | `0` | 是否关闭日志文件输出（`1` 为不写日志文件） |
| `-noprint` | `0` | 是否禁止控制台输出（`1` 为不打印到 stdout） |
| `-v` | `false` | 显示版本信息与构建时间 |

## 编译与测试

本项目提供了便捷的 `Makefile`：

```bash
# 编译二进制
make build

# 运行单元测试与竞态检测
make test
make race

# 构建 Docker 镜像
make docker

# 跨平台打包
make pack-fast   # 快速打包主流平台 (linux/darwin/windows, amd64/arm64)
make pack        # 打包所有支持的架构
```
