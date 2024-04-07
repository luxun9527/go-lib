# go程序-->docker私有仓库镜像-->k8s部署

https://www.ywbj.cc/?p=671

https://blog.51cto.com/u_1264026/7552228

https://www.cnblogs.com/guangdelw/p/16967841.html

https://blog.csdn.net/qq_43652666/article/details/132929348

https://blog.csdn.net/MssGuo/article/details/128149704

https://znunwm.top/archives/k8s-xiang-xi-jiao-cheng

简单实现go程序打包成docker镜像，推送到私有仓库，再到k8s部署，了解一下在k8s在我们的开发生产中的使用，主要是熟悉一下流程。

如果遇到问题欢迎留言，您的star就是我更新的动力 本文地址[go程序-docker-k8s](https://github.com/luxun9527/go-lib/tree/master/utils/k8s)

## 安装环境

**ubantu20，k8s版本v1.24.0，docker版本26.0.0 ，cri-dockerd版本0.3.1.3-0.ubuntu-focal，_amd64** **虚拟机网络桥接**

**master 192.168.2.199** 

**node    192.168.2.200**

#### 安装ks8预备工作

所有机器都执行,所有机器都要有root权限。

```shell
# 1、关闭防火墙
#ufw查看当前的防火墙状态：inactive状态是防火墙关闭状态 active是开启状态
ufw status
#启动、关闭防火墙
ufw disable

# 2、禁用selinux
#默认ubunt默认是不安装selinux的，如果没有selinux命令和配置文件则说明没有安装selinux，则下面步骤就不用做了
sed -ri 's/SELINUX=enforcing/SELINUX=disabled/g' /etc/selinux/config 
setenforce 0

#3、关闭swap分区（必须，因为k8s官网要求）
#注意：最好是安装虚拟机时就不要创建swap交换分区**
sed -ri 's/.*swap.*/#&/' /etc/fstab
swapoff -a



# 4、设置主机名
cat >> /etc/hosts <<EOF
192.168.2.199 master
192.168.2.200 node
EOF
cat >> /etc/hosts <<EOF
192.168.2.199 master
192.168.2.200 node
EOF
#199master上执行
vim /etc/hostname
master
#200node上执行
vim /etc/hostname
node

# 5、时间同步
#查看时区，时间
date
#先查看时区是否正常，不正确则替换为上海时区
timedatectl set-timezone Asia/Shanghai
#安装chrony，联网同步时间
apt install chrony -y && systemctl enable --now chronyd

# 6、将桥接的IPv4流量传递到iptables的链
#（有一些ipv4的流量不能走iptables链，因为linux内核的一个过滤器，每个流量都会经过他，然后再匹配是否可进入当前应用进程去处理，所以会导致流量丢失），配置k8s.conf文件（k8s.conf文件原来不存在，需要自己创建的）

touch /etc/sysctl.d/k8s.conf
cat >> /etc/sysctl.d/k8s.conf <<EOF
net.bridge.bridge-nf-call-ip6tables=1
net.bridge.bridge-nf-call-iptables=1
net.ipv4.ip_forward=1
vm.swappiness=0
EOF
sysctl -p
sysctl --system

# 7、设置服务器之间免密登陆(2台彼此之间均设置)
ssh-keygen -t rsa
ssh-copy-id -i /root/.ssh/id_rsa.pub root@192.168.2.199
ssh-copy-id -i /root/.ssh/id_rsa.pub root@192.168.2.200
ssh node1
ssh node2

# 8执行
modprobe  br_netfilter
#让配置生效
sysctl -p
```

#### 安装docker

https://cloud.tencent.com/developer/article/2309562

所有机器都要执行

```shell
apt-get remove docker docker-engine docker.io containerd runc

apt-get install ca-certificates curl gnupg lsb-release

curl -fsSL http://mirrors.aliyun.com/docker-ce/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] http://mirrors.aliyun.com/docker-ce/linux/ubuntu $(lsb_release -cs) stable"
apt-get install docker-ce docker-ce-cli containerd.io
#修改docker 配置
vim /etc/docker/daemon.json
{
  "registry-mirrors": ["https://v9nqzd2l.mirror.aliyuncs.com"], #镜像代理
    "exec-opts": ["native.cgroupdriver=systemd"], #指定cgroupdriver
    "insecure-registries": ["192.168.2.200:5000"] # 解决私有仓库走https的问题
}

{
  "registry-mirrors": ["https://v9nqzd2l.mirror.aliyuncs.com"],
    "exec-opts": ["native.cgroupdriver=systemd"], 
    "insecure-registries": ["192.168.2.200:5000"] 
}
#重启docker
systemctl restart docker
```



#### docker私有仓库

https://blog.csdn.net/weixin_38251332/article/details/129261314

https://blog.csdn.net/wqadxmm/article/details/127648237

```shell
mkdir -p /securitit/registry/certs/ ; \
mkdir -p /securitit/registry/auth/ ;\
mkdir -p /securitit/registry/conf/ ;\
mkdir -p /securitit/registry/db/ ;\
mkdir -p  /securitit/registry/data/registry/;
#执行这一句会报错，按照提示安装对于的软件即可。
htpasswd -Bbn admin admin  > /securitit/registry/auth/htpasswd
openssl req -new -newkey rsa:4096 -days 365 -subj "/CN=localhost" -nodes -x509 -keyout /securitit/registry/auth/auth.key -out /securitit/registry/auth/auth.cert
```



/securitit/registry/conf/registry-srv.yml

```yaml
version: 0.1    
log:
  fields:
    service: registry
storage:
  delete:
    enabled: true
  cache:
    blobdescriptor: inmemory
  filesystem:
    rootdirectory: /var/lib/registry
 
http:
  addr: 0.0.0.0:5000   
  headers:
    X-Content-Type-Options: [nosniff]
health:
  storagedriver:
    enabled: true
    interval: 10s
threshold: 3
auth:
  token:
    # registry-web的地址.
    realm: http://192.168.2.200:5050/api/auth
    # 私有仓库的配置地址.
    service: 192.168.2.200:5000
    # 需要与registry-web定义的名称一致.
    issuer: 'my issuer'
    # 容器内证书路径，容器启动时通过数据卷参数指定.
    rootcertbundle: /etc/docker/registry/auth.cert
```

/securitit/registry/conf/registry-web.yml

```yaml
registry:
  # 私有仓库地址.
  url: http://192.168.2.200:5000/v2
  # 私有仓库命名.
  name: 192.168.2.200:5000
  # 是否只读设置.
  readonly: false
  auth:
    # 是否进行鉴权处理.
    enabled: false
    # 需要与私有仓库定义的名称一致.
    issuer: 'my issuer'
    # 容器内私钥证书路径，容器启动时通过数据卷参数指定.
    key: /conf/auth.key
docker pull registry

docker run -d -p 5000:5000 --restart=always --name registry-srv \
		-v /securitit/registry/conf/registry-srv.yml:/etc/docker/registry/config.yml \
    -v /securitit/registry/data/registry:/var/lib/registry  \
		-v /securitit/registry/auth/auth.cert:/etc/docker/registry/auth.cert \
    -v /securitit/registry/auth/htpasswd:/etc/docker/registry/htpasswd \
		-e "REGISTRY_AUTH=htpasswd" \
		-e "REGISTRY_AUTH_HTPASSWD_REALM=Registry Realm" \
		-e  REGISTRY_AUTH_HTPASSWD_PATH=/etc/docker/registry/htpasswd  \
		registry
    
docker pull hyper/docker-registry-web

docker run -it -d -v /securitit/registry/conf/registry-web.yml:/conf/config.yml \
           -v /securitit/registry/auth/auth.key:/conf/auth.key \
           -v /securitit/registry/db:/data \
           -e REGISTRY_TRUST_ANY_SSL=false   \
           -e registry_url=http://192.168.2.200:5000/v2 \
           -e REGISTRY_BASIC_AUTH="" \
           -e REGISTRY_TRUST_ANY_SSL=false   \
            -e registry_auth_enabled=false \
            -e registry_readonly=false \
           -p 5050:8080 --name registry-web hyper/docker-registry-web
```



**验证**

```
docker login 192.168.2.200:5000
docker tag nginx 192.168.2.200:5000/nginx:1.0
docker push 192.168.2.200:5000/nginx:1.0
```

![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1712493070530-a73d278b-81ef-4ce1-83e8-4448a97c488b.png)

#### 安装cri-docker

所有机器都要执行

1.24.0开始使用containerd作为运行时环境，使用docker要安装cri-docker

```shell
# 下的比较慢，有代理执行一下代理  export https_proxy=http://192.168.2.109:7890
wget https://github.com/Mirantis/cri-dockerd/releases/download/v0.3.1/cri-dockerd_0.3.1.3-0.ubuntu-focal_amd64.deb
#安装插件
dpkg -i cri-dockerd_0.3.1.3-0.ubuntu-focal_amd64.deb
vim /lib/systemd/system/cri-docker.service
#修改启动命令，使用cni插件 指定仓库
ExecStart=/usr/bin/cri-dockerd --network-plugin=cni  --pod-infra-container-image=registry.aliyuncs.com/google_containers/pause:3.7

systemctl daemon-reload && systemctl restart cri-docker
```

#### 安装k8s相关命令

所有机器都要执行

```shell
apt-get update && apt-get install -y apt-transport-https
curl https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | apt-key add - 
cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main
EOF
apt-get update
apt install -y kubelet=1.24.0-00 kubeadm=1.24.0-00 kubectl=1.24.0-00
systemctl enable kubelet
```

#### 启动master启动node

```shell
#master执行，拉取k8s组件的相关镜像。
kubeadm config images pull --kubernetes-version=v1.24.0 --cri-socket unix:///run/cri-dockerd.sock --image-repository registry.aliyuncs.com/google_containers
#master执行 初始化
kubeadm init \
--kubernetes-version v1.24.0 \
--pod-network-cidr=10.244.0.0/16 \
--service-cidr=10.96.0.0/12 \
--cri-socket unix:///run/cri-dockerd.sock \
--image-repository registry.aliyuncs.com/google_containers 
--v 5
#成功后按照提示将config 文件复制到 /root/.kube/config

#在master节点执行 将master .kube中的config复制到node 200
#如果node中没有.kube文件夹
mkdir /root/.kube
scp -r $HOME/.kube/config root@192.168.2.200:/root/.kube/config


#node执行，子节点加入
kubeadm join 192.168.2.199:6443 --token jztt4p.60cvat29e6k7l2ud \
--cri-socket=unix:///run/cri-dockerd.sock \
        --discovery-token-ca-cert-hash sha256:086ca57da1b21a7d8fc805946b88aaab5661f0c339f63692842cbffcb5da3e75
# 如果token不记得执行下面的代码重新生成token
kubeadm token create --print-join-command
```


 安装cni插件 master， node 都要安装。

```shell
#下载calico
wget https://docs.projectcalico.org/manifests/calico.yaml
#编辑文件,找到下面这两句,去掉注释,修改ip为当前你设置的pod ip段
vim calico.yaml
- name: CALICO_IPV4POOL_CIDR
  value: "10.244.0.0/16"
#镜像拉取没有问题的话最好
kubectl apply -f calico.yaml 	
```





### 错误排除

#### node不可调度

https://stackoverflow.com/questions/55432764/my-worker-node-status-is-ready-schedulingdisabled

**kubectl uncordon node199**



#### node 处于 notready状态

tail -f /var/log/syslog 查看，没安装cni网络插件，按照上面安装插件



#### cni网络问题

 Failed to create pod sandbox: rpc error: code = Unknown desc = [failed to set up sandbox container "678b6cb1055849a659946bfed802d0066908651ef9c145f077ac4b797c1c1eee" network for pod "hello-world1-okteto-f54468447-mmpdg": networkPlugin cni failed to set up

 pod "hello-world1-okteto-f54468447-mmpdg_default" network: plugin type="calico" failed (add): error getting ClusterInformation: connection is unauthorized: Unauthorized, failed to clean up sandbox container "678b6cb1055849a659946bfed802d0066908651ef9c145f077

ac4b797c1c1eee" network for pod "hello-world1-okteto-f54468447-mmpdg": networkPlugin cni failed to teardown pod "hello-world1-okteto-f54468447-mmpdg_default" network: plugin type="calico" failed (delete): error getting ClusterInformation: connection is unauthorized: Unauthorized

卸载cni插件重新安装即可

```
kubectl delete -f calico.yaml 	&& kubectl delete -f calico.yaml 	
```



#### k8s卸载

https://www.orchome.com/16610

node和master都可以执行这个命令

```
kubeadm reset  --cri-socket unix:///run/cri-dockerd.sock 
```



```plain
rm -rf /etc/kubernetes/manifests/kube-apiserver.yaml  \                                                                                                                                                                                                                  
/etc/kubernetes/manifests/kube-controller-manager.yaml \                                                                                                                                                                                                          
/etc/kubernetes/manifests/kube-scheduler.yaml  \                                                                                                                                                                                                                  
/etc/kubernetes/manifests/etcd.yaml \  
/etc/kubernetes/kubelet.conf  \
/etc/kubernetes/admin.conf \
/etc/kubernetes/scheduler.conf \
/var/lib/kubelet/kubeadm-flags.env \
/etc/kubernetes/controller-manager.conf
```

### Kuboard 

k8s 网页管理工具

https://www.cnblogs.com/smj-7038/p/17098621.html

```shell
sudo docker run -d \
  --restart=unless-stopped \
  --name=kuboard \
  -p 8087:80/tcp \
  -p 10081:10081/tcp \
  -e KUBOARD_ENDPOINT="http://192.168.2.200:20" \
  -e KUBOARD_AGENT_SERVER_TCP_PORT="10081" \
  -e KUBOARD_ADMIN_DERAULT_PASSWORD="Kuboard123" \
  -v /home/docker-volumes/kubiard-data:/data \
  eipwork/kuboard:v3
```

将.kute/config填到这个地方

![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1712496833442-3c532335-1714-4c29-b117-370370abd055.png)

k8s的角色权限管理还没搞明白，无脑执行下面的，将admin和system:anonymous加入到群组中，你才能看到命名空间。

```
kubectl create clusterrolebinding admin --clusterrole=cluster-admin   --user=admin
kubectl create clusterrolebinding sanonymous --clusterrole=cluster-admin --user=system:anonymous
```

## go程序

这个要在200这台linux环境的机器中执行，推荐大家都配一个smb,将linux机器的文件映射为windows的一个磁盘。

```shell
package main

import (
	"github.com/gin-gonic/gin"
	"time"
)

func main() {

	route := gin.Default()
	route.GET("/api/time", func(c *gin.Context) {
		h := gin.H{"time": time.Now().Format(time.DateTime), "code": 200}
		c.JSON(200, h)
	})
	route.Run(":8089")
}
FROM debian:stretch-slim
WORKDIR /app
COPY dockerdemo /app/dockerdemo

ENTRYPOINT ["/app/dockerdemo"]
```



```yaml
buildDocker:
	docker build -t 192.168.2.200:5000/dockerdemo:1.0 .
push:
	docker push 192.168.2.200:5000/dockerdemo:1.0
apply:
	kubectl apply -f deployment
delete:
	kubectl delete -f deployment
	docker rmi 192.168.2.200:5000/dockerdemo:1.0
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -buildvcs=false -o dockerdemo
dockerRun:
	docker run -p 8089:8089 --name dockerdemo -it 192.168.2.200:5000/dockerdemo:1.0
start:
	make build
	make buildDocker
	make push
	make apply
```

如果你所有环境都配置好了，执行 `make start`即可

验证查看pod

![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1712501703006-1aaf1de4-9b43-4a03-a5cf-de475922b6da.png)

![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1712501659725-f7ef95bc-bd5f-4bf2-9c45-48914385090c.png)

访问192.168.2.200:30001

curl  http://192.168.2.200:30001/api/time   返回       {"code":200,"time":"2024-04-07 14:52:48"}                                                                                     

## pv pvc

https://kubernetes.io/zh-cn/docs/concepts/storage/persistent-volumes/

https://kubernetes.io/zh-cn/docs/tasks/configure-pod-container/configure-persistent-volume-storage/