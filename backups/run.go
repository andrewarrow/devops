package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func Run(db string) {
	now := time.Now().Unix()
	filename := fmt.Sprintf("%s_%d.sql", db, now)
	b, err := exec.Command("bash", "-c", fmt.Sprintf("pg_dump postgres://fred:fred@localhost:5432/%s > %s", db, filename)).CombinedOutput()
	fmt.Println(string(b), err)
	b, err = exec.Command("gzip", filename).CombinedOutput()
	fmt.Println(string(b), err)
	asBytes, _ := ioutil.ReadFile(filename)
	storeInBucket(asBytes, filename)
}

func storeInBucket(data []byte, filename string) {
	bucket := os.Getenv("STORAGE_BUCKET")
	keyPath := os.Getenv("KEY_PATH")

	gcsClient, err := storage.NewClient(context.Background(),
		option.WithCredentialsFile(keyPath))
	fmt.Println(err, bucket, keyPath)

	w := gcsClient.Bucket(bucket).Object(filename).NewWriter(context.Background())
	w.ContentType = "application/octet-stream"
	_, err = w.Write(data)
	fmt.Println("write", err)
	w.Close()
}
