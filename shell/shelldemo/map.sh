#!/bin/bash
declare -A site
site["google"]="www.google.com"
site["runoob"]="www.runoob.com"
site["taobao"]="www.taobao.com"

echo ${site["runoob"]}
echo "数组的键为: ${!site[*]}"
echo "数组的键为: ${!site[@]}"