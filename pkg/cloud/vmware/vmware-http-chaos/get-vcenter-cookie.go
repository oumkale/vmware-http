package http

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	experimentTypes "github.com/litmuschaos/litmus-go/pkg/vmware/vmware-http-chaos/types"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/guest/toolbox"
	"github.com/vmware/govmomi/vim25/types"
)

// Message contains attribute for message
type Message struct {
	MsgValue string `json:"value"`
}

//GetVcenterSessionID returns the vcenter sessionid
func GetVcenterSessionID(experimentsDetails *experimentTypes.ExperimentDetails) (string, error) {

	//Leverage Go's HTTP Post function to make request
	req, err := http.NewRequest("POST", "https://"+experimentsDetails.VcenterServer+"/rest/com/vmware/cis/session", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(experimentsDetails.VcenterUser, experimentsDetails.VcenterPass)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	//Handle Error
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var m Message
	json.Unmarshal([]byte(body), &m)

	login := "vmware-api-session-id=" + m.MsgValue + ";Path=/rest;Secure;HttpOnly"
	return login, nil
}

func GetVMWareToolClient(experimentDetails *experimentTypes.ExperimentDetails) (*toolbox.Client, error) {
	// Creating a connection context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Parsing URL
	url, err := url.Parse(`https://` + experimentDetails.VcenterUser + `:` + experimentDetails.VcenterPass + `@` + experimentDetails.VcenterServer + `/sdk`)
	if err != nil {
		return nil, err
	}

	// Connecting to vCenter
	client, err := govmomi.NewClient(ctx, url, true)
	if err != nil {
		return nil, err
	}

	// Selecting default datacenter
	vm, err := find.NewFinder(client.Client).VirtualMachine(ctx, experimentDetails.OperatingSystem)
	if err != nil {
		return nil, err
	}
	toolsClient, err := toolbox.NewClient(ctx, client.Client, vm, &types.NamePasswordAuthentication{
		Username: experimentDetails.VMUserName,
		Password: experimentDetails.VMPassword,
	})

	return toolsClient, nil
}

//GetVMStatus gets the current status of Vcenter VM
func GetVMStatus(experimentsDetails *experimentTypes.ExperimentDetails, cookie string) (string, error) {

	type Message struct {
		MsgValue struct {
			StateValue string `json:"state"`
		} `json:"value"`
	}

	//Leverage Go's HTTP Post function to make request
	req, err := http.NewRequest("GET", "https://"+experimentsDetails.VcenterServer+"/rest/vcenter/vm/"+experimentsDetails.AppVMMoid+"/power/", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", cookie)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	//Handle Error
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var m1 Message
	json.Unmarshal([]byte(body), &m1)
	return string(m1.MsgValue.StateValue), nil
}
