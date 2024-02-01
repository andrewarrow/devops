package main

import (
	"fmt"
	"os"
	"os/exec"
)

func InstallEtc() {

	//	tar -cf send.tar etc
	b, err := exec.Command("tar", "-cf", "send.tar", "etc").CombinedOutput()
	fmt.Println(string(b), err == nil)

	//scp -i ~/.ssh/aa-iot.pem devops/send.tar ec2-user@54.84.223.125:
	pem := os.Getenv("SSH_PEM")
	who := os.Getenv("SSH_USER")
	ip := os.Getenv("VM_IP")

	b, err = exec.Command("scp", "-i", pem, "send.tar", who+"@"+ip+":").CombinedOutput()
	fmt.Println(string(b), err == nil)
}
