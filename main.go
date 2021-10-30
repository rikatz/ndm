package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rikatz/ndm/pkg/cloudinit"
	"github.com/rikatz/ndm/pkg/machine"
	"github.com/rikatz/ndm/pkg/photon"
)

var (
	requiredFiles = map[string]string{
		"/Applications/VMware OVF Tool/ovftool":                 "VMware OVFTool does not exist. Please download it from https://developer.vmware.com/tool/ovf",
		"/Applications/VMware Fusion.app/Contents/Public/vmrun": "VMware Fusion not installed. Please download it from https://www.vmware.com/products/fusion.html",
	}
)

func main() {
	verifyBinaries()

	file, err := photon.CheckPhotonOVA("")
	if err != nil {
		log.Fatalf("Error getting Photon OVA: %s", err)
	}
	userdata, metadata, err := cloudinit.GenerateCloudInit()
	if err != nil {
		log.Fatalf("Error generating metadata: %s", err)
	}
	if !machine.MachineExist("runner") {
		log.Printf("Machine %s does not exists, creating", "runner")
		err = machine.CreateMachine(file, "runner", userdata, metadata)
		if err != nil {
			log.Fatalf("Failed creating the Runner VM: %s", err)
		}
	}
	ip, err := machine.MachineStart("runner")
	if err != nil {
		log.Fatalf("failed to start machine: %s", err)
	}
	fmt.Println("==========================")
	fmt.Println("Docker Machine Started")
	fmt.Println("* Please execute the following to add the SSH Key to the well-known hosts:") // TODO: Automate this
	fmt.Printf("ssh ndm@%s exit\n\n", strings.TrimSuffix(ip, "\n"))
	fmt.Println("* Now just use the docker daemon there:")
	fmt.Printf("export DOCKER_HOST=ssh://ndm@%s", ip)

}

func verifyBinaries() {
	var errMsg string
	for app, msg := range requiredFiles {
		if _, err := os.Stat(app); err != nil {
			errMsg = fmt.Sprintf("%s\n%s", errMsg, msg)
		}
	}
	if errMsg != "" {
		log.Fatalf(errMsg)
	}
}
