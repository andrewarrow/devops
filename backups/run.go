package main

import (
	"fmt"
	"os/exec"
	"time"
)

func Run(db string) {
	now := time.Now().Unix()
	filename := fmt.Sprintf("%s_%d.sql", db, now)
	b, err := exec.Command("bash", "-c", fmt.Sprintf("pg_dump postgres://fred:fred@localhost:5432/%s > %s", db, filename)).CombinedOutput()
	fmt.Println(string(b), err)
	b, err = exec.Command("gzip", filename).CombinedOutput()
	fmt.Println(string(b), err)
}
