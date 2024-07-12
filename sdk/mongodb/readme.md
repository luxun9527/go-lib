## 1、适用场景,数据类型。

### mongodb的适合存储的数据和使用场景

1、cms博客文档，评论消息，文本比较多的数据，数据相对较独立，关联较少。

2、im消息。

3、时序数据。iot记录，股票交易记录。

4、复杂嵌套数据，存储和查询复杂的JSON数据结构。  

5、每列字段不固定的数据。



### monogodb数据类型

https://cloud.tencent.com/developer/article/1511583

| Double                   | 1    | double              | shell中的数字类型     | 64位浮点数                     | { “x” : 3.14 } { “x” : 3 }                                   |
| ------------------------ | ---- | ------------------- | --------------------- | ------------------------------ | ------------------------------------------------------------ |
| String                   | 2    | string              |                       | 字符串类型                     |                                                              |
| Object                   | 3    | object              |                       | 对象类型，嵌套                 | {"address":{"street":"123 Main St","city":"New York","zipCode":"10001"},} |
| Array                    | 4    | array               |                       | 数组类型                       | { “x” : [“a”, “b”, “c”]}                                     |
| Binary data              | 5    | binData             | shell中不可用         | 二进制数据类型                 |                                                              |
| Undefined                | 6    | undefined           | 已过时                | 未定义类型                     |                                                              |
| ObjectId                 | 7    | objectId            |                       | 对象id类型                     | {“x” : objectId() }                                          |
| Boolean                  | 8    | bool                |                       | 布尔类型                       |                                                              |
| Date                     | 9    | date                |                       | 日期类型                       | ISODate("2021-06-12T19:20:30Z")   ISO标准的时间Z表示为UTC时间 |
| Null                     | 10   | null                |                       | 用于表示空值或者不存在的字段   |                                                              |
| Regular Expression       | 11   | regex               |                       | 正则表达式类型                 |                                                              |
| DBPointer                | 12   | dbPointer           | 已过时                |                                |                                                              |
| JavaScript               | 13   | javascript          |                       | JavaScript代码                 |                                                              |
| Symbol                   | 14   | symbol              | shell中不可用，已过时 |                                |                                                              |
| JavaScript（with scope） | 15   | javascriptWithScope |                       | 带作用域的JavaScript代码       |                                                              |
| 32-bit integer           | 16   | int                 | shell中不可用         | 32位整数                       |                                                              |
| Timestamp                | 17   | timestamp           |                       | 时间戳类型秒加同一秒的操作顺序 | Timestamp(1623515630,1)                                      |
| 64-bit integer           | 18   | long                | shell中不可用         | 64位整数                       |                                                              |
| Decimal128               | 19   | decimal             | 3.4版本新增           |                                |                                                              |
| Min key                  | -1   | minKey              | shell中无此类型       | 最小键                         |                                                              |
| Max key                  | 127  | maxKey              | shell中无此类型       | 最大键                         |                                                              |

### 常用命令

#### 插入操作

1. **插入单个文档**

```powershell
db.collection.insertOne({
   "name": "John Doe",
   "age": 30,
   "address": {
      "street": "123 Main St",
      "city": "New York"
   }
});
```

1. **插入多个文档**

```plain
db.collection.insertMany([
   {
      "name": "Jane Doe",
      "age": 25,
      "address": {
         "street": "456 Park Ave",
         "city": "New York"
      }
   },
   {
      "name": "Mike Smith",
      "age": 35,
      "address": {
         "street": "789 Broadway",
         "city": "New York"
      }
   }
]);
```

#### 查询操作

1. **查询所有文档**

```plain
db.collection.find({});
```

1. **查询匹配条件的文档**

```plain
db.collection.find({ "age": { "$gt": 25 } });
```

1. **查询特定字段**

```plain
db.collection.find(
   { "age": { "$gt": 25 } },
   { "name": 1, "address.city": 1, "_id": 0 }
);
```

1. **使用正则表达式进行查询**

```plain
db.collection.find({ "name": { "$regex": "^J", "$options": "i" } });
```

1. **计数查询结果**

```plain
javascript
Copy code
db.collection.find({ "age": { "$gt": 25 } }).count();
```

1. **排序查询结果**

