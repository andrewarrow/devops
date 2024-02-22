package main

import (
	"fmt"
)

func ReloadService(service, ip string) {
	who := "root"
	list := []string{`"systemctl daemon-reload"`,
		fmt.Sprintf(`"systemctl enable %s.service"`, service),
		fmt.Sprintf(`"systemctl start %s.service"`, service)}
	for _, item := range list {
		Run(who, ip, item)
	}
}
