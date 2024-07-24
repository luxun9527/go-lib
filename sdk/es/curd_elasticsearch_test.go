package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type Document struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

var es *elasticsearch.Client

func init() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://192.168.2.159:9200",
		},
		Username: "elastic",
		Password: "123456",
	}

	var err error
	es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
}

func TestIndexDocument(t *testing.T) {
	doc := Document{Title: "Test UserDoc", Content: "This is the content of the test document."}
	docJSON, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("Error marshaling document: %s", err)
	}

	req := esapi.IndexRequest{
		Index:      "test-index",
		DocumentID: "6",
		Body:       bytes.NewReader(docJSON),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		t.Fatalf("Error indexing document: %s", err)
	}
	defer res.Body.Close()
	log.Println(res)
}

func TestGetDocument(t *testing.T) {
	req := esapi.GetRequest{
		Index:      "test-index",
		DocumentID: "1",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		t.Fatalf("Error getting document: %s", err)
	}
	defer res.Body.Close()

	var getDoc map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&getDoc); err != nil {
		t.Fatalf("Error parsing the document: %s", err)
	}
	log.Println(getDoc)
}

func TestUpdateDocument(t *testing.T) {
	updateDoc := map[string]interface{}{
		"doc": map[string]interface{}{
			"content": "This is the updated content of the test document.",
		},
	}
	updateJSON, err := json.Marshal(updateDoc)
	if err != nil {
		t.Fatalf("Error marshaling update document: %s", err)
	}

	req := esapi.UpdateRequest{
		Index:      "test-index",
		DocumentID: "1",
		Body:       bytes.NewReader(updateJSON),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		t.Fatalf("Error updating document: %s", err)
	}
	defer res.Body.Close()
	log.Println(res)
}

func TestDeleteDocument(t *testing.T) {
	req := esapi.DeleteRequest{
		Index:      "test-index",
		DocumentID: "1",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		t.Fatalf("Error deleting document: %s", err)
	}
	defer res.Body.Close()
	log.Println(res)
}

func TestMatchQuery(t *testing.T) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "Test",
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		t.Fatalf("Error encoding query: %s", err)
	}

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("test-index"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		t.Fatalf("Error getting search response: %s", err)
	}

	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		t.Fatalf("Error parsing search response: %s", err)
	}
	log.Println(result)
}

func TestAnalyzedQuery(t *testing.T) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"content": "test document content",
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		t.Fatalf("Error encoding analyzed query: %s", err)
	}

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("test-index"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		t.Fatalf("Error getting analyzed search response: %s", err)
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		t.Fatalf("Error parsing analyzed search response: %s", err)
	}
	log.Println(result)
}
