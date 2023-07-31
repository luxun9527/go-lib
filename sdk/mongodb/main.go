package main

import (
	"context"
	"fmt"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

var cli *qmgo.QmgoClient

func main() {
	initDB()
	insert()
	//update()
	//findOne()
	//findMany()
}
func initDB() {

	// 设置客户端连接配置
	var err error
	cli, err = qmgo.Open(context.Background(), &qmgo.Config{Uri: "mongodb://192.168.2.99:27017", Database: "test", Coll: "student"})
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
