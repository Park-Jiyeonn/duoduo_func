#!/bin/bash
#这是一个awk命令，作用是输出行号不等于1的行中的第二个字段。

#具体来说，NR表示当前处理的行号，$2表示当前行的第二个字段，
# 'NR!=1 {print $2}'表示对于行号不等于1的行，输出其第二个字段。

#通常情况下，行号为1的行往往是表头，我们可以通过这个命令来跳过表头，只处理数据行。
# 查找占用x端口的进程并杀死
sudo lsof -i tcp:8888 | awk 'NR!=1 {print $2}' | xargs sudo kill
sudo lsof -i tcp:9000 | awk 'NR!=1 {print $2}' | xargs sudo kill
sudo lsof -i tcp:9001 | awk 'NR!=1 {print $2}' | xargs sudo kill
sudo lsof -i tcp:9002 | awk 'NR!=1 {print $2}' | xargs sudo kill

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