#!/bin/bash
for i in {4..20} ; do
   mdadm -E /dev/sd[a-z]$i >>/tmp/databack/disks.txt
  if [ $? -ne 0 ]; then
      break
  fi
 #ss
done


