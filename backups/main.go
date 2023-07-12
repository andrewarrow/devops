package main

import (
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if len(os.Args) == 1 {
		//PrintHelp()
		return
	}
	command := os.Args[1]

	if command == "cp" {
	} else if command == "run" {
		db := os.Args[2]
		Run(db)
	}
}
