package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
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
		key := os.Getenv("SSH_KEY")
		ip := os.Getenv("VM_IP")
		file := os.Args[2]
		dest := os.Args[3]
		who := os.Args[4]
		b, err := exec.Command("scp", "-i", "~/.ssh/"+key, file, who+"@"+ip+":"+dest).CombinedOutput()
		fmt.Println(string(b), err == nil)
	} else if command == "" {
	} else if command == "" {
	}
}
