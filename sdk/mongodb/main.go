package main

import (
	"context"
	"errors"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
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

func main() {
	//https://www.mongodb.com/docs/manual/crud/
	initDB()
	//insert()
	//update()
	//findOne()
	//findMany()
	//aggregate()
	transaction()
}
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
	_dbCli = _cli.Database("test")
	_collection = _dbCli.Collection("test111")
}

type Student struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func insert() {
	s1 := Student{"小红", 12}
	s2 := Student{"小兰", 10}
	s3 := Student{"小黄", 11}

	if _, err := _collection.InsertOne(context.TODO(), s1); err != nil {
		log.Printf("insert data to monodb  failed,err:%v", err)

	}

	students := []interface{}{s2, s3}
	insertManyResult, err := _collection.InsertMany(context.TODO(), students)
	if err != nil {
		log.Printf("insertMany data to monodb  failed,err:%v", err)
	}
	log.Printf("Inserted multiple documents: %v", insertManyResult.InsertedIDs)

}
func update() {
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
func findOne() {
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
func findMany() {

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
func transaction() {
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
