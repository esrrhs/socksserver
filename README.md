# socksserver

[<img src="https://img.shields.io/github/license/esrrhs/socksserver">](https://github.com/esrrhs/socksserver)
[<img src="https://img.shields.io/github/languages/top/esrrhs/socksserver">](https://github.com/esrrhs/socksserver)
[![Go Report Card](https://goreportcard.com/badge/github.com/esrrhs/socksserver)](https://goreportcard.com/report/github.com/esrrhs/socksserver)
[<img src="https://img.shields.io/github/v/release/esrrhs/socksserver">](https://github.com/esrrhs/socksserver/releases)
[<img src="https://img.shields.io/github/downloads/esrrhs/socksserver/total">](https://github.com/esrrhs/socksserver/releases)
[<img src="https://img.shields.io/docker/pulls/esrrhs/socksserver">](https://hub.docker.com/repository/docker/esrrhs/socksserver)
[<img src="https://img.shields.io/github/workflow/status/esrrhs/socksserver/Go">](https://github.com/esrrhs/socksserver/actions)

简单的socks5服务器
```
./socksserver -l :4455
```
docker
```
docker run --name socksserver -d --privileged --network host --restart=always esrrhs/socksserver ./socksserver -u yourusername -p yourpassword -l :1080
```
