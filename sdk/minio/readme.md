# 目录

1. minio 概念
2. minio安装、集群安装
3. 客户端常用命令
4. go sdk

# minio概念

[**https://juejin.cn/post/7108553532776644622**](https://juejin.cn/post/7108553532776644622)



**chatgpt:**

**MinIO 是一个高性能、分布式的对象存储系统，兼容 Amazon S3 API。它通常用于存储海量非结构化数据，如图片、视频、备份和日志文件等。MinIO 设计理念是简单、快速并且支持云原生架构。以下是 MinIO 的一些关键概念和特点：**

### **1.** **对象存储**

**对象存储是 MinIO 的核心概念，区别于传统的文件系统存储和块存储。数据被存储为对象，每个对象都包含数据、元数据和唯一的标识符。与文件存储不同，对象存储没有文件层次结构（如目录），而是通过对象键（Object Key）来标识数据。**

### **2.** **桶（Bucket）**

- **桶（Bucket）** **是存储对象的容器。MinIO 中的每个对象都必须存储在一个桶内，类似于 S3 中的 bucket。**
- **每个桶是全局唯一的，并且桶名需要遵循特定规则，通常必须是小写字母、数字、破折号或下划线的组合。**
- **你可以为不同的用途创建多个桶来组织数据。**

### **3.** **对象（Object）**

- **对象** **是 MinIO 存储的数据单位，包含三部分：**

- - **数据：实际存储的内容（例如图片、视频、文档等）。**
  - **元数据：描述对象的相关信息，如创建时间、修改时间、文件类型等。**
  - **对象键（Object Key）：唯一标识该对象的字符串，类似于文件系统中的路径。**

### **4.** **S3 兼容性**

- **MinIO 是兼容** **Amazon S3 API** **的，即支持使用 S3 的常见操作，如上传、下载、删除文件和创建桶等。通过 MinIO 提供的 API、CLI 或 SDK，可以像使用 AWS S3 一样管理对象。**
- **S3 API** **支持常见的对象存储操作，如：**`**PUT**`**,** `**GET**`**,** `**DELETE**`**,** `**LIST**`**,** `**HEAD**`**。**

### **5.** **分布式架构**

- **MinIO 支持** **分布式部署****，即可以通过多个节点组成一个分布式集群来提高存储的容错性和可扩展性。**
- **在分布式模式下，MinIO 会使用** **Erasure Coding****（纠删码）来实现数据的冗余存储，确保即使部分节点或磁盘出现故障，也能保证数据的安全性。**
- **分布式 MinIO 支持水平扩展，只需增加更多节点来提升存储容量和性能。**

### **6.** **版本控制（Object Versioning）**

- **MinIO 支持对象版本控制，允许你为同一个对象存储多个版本。每次修改对象时，MinIO 会保留旧版本，并允许用户根据需要恢复到任何历史版本。**

### **7.** **访问控制与权限**

- **MinIO 提供了** **基于策略的访问控制****，可以通过** **桶策略** **来设置权限。**
- **策略允许定义哪些用户或应用可以访问特定的桶或对象。MinIO 支持** **多种认证方式****：**

- - **Access Key 和 Secret Key****：类似于 AWS 的访问密钥。**
  - **STSToken（临时令牌）****：支持动态生成的临时访问令牌。**
  - **IAM Policies****：可以通过策略来限制对桶或对象的访问权限。**

### **8.** **镜像和备份**

- **MinIO 提供了** **复制** **和** **数据同步** **功能，可以在多个 MinIO 实例之间进行数据复制，实现数据的高可用性和灾难恢复。**
- **还支持将数据备份到其他对象存储系统，如 Amazon S3、Google Cloud Storage 等，或者其他 MinIO 集群。**

### **9.** **性能与扩展性**

- **MinIO 在设计上强调高性能，采用并行化的 IO 操作来实现极高的吞吐量，适用于大规模数据存储需求。**
- **支持** **分布式存储集群****，可以横向扩展，以适应不断增加的数据存储需求。**

### **10.** **MinIO 客户端与 SDK**

- **MinIO 提供了多种客户端工具和 SDK，支持常见的编程语言，如 Go、Java、Python、JavaScript 等，可以直接进行对象存储的访问和操作。**
- **MinIO 还提供了** `**mc**`**（MinIO Client）工具，支持命令行操作，如上传、下载、列出、删除对象等。**

### **11.** **事件通知**

- **MinIO 支持基于对象操作的** **事件通知****，例如，当对象被上传或删除时，MinIO 可以触发 HTTP 请求、AWS SQS 消息或 Kafka 消息等通知。**

### **12.** **安全性**

- **MinIO 提供了** **TLS 加密** **来确保数据传输过程中的安全性。**
- **对象数据支持** **服务器端加密****，可以在上传或存储时自动加密数据，支持多种加密方式，如 AES-256。**
- **支持** **静态数据加密****，加密存储在硬盘上的数据。**

### 

# minio部署

## 单节点

[**https://github.com/bitnami/containers/tree/main/bitnami/minio**](https://github.com/bitnami/containers/tree/main/bitnami/minio)

```yaml
# Copyright Broadcom, Inc. All Rights Reserved.
# SPDX-License-Identifier: APACHE-2.0
#https://github.com/bitnami/containers/blob/main/bitnami/minio/docker-compose-distributed.yml
services:
  minio:
    image: docker.fxxk.dedyn.io/bitnami/minio:2024
    ports:
      - '9000:9000'
      - '9001:9001'
    volumes:
      - './data:/bitnami/minio/data'
    environment:
      - MINIO_ROOT_USER=minio
      - MINIO_ROOT_PASSWORD=miniosecret
```

# 常用命令

# go sdk

minio是s3协议，可以使用aws提供的sdk来操作。

```go
package minio

import (
    "bytes"
    "context"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/gookit/goutil/byteutil"
    "io"
    "log"
    "net/http"
    "os"
    "testing"
)

var (
    s3Client        *s3.Client
    bucketName      = "test1"
    filePath        = "./test.txt"
    constEndPoint   = "http://192.168.2.159:9000"
    region          = "us-east-1"
    accessKeyID     = "minio"
    secretAccessKey = "miniosecret"
)

// minio 兼容s3协议，可以使用aws sdk进行操作
func init() {
    var f1 aws.CredentialsProviderFunc = func(ctx context.Context) (aws.Credentials, error) {
        return aws.Credentials{
            AccessKeyID:     accessKeyID,
            SecretAccessKey: "secretAccessKey",
        }, nil
    }

    s3Client = s3.NewFromConfig(aws.Config{
        Region:       region,
        Credentials:  f1,
        BaseEndpoint: aws.String(constEndPoint),
    })
}
func TestCreateBucket(t *testing.T) {

    _, err := s3Client.CreateBucket(context.Background(), &s3.CreateBucketInput{
        Bucket: aws.String(bucketName),
    })
    if err != nil {
        log.Printf("%v", err)
    }
    log.Printf("[INFO] create bucket %s success", bucketName)
}

func TestUploadFile(t *testing.T) {
    file, err := os.Open(filePath)
    if err != nil {
        log.Printf("[ERROR] Failed to read file %s", err.Error())
        return
    }
    objectName := byteutil.Md5(filePath)

    if _, err := s3Client.PutObject(context.Background(), &s3.PutObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(string(objectName)),
        Body:   file,
    }); err != nil {
        log.Printf("[ERROR] Failed to upload %s to, ", err.Error())
        return
    }
    log.Printf("[INFO] Uploaded %s to %s", objectName, bucketName)
}
func TestDeleteFile(t *testing.T) {

    objectName := byteutil.Md5(filePath)

    if _, err := s3Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(string(objectName)),
    }); err != nil {
        log.Printf("[ERROR] Failed to delete %s from, ", err.Error())
        return
    }
    log.Printf("[INFO] Deleted %s from %s", objectName, bucketName)

}
func TestDeleteBucket(t *testing.T) {

    if _, err := s3Client.DeleteBucket(context.Background(), &s3.DeleteBucketInput{
        Bucket:              aws.String(bucketName),
        ExpectedBucketOwner: nil,
    }); err != nil {
        log.Printf("[ERROR] Failed to delete %s from, ", err.Error())
        return
    }
    log.Printf("[INFO] Deleted %s from %s", bucketName, bucketName)

}

// 生成预上传的url
func TestGetuploadUrl(t *testing.T) {
    var s PreSigner = NewPreSigner(s3Client)
    //生成预签名上传url，url的有效期为1小时。
    url, err := s.PutObject(bucketName, "presignerUploadData", 3600)
    if err != nil {
        log.Printf("[ERROR] Failed to get upload url, err %v", err.Error())
    }
    log.Printf("[INFO] Get upload url success, url %s", url.URL)
    req, err := http.NewRequest("PUT", url.URL, bytes.NewReader([]byte("test")))
    if err != nil {
        log.Printf("[ERROR] Failed to create request, err %v", err.Error())
        return
    }
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Printf("[ERROR] Failed to send request, err %v", err.Error())
        return
    }
    if resp.StatusCode != http.StatusOK {
        log.Printf("[ERROR] Failed to upload file, status code %d", resp.StatusCode)
        return
    }
    data, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("[ERROR] Failed to read request body, err %v", err.Error())
        return
    }
    log.Printf("[INFO] Request body %s", data)

}
```