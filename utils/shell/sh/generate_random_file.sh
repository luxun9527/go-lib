#!/bin/bash
count=$1
filePath=$2
for((i=1;i<=$count;i++));
do
uuid=$(cat /proc/sys/kernel/random/uuid)
p="$filePath/$uuid"
echo $uuid > $p > /dev/null 2>&1
done