#!/bin/bash

for i in `ls`;
do
echo $i is file name\! ;
done

list="rootfs
usr
 data
  data2"
for i in $list;
do
echo $i is appoint ;
done


sum=0
for ((i=1; i<=100; i++))
do
    ((sum += i))
done
echo "The sum is: $sum"

