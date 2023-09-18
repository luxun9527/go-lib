#!/bin/bash


# 返回$1指定的git项目的当前分支(branch)或标签名(tag)
# $1 git项目源码位置,为空获则默认为当前文件夹
function current_branch () {
    local folder="$(pwd)"
    [ -n "$1" ] && folder="$1"
    git -C "$folder" describe --tags HEAD || \
    git -C "$folder" rev-parse --abbrev-ref HEAD | grep -v HEAD || \
    git -C "$folder" rev-parse HEAD
}
current_branch