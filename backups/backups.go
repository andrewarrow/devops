package main

import (
	"fmt"
	"os/exec"
)

func Run(db string) {
	b, err := exec.Command("pg_dump",
		"postgres://fred:fred@localhost:5432/"+db,
		">",
		db+".sql").CombinedOutput()
	fmt.Println(string(b), err)
	//s = `"gzip %s.sql"`
}
