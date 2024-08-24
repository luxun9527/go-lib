docker:
	docker build -t 192.168.2.200:5000/grpcdemoclient:latest -f Dockerfile-client .
	docker build -t 192.168.2.200:5000/grpcdemoserver:latest -f Dockerfile-server .
	docker push 192.168.2.200:5000/grpcdemoserver:latest
	docker push 192.168.2.200:5000/grpcdemoclient:latest
	kubectl apply -f client.yaml
	kubectl apply -f server.yaml
apply:
	kubectl apply -f deployment
delete:
	#kubectl delete -f deployment
	kubectl delete -f client.yaml
	kubectl delete -f server.yaml
	docker rm -f  grpcdemoserver
	docker rm -f  grpcdemoclient
	docker rmi -f 192.168.2.200:5000/grpcdemoclient:latest
	docker rmi -f 192.168.2.200:5000/grpcdemoserver:latest
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -buildvcs=false -o grpcdemo
dockerRun:
	docker run -p 8089:8089 --name grpcdemoserver -it 192.168.2.200:5000/grpcdemoserver:latest
start:
	make build
	make buildDocker
	make push
	make apply