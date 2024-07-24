package main

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
	"testing"
)

type UserDoc struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

var esClient *elastic.Client

func init() {
	var err error
	esClient, err = elastic.NewClient(
		elastic.SetURL("http://192.168.2.159:9200"),
		elastic.SetBasicAuth("elastic", "123456"),
		elastic.SetSniff(false), // 禁用 Sniffing
		//elastic.SetHealthcheck(false),                                      // 禁用健康检查
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)), // 启用错误日志
		elastic.SetInfoLog(log.New(os.Stdout, "ELASTIC ", log.LstdFlags)),  // 启用信息日志
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

func TestGetDocument2(t *testing.T) {
	getResult, err := esClient.Get().
		Index("test-index").
		Id("1").
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

func TestMatchQuery2(t *testing.T) {
	query := elastic.NewMatchQuery("title", "Test")
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
