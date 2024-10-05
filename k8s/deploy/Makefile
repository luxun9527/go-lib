buildDocker:
	docker build -t 192.168.2.200:5000/demosrv:latest .
push:
	 docker push 192.168.2.200:5000/demosrv:latest
apply:
	kubectl apply -f deployment
delete:
	kubectl delete -f deployment
	docker rmi 192.168.2.200:5000/demosvr:latest
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -buildvcs=false -o k8sdemo
dockerRun:
	docker run -p 8089:8089 --name demosvr -it 192.168.2.200:5000/demosvr:latest
start:
	make build
	make buildDocker
	make push
	make apply