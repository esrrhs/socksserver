# socksserver

[<img src="https://img.shields.io/github/license/esrrhs/socksserver">](https://github.com/esrrhs/socksserver)
[<img src="https://img.shields.io/github/languages/top/esrrhs/socksserver">](https://github.com/esrrhs/socksserver)
[<img src="https://img.shields.io/github/v/release/esrrhs/socksserver">](https://github.com/esrrhs/socksserver/releases)
[<img src="https://img.shields.io/github/downloads/esrrhs/socksserver/total">](https://github.com/esrrhs/socksserver/releases)
[<img src="https://img.shields.io/docker/pulls/esrrhs/socksserver">](https://hub.docker.com/repository/docker/esrrhs/socksserver)
[<img src="https://img.shields.io/github/actions/workflow/status/esrrhs/socksserver/go.yml?branch=master">](https://github.com/esrrhs/socksserver/actions)

A lightweight, high-performance SOCKS5 proxy server with authentication and graceful shutdown support.

## Features

- **Lightweight & High Performance**: Powered by Go standard networking with zero-copy stream splicing.
- **Authentication**: Supports both unauthenticated mode and RFC 1929 username/password authentication.
- **Robust & Resilient**: Dial timeout protection and automatic connection resource reclamation to avoid leaks.
- **Graceful Shutdown**: Captures system signals (`SIGINT`/`SIGTERM`) for smooth teardown.
- **Cross-Platform Packaging**: One-command cross-compilation for major OSes and architectures (Linux, macOS, Windows, ARM, etc.).
- **Minimal Docker Image**: Multi-stage Alpine build resulting in a compact ~25 MB image.

## Quick Start

### Local Execution

Without authentication:
```bash
./socksserver -l :4455
```

With username and password authentication:
```bash
./socksserver -l :1080 -u myuser -p mypassword
```

### Docker

```bash
docker run --name socksserver -d --privileged --network host --restart=always esrrhs/socksserver ./socksserver -u yourusername -p yourpassword -l :1080
```

Or pass arguments directly:
```bash
docker run --name socksserver -d --net=host --restart=always esrrhs/socksserver -u yourusername -p yourpassword -l :1080
```

## Command Line Options

| Flag | Default | Description |
| :--- | :--- | :--- |
| `-l` | (Required) | Listen address and port, e.g. `:1080` or `0.0.0.0:1080` |
| `-u` | `""` | Username for authentication (empty for no auth) |
| `-p` | `""` | Password for authentication |
| `-loglevel` | `info` | Log level (`debug`, `info`, `warn`, `error`) |
| `-nolog` | `0` | Disable log file output (`1` to disable writing to log files) |
| `-noprint` | `0` | Disable stdout printing (`1` to disable printing to stdout) |
| `-v` | `false` | Show version and build date |

## Build & Test

A handy `Makefile` is provided for standard workflows:

```bash
# Build binary
make build

# Run tests and race detector
make test
make race

# Build Docker image
make docker

# Cross-platform packaging
make pack-fast   # Quickly pack major platforms (linux/darwin/windows, amd64/arm64)
make pack        # Pack all supported architectures
```
