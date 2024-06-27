# tor-control协议服务端模拟工具

本项目是一个使用 Go 实现的 Tor 控制协议模拟器。它模拟了一个 Tor 控制协议服务器，使您可以通过模拟服务器-客户端交互来测试扫描器或客户端工具。模拟器支持通过配置文件动态配置认证方法和 Tor 版本

## 功能介绍

- 模拟 Tor 控制协议服务器
- 通过命令行或配置文件进行配置
- 记录客户端交互日志（连接握手、身份验证等）
- 日志以 [Tor-Control] 为前缀
- 支持各种扫描工具（例如 nmap、fofa）
- 通过配置文件动态配置 AUTH METHODS 和 VERSION

## 使用方法

### 1. 安装依赖
首先，安装所需的 Go 库
- go get -u github.com/sirupsen/logrus

### 2. 对配置文件进行配置
```go
{
    "port": "9051",
    "password": "your_password",
    "logfile": "tor_control.log",
    "auth_methods": "HASHEDPASSWORD",
    "version": "0.4.6.5"
}
```
#### 参数说明
- port：服务器监听的端口。
- password：身份验证密码。
- logfile：日志文件名称。
- auth_methods：服务器支持的认证方法。
- version：服务器报告的 Tor 版本。
### 3.运行模拟器
```go 
go run main.go -config config.json
```


## 效果展示

模拟器运行后，将在指定端口监听传入连接。您可以使用各种工具和命令来测试模拟器。以下是一些示例交互：
### 1.身份验证
`AUTHENTICATE "your_password"`
### 响应：
`250 OK`
### 2.获取协议信息
`PROTOCOLINFO 1`
### 响应：
```
250-PROTOCOLINFO 1
250-AUTH METHODS=HASHEDPASSWORD
250-VERSION Tor="0.4.6.5"
250 OK
```
### 3.获取版本信息
`GETINFO version`
### 响应：
``` 
250-version=0.4.6.5
250 OK
```
### 4.重新加载信号
`SIGNAL RELOAD`
### 响应：
``` 
250 OK
```

## 日志
模拟器记录所有客户端交互。日志以 [Tor-Control] 为前缀，并包含关于接收命令和发送响应的详细信息。日志将写入配置文件中指定的文件（例如示例中的 tor_control.log）
### 日志输出示例：
```
[Tor-Control] Server started on port 9051
[Tor-Control] Client connected: 127.0.0.1:56789
[Tor-Control] Received: AUTHENTICATE "your_password"
[Tor-Control] Authentication successful
[Tor-Control] Received: PROTOCOLINFO 1
[Tor-Control] Sent protocol info
[Tor-Control] Received: GETINFO version
[Tor-Control] Sent version info
[Tor-Control] Received: SIGNAL RELOAD
[Tor-Control] Signal RELOAD received
[Tor-Control] Client disconnected
```
## 贡献
如果您发现任何问题或有改进建议，请随时提交问题或拉取请求。

## 作者
sunhanfei

