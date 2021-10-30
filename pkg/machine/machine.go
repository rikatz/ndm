package machine

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
)

func CreateMachine(template, machine string, userdata, metadata []byte) error {
	userDataBase64 := base64.StdEncoding.EncodeToString(userdata)
	metaDataBase64 := base64.StdEncoding.EncodeToString(metadata)

	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	machineDir := fmt.Sprintf("%s/Virtual Machines.localized/", homedir)

	command := exec.Cmd{
		Path:   "/Applications/VMware OVF Tool/ovftool",
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	args := []string{
		"--allowExtraConfig",
		"--X:enableHiddenProperties",
		"--X:injectOvfEnv",
		"--acceptAllEulas",
	}

	nameArg := fmt.Sprintf("--name=%s", machine)
	args = append(args, nameArg, template, machineDir)

	command.Args = args
	if err := command.Run(); err != nil {
		return err
	}
	// TODO: Turn memory configurable
	vmxfile := fmt.Sprintf("%s/%s.vmwarevm/%s.vmx", machineDir, machine, machine)
	if err := fixParametersInVMX(vmxfile, 2, 8192, userDataBase64, metaDataBase64); err != nil {
		return fmt.Errorf("error setting memory size: %s", err)
	}
	return nil
}

func MachineExist(name string) bool {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	machineDir := fmt.Sprintf("%s/Virtual Machines.localized/%s.vmwarevm/%s.vmx", homedir, name, name)
	if _, err := os.Stat(machineDir); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// This function could be a sed...but....spaces in directories...
func fixParametersInVMX(vmxfile string, cpu, memory int, userdata, metadata string) error {
	// open original file
	f, err := os.Open(vmxfile)
	if err != nil {
		return fmt.Errorf("failed opening file: %s", err)
	}
	defer f.Close()
	// create temp file
	tmp, err := os.CreateTemp("", "template")
	if err != nil {
		return fmt.Errorf("failed creating temporary file: %s", err)
	}
	defer tmp.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()

		if ok, _ := regexp.MatchString("memsize.*", line); ok {
			line = fmt.Sprintf("memsize = \"%d\"", memory)
		}
		if ok, _ := regexp.MatchString("numvcpus.*", line); ok {
			line = fmt.Sprintf("numvcpus = \"%d\"", cpu)
		}
		if ok, _ := regexp.MatchString("ethernet0.connectionType.*", line); ok {
			line = "ethernet0.connectionType = \"nat\""
		}
		if _, err := io.WriteString(tmp, line+"\n"); err != nil {
			return err
		}
	}
	if sc.Err() != nil {
		return sc.Err()
	}
	cloudinitString := "guestinfo.metadata.encoding = \"base64\"\n"
	cloudinitString = cloudinitString + "guestinfo.userdata.encoding = \"base64\"\n"
	cloudinitString = cloudinitString + fmt.Sprintf("guestinfo.userdata = \"%s\"\n", userdata)
	cloudinitString = cloudinitString + fmt.Sprintf("guestinfo.metadata = \"%s\"", metadata)
	if _, err := io.WriteString(tmp, cloudinitString+"\n"); err != nil {
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmp.Name(), vmxfile); err != nil {
		return err
	}
	return nil
}

func MachineStart(name string) (ip string, err error) {
	log.Printf("Starting Virtual Machine %s", name)
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	machineDir := fmt.Sprintf("%s/Virtual Machines.localized/%s.vmwarevm/%s.vmx", homedir, name, name)
	if err := exec.Command("/Applications/VMware Fusion.app/Contents/Public/vmrun", "upgradevm", machineDir).Run(); err != nil {
		return "", fmt.Errorf("failed to upgrade machine virtual hardware: %s", err)
	}

	if err := exec.Command("/Applications/VMware Fusion.app/Contents/Public/vmrun", "start", machineDir, "nogui").Run(); err != nil {
		return "", fmt.Errorf("failed to start machine virtual: %s", err)
	}

	ipByte, err := exec.Command("/Applications/VMware Fusion.app/Contents/Public/vmrun", "getGuestIPAddress", machineDir, "-wait").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to obtain virtual machine IP: %s", err)
	}
	log.Printf("Machine %s started", name)
	return string(ipByte), nil
}
