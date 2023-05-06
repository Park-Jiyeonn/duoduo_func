#!/bin/bash

sh ./stop.sh
#kill $(pidof api)
#kill $(pidof base)
#kill $(pidof interact)
#kill $(pidof social)

# 进入第一个文件夹
cd ./cmd/api

# 在当前文件夹下执行命令，并将其放入后台运行
go run . &

# 进入第二个文件夹
cd ../base

# 在当前文件夹下执行命令，并将其放入后台运行
sh build.sh
sh output/bootstrap.sh &

# 进入第三个文件夹
cd ../interact

# 在当前文件夹下执行命令，并将其放入后台运行
sh build.sh
sh output/bootstrap.sh &

# 进入第四个文件夹
cd ../social

# 在当前文件夹下执行命令，并将其放入后台运行
sh build.sh
sh output/bootstrap.sh &

# netsh interface portproxy add v4tov4 listenport=8888 connectaddress=172.17.161.115 connectport=8888 listenaddress=0.0.0.0 protocol=tcp
# netsh interface portproxy show all