```plain
javascript
Copy code
db.collection.find({ "age": { "$gt": 25 } }).sort({ "age": -1 });
```

1. **限制和跳过结果**

```plain
javascript
Copy code
db.collection.find({ "age": { "$gt": 25 } }).limit(5).skip(10);
```

#### 更新操作

1. **更新单个文档**

```plain
db.collection.updateOne(
   { "name": "John Doe" },
   { "$set": { "age": 31 } }
);
```

1. **更新多个文档**

```plain
db.collection.updateMany(
   { "address.city": "New York" },
   { "$set": { "status": "Active" } }
);
```

1. **替换文档**

```plain
db.collection.replaceOne(
   { "name": "John Doe" },
   { "name": "John Doe", "age": 31, "status": "Active" }
);
```

#### 删除操作

1. **删除单个文档**

```plain
db.collection.deleteOne({ "name": "John Doe" });
```

1. **删除多个文档**

```plain
db.collection.deleteMany({ "status": "Inactive" });
```



## 2、monogdb存储引擎

底层存储引擎.wiredTiger做为存储引擎。

#### wiredTiger的数据结构

https://cloud.tencent.com/developer/article/1815262
 通过查阅资料，我从 MongoDb 的官网和 WiredTiger 官网找到了答案。MongoDb 官网关于存储引擎（Storage Engine）的描述写道：从 MongoDb 3.2 版本开始，其使用了 WiredTiger 作为其默认的存储引擎。而从 WiredTiger 官网文档，我们可以知道：**WiredTiger 使用的是 B+ 树作为其存储结构。**

