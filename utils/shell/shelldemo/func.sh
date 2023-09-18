#!/bin/bash
demoFun(){
    echo "这是我的第一个 shell 函数!"
}
echo "-----函数开始执行-----"
demoFun
echo "-----函数执行完毕-----"

funWithReturn(){
    echo "这个函数会对输入的两个数字进行相加运算..."
    echo "输入第一个数字: "
    read aNum
    echo "输入第二个数字: "
    read anotherNum
    echo "两个数字分别为 $aNum 和 $anotherNum !"
    return $(($aNum+$anotherNum))
}
funWithReturn
echo "输入的两个数字之和为 $? !"

# 返回$1指定的git项目的当前分支(branch)或标签名(tag)
# $1 git项目源码位置,为空获则默认为当前文件夹
function current_branch () {
    local folder="$(pwd)"
    [ -n "$1" ] && folder="$1"
    git -C "$folder" rev-parse --abbrev-ref HEAD | grep -v HEAD || \
    git -C "$folder" describe --tags HEAD || \
    git -C "$folder" rev-parse HEAD
}
current_branch