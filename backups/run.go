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

	gcsClient, _ := storage.NewClient(context.Background(),
		option.WithCredentialsFile(keyPath))

	w := gcsClient.Bucket(bucket).Object(filename).NewWriter(context.Background())
	w.ContentType = "application/octet-stream"
	w.Write(data)
	w.Close()
}
