package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"google.golang.org/api/option"
)

func Run(db string) {
	for {
		fmt.Println("Running")
		RunOnce(db)
		time.Sleep(time.Hour * 24)
	}
}

func RunOnce(db string) {
	now := time.Now().Unix()
	filename := fmt.Sprintf("%s_%d.sql", db, now)
	b, err := exec.Command("bash", "-c", fmt.Sprintf("pg_dump --inserts --table=stores --table=flex_stores --table=configs --table=searches --table=admins --table=links --table=users postgres://fred:fred@localhost:5432/%s > %s", db, filename)).CombinedOutput()
	fmt.Println(string(b), err)
	b, err = exec.Command("gzip", filename).CombinedOutput()
	fmt.Println(string(b), err)
	filename = filename + ".gz"
	asBytes, _ := ioutil.ReadFile(filename)
	storeInGoogleBucket(asBytes, filename)
	//storeInAwsBucket(asBytes, filename)
}

func storeInAwsBucket(data []byte, filename string) {

	bucketName := os.Getenv("STORAGE_BUCKET")

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("us-west-2"),
	)
	fmt.Println(err)

	client := s3.NewFromConfig(cfg)

	dataReader := bytes.NewReader(data)

	putObjectInput := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
		Body:   dataReader,
	}

	_, err = client.PutObject(context.Background(), putObjectInput)
	fmt.Println(err)
}

func storeInGoogleBucket(data []byte, filename string) {
	bucket := os.Getenv("STORAGE_BUCKET")
	keyPath := os.Getenv("KEY_PATH")

	gcsClient, err := storage.NewClient(context.Background(),
		option.WithCredentialsFile(keyPath))
	fmt.Println(err, bucket, keyPath, len(data), filename)

	w := gcsClient.Bucket(bucket).Object(filename).NewWriter(context.Background())
	w.ContentType = "application/octet-stream"
	_, err = w.Write(data)
	fmt.Println("write", err)
	w.Close()
}
