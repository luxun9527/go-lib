#!/bin/bash
#export http_proxy=http://10.18.13.129:7890
#export https_proxy=http://10.18.13.129:7890
echo 'start'
source /etc/profile
source /home/deng/.bashrc

go build -gcflags='all=-N -l' -o /home/deng/smb/go-lib/rclone/rclone_temp
nohup /home/deng/gopath/bin/dlv --listen=:2348 --headless=true --api-version=2 --accept-multiclient exec /home/deng/smb/go-lib/rclone/rclone_temp > /home/deng/smb/go-lib/rclone/test.log 2>&1 &
echo 'finish'