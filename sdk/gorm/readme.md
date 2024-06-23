1、gorm的常用用法。自定义类型。clause

2、gorm+gen 提高crud的速度

3、trace链路追踪。

4、gorm集群和分库分表。

5、日志。

6、gorm.DB链式操作需要注意的地方

https://gorm.io/zh_CN/docs/index.html

gorm提供的文档非常齐全了，总结一下在工作中的一些实践，如果对您有帮我，您的star就是对我的鼓励https://github.com/luxun9527/go-lib 如有任何问题也欢迎留言。



**安装mysql**

```powershell
mkdir /root/docker/mysql/conf/ && mkdir /root/docker/mysql/data
docker run -p 3306:3306 --name mysql8  --restart always \
-e MYSQL_ROOT_PASSWORD=root \
-v /root/docker/mysql/conf/my.cnf:/etc/mysql/my.cnf \
-v /root/docker/mysql/data:/var/lib/mysql \
-d mysql:8.0
```



## 1、gorm的常用用法,curd,自定义类型，clause子句。

### 新增

```go
func TestCreate(t *testing.T) {
    user := &User{
        Username:  "",
        Age:       0,
        Fav:       "",
        CreatedAt: 0,
        UpdatedAt: 0,
    }
    //INSERT INTO `user` (`username`,`age`,`fav`,`created_at`,`updated_at`) VALUES ('',0,'',1692947238,1692947238)
    if err := db.Create(user).Error; err != nil {
        log.Panicf("create no select  err %v", err)
    }
    log.Printf("user1 = %+v", user)
    //INSERT INTO `user` (`fav`,`created_at`,`updated_at`) VALUES ('',1692947405,1692947405)
    if err := db.Select("fav").Create(user).Error; err != nil {
        log.Panicf("create select field err %v", err)
    }
    log.Printf("user2 = %+v", user)
    //INSERT INTO `user` (`username`,`age`,`created_at`,`updated_at`,`id`) VALUES ('',0,1692947813,1692947813,14)
    if err := db.Omit("fav").Create(user).Error; err != nil {
        log.Panicf("create Omit field err %v", err)
    }
    log.Printf("user2 = %+v", user)

    //https://gorm.io/zh_CN/docs/models.html https://gorm.io/zh_CN/docs/create.html#%E9%BB%98%E8%AE%A4%E5%80%BC
    //官方的文档很全主要是验证一些模糊的地方
    //1、不使用select指定，零值字段也会被插入。使用select只会插入指定的字段。可以使用tag指定默认值 `gorm:"default:18"`
    //2、created_at updated_at会被填充当前时间插入。 https://gorm.io/zh_CN/docs/models.html#%E5%88%9B%E5%BB%BA-x2F-%E6%9B%B4%E6%96%B0%E6%97%B6%E9%97%B4%E8%BF%BD%E8%B8%AA%EF%BC%88%E7%BA%B3%E7%A7%92%E3%80%81%E6%AF%AB%E7%A7%92%E3%80%81%E7%A7%92%E3%80%81Time%EF%BC%89
    //3、插入后会将主键赋值回来，使用select的方式插入不会赋值回来。
}

2024/04/21 23:24:21 E:/demoproject/go-lib/sdk/gorm/curd_test.go:74
[3.177ms] [rows:1] INSERT INTO `user` (`username`,`age`,`fav`,`created_at`,`updated_at`,`deleted_at`) VALUES ('',0,'',1713713061,1713713061,'0')
2024/04/21 23:24:21 user1 = &{ID:7 Username: Age:0 Fav: CreatedAt:1713713061 UpdatedAt:1713713061 DeletedAt:0}

2024/04/21 23:24:21 E:/demoproject/go-lib/sdk/gorm/curd_test.go:79
[3.227ms] [rows:1] INSERT INTO `user` (`fav`,`created_at`,`updated_at`) VALUES ('',1713713061,1713713061)
2024/04/21 23:24:21 user2 = &{ID:7 Username: Age:0 Fav: CreatedAt:1713713061 UpdatedAt:1713713061 DeletedAt:0}

2024/04/21 23:24:21 E:/demoproject/go-lib/sdk/gorm/curd_test.go:84 Error 1062 (23000): Duplicate entry '7' for key 'user.PRIMARY'
[0.534ms] [rows:0] INSERT INTO `user` (`username`,`age`,`created_at`,`updated_at`,`deleted_at`,`id`) VALUES ('',0,1713713061,1713713061,'0',7)
2024/04/21 23:24:21 create Omit field err Error 1062 (23000): Duplicate entry '7' for key 'user.PRIMARY'
--- FAIL: TestCreate (0.01s)
panic: create Omit field err Error 1062 (23000): Duplicate entry '7' for key 'user.PRIMARY' [recovered]
panic: create Omit field err Error 1062 (23000): Duplicate entry '7' for key 'user.PRIMARY'
```

