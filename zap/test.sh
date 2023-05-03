#!/bin/sh
# 执行脚本后重命令防止其他程序执行disable
echo "rckwake唤醒时间$(rtcwake -m show)"
hwclock -w
rtcwake -l -m no -s 300
echo "rckwake唤醒时间 /proc/driver/rtc文件 $(cat /proc/driver/rtc)"
echo "rckwake唤醒时间$(rtcwake -m show)"
poweroff