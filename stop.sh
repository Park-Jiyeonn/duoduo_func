#这是一个awk命令，作用是输出行号不等于1的行中的第二个字段。

#具体来说，NR表示当前处理的行号，$2表示当前行的第二个字段，
# 'NR!=1 {print $2}'表示对于行号不等于1的行，输出其第二个字段。

#通常情况下，行号为1的行往往是表头，我们可以通过这个命令来跳过表头，只处理数据行。
# 查找占用x端口的进程并杀死
sudo lsof -i tcp:8888 | awk 'NR!=1 {print $2}' | xargs sudo kill
sudo lsof -i tcp:10086 | awk 'NR!=1 {print $2}' | xargs sudo kill
sudo lsof -i tcp:10087 | awk 'NR!=1 {print $2}' | xargs sudo kill
sudo lsof -i tcp:10088 | awk 'NR!=1 {print $2}' | xargs sudo kill
