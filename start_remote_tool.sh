#!/usr/bin/env bash

# 需要替换为自己指定软件路径
cd ~/software/

# 启动redis服务
nohup ./redis/bin/redis-server ./redis/bin/redis.conf &
# 启动zookeeper
nohup ./kafka/bin/zookeeper-server-start.sh ./kafka/config/zookeeper.properties &

# 启动mongodb
mongod --dbpath /usr/local/mongodb/data --logpath /usr/local/mongodb/log/mongo.log --fork

## 启动服务

# 替换成自己的路径
cd ~/go/src/github.com/OpenIMSDK/Open-IM-Server/cmd/

# 先杀死进程
killall openimapi
killall openimrpcuser

# 构建服务
# 启动api服务
go build -o openimapi openim-api/main.go
# 构建rpc服务
go build -o openimrpcuser ./openim-rpc/openim-rpc-user/main.go

#后台构建运行服务
nohup ./openimapi &
nohup ./openimrpcuser &


