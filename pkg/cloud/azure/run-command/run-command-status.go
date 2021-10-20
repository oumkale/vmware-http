package runcommand

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/litmuschaos/litmus-go/pkg/log"
)

func GetRunCommandResult(result *compute.RunCommandResult) {
	message := strings.Split(strings.TrimSuffix(*(*result.Value)[0].Message, "\n"), "\n")

	i := 0

	for ; i < len(message) && message[i] != "[stdout]"; i++ {
	}

	if i < len(message)-1 && message[i+1] != "" {
		i++
		for ; i < len(message) && message[i] != "[stderr]"; i++ {
			if message[i] != "" {
				log.Infof("[stdout]: %v\n", message[i])
			}
		}
	}

	if i < len(message)-1 && message[i+1] != "" {
		i++
		for ; i < len(message) && message[i] != ""; i++ {
			log.Errorf("[stderr]: %v\n", message[i])
		}
	}
}