### 修改

```go
func TestUpdate(t *testing.T) {
    u := &User{ID: 1}
    //更新单个列 只有非空的字段会被修改
    //UPDATE `user` SET `username`='hello',`updated_at`=1692948530 WHERE `id` = 1
    db.Model(&u).Update("username", "hello")
    //更新多个列,只有非空的字段会被修改
    // UPDATE `user` SET `username`='hello12',`updated_at`=1692948530 WHERE `id` = 1
    db.Model(&u).Updates(&User{Username: "hello12"})
    //select 更新指定了列,不会忽略零值
    // UPDATE `user` SET `username`='',`updated_at`=1692948530 WHERE `id` = 1
    db.Model(&u).Select("username").Updates(&User{
        Username: "",
        Fav:      "",
    })
    //忽略指定的列，会忽略空值。
    //UPDATE `user` SET `updated_at`=1695796004 WHERE `id` = 1
    db.Model(&u).Omit("fav").Updates(&User{
        Username: "",
        Fav:      "1",
    })
    //UPDATE `user` SET `id`=0,`username`='',`age`=0,`created_at`=0,`updated_at`=1695796443,`deleted_at`=0 WHERE `id` = 1
    //select 所有再忽略指定的列，不会忽略零值
    db.Select("*").Omit("fav").Updates(&User{
        ID:       1,
        Username: "",
        Fav:      "1",
    })
    //指定条件更新，会忽略空值
    // UPDATE `user` SET `id`=1,`username`='1',`updated_at`=1692948928 WHERE username = 'admin' AND `id` = 1
    db.Where("username = ?", "admin").Updates(User{Username: "1", ID: 1})
    //默认没有指定条件不会全局更新。db.Model(&User{}).Update("name", "jinzhu").Error gorm.ErrMissingWhereClause

    //如果要更新零字段，要使用map或者select指定字段。 select 所有不会忽略零值
    // UPDATE `user` SET `id`=1,`username`='',`age`=0,`fav`='',`created_at`=0,`updated_at`=1692948946 WHERE `id` = 1
    db.Select("*").Updates(User{
        ID:        1,
        Username:  "",
        Age:       0,
        Fav:       "",
        CreatedAt: 0,
        UpdatedAt: 0,
    })
}
```

### 删除

```go
func TestDelete(t *testing.T) {
	//UPDATE `user` SET `deleted_at`=1713714861 WHERE id = 1 AND `user`.`deleted_at` = 0
	db.Where("id = ?", 3).Delete(&User{})

	//DELETE FROM `user` WHERE id = 3
	db.Unscoped().Where("id = ?", 3).Delete(&User{})
}
```

### 查询

