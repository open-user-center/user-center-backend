# 用户中心后台

## 技术选型
Golang + Kratos + Gorm + Mysql + Redis + Grpc + Wire
## 启动前准备（推荐用容器，很方便！！！）
### 安装redis
```
docker run --name user-center-redis -p 6379:6379 -itd redis:6.0
```
### 安装mysql
```
docker pull mysql:5.7
docker run --name user-center-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -d mysql:5.7
```
### 运行以下sql脚本注入数据
```
# 建表脚本
drop database if exists user_center;
create database user_center;
use user_center;

DROP TABLE IF EXISTS user;
create table if not exists user
(
    id           bigint auto_increment comment 'id'
        primary key,
    username     varchar(256)                       null comment '用户昵称',
    userAccount  varchar(256)                       null comment '账号',
    avatarUrl    varchar(1024)                      null comment '用户头像',
    gender       tinyint                            null comment '性别',
    userPassword varchar(512)                       not null comment '密码',
    phone        varchar(128)                       null comment '电话',
    email        varchar(512)                       null comment '邮箱',
    userStatus   int      default 0                 not null comment '用户状态 0-正常',
    createTime   datetime default CURRENT_TIMESTAMP null comment '创建时间',
    updateTime   datetime default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP comment '更新时间',
    isDelete     tinyint  default 0                 not null comment '是否删除',
    role         int      default 0                 not null comment '用户角色 0-普通用户 1-管理员'
)
    comment '用户';

insert into user value(null, 'Tom', 'admin', 'http://cdn.u2.huluxia.com/g3/M00/36/56/wKgBOVwPmcmAB2cnAACcXKrjLlw989.jpg',
                       0, '25d55ad283aa400af464c76d713c07ad', null, null, 0, null, null, 0, 1);
                    
```
### 安装相应的依赖
```
make init
make tidy
```
### 项目编译
```
make build
```
### 项目运行
```
./bin/main -conf ./app/user/service/configs/config.yaml
```





