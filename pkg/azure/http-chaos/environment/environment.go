package environment

import (
	"strconv"

	clientTypes "k8s.io/apimachinery/pkg/types"

	experimentTypes "github.com/chaosnative/litmus-go/pkg/azure/http-chaos/types"
	"github.com/litmuschaos/litmus-go/pkg/types"
	"github.com/litmuschaos/litmus-go/pkg/utils/common"
)

// STEPS TO GETENV OF YOUR CHOICE HERE
// ADDED FOR FEW MANDATORY FIELD

//GetENV fetches all the env variables from the runner pod
func GetENV(experimentDetails *experimentTypes.ExperimentDetails) {
	experimentDetails.ExperimentName = common.Getenv("EXPERIMENT_NAME", "azure-http-chaos")
	experimentDetails.ChaosNamespace = common.Getenv("CHAOS_NAMESPACE", "litmus")
	experimentDetails.EngineName = common.Getenv("CHAOSENGINE", "")
	experimentDetails.ChaosDuration, _ = strconv.Atoi(common.Getenv("TOTAL_CHAOS_DURATION", "30"))
	experimentDetails.ChaosInterval, _ = strconv.Atoi(common.Getenv("CHAOS_INTERVAL", "30"))
	experimentDetails.RampTime, _ = strconv.Atoi(common.Getenv("RAMP_TIME", "0"))
	experimentDetails.ChaosLib = common.Getenv("LIB", "litmus")
	experimentDetails.AppNS = common.Getenv("APP_NAMESPACE", "")
	experimentDetails.AppLabel = common.Getenv("APP_LABEL", "")
	experimentDetails.AppKind = common.Getenv("APP_KIND", "")
	experimentDetails.ChaosUID = clientTypes.UID(common.Getenv("CHAOS_UID", ""))
	experimentDetails.InstanceID = common.Getenv("INSTANCE_ID", "")
	experimentDetails.Delay, _ = strconv.Atoi(common.Getenv("STATUS_CHECK_DELAY", "2"))
	experimentDetails.Timeout, _ = strconv.Atoi(common.Getenv("STATUS_CHECK_TIMEOUT", "15"))
	experimentDetails.AzureInstanceNames = common.Getenv("AZURE_INSTANCE_NAMES", "akash-run-command,akash-chaos-test")
	experimentDetails.ResourceGroup = common.Getenv("RESOURCE_GROUP", "akash-litmus-test")
	experimentDetails.ScaleSet = common.Getenv("SCALE_SET", "disable")
	experimentDetails.Sequence = common.Getenv("SEQUENCE", "serial")
	experimentDetails.StreamPort = common.Getenv("STREAM_PORT", "6379")
	experimentDetails.StreamType = common.Getenv("STREAM_TYPE", "upstream")
	experimentDetails.ListenPort = common.Getenv("LISTEN_PORT", "20000")
	experimentDetails.HttpChaosType = common.Getenv("HTTP_CHAOS_TYPE", "latency")
	experimentDetails.InstallDependency = common.Getenv("INSTALL_DEPENDENCY", "True")
	experimentDetails.OperatingSystem = common.Getenv("OPERATING_SYSTEM", "linux")
	experimentDetails.Latency, _ = strconv.Atoi(common.Getenv("LATENCY", "2000"))
	experimentDetails.RateLimit, _ = strconv.Atoi(common.Getenv("RATE_LIMIT", "100"))
	experimentDetails.DataLimit, _ = strconv.Atoi(common.Getenv("DATA_LIMIT", "10000"))
	experimentDetails.RequestTimeout, _ = strconv.Atoi(common.Getenv("REQUEST_TIMEOUT", "1000"))
	experimentDetails.ScriptPath = common.Getenv("SCRIPT_PATH", "pkg/azure/http-chaos/scripts/run-script.sh")
	experimentDetails.AbortScriptPath = common.Getenv("ABORT_SCRIPT_PATH", "pkg/azure/http-chaos/scripts/stop-script.sh")
}

//InitialiseChaosVariables initialise all the global variables
func InitialiseChaosVariables(chaosDetails *types.ChaosDetails, experimentDetails *experimentTypes.ExperimentDetails) {
	appDetails := types.AppDetails{}
	appDetails.AnnotationCheck, _ = strconv.ParseBool(common.Getenv("ANNOTATION_CHECK", "false"))
	appDetails.AnnotationKey = common.Getenv("ANNOTATION_KEY", "litmuschaos.io/chaos")
	appDetails.AnnotationValue = "true"
	appDetails.Kind = experimentDetails.AppKind
	appDetails.Label = experimentDetails.AppLabel
	appDetails.Namespace = experimentDetails.AppNS

	chaosDetails.ChaosNamespace = experimentDetails.ChaosNamespace
	chaosDetails.ChaosPodName = experimentDetails.ChaosPodName
	chaosDetails.ChaosUID = experimentDetails.ChaosUID
	chaosDetails.EngineName = experimentDetails.EngineName
	chaosDetails.ExperimentName = experimentDetails.ExperimentName
	chaosDetails.InstanceID = experimentDetails.InstanceID
	chaosDetails.Timeout = experimentDetails.Timeout
	chaosDetails.Delay = experimentDetails.Delay
	chaosDetails.AppDetail = appDetails
	chaosDetails.JobCleanupPolicy = common.Getenv("JOB_CLEANUP_POLICY", "retain")
	chaosDetails.ProbeImagePullPolicy = experimentDetails.LIBImagePullPolicy
}
