package main

import (
	"fmt"
)

func ReloadService(service, ip string) {
	who := "root"
	list := []string{"systemctl daemon-reload",
		fmt.Sprintf("systemctl enable --now %s.service", service),
		fmt.Sprintf("systemctl restart %s.service", service)}
	for _, item := range list {
		Run(who, ip, item)
	}
}
