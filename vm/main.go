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

	ip := os.Getenv("VM_IP")
	who := "root"

	if command == "cp" {
		file := os.Args[2]
		dest := os.Args[3]
		who = os.Args[4]
		b, err := exec.Command("scp", "-i", "~/.ssh/"+who, file, who+"@"+ip+":"+dest).CombinedOutput()
		fmt.Println(string(b), err == nil)
	} else if command == "reload" {
		service := os.Args[2]
		// systemctl daemon-reload
		// systemctl enable --now web.service
		// systemctl restart web.service
		list := []string{"systemctl daemon-reload",
			fmt.Sprintf("systemctl enable --now %s.service", service),
			fmt.Sprintf("systemctl restart %s.service", service)}
		for _, item := range list {
			Run(who, ip, item)
		}
	} else if command == "env" {
		guid := PseudoUuid()
		email := os.Args[2]
		domains := os.Args[3]
		list := []string{fmt.Sprintf(`BALANCER_GUID=%s`, guid),
			fmt.Sprintf(`BALANCER_EMAIL=%s`, email),
			fmt.Sprintf(`BALANCER_DOMAINS=%s`, domains)}
		for _, item := range list {
			run := fmt.Sprintf(`"echo '%s' >> /etc/systemd/system/aa.conf"`, item)
			Run(who, ip, run)
		}
	}
}

func Run(who, ip, item string) {
	b, err := exec.Command("ssh", "-i", "~/.ssh/"+who, who+"@"+ip,
		"bash -s", "<<<", item).CombinedOutput()
	fmt.Println(string(b), err == nil)
}