```go
func TestSelect(t *testing.T) {
    //普通的查询，官方文档上写的比较全,主要是显示预加载。
    //========================================一对一 一个用户有一个profile======================
    var users1 []*User
    if err := db.Model(&User{}).Where("id = ?", 6).Preload("Profile").Find(&users1).Error; err != nil {
        log.Println(err)
    }
    //SELECT * FROM `card` WHERE `card`.`user_id` = 6
    //SELECT * FROM `user` WHERE id = 6
    //SELECT * FROM `profile` WHERE `profile`.`user_id` = 6
    //===================================================一对多，一个用户有多个cards=================================
    var users2 []*User
    if err := db.Model(&User{}).Where("id = ?", 6).Preload("Cards").Preload("Profile").Find(&users2).Error; err != nil {
        log.Println(err)

    }
    for _, v := range users2 {
        log.Printf("%+v", v)
    }
    /*
		[0.529ms] [rows:0] SELECT * FROM `card` WHERE `card`.`user_id` = 6

		2024/04/23 22:57:55 E:/demoproject/go-lib/sdk/gorm/curd_test.go:215
		[1.098ms] [rows:0] SELECT * FROM `profile` WHERE `profile`.`user_id` = 6

		2024/04/23 22:57:55 E:/demoproject/go-lib/sdk/gorm/curd_test.go:215
		[3.268ms] [rows:1] SELECT * FROM `user` WHERE id = 6 AND `user`.`deleted_at` = 0
		2024/04/23 22:57:55 &{ID:6 Username: Age:0 Fav: CompanyID:1 CreatedAt:1713713017 UpdatedAt:1713713017 Cards:[] Profile:{ID:0 UserID:0 Nickname: Desc: CreatedAt:0 UpdatedAt:0 DeletedAt:0} Company:{ID:0 Name:} DeletedAt:0}

	*/
    //===============================belong to 一个用户属于一家公司=====================
    var (
        user    User
        company Company
    )
    //查用户所属的公司。
    //SELECT * FROM `user` WHERE id = 1 AND `user`.`deleted_at` = 0 LIMIT 1
    db.Where("id = ?", 6).Take(&user)
    //SELECT * FROM `company` WHERE `company`.`id` = 1
    if err := db.Debug().Model(&user).Association("Company").Find(&company); err != nil {
        log.Println("belong to error", err)
    }
    user.Company = company
    log.Printf("user =%+v", user)
    //user ={ID:6 Username: Age:0 Fav: CompanyID:0 CreatedAt:1713713017 UpdatedAt:1713713017  Company:{ID:1 Name:test} DeletedAt:0}

    //=====================================join 预加载==========================================================
    printLine("join 预加载")
    var u User
    if err := db.Joins("Company").Take(&u, 6).Error; err != nil {
        log.Panicf("join error %v", err)
    }
    //SELECT `user`.`id`,`user`.`username`,`user`.`age`,`user`.`fav`,`user`.`company_id`,`user`.`created_at`,`user`.`updated_at`,`user`.`deleted_at`,
    //`Company`.`id` AS `Company__id`,`Company`.`name` AS `Company__name` FROM `user`
    //LEFT JOIN `company` `Company` ON `user`.`company_id` = `Company`.`id` WHERE `user`.`id` = 4 AND `user`.`deleted_at` = 0 LIMIT 1
    log.Printf("join u =%+v", u)
}
```

### 自定义类型

https://github.com/go-gorm/datatypes

https://gorm.io/zh_CN/docs/data_types.html#

