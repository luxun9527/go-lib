#!/bin/bash

paths=$*
for i in $paths ; do
    path=$i
    echo $path
    btrfs property set -ts $path ro false
    chattr -i $path
    base=$(echo $(pwd) |awk  -F/ '{print $2}')
    base="/$base/"
    arr=()
    index=0
    function foo() {
      for j in $1;
      do
        arr[index]=$j
        index=`expr $index + 1`;
        p1=$(btrfs subvolume list $j -o |awk -v base=$base  '{print base$9}'  )
    #    if [ ! "$p1" ]; then
    #      echo 'finish'
    #        return
    #     fi
        foo "$p1"
      done
    }
    foo $path
    len=${#arr[@]}
    for ((k=$len-1; k>=0; k--))
    do
       btrfs subvolume delete ${arr[$k]}
       #btrfs subvolume sync  ${arr[$k]}
    done
done
