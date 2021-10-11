package main

import (
	"flag"
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"

	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"

	azureHttpChaos "github.com/chaosnative/litmus-go/experiments/azure/azure-http-chaos/experiment"
	azureStressChaos "github.com/chaosnative/litmus-go/experiments/azure/azure-stress-chaos/experiment"

	"github.com/chaosnative/litmus-go/pkg/clients"
	"github.com/chaosnative/litmus-go/pkg/log"
	"github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		DisableSorting:         true,
		DisableLevelTruncation: true,
	})
}

func main() {

	clients := clients.ClientSets{}

	// parse the experiment name
	experimentName := flag.String("name", "azure-http-chaos", "name of the chaos experiment")

	//Getting kubeConfig and Generate ClientSets
	if err := clients.GenerateClientSetFromKubeConfig(); err != nil {
		log.Errorf("Unable to Get the kubeconfig, err: %v", err)
		return
	}

	log.Infof("Experiment Name: %v", *experimentName)

	// invoke the corresponding experiment based on the the (-name) flag
	switch *experimentName {
	case "azure-stress-chaos":
		azureStressChaos.AzureStressChaosExperiment(clients)
	case "azure-http-chaos":
		azureHttpChaos.AzureHttpChaosExperiment(clients)
	default:
		log.Errorf("Unsupported -name %v, please provide the correct value of -name args", *experimentName)
		return
	}
}