WiredTiger（以下简称WT）是一个优秀的单机[数据库存储](https://cloud.tencent.com/product/crs?from_column=20065&from=20065)引擎，它拥有诸多的特性，既支持BTree索引，也支持LSM Tree索引，支持行存储和列存储，**实现ACID级别事务**、支持大到4G的记录等。WT的产生不是因为这些特性，而是和计算机发展的现状息息相关。

为啥选择b+树，

**b+树在数据存在叶子节点，每个节点能存储的数据相对更多，树高相对更低，磁盘io相对更少。**

**b+树，叶子节点的数据使用单链表连接，范围查询效率更高。**

#### 缓存

https://mongoing.com/archives/78895



cacheSizeGB默认值是(RAM – 1GB) / 2。这个限制的出发点是防止OOM，因为MongoDB使用的总内存不仅是WT Cache，还包括了一些额外的消耗，

## 3、索引

### 理论知识

monodb的索引和mysql的索引类似

1、WiredTiger索引的数据结构是b+树

2、主键索引和普通索引。 

每个文档必须有一个唯一的_id字段，它作为主键索引，叶子节点存储的是真实的文档记录。

普通索引：数据节点存的是文档的引用即_id字段



索引相关命令

创在MongoDB中，创建索引和使用explain进行慢查询分析是优化数据库性能的重要手段。下面是详细步骤和示例。

### 相关命令

创建单字段索引

```plain
db.collection.createIndex({ fieldName: 1 })  // 升序
db.collection.createIndex({ fieldName: -1 }) // 降序
```

创建复合索引

```plain
db.collection.createIndex({ field1: 1, field2: -1 })
```

创建全文索引

```plain
db.collection.createIndex({ fieldName: "text" })
```

创建唯一索引

```plain
db.collection.createIndex({ fieldName: 1 }, { unique: true })
```

查看已有索引

```plain
db.collection.getIndexes()
```

删除索引

```plain
db.collection.dropIndex("indexName")
```



### 慢查询分析





在MongoDB中，使用explain方法来分析查询的执行计划和性能。要判断查询是否命中索引，可以关注explain结果中的以下几个关键参数：

1. **winningPlan**: 这是MongoDB选择的执行计划。如果查询使用了索引，winningPlan中会包含IXSCAN（索引扫描）阶段。
2. **totalKeysExamined**: 表示查询过程中检查的索引键数量。命中索引时，这个值通常比totalDocsExamined小。
3. **totalDocsExamined**: 表示查询过程中扫描的文档数量。如果命中索引，这个值会显著减少。
4. **stage**: 查看查询计划中的阶段，是否包含IXSCAN或FETCH。

```
db.users.find({ age: { $eq: 25 } }).explain("executionStats")
 "executionStats": {
        "executionSuccess": true,
        "nReturned": NumberInt("112"),
        "executionTimeMillis": NumberInt("0"),
        "totalKeysExamined": NumberInt("112"),
        "totalDocsExamined": NumberInt("112"),
        "executionStages": {
            "stage": "FETCH",
            "nReturned": NumberInt("112"),
            "executionTimeMillisEstimate": NumberInt("0"),
            "works": NumberInt("113"),
            "advanced": NumberInt("112"),
            "needTime": NumberInt("0"),
            "needYield": NumberInt("0"),
            "saveState": NumberInt("0"),
            "restoreState": NumberInt("0"),
            "isEOF": NumberInt("1"),
            "docsExamined": NumberInt("112"),
            "alreadyHasObj": NumberInt("0"),
            "inputStage": {
                "stage": "IXSCAN",
                "nReturned": NumberInt("112"),
                "executionTimeMillisEstimate": NumberInt("0"),
                "works": NumberInt("113"),
                "advanced": NumberInt("112"),
                "needTime": NumberInt("0"),
                "needYield": NumberInt("0"),
                "saveState": NumberInt("0"),
                "restoreState": NumberInt("0"),
                "isEOF": NumberInt("1"),
                "keyPattern": {
                    "age": 1
                },
                "indexName": "age_1",
                "isMultiKey": false,
                "multiKeyPaths": {
                    "age": [ ]
                },
                "isUnique": false,
                "isSparse": false,
                "isPartial": false,
                "indexVersion": NumberInt("2"),
                "direction": "forward",
                "indexBounds": {
                    "age": [
                        "[25.0, 25.0]"
                    ]
                },
                "keysExamined": NumberInt("112"),
                "seeks": NumberInt("1"),
                "dupsTested": NumberInt("0"),
                "dupsDropped": NumberInt("0")
            }
        }
    },
```

## 4、mongodb安装和mongodb集群

https://juejin.cn/post/7142753920043974663

MongoDB 有三种集群部署模式，主从复制（Master-Slaver）、副本集（Replica Set）和分片（Sharding）模式。

#### 1. 主从复制（Master-Slaver）

基本的设置方式是建立一个主节点（Primary）和一个或者多个从节点（Secondary）。

主从复制模式的集群中只能有一个主节点，主节点提供所有的增、删、查、改服务，从节点不提供任何服务。但是可以通过设置使从节点提供查询服务，这样可以减少主节点的压力。 每个从节点要知道主节点的地址，主节点记录在其上的所有操作，从节点定期轮询主节点获取这些操作，然后对自己的数据副本执行这些操作，从而保证从节点的与主节点的数据一致。

**弊端**：主从复制的集群中，当主节点出现故障使，只能人工介入指定新的主节点，从节点不会自动升级为主节点。同时，在这段时间内，该集群架构只能处于只读状态。

#### 2. 副本集（Replica Set）

当集群中的主节点发生故障时，副本集可以自动投票，选举出新的主节点，而且这个过程对应用时透明的。主要为了解决主从复制模式不具备高可用的问题。

#### 3. 分片（Sharding）模式

副本集可以解决主节点发生故障导致数据丢失或不可用的问题，但遇到需要存储海量数据的情况时，副本集中的一台机器不足以存储数据，或者说集群不足以提供可接受的读写吞吐量。这就需要用到 MongoDB 的分片（Sharding）技术，这也是 MongoDB 的另外一种集群部署模式。
MongoDB 的分片机制允许创建一个包含许多台机器的集群，将数据子集分散在集群中，每个分片维护着一个数据集合的子集。与副本集相比，使用集群架构可以使应用程序具有更强大的数据处理能力。 构建一个MongoDB的分片集群，需要三个重要组件，分别是：分片服务器（Shard Server）、配置服务器（Config Server）和路由服务器（Route Server）。

####  4、使用docker部署副本集

https://chenyejun.github.io/blog/mongoDB/mongodbAddUser.html

https://www.cnblogs.com/linmt/p/17365572.html docker安装单节点。

[https://outmanzzq.github.io/2019/01/30/docker-mongo-replica/#%E4%BA%8C%E5%AE%9E%E6%88%98%E5%9C%A8-docker-%E4%B8%AD%E5%88%9B%E5%BB%BA%E5%89%AF%E6%9C%AC%E9%9B%86](https://outmanzzq.github.io/2019/01/30/docker-mongo-replica/#二实战在-docker-中创建副本集)

https://medium.com/workleap/the-only-local-mongodb-replica-set-with-docker-compose-guide-youll-ever-need-2f0b74dd8384

MongoDB的副本集（Replica Set）是MongoDB高可用性和数据冗余的核心机制。副本集是一组MongoDB实例，它们维护相同的数据集，通过选举机制实现自动故障转移和数据同步。副本集模式的集群能够确保在单个节点发生故障时，系统依然能够继续提供服务。

##### 副本集的组成

一个典型的MongoDB副本集由以下几部分组成：

1. **主节点（Primary）**：

- - 负责处理所有的写操作。
  - 每个副本集只能有一个主节点。
  - 主节点会定期将数据变更复制到从节点。

1. **从节点（Secondary）**：

- - 负责复制主节点的数据。
  - 可以提供读操作（如果启用了读偏好）。
  - 在主节点发生故障时，从节点可以参与选举，成为新的主节点。

1. **仲裁节点（Arbiter）**：

- - 不存储数据，只参与选举投票。
  - 用于维持副本集的奇数节点，以避免选举投票中的平票情况。

##### 工作机制

1. **数据复制**：

- - 主节点接收到写操作请求后，将数据变更记录在操作日志（oplog）中。
  - 从节点通过不断拉取主节点的oplog来复制数据变更，保持数据一致性。

1. **自动故障转移**：

- - 如果主节点不可用，从节点会通过选举机制选择一个新的主节点。
  - 选举机制依赖于Raft算法，保证选举过程快速可靠。
  - 仲裁节点帮助维持选举的奇数票数，防止平票。

1. **读操作**：

- - 默认情况下，读操作由主节点处理。
  - 通过设置读偏好，可以配置从节点也处理读操作，以提高读取性能和分担主节点的负载。



**副本集 oplog**

- oplog：操作日志，一个特殊的固定集合，保持所有修改数据库上数据的操作的滚动记录。
- 副本集的所有成员都包含一份 oplog 的副本，在 local.oplog.rs 中，这允许它们维护数据库的当前状态。
- 为了减轻复制的困难，所有的副本集成员成员都发送心跳（ping）到所有其它成员。
- 任何成员都可以从其它成员那里导入 oplog。
- oplog 中的每个操作都是幂等（idempotent）的。即，oplog 对目标资料应用不管一次或多次操作，都产生相同的结果。

**oplog 大小**

- 当第一次开始一个副本集成员时，MongoDB 以默认大小创建一个 oplog。
- 如果可以预料副本集工作量的大小，可以将 oplog 设置为比默认值大些。相反，如果应用主要用来执行读操作和少量写操作，一个小的 oplog 可能就足够了。
- 以下情况可能需要较大的 oplog：

- - 一次更新多个文档
  - 删除与插入量大致相同的数据时
  - 大量的就地更新

**oplog 状态**

使用 rs.printReplicationInfo() 方法来查看 oplog 状态，包括其大小和操作的时间范围。









## 4、go sdk

go sdk增删改查

```go
package mongodb

import (
	"context"
	"errors"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"math/rand"
	"testing"
	"time"
)

var (
	_cli        *mongo.Client
	_dbCli      *mongo.Database
	_collection *mongo.Collection
)

const (
	dbName         = "test"
	collectionName = "users"
)

func initDB() {
	var err error
	uri := "mongodb://root:example@192.168.2.159:30011,192.168.2.159:30012,192.168.2.159:30013/?replicaSet=rs0"

	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(1000).
		SetMinPoolSize(100).
		SetReadPreference(readpref.SecondaryPreferred())
	_cli, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Printf("connect mongodb error: %v", err)
	}
	_dbCli = _cli.Database(dbName)
	_collection = _dbCli.Collection(collectionName)
}

type Student struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

// 定义嵌套的地址结构体
type Address struct {
	Street string `bson:"street"`
	City   string `bson:"city"`
}

type User struct {
	Name      string              `bson:"name"`
	Age       int                 `bson:"age"`
	Email     string              `bson:"email"`
	IsStudent bool                `bson:"isStudent"`
	Scores    []int               `bson:"scores"`
	Address   Address             `bson:"address"`
	CreatedAt primitive.Timestamp `bson:"createdAt"`
	UpdatedAt time.Time           `bson:"updatedAt"`
}

func TestInsert(t *testing.T) {
	initDB()
	s1 := Student{"小红", 12}
	s2 := Student{"小兰", 10}
	s3 := Student{"小黄", 11}

	if _, err := _collection.InsertOne(context.TODO(), s1); err != nil {
		log.Printf("insert data to monodb  failed,err:%v", err)
		return
	}

	students := []interface{}{s2, s3}
	insertManyResult, err := _collection.InsertMany(context.TODO(), students)
	if err != nil {
		log.Printf("insertMany data to monodb  failed,err:%v", err)
	}

	log.Printf("Inserted multiple documents: %v", insertManyResult.InsertedIDs)

	// Create a document with a Timestamp
	timestamp := primitive.Timestamp{T: uint32(time.Now().Unix()), I: 1}
	date := primitive.NewDateTimeFromTime(time.Now())

	document := bson.D{
		{"name", "Alice"},
		{"age", 30},
		{"email", "alice@example.com"},
		{"isStudent", false},
		{"scores", bson.A{85, 90, 88}},
		{"address", bson.D{
			{"street", "123 Main St"},
			{"city", "Wonderland"},
		}},
		{"createdAt", timestamp},
		{"updatedAt", date},
	}
	if _, err := _collection.InsertOne(context.TODO(), document); err != nil {
		log.Printf("insert data to monodb  failed,err:%v", err)
		return
	}
	// 创建时间戳和日期
	createdAt := primitive.Timestamp{T: uint32(time.Now().Unix()), I: 1}
	updatedAt := time.Now()

	// 直接给结构体赋值
	user := &User{
		Name:      "Alice",
		Age:       30,
		Email:     "alice@example.com",
		IsStudent: false,
		Scores:    []int{85, 90, 88},
		Address: Address{
			Street: "123 Main St",
			City:   "Wonderland",
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	for i := 0; i < 5000; i++ {
		user.Age = int(rand.Int31n(100))
		_collection.InsertOne(context.TODO(), user)
	}
}
func TestUpdate(t *testing.T) {
	initDB()
	filter := bson.D{{"name", "小兰"}}
	u := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}

	if _, err := _collection.UpdateOne(context.TODO(), filter, u); err != nil {
		log.Printf("insertMany data to monodb  failed,err:%v", err)
	}

}

func TestFindOne(t *testing.T) {
	initDB()
	// 创建一个Student变量用来接收查询的结果
	var result Student
	filter := bson.D{{"name", "小1兰"}}
	if err := _collection.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		if errors.Is(err, qmgo.ErrNoSuchDocuments) {
			log.Printf("no data found ")
		} else {
			log.Printf("find one data failed %v", err)
		}
	}
	log.Printf("Found a single document: %+v\n", result)
}
func TestFindMany(t *testing.T) {
	initDB()
	filter := bson.D{{"name", "小兰"}}
	result := make([]Student, 0, 10)
	cur, err := _collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("find many data failed %v", err)
	}
	defer cur.Close(context.TODO())
	if err := cur.All(context.TODO(), &result); err != nil {
		log.Printf("decode data failed %v", err)
	}
	log.Printf("Found a many document: %+v", result)
}
func TestTransaction(t *testing.T) {
	initDB()
	session, err := _cli.StartSession()
	if err != nil {
		log.Panicf("start session failed,err:%v", err)
	}
	defer session.EndSession(context.Background())
	_, err = session.WithTransaction(context.Background(), func(sessCtx mongo.SessionContext) (interface{}, error) {
		collection := _cli.Database("testdb").Collection("testcol")
		if _, err := collection.InsertOne(sessCtx, bson.D{{"name", "Alice"}}); err != nil {
			return nil, err
		}

		//return nil, errors.New("test")
		if _, err := collection.InsertOne(sessCtx, bson.D{{"name", "Bob"}}); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		log.Panicf("Transaction failed: %v", err)
	}

	log.Printf("Transaction success")
}
```