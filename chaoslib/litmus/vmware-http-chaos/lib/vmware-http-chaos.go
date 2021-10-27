package lib

import (
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	vmware "github.com/chaosnative/litmus-go/pkg/cloud/vmware/vmware-http-chaos/http"
	commonExperimentTypes "github.com/chaosnative/litmus-go/pkg/vmware/vmware-http-chaos/types"
	experimentTypes "github.com/chaosnative/litmus-go/pkg/vmware/vmware-http-chaos/types"
	clients "github.com/litmuschaos/litmus-go/pkg/clients"
	"github.com/litmuschaos/litmus-go/pkg/events"
	"github.com/litmuschaos/litmus-go/pkg/log"
	"github.com/litmuschaos/litmus-go/pkg/probe"
	"github.com/litmuschaos/litmus-go/pkg/types"
	"github.com/litmuschaos/litmus-go/pkg/utils/common"
	"github.com/pkg/errors"
	"github.com/vmware/govmomi/guest/toolbox"
)

var (
	err           error
	inject, abort chan os.Signal
)

func InjectVMHttpChaosChaos(experimentsDetails *experimentTypes.ExperimentDetails, clients clients.ClientSets, resultDetails *types.ResultDetails, eventsDetails *types.EventDetails, chaosDetails *types.ChaosDetails) error {

	// inject channel is used to transmit signal notifications
	inject = make(chan os.Signal, 1)
	// Catch and relay certain signal(s) to inject channel
	signal.Notify(inject, os.Interrupt, syscall.SIGTERM)

	// abort channel is used to transmit signal notifications.
	abort = make(chan os.Signal, 1)
	signal.Notify(abort, os.Interrupt, syscall.SIGTERM)

	// Waiting for the ramp time before chaos injection
	if experimentsDetails.RampTime != 0 {
		log.Infof("[Ramp]: Waiting for the %vs ramp time before injecting chaos", experimentsDetails.RampTime)
		common.WaitForDuration(experimentsDetails.RampTime)
	}

	//  get the instance name or list of instance names
	instanceNameList := strings.Split(experimentsDetails.VMWareInstanceNames, ",")
	if len(instanceNameList) == 0 {
		return errors.Errorf("no instance name found")
	}

	// prepare run command parameters
	parameters, err := vmware.PrepareInputParameters(experimentsDetails)
	if err != nil {
		return errors.Errorf("failed to prepare chaos parameters, err: %v", err)
	}

	toolClient, err := vmware.GetVMWareToolClient(experimentsDetails)
	if err != nil {
		return errors.Errorf("failed to prepare tool client, err: %v", err)
	}

	runCommandInput, err := vmware.PrepareRunCommandInput(experimentsDetails, &parameters)
	if err != nil {
		return errors.Errorf("failed to prepare run command input, err: %v", err)
	}

	err = vmware.UploadFiles(toolClient, experimentsDetails, &parameters)
	if err != nil {
		return errors.Errorf("failed to upload file to target vm, err: %v", err)
	}

	abortRunCommandInput, err := vmware.PrepareAbortRunCommandInput(experimentsDetails, &parameters)
	if err != nil {
		return errors.Errorf("failed to prepare abort command input, err: %v", err)
	}

	// watching for the abort signal and revert the chaos
	go abortWatcher(experimentsDetails, instanceNameList, abortRunCommandInput, toolClient)

	switch strings.ToLower(experimentsDetails.Sequence) {
	case "serial":
		if err = injectChaosInSerialMode(experimentsDetails, instanceNameList, runCommandInput, clients, resultDetails, eventsDetails, chaosDetails, toolClient); err != nil {
			return err
		}
	case "parallel":
		if err = injectChaosInParallelMode(experimentsDetails, instanceNameList, runCommandInput, clients, resultDetails, eventsDetails, chaosDetails, toolClient); err != nil {
			return err
		}
	default:
		return errors.Errorf("%v sequence is not supported", experimentsDetails.Sequence)
	}

	// Waiting for the ramp time after chaos injection
	if experimentsDetails.RampTime != 0 {
		log.Infof("[Ramp]: Waiting for the %vs ramp time after injecting chaos", experimentsDetails.RampTime)
		common.WaitForDuration(experimentsDetails.RampTime)
	}
	return nil
}

