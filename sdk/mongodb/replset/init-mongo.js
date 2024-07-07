// 等待 MongoDB 实例启动
sleep(10000);
//openssl rand -base64 756 > ${PWD}/rs_keyfile
// 连接到 MongoDB 实例并进行身份验证
db = connect("mongo1:27017/admin");
db.auth("root", "example");

// 初始化副本集
rs.initiate({
    _id: "rs0",
    members: [
        { _id: 0, host: "192.168.2.159:30011" },
        { _id: 1, host: "192.168.2.159:30012" },
        { _id: 2, host: "192.168.2.159:30013" }
    ]
});
