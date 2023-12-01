.PHONY: proto
proto:
	protoc   -I./pb -I./pb/googleapis  --grpc-gateway_out=./pb \
			--openapiv2_out=. \
			--openapiv2_opt json_names_for_fields=false\
			--openapiv2_opt generate_unbound_methods=true \
			--openapiv2_opt output_format=yaml \
			--grpc-gateway_opt generate_unbound_methods=true \
			--grpc-gateway_opt logtostderr=true \
			  --go_out=./pb --go-grpc_out=./pb grpcdemo/grpcdemo.proto
	#protoc   -I./pb --go_out=./pb --openapiv2_out=.  grpcdemo/custom.proto
	#protoc   -I./pb --go_out=../../../ --openapiv2_out=. grpcdemo/folder/imported.proto


   #当要生成的pb文件有复杂的导入， 1.定义-I+grpcdemo/folder/imported.proto 指定要生成的proto文件的位置。