func injectChaosInParallelMode(experimentsDetails *experimentTypes.ExperimentDetails, instanceNameList []string, runCommandInput *exec.Cmd, clients clients.ClientSets, resultDetails *types.ResultDetails, eventsDetails *types.EventDetails, chaosDetails *types.ChaosDetails, tools *toolbox.Client) error {
	select {
	case <-inject:
		// Stopping the chaos execution, if abort signal received
		os.Exit(0)
	default:
		// ChaosStartTimeStamp contains the start timestamp, when the chaos injection begin
		ChaosStartTimeStamp := time.Now()
		duration := int(time.Since(ChaosStartTimeStamp).Seconds())

		for duration < experimentsDetails.ChaosDuration {

			log.Infof("[Info]: Target instanceName list, %v", instanceNameList)

			if experimentsDetails.EngineName != "" {
				msg := "Injecting " + experimentsDetails.ExperimentName + " chaos on VMWare instance"
				types.SetEngineEventAttributes(eventsDetails, types.ChaosInject, msg, "Normal", chaosDetails)
				events.GenerateEvents(eventsDetails, clients, chaosDetails, "ChaosEngine")
			}

			// Running scripts parallely
			for _, vmName := range instanceNameList {
				log.Infof("[Chaos]: Running script on the VMWare instance: %v", vmName)
				if err := vmware.Run(runCommandInput, tools); err != nil {
					return errors.Errorf("unable to run script on vmware instance, err: %v", err)
				}
			}

			// Run probes during chaos
			if len(resultDetails.ProbeDetails) != 0 {
				if err = probe.RunProbes(chaosDetails, clients, resultDetails, "DuringChaos", eventsDetails); err != nil {
					return err
				}
			}

			// Wait for Chaos interval
			log.Infof("[Wait]: Waiting for chaos interval of %vs", experimentsDetails.ChaosInterval)
			common.WaitForDuration(experimentsDetails.ChaosInterval)

			duration = int(time.Since(ChaosStartTimeStamp).Seconds())
		}
	}
	return nil
}

func injectChaosInSerialMode(experimentsDetails *experimentTypes.ExperimentDetails, instanceNameList []string, runCommandInput *exec.Cmd, clients clients.ClientSets, resultDetails *types.ResultDetails, eventsDetails *types.EventDetails, chaosDetails *types.ChaosDetails, tools *toolbox.Client) error {
	select {
	case <-inject:
		// Stopping the chaos execution, if abort signal received
		os.Exit(0)
	default:
		// ChaosStartTimeStamp contains the start timestamp, when the chaos injection begin
		ChaosStartTimeStamp := time.Now()
		duration := int(time.Since(ChaosStartTimeStamp).Seconds())

		for duration < experimentsDetails.ChaosDuration {

			log.Infof("[Info]: Target instanceName list, %v", instanceNameList)

			if experimentsDetails.EngineName != "" {
				msg := "Injecting " + experimentsDetails.ExperimentName + " chaos on VMWare instance"
				types.SetEngineEventAttributes(eventsDetails, types.ChaosInject, msg, "Normal", chaosDetails)
				events.GenerateEvents(eventsDetails, clients, chaosDetails, "ChaosEngine")
			}

			// Running scripts serially
			for _, vmName := range instanceNameList {

				log.Infof("[Chaos]: Running script on the VMWare instance: %v", vmName)
				if err := vmware.Run(runCommandInput, tools); err != nil {
					return errors.Errorf("unable to run script on vmware instance, err: %v", err)
				}

				// Run probes during chaos
				if len(resultDetails.ProbeDetails) != 0 {
					if err = probe.RunProbes(chaosDetails, clients, resultDetails, "DuringChaos", eventsDetails); err != nil {
						return err
					}
				}

				// Wait for Chaos interval
				log.Infof("[Wait]: Waiting for chaos interval of %vs", experimentsDetails.ChaosInterval)
				common.WaitForDuration(experimentsDetails.ChaosInterval)
			}
			duration = int(time.Since(ChaosStartTimeStamp).Seconds())
		}
	}
	return nil
}

func abortWatcher(experimentsDetails *experimentTypes.ExperimentDetails, instanceNameList []string, abortRunCommandInput *exec.Cmd, tools *toolbox.Client) {
	<-abort

	log.Info("[Abort]: Chaos Revert Started")
	runCommandFutures := []commonExperimentTypes.RunCommandFuture{}

	for _, vmName := range instanceNameList {
		log.Infof("[Chaos]: Running abort script on the VMWare instance: %v", vmName)
		runCommandFuture := commonExperimentTypes.RunCommandFuture{}
		if err := vmware.Run(abortRunCommandInput, tools); err != nil {
			log.Errorf("unable to run abort script on vmware instance, err: %v", err)
		}
		runCommandFutures = append(runCommandFutures, runCommandFuture)
	}

	log.Infof("[Abort]: Chaos Revert Completed")
	os.Exit(1)
}
