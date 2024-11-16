package minio

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gookit/goutil/byteutil"
	"github.com/zeromicro/go-zero/core/stringx"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
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
			SecretAccessKey: secretAccessKey,
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
	//生成预签名上传url，url的有效期为1小时。
	key := stringx.Randn(10)
	url, err := s3.NewPresignClient(s3Client).PresignPutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(3600 * int64(time.Second))
	})
	if err != nil {
		log.Printf("[ERROR] Failed to create request, err %v", err.Error())
		return
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

	log.Printf("[INFO] Request body %s key %s", data, key)

}
