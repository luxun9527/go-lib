.PHONY:genmodel
genmodel:
	gentool --dsn="root:root@tcp(192.168.2.200:33606)/test?charset=utf8mb4&parseTime=True&loc=Local" --db=mysql --tables=card \
   -outPath=gen/query -fieldMap="decimal:string;tinyint:int8;int:int32;bigint:int64;smallint:int16" -fieldSignable=true
