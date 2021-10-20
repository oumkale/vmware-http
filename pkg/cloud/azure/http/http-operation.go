package azure

import (
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	experimentTypes "github.com/chaosnative/litmus-go/pkg/azure/http-chaos/types"
	"github.com/litmuschaos/litmus-go/pkg/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//PrepareInputParameters will set the required parameters for the http chaos experiment
func PrepareInputParameters(experimentDetails *experimentTypes.ExperimentDetails) ([]compute.RunCommandInputParameter, error) {

	parameters := []compute.RunCommandInputParameter{}

	// Setting up toxic name
	toxicName := experimentDetails.StreamType + "_chaos"
	// Initialising the input parameters
	parameterName := []string{"InstallDependency", "ToxicName", "ListenPort", "StreamType", "StreamPort", "ToxicType", "ToxicValue"}
	parameterValues := []string{experimentDetails.InstallDependency, toxicName, experimentDetails.ListenPort, experimentDetails.StreamType, experimentDetails.StreamPort, "( ", "( "}

	//  get the toxic name or list of toxic names
	toxicsTypeList := strings.Split(experimentDetails.HttpChaosType, ",")
	if len(toxicsTypeList) == 0 {
		return nil, errors.Errorf("no toxics found")
	}

	toxicTypes := ""
	toxicValues := ""

	// Adding experiment args to parameter list
	for _, toxic := range toxicsTypeList {
		toxicTypes = toxicTypes + toxic + ","
		switch toxic {
		case "latency":

			log.InfoWithValues("[Info]: Details of Http Chaos:", logrus.Fields{
				"Chaos Type":  toxic,
				"Latency":     experimentDetails.Latency,
				"Listen Port": experimentDetails.ListenPort,
				"Stream Type": experimentDetails.StreamType,
				"Stream Port": experimentDetails.StreamPort,
			})
			toxicValues = toxicValues + strconv.Itoa(experimentDetails.Latency) + ","

		case "timeout":

			log.InfoWithValues("[Info]: Details of Http Chaos:", logrus.Fields{
				"Chaos Type":  toxic,
				"Timeout":     experimentDetails.RequestTimeout,
				"Listen Port": experimentDetails.ListenPort,
				"Stream Type": experimentDetails.StreamType,
				"Stream Port": experimentDetails.StreamPort,
			})

			toxicValues = toxicValues + strconv.Itoa(experimentDetails.RequestTimeout) + ","

		case "rate-limit":

			log.InfoWithValues("[Info]: Details of Http Chaos:", logrus.Fields{
				"Chaos Type":  toxic,
				"Rate Limit":  experimentDetails.RateLimit,
				"Listen Port": experimentDetails.ListenPort,
				"Stream Type": experimentDetails.StreamType,
				"Stream Port": experimentDetails.StreamPort,
			})
			toxicValues = toxicValues + strconv.Itoa(experimentDetails.RateLimit) + ","

		case "data-limit":

			log.InfoWithValues("[Info]: Details of Http Chaos:", logrus.Fields{
				"Chaos Type":  toxic,
				"Data Limit":  experimentDetails.DataLimit,
				"Listen Port": experimentDetails.ListenPort,
				"Stream Type": experimentDetails.StreamType,
				"Stream Port": experimentDetails.StreamPort,
			})
			toxicValues = toxicValues + strconv.Itoa(experimentDetails.DataLimit) + ","

		default:
			return nil, errors.Errorf("Http chaos for type: %v is not supported", toxic)
		}
	}

	toxicValues = strings.TrimSuffix(toxicValues, ",")
	toxicTypes = strings.TrimSuffix(toxicTypes, ",")

	parameterValues[5] = toxicTypes
	parameterValues[6] = toxicValues

	// appending values to parameters
	for i := range parameterValues {
		parameters = append(parameters, compute.RunCommandInputParameter{
			Name:  &parameterName[i],
			Value: &parameterValues[i],
		})
	}

	return parameters, nil
}

//PrepareAbortInputParameters will set the required parameters for the abort of http chaos experiment
func PrepareAbortInputParameters(experimentDetails *experimentTypes.ExperimentDetails) ([]compute.RunCommandInputParameter, error) {

	parameters := []compute.RunCommandInputParameter{}

	// Setting up toxic name
	toxicName := experimentDetails.StreamType + "_chaos"

	// Initialising the input parameters
	parameterName := []string{"ToxicName"}
	parameterValues := []string{toxicName}

	// appending values to parameters
	for i := range parameterValues {
		parameters = append(parameters, compute.RunCommandInputParameter{
			Name:  &parameterName[i],
			Value: &parameterValues[i],
		})
	}

	return parameters, nil
}

func CheckRunCommandResultError(result *compute.RunCommandResult) error {
	message := strings.Split(strings.TrimSuffix(*(*result.Value)[0].Message, "\n"), "\n")
	i := 0

	for ; i < len(message) && message[i] != "[stderr]"; i++ {
	}

	if i < len(message)-1 && message[i+1] != "" {
		return errors.Errorf("Script failed due to %v. Check logs!", message[i+1])
	}
	return nil
}
