package main

import (
	"fmt"
	"io/ioutil"
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
		Scp(who, file, ip, dest)
	} else if command == "reload" {
		service := os.Args[2]
		ReloadService(service, ip)
		Run(who, ip, `"mkdir /certs"`)
	} else if command == "web" {
		file := `[Unit]
Description=web%s
After=network.target network-online.target
Requires=network-online.target

[Service]
User=aa
Group=aa
EnvironmentFile=/etc/systemd/system/aa.conf
ExecStart=/home/aa/web-%s run %s
Restart=on-failure
RestartSec=1s

[Install]
WantedBy=multi-user.target
`
		ports := []string{"3000", "3001"}
		for _, port := range ports {
			file := fmt.Sprintf(file, port, port, port)
			ioutil.WriteFile(port, []byte(file), 0644)
			Scp(who, port, ip, fmt.Sprintf("/etc/systemd/system/web-%s.service", port))
			os.Remove(port)
		}
	} else if command == "psql" {
		//Run("root", ip, "apt install -y postgresql")
		Run("root", ip, `'psql --user=postgres -c "CREATE USER fred WITH SUPERUSER PASSWORD 'fred'"'`)
		Run("root", ip, `'psql --user=postgres -c "CREATE database feedback"'`)
		Run("root", ip, `'psql --user=postgres -c "CREATE EXTENSION IF NOT EXISTS citext"'`)
	} else if command == "deploy-web" {
		domain := os.Args[2]
		DeployWeb(domain, ip)
	} else if command == "deploy-balancer" {
		Scp("aa", "../balancer/balancer", ip, "/home/aa/balancer2")
		Run("root", ip, "systemctl stop balancer.service")
		Run("aa", ip, `"mv /home/aa/balancer2 /home/aa/balancer"`)
		Run("root", ip, "systemctl start balancer.service")
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
		fmt.Println("export BALANCER_GUID=" + guid)
	}
}

func Run(who, ip, item string) {
	b, err := exec.Command("ssh", "-i", "~/.ssh/"+who, who+"@"+ip,
		"bash -s", "<<<", item).CombinedOutput()
	fmt.Println(string(b), err == nil)
}

func Scp(who, file, ip, dest string) {
	b, err := exec.Command("scp", "-i", "~/.ssh/"+who, file, who+"@"+ip+":"+dest).CombinedOutput()
	fmt.Println(string(b), err == nil)
}
