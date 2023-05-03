#!/bin/bash
my_array=(A B "C" D)
my_array[1]=value1
my_array[4]=value1
echo "第一个元素为: ${my_array[0]}"
echo "第二个元素为: ${my_array[1]}"
echo "第三个元素为: ${my_array[2]}"
echo "第四个元素为: ${my_array[3]}"
echo "第四个元素为: ${my_array[4]}"
#获取所有
echo "数组的元素为: ${my_array[*]}"
echo "数组的元素为: ${my_array[@]}"
echo ${#my_array[@]}