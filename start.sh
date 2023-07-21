#!/bin/bash

# 编译所有服务
cd ./cmd/user
sh build.sh
cd ../relation
sh build.sh
cd ../favorite
sh build.sh
cd ../comment
sh build.sh
cd ../video
sh build.sh
cd ../api
sh build.sh

# 启动所有服务
cd ../..
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
cd ../api
./output/bin/api
