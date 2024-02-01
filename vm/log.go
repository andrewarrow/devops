package main

import (
	"fmt"
	"os"
	"os/exec"
)

func GetLog(service string) {
	pem := os.Getenv("SSH_PEM")
	who := os.Getenv("SSH_USER")
	ip := os.Getenv("VM_IP")
	item := `sudo journalctl -u %s.service --since "1 hour ago" > /tmp/t; cat /tmp/t`
	send := fmt.Sprintf(item, service)
	b, err := exec.Command("ssh", "-i", pem, who+"@"+ip,
		"bash -s", "<<<", send).CombinedOutput()
	fmt.Println(string(b), err == nil)
}
