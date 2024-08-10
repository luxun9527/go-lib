1、相关概念，数据类型

2、docker安装集群

3、常用命令

4、gosdk 





## 相关概念

### 核心概念

- **Index（索引）**: Elasticsearch中的索引相当于关系型数据库中的表。一个索引包含多个文档，并且每个文档有相同的结构。
- **Document（文档）**: 文档是Elasticsearch中的最小数据单元，相当于关系型数据库中的一行数据。每个文档是一个JSON对象。
- **Shards（分片）**: 为了支持大规模数据，Elasticsearch将索引分割成多个分片，每个分片可以单独存储和搜索。
- **Replica（副本）**: 副本是分片的拷贝，存在于其他节点上，以提高数据的冗余和可用性。
- **Cluster（集群）**: Elasticsearch集群是由一个或多个节点组成的，并且每个节点是一个Elasticsearch实例。集群有一个唯一的名字，集群中的每个节点共享数据。
- **Node（节点）**: 集群中的每个节点是一个独立的Elasticsearch实例。节点负责存储数据并参与集群的索引和搜索能力。

### 数据类型

- **Text**: 适合全文检索，字段的内容会被分析器（Analyzer）处理，用于搜索和索引。
- **Keyword**: 用于精确匹配，不会被分析器处理，适合排序、聚合和过滤。
- **Numeric**: 包括`long`、`integer`、`short`、`byte`、`double`、`float`、`half_float`等，用于存储数值类型。
- **Date**: 用于存储日期数据，可以以多种格式存储和查询。
- **Boolean**: 用于存储`true`或`false`值。
- **Range**: 用于存储数值、日期、IP地址等的范围值。
- 

## docker安装

https://gitee.com/zhengqingya/docker-compose/tree/master/Linux/elasticsearch

单节点安装参考，使用docker-compose安装。

集群安装待定。

## 常用命令

在kibana中输入命令

点击开发工具进入控制台

https://www.cnblogs.com/machangwei-8/p/14979956.html#_label1



### 新增一个索引

```
 PUT /my_index
```



新增索引指定 mapping  ,mapping类似mysql中字段的类型定义

```json
PUT /my_index
{
  "mappings": {
    "properties": {
      "class_name": {
        "type": "keyword"
      },
      "student_id": {
        "type": "integer"
      },
      "age": {
        "type": "integer"
      }
    }
  }
}
```

也可以不指定mapping在插入的时候会自动推断，字符串类型的默认为text，如何想让字符串类型为keyword，需要提前定义，或者使用动态模板



```json
PUT /my_index
{
  "mappings": {
    "dynamic_templates": [
      {
        "class_name_as_keyword": {
          "match": "class_name",
          "mapping": {
            "type": "keyword"
          }
        }
      }
    ]
  }
}
```

### 获取所有索引 

GET _cat/indices?v

### 获取单个索引 

POST demo-2024.08.03/_search     

### 删除单个索引

```
 DELETE /my_index
```

### 查看所有文档

```
POST /test-index/_search
```

### 插入单个文档

```json
POST /my_index1/_doc/1
 {"class_name": "c1", "student_id": "1","age":12}  
```

### 查看单个文档

```
GET /test-index/_doc/1
```

### 删除单个文档

```
DELETE /test-index/_doc/1
```

### 在单个索引中分词，使用match

```go
POST /test-index/_search
{
  "query": {
    "match_phrase": {
      "content": "中国"
    }
  }
}
```

### 单个文档全等查询,使用term查询

如果查询的字段是text类型分词过的，需要加上,keyword后缀

```go
{
  "query": {
    "term": {
      "content.keyword": "中国北京"
    }
  }
}
```

### 聚合查询

获取文档数量

```go
POST test-index/_search     
{
  "size": 0
}
```

##### 分组查询

统计某个分组某个字段之和

假设你的Elasticsearch索引中包含以下学生数据：

| student_id | class_id | age  |
| ---------- | -------- | ---- |
| 1          | class_1  | 15   |
| 2          | class_1  | 14   |
| 3          | class_1  | 16   |
| 4          | class_1  | 13   |
| 5          | class_2  | 14   |
| 6          | class_2  | 15   |
| 7          | class_2  | 16   |
| 8          | class_3  | 14   |
| 9          | class_3  | 15   |
| 10         | class_3  | 15   |

Elasticsearch聚合查询

我们要统计每个班级的学生年龄之和，使用如下的Elasticsearch查询：

```json
POST /your_index/_search
{
  "size": 0,
  "aggs": {
    "by_class": {
      "terms": {
        "field": "class_id",
        "size": 10
      },
      "aggs": {
        "total_age": {
          "sum": {
            "field": "age"
          }
        }
      }
    }
  }
}
```

查询响应示例

执行上述查询后的响应可能如下：

```json
{
  "took": 5,
  "timed_out": false,
  "_shards": {
    "total": 1,
    "successful": 1,
    "skipped": 0,
    "failed": 0
  },
  "hits": {
    "total": {
      "value": 10,
      "relation": "eq"
    },
    "max_score": null,
    "hits": []
  },
  "aggregations": {
    "by_class": {
      "doc_count_error_upper_bound": 0,
      "sum_other_doc_count": 0,
      "buckets": [
        {
          "key": "class_1",
          "doc_count": 4,
          "total_age": {
            "value": 58.0
          }
        },
        {
          "key": "class_2",
          "doc_count": 3,
          "total_age": {
            "value": 45.0
          }
        },
        {
          "key": "class_3",
          "doc_count": 3,
          "total_age": {
            "value": 44.0
          }
        }
      ]
    }
  }
}
```

响应数据表格化

根据上述响应，查询结果如下表格所示：

| class_id | 学生数量 (doc_count) | 年龄总和 (total_age.value) |
| -------- | -------------------- | -------------------------- |
| class_1  | 4                    | 58                         |
| class_2  | 3                    | 45                         |
| class_3  | 3                    | 44                         |

解释

- `class_id`：班级的标识符。
- `doc_count`：该班级的学生数量。
- `total_age.value`：该班级中所有学生年龄的总和。

这些结果展示了每个班级的学生数量以及他们年龄的总和。这种查询方式可以帮助你快速获取每个分组（例如班级）的统计数据。



## go sdk

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/olivere/elastic/v7"
    "github.com/spf13/cast"
    "log"
    "math/rand"
    "os"
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
```