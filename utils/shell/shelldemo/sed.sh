#!/bin/bash
#在 testfile 文件的第四行后添加一行，并将结果输出到标准输出，在命令行提示符下输入如下命令：
 sed -e 4a\newLine testfile
 #将 testfile 的内容列出并且列印行号，同时，请将第 2~5 行删除！
 nl testfile | sed '2,5d'