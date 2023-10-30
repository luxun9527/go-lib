#!/bin/bash
#在 testfile 文件的第四行后添加一行，并将结果输出到标准输出，在命令行提示符下输入如下命令：
 sed -e 4a\newLine testfile
 #将 testfile 的内容列出并且列印行号，同时，请将第 2~5 行删除！
 nl testfile | sed '2,5d'

 sed '1i deb http://mirrors.aliyun.com/debian/ buster main non-free contrib\ndeb-src http://mirrors.aliyun.com/debian/ buster main non-free contrib\ndeb http://mirrors.aliyun.com/debian-security buster/updates main\ndeb-src http://mirrors.aliyun.com/debian-security buster/updates main\ndeb http://mirrors.aliyun.com/debian/ buster-updates main non-free contrib\ndeb-src http://mirrors.aliyun.com/debian/ buster-updates main non-free contrib\ndeb http://mirrors.aliyun.com/debian/ buster-backports main non-free contrib\ndeb-src http://mirrors.aliyun.com/debian/ buster-backports main non-free contrib' file  #这是在第一行前添加字符串
 sed '$i 添加的内容' file  #这是在最后一行行前添加字符串
 sed '$a添加的内容' file  #这是在最后一行行后添加字符串