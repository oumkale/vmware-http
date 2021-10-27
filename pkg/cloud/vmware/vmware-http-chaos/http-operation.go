package http

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/litmuschaos/litmus-go/pkg/log"
	experimentTypes "github.com/litmuschaos/litmus-go/pkg/vmware/vmware-http-chaos/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//PrepareInputParameters will set the required parameters for the http chaos experiment
func PrepareInputParameters(experimentDetails *experimentTypes.ExperimentDetails) (map[string]string, error) {

	parameters := map[string]string{}
	// Setting up toxic name
	toxicName := experimentDetails.StreamType + "_" + experimentDetails.HttpChaosType

	// Initialising the input parameters
	parameters["InstallDependency"] = experimentDetails.InstallDependency
	parameters["ToxicName"] = toxicName
	parameters["ListenPort"] = experimentDetails.ListenPort
	parameters["StreamType"] = experimentDetails.StreamType
	parameters["StreamPort"] = experimentDetails.StreamPort
	parameters["ToxicType"] = experimentDetails.HttpChaosType
	parameters["ToxicValue"] = ""

	// Adding experiment args to parameter list
	switch experimentDetails.HttpChaosType {
	case "latency":

		log.InfoWithValues("[Info]: Details of Http Chaos:", logrus.Fields{
			"Chaos Type":  experimentDetails.HttpChaosType,
			"Latency":     experimentDetails.Latency,
			"Listen Port": experimentDetails.ListenPort,
			"Stream Type": experimentDetails.StreamType,
			"Stream Port": experimentDetails.StreamPort,
		})

		parameters["ToxicValue"] = strconv.Itoa(experimentDetails.Latency)

	case "timeout":

		log.InfoWithValues("[Info]: Details of Http Chaos:", logrus.Fields{
			"Chaos Type":  experimentDetails.HttpChaosType,
			"Timeout":     experimentDetails.RequestTimeout,
			"Listen Port": experimentDetails.ListenPort,
			"Stream Type": experimentDetails.StreamType,
			"Stream Port": experimentDetails.StreamPort,
		})

		parameters["ToxicValue"] = strconv.Itoa(experimentDetails.RequestTimeout)

	case "rate-limit":

		log.InfoWithValues("[Info]: Details of Http Chaos:", logrus.Fields{
			"Chaos Type":  experimentDetails.HttpChaosType,
			"Rate Limit":  experimentDetails.RateLimit,
			"Listen Port": experimentDetails.ListenPort,
			"Stream Type": experimentDetails.StreamType,
			"Stream Port": experimentDetails.StreamPort,
		})

		parameters["ToxicValue"] = strconv.Itoa(experimentDetails.RateLimit)

	case "data-limit":

		log.InfoWithValues("[Info]: Details of Http Chaos:", logrus.Fields{
			"Chaos Type":  experimentDetails.HttpChaosType,
			"Data Limit":  experimentDetails.DataLimit,
			"Listen Port": experimentDetails.ListenPort,
			"Stream Type": experimentDetails.StreamType,
			"Stream Port": experimentDetails.StreamPort,
		})

		parameters["ToxicValue"] = strconv.Itoa(experimentDetails.DataLimit)

	default:
		return nil, errors.Errorf("Http chaos for type: %v is not supported", experimentDetails.HttpChaosType)
	}

	// Adding " to start and end of strings
	parameters["ToxicValue"] = "\"" + parameters["ToxicValue"] + "\""

	return parameters, nil
}

func CheckRunCommandResultError(result *compute.RunCommandResult) error {
	message := strings.Split(strings.TrimSuffix(*(*result.Value)[0].Message, "\n"), "\n")
	i := 0

	for ; i < len(message) && message[i] != "[stderr]"; i++ {
	}
	// errorCodes := make([][]int)
	var errorCode []int
	errorCode = nil

	if message[i+1] != "" {
		exitCodeRegex := regexp.MustCompile("error:")
		for ; i < len(message); i++ {
			// errorCodes = append(errorCodes, exitCodeRegex.FindStringIndex())
			errorCode = exitCodeRegex.FindStringIndex(message[i])
			break
		}
	}
	if errorCode != nil {
		return errors.Errorf("Script failed due to %v", message[errorCode[0]:])
	}
	return nil
}
