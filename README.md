# Go App

一个基于 Gin 框架的用户相关应用。

## 功能特性

- 用户注册
- 邮箱登录
- 密码修改
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

## 目录结构
- `config`: 配置文件
- `controllers`: 控制器
- `models`: 数据模型
- `routes`: 路由配置
- `utils`: 工具函数
- `main.go`: 入口文件   
- `go.mod`: Go 模块配置文件
- `go.sum`: Go 依赖的哈希值 
- `README.md`: 项目说明文件

## 接口文档
> 访问 Swagger 文档：
```swagger
URL_ADDRESS:8080/swagger/index.html
```

## 快速开始

### 前置要求

- Go 1.21 或更高版本
- MySQL 8.0
- 设置好 GOPATH 和 GOROOT


### 安装

1. 克隆项目
```bash
git clone https://github.com/developerjet/go-app.git
```

2. 安装依赖
```bash
cd go-app
go mod download
```

3. 配置数据库
```
1. 确保 MySQL 服务已启动
2. 创建一个名为 go_app 的数据库
3. 在 config/config.go 文件中配置数据库连接信息
```
  
4. 安装依赖         
```bash
go mod download
go mod tidy
```

### 运行项目 
1. 运行项目
```bash
go run main.go
```

2. 运行项目须注意⚠️
- 查找所有运行中的 Go 进程：
```bash
ps aux | grep go
```

- 找到并终止所有相关的 Go 进程
```bash
pkill -f "go"
```

- 如果想要更精确地只终止特定的项目进程：
```bash
pkill -f "go_app"
```

- 如果上述命令不能完全终止进程，可以使用更强制的方式
```bash
pkill -9 -f "go_app"
```

3. 重新启动项目
```bash
cd /Users/edy/Documents/Github/go-app/go_app
go run main.go
```

> 服务将在 http://localhost:8080 启动


## 数据库配置

### 连接信息
在 config/config.go 文件中配置数据库连接信息。
```go
dsn := "root:123456@tcp(127.0.0.1:3306)/go_app?charset=utf8mb4&parseTime=True&loc=Local"
```

### 默认配置：
- 主机：localhost (127.0.0.1)
- 端口：3306
- 用户名：root
- 密码：123456
- 数据库：go_app
- 字符集：utf8mb4

## 数据库使用指南

### MySQL 服务管理
1. 启动 MySQL 服务
```bash
brew services start mysql@8.0
```

2. 停止 MySQL 服务
```bash
brew services stop mysql@8.0
```

3. 查看服务状态
```bash
brew services list | grep mysql
```

### 数据库操作
1. 连接数据库
```bash
mysql -u root -p
```
2. 查看数据库列表
```sql
SHOW DATABASES;
```
3. 切换数据库
```sql
USE go_app;
```
4. 查看表列表
```sql
SHOW TABLES;
```
5. 查看表结构
```sql
DESCRIBE users;
```
6. 查询数据
```sql
SELECT * FROM users;
```
7. 插入数据
```sql
INSERT INTO users (name, email, password) VALUES ('John Doe', '
```sql
SELECT * FROM users;
```
8. 插入数据
```sql
INSERT INTO users (name, email, password) VALUES ('John Doe', 'EMAIL', 'password123');
```
9. 更新数据
```sql
UPDATE users SET name = 'Jane Doe' WHERE id = 1;
```
10. 删除数据
```sql
DELETE FROM users WHERE id = 1;
```                 


