package main

import (
	"context"
	"fmt"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

var cli *qmgo.QmgoClient

func main() {
	//https://www.mongodb.com/docs/manual/crud/
	initDB()
	//insert()
	//update()
	//findOne()
	//findMany()
	aggregate()
}
func initDB() {

	// 设置客户端连接配置
	var err error
	cli, err = qmgo.Open(context.Background(), &qmgo.Config{Uri: "mongodb://admin:admin123@192.168.254.99:27018", Database: "test", Coll: "bc_chat_box"})
	// 检查连接
	err = cli.Ping(3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
}

type Student struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func insert() {
	s1 := Student{"小红", 12}
	s2 := Student{"小兰", 10}
	s3 := Student{"小黄", 11}
	_, err := cli.InsertOne(context.TODO(), s1)
	if err != nil {
		log.Fatal(err)
	}
	students := []interface{}{s2, s3}
	insertManyResult, err := cli.InsertMany(context.TODO(), students)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

}
func update() {
	filter := bson.D{{"name", "小兰"}}
	u := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}
	err := cli.UpdateOne(context.TODO(), filter, u)
	if err != nil {
		log.Fatal(err)
	}
}
func findOne() {
	// 创建一个Student变量用来接收查询的结果
	var result Student
	filter := bson.D{{"name", "小兰"}}
	if err := cli.Find(context.TODO(), filter).One(&result); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found a single document: %+v\n", result)
}
func findMany() {
	filter := bson.D{{"name", "小兰"}}
	result := make([]Student, 0, 10)
	if err := cli.Find(context.TODO(), filter).All(&result); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)
}

type LastMessage struct {
	SessionID   int64  `bson:"_id,omitempty" json:"id,omitempty"`
	Type        int    `bson:"type" json:"type"`
	UserId      int64  `bson:"user_id" json:"user_id"`
	Content     string `bson:"content" json:"content"`
	CreatedAt   int64  `bson:"created_at,omitempty" json:"updated_at,omitempty"`
	UnReadCount int64  `bson:"un_read_count,omitempty" json:"un_read_count,omitempty"`
}

func aggregate() {
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"user_id", 445299875098629}}}},
		{{"$group", bson.D{{"_id", "$chat_id"},
			{"created_at", bson.D{{"$last", "$created_at"}}},
			{"type", bson.D{{"$last", "$type"}}},
			{"user_id", bson.D{{"$last", "$user_id"}}},
			{"un_read_count", bson.D{{"$sum", bson.D{{"$cond", bson.D{{"if", bson.D{{"$eq", bson.A{"$is_read", 0}}}}, {"then", 1}, {"else", 0}}}}}}},
			{"content", bson.D{{"$last", "$content"}}}}}},
	}
	var msg []*LastMessage
	if err := cli.Aggregate(context.Background(), pipeline).All(&msg); err != nil {
		log.Fatal(err)
	}
	for _, v := range msg {
		log.Printf("%+v", v)
	}
}
