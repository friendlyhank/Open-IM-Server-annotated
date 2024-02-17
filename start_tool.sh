#!/usr/bin/env bash

cd ~/software/

# 启动redis服务
nohup ./redis/bin/redis-server ./redis/bin/redis.conf &
# 启动zookeeper
nohup ./kafka/bin/zookeeper-server-start.sh ./kafka/config/zookeeper.properties &
