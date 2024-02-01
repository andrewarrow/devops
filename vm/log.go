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
	item := `sudo journalctl -u %s.service --since=today > /tmp/t
	sudo cat /tmp/t`
	send := fmt.Sprintf(item, service)
	b, err := exec.Command("ssh", "-i", pem, who+"@"+ip,
		"<<<", send).CombinedOutput()
	fmt.Println(string(b), err == nil)
}
