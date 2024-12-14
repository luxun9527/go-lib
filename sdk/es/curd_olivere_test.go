package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/luxun9527/zlog"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/cast"
	"go.uber.org/zap/zapcore"
	"log"
	"math/rand"
	"testing"
	"time"
)

type UserDoc struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Student struct {
	Id        int64
	ClassName string
	Age       int32
}

var esClient *elastic.Client

func init() {
	var err error
	zlog.DevConfig.UpdateLevel(zapcore.DebugLevel)
	esClient, err = elastic.NewClient(
		elastic.SetURL("http://192.168.2.159:9200"),
		elastic.SetBasicAuth("elastic", "123456"),
		elastic.SetSniff(false), // 禁用 Sniffing
		//elastic.SetHealthcheck(false),                                      // 禁用健康检查
		elastic.SetErrorLog(zlog.ErrorEsOlivereLogger), // 启用错误日志
		elastic.SetInfoLog(zlog.InfoEsOlivereLogger),   // 启用信息日志
	)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
}

func TestIndexDocument2(t *testing.T) {
	doc := UserDoc{Title: "Test UserDoc", Content: "This is the content of the test document."}
	_, err := esClient.Index().
		Index("test-index").
		Id("4").
		BodyJson(doc).
		Refresh("true").
		Do(context.Background())
	if err != nil {
		t.Fatalf("Error indexing document: %s", err)
	}
	log.Println("UserDoc indexed successfully")
}

// 插入一条分词数据
func TestIndexDocument3(t *testing.T) {
	doc := UserDoc{Title: "地址", Content: "中国北京"}
	_, err := esClient.Index().
		Index("test-index").
		Id("4").
		BodyJson(doc).
		Refresh("true").
		Do(context.Background())
	if err != nil {
		t.Fatalf("Error indexing document: %s", err)
	}
	log.Println("UserDoc indexed successfully")
}

func TestGetDocument2(t *testing.T) {
	getResult, err := esClient.Get().
		Index("test-index").
		Id("4").
		Do(context.Background())
	if err != nil {
		t.Fatalf("Error getting document: %s", err)
	}
	if getResult.Found {
		var doc UserDoc
		if err := json.Unmarshal(getResult.Source, &doc); err != nil {
			t.Fatalf("Error parsing the document: %s", err)
		}
		log.Printf("Got document: %+v", doc)
	} else {
		t.Fatalf("UserDoc not found")
	}
}

func TestUpdateDocument2(t *testing.T) {
	updateDoc := map[string]interface{}{
		"doc": map[string]interface{}{
			"content": "This is the updated content of the test document.",
		},
	}
	_, err := esClient.Update().
		Index("test-index").
		Id("1").
		Doc(updateDoc).
		Refresh("true").
		Do(context.Background())
	if err != nil {
		t.Fatalf("Error updating document: %s", err)
	}
	log.Println("UserDoc updated successfully")
}

func TestDeleteDocument2(t *testing.T) {
	_, err := esClient.Delete().
		Index("test-index").
		Id("1").
		Refresh("true").
		Do(context.Background())
	if err != nil {
		t.Fatalf("Error deleting document: %s", err)
	}
	log.Println("UserDoc deleted successfully")
}

func TestTermQuery2(t *testing.T) {
	query := elastic.NewTermQuery("content.keyword", "中国北京")
	searchResult, err := esClient.Search().
		Index("test-index").
		Query(query).
		Do(context.Background())
	if err != nil {
		t.Fatalf("Error getting search response: %s", err)
	}
	log.Printf("Found %d documents", searchResult.TotalHits())
	for _, hit := range searchResult.Hits.Hits {
		var doc UserDoc
		if err := json.Unmarshal(hit.Source, &doc); err != nil {
			t.Fatalf("Error parsing search response: %s", err)
		}
		log.Printf("Found document: %+v", doc)
	}
}

func TestMatchQuery2(t *testing.T) {
	query := elastic.NewMatchQuery("content", "中国")
	searchResult, err := esClient.Search().
		Index("test-index").
		Query(query).
		Do(context.Background())
	if err != nil {
		t.Fatalf("Error getting search response: %s", err)
	}
	log.Printf("Found %d documents", searchResult.TotalHits())
	for _, hit := range searchResult.Hits.Hits {
		var doc UserDoc
		if err := json.Unmarshal(hit.Source, &doc); err != nil {
			t.Fatalf("Error parsing search response: %s", err)
		}
		log.Printf("Found document: %+v", doc)
	}
}

func TestAnalyzedQuery2(t *testing.T) {
	query := elastic.NewMatchQuery("content", "test document content")
	searchResult, err := esClient.Search().
		Index("test-index").
		Query(query).
		Do(context.Background())
	if err != nil {
		t.Fatalf("Error getting search response: %s", err)
	}
	log.Printf("Found %d documents", searchResult.TotalHits())
	for _, hit := range searchResult.Hits.Hits {
		var doc UserDoc
		if err := json.Unmarshal(hit.Source, &doc); err != nil {
			t.Fatalf("Error parsing search response: %s", err)
		}
		log.Printf("Found document: %+v", doc)
	}
}

func TestBlukInsert(t *testing.T) {
	for i := 0; i < 500; i++ {
		rand.Seed(time.Now().UnixNano())

		// 生成0到99之间的随机整数
		age := rand.Intn(100)
		className := fmt.Sprintf("class_%v", i/100)
		stu := Student{
			Id:        time.Now().UnixNano(),
			ClassName: className,
			Age:       int32(age),
		}
		_, err := esClient.Index().
			Index("test-index").
			Id(cast.ToString(i + 1)).
			BodyJson(stu).
			Refresh("true").
			Do(context.Background())
		if err != nil {
			t.Fatalf("Error indexing document: %s", err)
		}
		log.Println("Indexed document: ", stu.Id)
	}

}

// 测试聚合查询
/*
1、分组聚合
*/
func TestAggQuery1(t *testing.T) {
	agg := elastic.NewTermsAggregation().
		Field("ClassName.keyword").
		Size(3).
		SubAggregation("total_age", elastic.NewSumAggregation().Field("Age"))
	searchResult, err := esClient.Search().
		Index("test-index"). // 替换为你的索引名称
		Size(0).
		Aggregation("by_class", agg).
		Do(context.Background())
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	// 解析聚合结果
	byClassAgg, found := searchResult.Aggregations.Terms("by_class")
	if !found {
		log.Fatalf("Aggregation 'by_class' not found")
	}
	for _, bucket := range byClassAgg.Buckets {
		// 获取班级名称
		className := bucket.Key.(string)
		// 获取年龄之和
		totalAge, found := bucket.Sum("total_age")
		if found {
			fmt.Printf("Class: %s, Total Age: %f\n", className, *totalAge.Value)
		} else {
			fmt.Printf("Class: %s, Total Age not found\n", className)
		}
	}

}
