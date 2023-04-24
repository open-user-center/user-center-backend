# 用户中心后台

## 技术选型
Golang + Kratos + Gorm + mysql
## 启动前准备（推荐用容器，很方便！！！）
### 安装mysql
```
docker pull mysql:5.7
docker run --name user-center-mysql -e MYSQL_ROOT_PASSWORD=123456 -d mysql:5.7
```
