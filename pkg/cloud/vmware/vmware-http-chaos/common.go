package http

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/litmuschaos/litmus-go/pkg/log"
	"github.com/litmuschaos/litmus-go/pkg/utils/retry"
	experimentTypes "github.com/litmuschaos/litmus-go/pkg/vmware/vmware-http-chaos/types"
	"github.com/pkg/errors"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/guest/toolbox"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

// StringInSlice will check and return whether a string is present inside a slice or not
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// GetVMDetails fetch the details for VMWare instance
func GetVMDetails(experimentsDetails *experimentTypes.ExperimentDetails) (*object.VirtualMachine, error) {

	// Creating a connection context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Parsing URL
	url, err := url.Parse("https://administrator@vsphere.local:Chaosnative@1658@103.195.244.216/sdk")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	// Connecting to vCenter
	client, err := govmomi.NewClient(ctx, url, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	// Selecting default datacenter
	vm, err := find.NewFinder(client.Client).VirtualMachine(ctx, "win-vm")
	if err != nil {
		fmt.Fprintf(os.Stderr, "1Error: %s\n", err)
		os.Exit(1)
	}
	return vm, nil
}

// GetScaleSetNameAndInstanceId extracts the scale set name and VM id from the instance name
func GetScaleSetNameAndInstanceId(instanceName string) (string, string) {
	scaleSetAndInstanceId := strings.Split(instanceName, "_")
	return scaleSetAndInstanceId[0], scaleSetAndInstanceId[1]
}

func UploadFiles(tools *toolbox.Client, experimentsDetails *experimentTypes.ExperimentDetails, parameters *map[string]string) error {

	err := Upload(tools, experimentsDetails.ScriptPath, parameters)
	if err != nil {
		fmt.Println(err)
	}
	err = Upload(tools, experimentsDetails.AbortScriptPath, parameters)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func Upload(tools *toolbox.Client, scriptPath string, parameters *map[string]string) error {
	// Creating a connection context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	data, err := ioutil.ReadFile(scriptPath)
	if err != nil {
		fmt.Println(err)
	}
	p := soap.DefaultUpload
	file, err := os.Open(filepath.Clean(scriptPath))
	if err != nil {
		fmt.Fprintf(os.Stderr, "run :Error: %s\n", err)
		return err
	}
	defer file.Close()

	// Write some text line-by-line to file.
	for key, element := range *parameters {
		file.WriteString("\n" + key + "=" + element + "\n")
	}
	file.WriteString(string(data))
	err = tools.Upload(ctx, file, experimentsDetails.ScriptPathDesVM, p, &types.GuestFileAttributes{}, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "upload :Error: %s\n", err)
		return err
	}
	return nil
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func Run(cmd *exec.Cmd, tools *toolbox.Client) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := tools.Run(ctx, cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "run :Error: %s\n", err)
		os.Exit(1)
	}
	return nil
}

// WaitForServiceStop will wait for the service to completely stop
func WaitForCompleion(cmd *exec.Cmd, serviceName, delay, timeout int) error {

	log.Infof("[Status]: Checking %s toxiproxy server status for running", serviceName)
	return retry.
		Times(uint(timeout / delay)).
		Wait(time.Duration(delay) * time.Second).
		Try(func(attempt uint) error {

			serviceState := CheckRunning(cmd)
			if serviceState != true {
				log.Infof("[Info]: The toxiproxy server state is %s", serviceState)
				return errors.Errorf("%s toxiproxy server is not yet in inactive state", serviceName)
			}

			log.Infof("[Info]: The toxiproxy server state is %s", serviceState)
			return nil
		})
}

func CheckRunning(cmd *exec.Cmd) bool {
	if cmd == nil || cmd.ProcessState != nil && cmd.ProcessState.Exited() || cmd.Process == nil {
		return false
	}

	return true
}
