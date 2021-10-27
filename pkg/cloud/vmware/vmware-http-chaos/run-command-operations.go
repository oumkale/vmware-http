package runcommand

import (
	"context"
	"os/exec"
	"strings"

	experimentTypes "github.com/litmuschaos/litmus-go/pkg/vmware/vmware-http-chaos/types"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"

	"github.com/pkg/errors"
)

func PerformRunCommand(runCommandFuture *experimentTypes.RunCommandFuture, runCommandInput compute.RunCommandInput, subscriptionID, resourceGroup, azureInstanceName, scaleSet string) error {

	future, err := vmssClient.RunCommand(context.TODO(), resourceGroup, scaleSetName, vmId, runCommandInput)
	if err != nil {
		return errors.Errorf("failed to perform run command, err: %v", err)
	}
	return nil
}

// GetScaleSetNameAndInstanceId extracts the scale set name and VM id from the instance name
func GetScaleSetNameAndInstanceId(instanceName string) (string, string) {
	scaleSetAndInstanceId := strings.Split(instanceName, "_")
	return scaleSetAndInstanceId[0], scaleSetAndInstanceId[1]
}

// prepareRunCommandInput prepares the run command with the script, parameters and commandID
func PrepareRunCommandInput(experimentsDetails *experimentTypes.ExperimentDetails, parameters *map[string]string) (*exec.Cmd, error) {

	var dir string
	// Setting up command id based on operating system of VM
	if experimentsDetails.OperatingSystem == "windows" {
		experimentsDetails.ScriptPathDesVM = "C:\\Users\\inject.ps1"
		dir = "C:\\Users\\"
	} else {
		experimentsDetails.ScriptPathDesVM = "/usr/bin/inject.sh"
		dir = "/usr/bin/"
	}

	if strings.TrimSpace(experimentsDetails.ScriptPath) == "" {
		return &exec.Cmd{}, errors.Errorf("no script provided, either provide a custom script or use default script")
	}

	cmd := &exec.Cmd{
		Path:  "",
		Args:  []string{`echo ` + experimentsDetails.VMPassword + ` | sudo -S bash inject.sh`},
		Env:   []string{`SUDO_AKSPASS=` + experimentsDetails.VMPassword},
		Dir:   dir,
		Stdin: nil,
	}

	return cmd, nil
}

// prepareAbortRunCommandInput prepares abort script command
func PrepareAbortRunCommandInput(experimentsDetails *experimentTypes.ExperimentDetails, parameters *map[string]string) (*exec.Cmd, error) {

	var dir string
	// Setting up command id based on operating system of VM
	if experimentsDetails.OperatingSystem == "windows" {
		experimentsDetails.AbortScriptPath = "C:\\Users\\abort.ps1"
		dir = "C:\\Users\\"
	} else {
		experimentsDetails.AbortScriptPath = "/usr/bin/abort.sh"
		dir = "/usr/bin/"
	}
	// Checking if custom script path is provided or using the default if not
	if strings.TrimSpace(experimentsDetails.ScriptPath) != "" {
		return &exec.Cmd{}, errors.Errorf("no script provided, either provide a custom script or use default script")
	}

	cmd := &exec.Cmd{
		Path:  "",
		Args:  []string{`echo ` + experimentsDetails.VMPassword + ` | sudo -S bash abort.sh`},
		Env:   []string{`SUDO_AKSPASS=` + experimentsDetails.VMPassword},
		Dir:   dir,
		Stdin: nil,
	}

	return cmd, nil
}
