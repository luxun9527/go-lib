



## 单机启动

```bash
##先将配置文件中授权先关闭auth=false
docker-compose up;

docker exec -it mongodb /bin/bash

##进入mongodb shell##
mongo

##切换到admin库##
> use admin

##创建账号/密码##
db.createUser({ user: 'admin', pwd: 'admin', roles: [ { role: "userAdminAnyDatabase", db: "admin" } ] });

docker-compose up -d;
```

## 集群副本集模式启动