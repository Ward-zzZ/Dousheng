#!/bin/bash

# 编译所有服务
cd ./cmd/user
sh build.sh &
cd ../relation
sh build.sh &
cd ../favorite
sh build.sh &
cd ../comment
sh build.sh &
cd ../video
sh build.sh &
cd ../message
sh build.sh &
cd ../api
sh build.sh &

# 等待所有服务编译完成
wait

# 加载配置
cd ../..
go run "./data/configs/consul/uploadKV.go"

# 启动所有服务
cd ./cmd/user
sh output/bootstrap.sh &
cd ../relation
sh output/bootstrap.sh &
cd ../favorite
sh output/bootstrap.sh &
cd ../comment
sh output/bootstrap.sh &
cd ../video
sh output/bootstrap.sh &
cd ../message
sh output/bootstrap.sh &
cd ../api
sleep 2
echo "api service start success"
./output/bin/api
