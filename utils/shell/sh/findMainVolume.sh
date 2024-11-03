#!/bin/bash
#删除更新包
#https://www.cnblogs.com/huchong/p/10069521.html
volumes=$(df-json |awk '{print $8}'|  egrep '^/Volume[1-9]{1,}$')
for v in $volumes;
do
#找主卷
p="$v/.main.inc"
if [  -e $p ];then
  rm -rf "$v/.update_6cf63d40d9af5463182da96e489889c1"
  rm -rf "$v/.deCompressDir_0b812568da4ba30bc002ace26544dc24"
fi
done
