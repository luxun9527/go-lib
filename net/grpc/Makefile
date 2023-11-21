.PHONY: proto
proto:
	protoc   -I./pb -I./pb/googleapis --go_out=./pb --go-grpc_out=./pb grpcdemo/grpcdemo.proto
	protoc   -I./pb --go_out=./pb  grpcdemo/custom.proto
	protoc   -I./pb --go_out=../../../  grpcdemo/folder/imported.proto


   #当要生成的pb文件有复杂的导入， 1.定义-I+grpcdemo/folder/imported.proto 指定要生成的proto文件的位置。