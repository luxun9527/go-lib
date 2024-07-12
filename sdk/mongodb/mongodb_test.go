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