自定义的数据类型必须实现 [Scanner](https://pkg.go.dev/database/sql#Scanner) 和 [Valuer](https://pkg.go.dev/database/sql/driver#Valuer) 接口，以便让 GORM 知道如何将该类型接收、保存到数据库https://github.com/go-gorm/datatypes/blob/master/date.go 示例 

其他高级的用法参考https://gorm.io/zh_CN/docs/data_types.html#

```go
package main

import (
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"
)

type CustomDateModel struct {
	ID         int32 `gorm:"column:id;not null;" json:"user_id"`
	CustomDate Date  `gorm:"column:custom_date;not null;" json:"custom_date"`
}

// TableName Card's table name
func (*CustomDateModel) TableName() string {
	return "custom_date"
}

func TestCustomType(t *testing.T) {
	db.Create(&CustomDateModel{CustomDate: Date(time.Now())})
	//INSERT INTO `custom_date` (`custom_date`) VALUES ('2024-05-05 00:00:00')
}

type Date time.Time

//将数据库中的数据转换为Date类型
func (date *Date) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*date = Date(nullTime.Time)
	return
}
//将Date类型转换为数据库中的数据
func (date Date) Value() (driver.Value, error) {
	y, m, d := time.Time(date).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Time(date).Location()), nil
}

// GormDataType gorm common data type
func (date Date) GormDataType() string {
	return "date"
}

func (date Date) GobEncode() ([]byte, error) {
	return time.Time(date).GobEncode()
}

func (date *Date) GobDecode(b []byte) error {
	return (*time.Time)(date).GobDecode(b)
}

func (date Date) MarshalJSON() ([]byte, error) {
	return time.Time(date).MarshalJSON()
}

func (date *Date) UnmarshalJSON(b []byte) error {
	return (*time.Time)(date).UnmarshalJSON(b)
}
```




### clause子句

https://blog.csdn.net/dorlolo/article/details/127238144

[https://gorm.io/zh_CN/docs/create.html#Upsert-%E5%8F%8A%E5%86%B2%E7%AA%81](https://gorm.io/zh_CN/docs/create.html#Upsert-及冲突)

gorm与子句生成器有关的类，按父级到子集排列为 [DB](https://github.com/go-gorm/gorm/blob/master/gorm.go#L89) **-->** [Statement](https://github.com/go-gorm/gorm/blob/master/statement.go) **-->** [Clause](https://github.com/go-gorm/gorm/blob/master/clause/clause.go) **-->** [Expression](https://github.com/go-gorm/gorm/blob/master/clause/expression.go) (分别对应 数据库连接对象–> 语句 --> 子句 --> 表达式),它们都是以属性形式保存在父类中。只要知道这个结构，看源码就会轻松很多。

![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1714919325706-db595e7b-bd0e-407e-9b48-7eab3a2f9a0d.png)

```go
user := &User{}
db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)
//INSERT INTO `user` (`username`,`age`,`fav`,`company_id`,`created_at`,`updated_at`,`deleted_at`) VALUES ('',0,'',0,1714919066,1714919066,0) ON DUPLICATE KEY UPDATE `id`=`id`
```



**2. 为单个字段实现开发构造器**

2.1 实现方法

只要自定义的模型字段中包含以下的其中一种方法即可，

```go
CreateClauses(*Field) []clause.Interface
QueryClauses(*Field) []clause.Interface
UpdateClauses(*Field) []clause.Interface
DeleteClauses(*Field) []clause.Interface
```

这个返回值类型可以实现StatementModifier接口 或者 clause.Interface接口都行。使用实例。

```go
package main

import (
	"go-lib/sdk/gorm/gen/dao/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"testing"
)

type ClauseUser struct {
	model.User
	C CustomClause
}

func TestClauses(t *testing.T) {
	user := &User{}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user)
	//INSERT INTO `user` (`username`,`age`,`fav`,`company_id`,`created_at`,`updated_at`,`deleted_at`) VALUES ('',0,'',0,1714919066,1714919066,0) ON DUPLICATE KEY UPDATE `id`=`id`
	db.Clauses(clause.Eq{
		Column: "id",
		Value:  1,
	}).Find(&user)
	//SELECT * FROM `user` WHERE `id` = 1 AND `user`.`deleted_at` = 0 AND `user`.`id` = 19
	db.Find(&ClauseUser{})
	//SELECT * FROM `user` WHERE `user`.`c` IS NULL
}

type CustomClause int32

func (CustomClause) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{CustomClauseQuery{Field: f}}
}

type CustomClauseQuery struct {
	Field *schema.Field
}

func (sd CustomClauseQuery) Name() string {
	return ""
}

func (sd CustomClauseQuery) Build(builder clause.Builder) {

}

func (sd CustomClauseQuery) MergeClause(c *clause.Clause) {

}

func (sd CustomClauseQuery) ModifyStatement(stmt *gorm.Statement) {
	stmt.AddClause(clause.Where{Exprs: []clause.Expression{
		clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Value: nil},
	}})
}
```

详情用法可以参考https://github.com/go-gorm/soft_delete 软删除的做法。



## 2、gorm + gen

https://gorm.io/zh_CN/gen/dao.html

https://gorm.io/zh_CN/gen/database_to_structs.html

相对于gorm gen提供对表的curd进一步的封装，能提高我们的开发效率，具体有如下优势。

**1、直接使用表结构生成model，不用自己手写model**

```
 go install gorm.io/gen/tools/gentool@latest
gentool --dsn="root:root@tcp(192.168.2.200:33606)/test?charset=utf8mb4&parseTime=True&loc=Local" --db=mysql  -outPath=gen/dao/query -fieldSignable=true
```

**如果有定制需求可以参考**。https://gorm.io/zh_CN/gen/database_to_structs.html 

**有一些坑**

1、https://github.com/go-gorm/gen/issues/755自定义方法CommonMethod 需要和生成代码不在同一个包

[https://gorm.io/zh_CN/gen/database_to_structs.html#%E6%A8%A1%E6%9D%BF%E6%96%B9%E6%B3%95](https://gorm.io/zh_CN/gen/database_to_structs.html#模板方法) 

```go
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameUser = "user"

// User mapped from table <user>
type User struct {
    ID        int32  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
    Username  string `gorm:"column:username;not null;comment:用户名" json:"username"`
    Age       uint32 `gorm:"column:age;not null;comment:年龄" json:"age"`
    Fav       string `gorm:"column:fav;not null;comment:爱好" json:"fav"`
    CompanyID int32  `gorm:"column:company_id;not null;comment:公司Id" json:"company_id"`
    CreatedAt uint64 `gorm:"column:created_at;not null;comment:创建时间" json:"created_at"`
    UpdatedAt uint64 `gorm:"column:updated_at;not null;comment:修改时间" json:"updated_at"`
    DeletedAt uint64 `gorm:"column:deleted_at;not null" json:"deleted_at"`
}

// TableName User's table name
func (*User) TableName() string {
    return TableNameUser
}
```

**2、使用变量名而不是直接表字段名，避免写错字段。**

```go
gen ====> users, err := u.WithContext(ctx).Where(u.Name.Neq("modi")).Find()
gorm====> db.Where("name=?","modi").Find()
// SELECT * FROM users WHERE name <> 'modi';
```

**3、对查询的进一步封装,在查询上更加方便。**

```go
dao := query.Use(db)
//card := dao.Card
user := dao.User
company := dao.Company
ctx := context.Background()

_, err := user.WithContext(ctx).Find()
if err != nil {
    log.Panicf("err:%v", err)
}
var u []*User
if db.Find(&u).Error != nil {
    log.Panicf("err:%v", err)
}

//gen join
dao.User.WithContext(ctx).Join(company, user.CompanyID.EqCol(company.ID)).Find()
//gorm join
db.Joins("LEFT JOIN Company c ON c.id =  u.company_id").Find(&u)

// gen分页
dao.User.WithContext(ctx).FindByPage(0, 10)
//gorm 分页
db.Offset(0).Limit(10).Find(&u)

	// gen子查询
	subQuery := company.WithContext(ctx).
		Select(company.ID)
	//
	_, err = user.WithContext(ctx).
		Where(user.Columns(user.CompanyID).In(subQuery)).
		Find()
//SELECT * FROM `user` WHERE `user`.`company_id` IN (SELECT `company`.`id` FROM `company`)
```

**需要注意的的是：**

https://github.com/go-gorm/gen/issues/900 gen的Select方法和Where都是不能自定义的，当进行一些复杂查询的时候，需要自定义的时候如使用json函数`db.Select("JSON_OBJECT(key1,val1,key2,val2...)")` gen是无法做到的。但是我们可以使用ast来生成对应的代码。具体可以参考。https://github.com/luxun9527/go-lib/blob/master/utils/ast_apply/add_rawselect_rawmethod.go

## 3、trace链路追踪

将gorm加入到链路追踪中。使用WithContext传递。

```go
package main

import (
	"context"
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"log"
	"testing"
)

const (
	service     = "trace-demo1" // 服务名
	environment = "production"  // 环境
	id          = 1             // id
)

func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	// 创建 Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return tp, nil
}
func TestGormTrace(t *testing.T) {
	ctx := context.Background()

	tp, err := tracerProvider("http://192.168.11.185:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Use(otelgorm.NewPlugin(
		otelgorm.WithDBName("test"),
		otelgorm.WithTracerProvider(tp),
	)); err != nil {
		log.Fatal(err)
	}
	otel.SetTracerProvider(tp)

	tracer := otel.Tracer("gormtracer")

	ctx, span := tracer.Start(ctx, "gormtest")
	defer span.End()
	user := &User{
		Username: "test1",
		Age:      1,
		Fav:      "1",
	}
	//INSERT INTO `user` (`username`,`age`,`fav`,`created_at`,`updated_at`) VALUES ('',0,'',1692947238,1692947238)
	if err := db.WithContext(ctx).Create(user).Error; err != nil {
		fmt.Println("create err", err)
	}

	//otelplay.PrintTraceID(ctx)
	tp.Shutdown(ctx)
}
```

![img](https://cdn.nlark.com/yuque/0/2024/png/12466223/1714832378749-287c1857-8b4f-46a5-9ee3-335f988ab481.png)

## 4、主从集群，分库分表

refer https://www.freesion.com/article/92701394177/

https://www.cnblogs.com/cao-lei/p/13603043.html

主从延迟 https://jishuin.proginn.com/p/763bfbd589f7

https://blog.csdn.net/JesseYoung/article/details/40585809

### mysql使用docker安装mysql主从集群



进入Master库mysql客户端 `mysql -uroot -h127.0.0.1 -P33307 -proot`。

如果本地没有mysql客户端则进入容器中

`docker exec -it mysql-master /bin/bash `在容器中的端口为 3306

输入`show master status`查看Master状态：

![img](https://cdn.nlark.com/yuque/0/2022/png/12466223/1661091602934-b7fad82d-5d74-440c-85a5-49822c41913b.png)

记住File和Position，后面需要用到。此时一定不要操作Master库，否则将会引起Master状态的变化，File和Position字段也将会进行变化。

#### 4.2 slave配置

```powershell
docker run -p 33307:3306 --name mysql8s  --restart always \
-e MYSQL_ROOT_PASSWORD=root \
-v /root/docker/mysql/slave/conf/my.cnf:/etc/mysql/my.cnf \
-v /root/docker/mysql/slave/data:/var/lib/mysql \
-d mysql:8.0
```

进入mysql客户端执行命令

```
change master to master_host='192.168.2.200', master_user='root', master_password='root', master_port=33606, master_log_file='mysql-bin.000009', master_log_pos=814, master_connect_retry=30;
```

master_host ：Master库的地址局域网的ip
master_port：Master的端口号，指的是容器的端口号
master_user：用于数据同步的用户
master_password：用于同步的用户的密码
master_log_file：指定 Slave 从哪个日志文件开始复制数据，即上文中提到的 File 字段的值
master_log_pos：从哪个 Position 开始读，即上文中提到的 Position 字段的值
master_connect_retry：如果连接失败，重试的时间间隔，单位是秒，默认是60秒

2、 启动slave

```
start slave;
```

3、查看slave

```
show slave status \G;
```

![img](https://cdn.nlark.com/yuque/0/2022/png/12466223/1661092118656-e79f4ede-b693-476c-98b0-1612688c40a7.png)



### 4.3 slave常见错误

 1、日志设置错误![img](https://cdn.nlark.com/yuque/0/2022/png/12466223/1661092380180-bf21cc9b-29a1-4ad2-83b6-c9dc5526c8e4.png)

```
master_log_file='mysql-bin1.000004', master_log_pos=157
```

![img](https://cdn.nlark.com/yuque/0/2022/png/12466223/1661092437579-b3d0d006-ed30-4957-b956-927d84c03ab6.png)

检查日志是否和master的一致

2、master ip端口配置错误配置错误![img](https://cdn.nlark.com/yuque/0/2022/png/12466223/1661092499187-53369e2a-3979-4526-98f1-ba7a77b745ae.png)

### 其他命令

`stop slave;   `停止slave

`reset master;` 重置master

### 测试是否成功

```plsql
#创建数据库
create DATABASE test1;

use test1;
#创建数据表
CREATE TABLE IF NOT EXISTS `runoob_tbl`(
  `runoob_id` INT UNSIGNED AUTO_INCREMENT,
  `runoob_title` VARCHAR(100) NOT NULL,
  `runoob_author` VARCHAR(40) NOT NULL,
  `submission_date` DATE,
  PRIMARY KEY ( `runoob_id` )
);
#插入数据
INSERT INTO runoob_tbl (runoob_title, runoob_author, submission_date)  
VALUES("学习 PHP", "菜鸟教程", NOW());
docker exec -it mysql8s mysql -uroot -proot
use test1;
select * from runoob_tbl;
```

### `gorm mysql主从集群`

https://gorm.io/zh_CN/docs/dbresolver.html

```go
package main

import (
    "github.com/spf13/cast"
    "github.com/zeromicro/go-zero/core/stringx"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "gorm.io/plugin/dbresolver"
    "log"
    "os"
    "testing"
    "time"
)

const TableNameRunoobTbl = "runoob_tbl"

// RunoobTbl mapped from table <runoob_tbl>
type RunoobTbl struct {
    RunoobID       int32     `gorm:"column:runoob_id;primaryKey;autoIncrement:true" json:"runoob_id"`
    RunoobTitle    string    `gorm:"column:runoob_title;not null" json:"runoob_title"`
    RunoobAuthor   string    `gorm:"column:runoob_author;not null" json:"runoob_author"`
    SubmissionDate time.Time `gorm:"column:submission_date" json:"submission_date"`
}

// TableName RunoobTbl's table name
func (*RunoobTbl) TableName() string {
    return TableNameRunoobTbl
}

const (
    masterDsn = "root:root@tcp(192.168.2.200:33606)/test1?charset=utf8mb4&parseTime=True&loc=Local"
    slaveDsn  = "root:root@tcp(192.168.2.200:33307)/test1?charset=utf8mb4&parseTime=True&loc=Local"
)

// gentool --dsn="root:root@tcp(192.168.2.99:33307)/test1?charset=utf8mb4&parseTime=True&loc=Local" --onlyModel=true --db=mysql --tables=runoob_tbl -outPath=./ -fieldMap="decimal:string;tinyint:int32;"
func TestDbresolve(t *testing.T) {
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
        logger.Config{
            SlowThreshold:             time.Second, // Slow SQL threshold
            LogLevel:                  logger.Info, // Log level
            IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
            Colorful:                  true,        // Disable color

        },
    )
    masterDB, err := gorm.Open(mysql.New(mysql.Config{
        DSN: masterDsn,
    }), &gorm.Config{
        Logger: newLogger,
    })
    if err != nil {
        log.Fatal(err)
    }
    if err := masterDB.Use(
        dbresolver.Register(dbresolver.Config{
            Sources: []gorm.Dialector{mysql.New(mysql.Config{
                DSN: masterDsn,
            })},
            Replicas: []gorm.Dialector{mysql.New(mysql.Config{
                DSN: slaveDsn,
            })},
            Policy:            dbresolver.RandomPolicy{},
            TraceResolverMode: true,
        }).
        SetMaxIdleConns(10).
        SetConnMaxLifetime(time.Hour).
        SetMaxOpenConns(200),
    ); err != nil {
        log.Fatal(err)
    }
    d := &RunoobTbl{
        RunoobTitle:    cast.ToString(time.Now().Unix()),
        RunoobAuthor:   stringx.Randn(10),
        SubmissionDate: time.Now(),
    }
    masterDB.Create(d)

    if masterDB.Where("runoob_id=?", 2).Find(d).Error != nil {
        log.Println(err)
    }
    log.Println(d)
    /*
    2024/05/05 21:48:01 E:/demoproject/go-lib/sdk/gorm/dbresolver_test.go:78
    [9.653ms] [rows:1] [source] INSERT INTO `runoob_tbl` (`runoob_title`,`runoob_author`,`submission_date`) VALUES ('1714916881','XVG0Iay48N','2024-05-05 21:48:01.959')

    2024/05/05 21:48:01 E:/demoproject/go-lib/sdk/gorm/dbresolver_test.go:80
    [1.594ms] [rows:0] [replica] SELECT * FROM `runoob_tbl` WHERE runoob_id=2 AND `runoob_tbl`.`runoob_id` = 2
    */
}
```

### gorm 水平分表

https://gorm.io/zh_CN/docs/sharding.html

这个分表暂时不支持mysql如果mysql进行水平分表的话，要在业务代码中实现。

## 5、日志

https://gorm.io/zh_CN/docs/logger.html 自定义日志需要实现这几个接口。

```go
type Interface interface {
    LogMode(LogLevel) Interface
    Info(context.Context, string, ...interface{})
    Warn(context.Context, string, ...interface{})
    Error(context.Context, string, ...interface{})
    Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}
```

在实际的开发中常常使的比较多的是 `**moul.io/zapgorm2**` 这个库

## 6、gorm.DB链式操作需要注意的地方

**1、gorm 并发不安全。**

https://juejin.cn/post/7134002645651439630

首先，我们需要先去理解几乎每个方法中都会调用的函数：tx = db.getInstance()。

```go
 func (db *DB) getInstance() *DB {
     if db.clone > 0 {
         tx := &DB{Config: db.Config, Error: db.Error}
 
         if db.clone == 1 {
             // clone with new statement
             tx.Statement = &Statement{
                 DB:       tx,
                 ConnPool: db.Statement.ConnPool,
                 Context:  db.Statement.Context,
                 Clauses:  map[string]clause.Clause{},
                 Vars:     make([]interface{}, 0, 8),
             }
         } else {
             // with clone statement
             tx.Statement = db.Statement.clone()
             tx.Statement.DB = tx
         }
         return tx
     }

     return db
 }
```

将上述改写并简化一下，大概是这么个逻辑：

```go
func (db *DB) getInstance() *DB {
    switch db.clone:
    case 0:
    return db
    case 1:
    return newStatement() // 一个全新的，空白的Statement
    case 2:
    return db.cloneStatement() // 将之前的Statement复制一份
}
```

当clone=1时，这个*gorm.DB 实例总是并发安全的，因为它总是会返回一个全新的*gorm.DB 实例，不会对老*gorm.DB 实例有什么读写。

当clone=2时，这个*gorm.DB 实例也总是并发安全的，因为任何的 Chain Method 和 Finisher Method 都只会去读和复制当前*gorm.DB 实例的值，而不会修改，因此只会对这个*gorm.DB 实例并发读，那么当然是并发安全的。

当clone=0时，这个*gorm.DB 实例就**不并发安全**。

那clone字段分别会在什么情况下等于0、1、2呢？

- **在使用gorm.Open()之后，新建出来的\*gorm.DB 实例clone字段总是1。**
- 在调用(*gorm.Gorm).Session()时，如果Session{}.NewDB为false，则为返回的*gorm.DB 实例clone字段是2，如果为true，则为1。
- 在调用(*gorm.Gorm).Session()时，如果Session{}.Initialized为true，则返回的*gorm.DB 实例clone字段是0。这条规则优先级高于Session.NewDB。
- **在调用了任意Chain Method、Finisher Method之后，返回的Gorm对象clone字段是0。**



在日常的开发中，我们喜欢将db定义为一个全局变量，**要注意的是绝对不要将一个db.clone=0的db赋值回全局变量下面是一个错误用法。**

```go
	var users1 []*User
    //初始的db.clone为1 执行where db.clone为0
	db = db.Where("id = ?", 6)

	db.Where("name = ?", "lisi").Find(&users1)
```



**下面也是链式操作不注意db.clone造成的问题。**

https://gorm.io/zh_CN/docs/method_chaining.html#Reusability-and-Safety

Reusability and Safety

A critical aspect of GORM is understanding when a *gorm.DB instance is safe to reuse. Following a Chain Method or Finisher Method, GORM returns an initialized *gorm.DB instance. This instance is not safe for reuse as it may carry over conditions from previous operations, potentially leading to contaminated SQL queries. For example:

### Example of Unsafe Reuse

```
queryDB := DB.Where("name = ?", "jinzhu") // First query queryDB.Where("age > ?", 10).First(&user) // SQL: SELECT * FROM users WHERE name = "jinzhu" AND age > 10 // Second query with unintended compounded condition queryDB.Where("age > ?", 20).First(&user2) // SQL: SELECT * FROM users WHERE name = "jinzhu" AND age > 10 AND age > 20
```

### Example of Safe Reuse

To safely reuse a *gorm.DB instance, use a New Session Method:

```
queryDB := DB.Where("name = ?", "jinzhu").Session(&gorm.Session{}) // First query queryDB.Where("age > ?", 10).First(&user) // SQL: SELECT * FROM users WHERE name = "jinzhu" AND age > 10 // Second query, safely isolated queryDB.Where("age > ?", 20).First(&user2) // SQL: SELECT * FROM users WHERE name = "jinzhu" AND age > 20
```

In this scenario, using Session(&gorm.Session{}) ensures that each query starts with a fresh context, preventing the pollution of SQL queries with conditions from previous operations. This is crucial for maintaining the integrity and accuracy of your database interactions.



**核心就是当你进行了一次链式操作，要注意这个db的db.clone字段已经为零，可以一些其他的方法去重置如session**