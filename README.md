# Go App

一个基于 Gin 框架的用户管理系统。

## 功能特性

- 用户注册
- 邮箱登录
- JWT 认证
- 用户管理（CRUD）
- Swagger API 文档

## 技术栈

- Go 1.21
- Gin Web Framework
- GORM
- MySQL
- JWT
- Swagger

## 快速开始

### 前置要求

- Go 1.21 或更高版本
- MySQL 8.0
- 设置好 GOPATH 和 GOROOT

### 安装

1. 克隆项目
```bash
git clone https://github.com/yourusername/go-app.git
cd go-app

2.安装依赖
```bash
go mod tidy

3.运行项目
```bash
go run main.go

4.运行项目须注意⚠️
- 查找所有运行中的 Go 进程：
```bash
ps aux | grep go

- 找到并终止所有相关的 Go 进程
```bash
pkill -f "go"

- 如果想要更精确地只终止特定的项目进程：
```bash
pkill -f "go_app"

- 如果上述命令不能完全终止进程，可以使用更强制的方式
```bash
pkill -9 -f "go_app"

5.重新启动项目
```bash
go run main.go

> 服务将在 http://localhost:8080 启动


## 数据库配置

### 连接信息
dsn := "root:123456@tcp(127.0.0.1:3306)/go_app?charset=utf8mb4&parseTime=True&loc=Local"

### 默认配置：

- 主机：localhost (127.0.0.1)
- 端口：3306
- 用户名：root
- 密码：123456
- 数据库：go_app
- 字符集：utf8mb4


## 接口文档
> 访问 Swagger 文档： http://localhost:8080/swagger/index.html
