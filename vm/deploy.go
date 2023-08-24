package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func DeployWeb2(domain, ip string) {
	guid := os.Getenv("BALANCER_GUID")
	script := fmt.Sprintf(deployScript, domain, guid, domain, guid, domain, guid)
	fmt.Println(script)
}
func DeployWeb(domain, ip string) {
	guid := os.Getenv("BALANCER_GUID")
	script := fmt.Sprintf(deployScript, domain, guid, domain, guid, domain, guid)
	//fmt.Println(script)

	deploy3000 := fmt.Sprintf(runScriptDeploy, ip, "3000")
	deploy3001 := fmt.Sprintf(runScriptDeploy, ip, "3001")
	script3000 := fmt.Sprintf(runScript, "3000", "3000", "3000")
	script3001 := fmt.Sprintf(runScript, "3001", "3001", "3001")

	ioutil.WriteFile("deploy.sh", []byte(script), 0755)
	ioutil.WriteFile("deploy-3000.sh", []byte(deploy3000), 0755)
	ioutil.WriteFile("deploy-3001.sh", []byte(deploy3001), 0755)
	ioutil.WriteFile("script-3000.sh", []byte(script3000), 0755)
	ioutil.WriteFile("script-3001.sh", []byte(script3001), 0755)
	Scp("aa", "../web/web", ip, "/home/aa/web")
	b, err := exec.Command("./deploy.sh").CombinedOutput()
	fmt.Println(string(b), err == nil)

	os.Remove("deploy.sh")
	os.Remove("deploy-3000.sh")
	os.Remove("deploy-3001.sh")
	os.Remove("script-3000.sh")
	os.Remove("script-3001.sh")
}

var runScriptDeploy = `ssh -i ~/.ssh/aa aa@%s 'bash -s' < script-%s.sh`
var runScript = `sudo systemctl stop web-%s.service
mv /home/aa/web /home/aa/web-%s
sudo systemctl start web-%s.service
`

var deployScript = `#!/bin/bash

check_result() {
  if [ "$1" == "3000" ]; then
    "$script_dir/deploy-3001.sh"
  elif [ "$1" == "3001" ]; then
    "$script_dir/deploy-3000.sh"
  else
    "$script_dir/deploy-3000.sh"
  fi
}

script_dir=$(dirname "$0")

result=$(curl -s "https://%s/%s/web")

check_result "$result"

echo "switch"
curl -s "https://%s/%s/3000"

new_result=$(curl -s "https://%s/%s/web")
echo "$new_result"
`
