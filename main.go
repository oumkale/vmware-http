package main

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/event"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func handleEvent(ref types.ManagedObjectReference, events []types.BaseEvent) (err error) {
	// for _, event := range events {
	// 	eventType := reflect.TypeOf(event).String()
	// 	fmt.Printf("Event found of type %s\n", eventType)
	// }

	return nil
}

var (
	ServiceInstance = types.ManagedObjectReference{
		Type:  "VslmServiceInstance",
		Value: "ServiceInstance",
	}
)

func main() {

	// Creating a connection context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Parsing URL
	url, err := url.Parse("https://administrator@vsphere.local:Chaosnative@1658@103.195.244.216/sdk")
	//url, err := url.Parse("https://administrator@vsphere.local:@Chaosnative@1658@103.195.244.216/sdk")

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

	// vCenter version
	info := client.ServiceContent.About
	fmt.Println("Connected to vCenter version \n", client.ServiceContent.RootFolder)
	fmt.Print("Connected to vCenter version \n", info, client.ServiceContent.HostProfileManager)

	// Selecting default datacenter
	finder := find.NewFinder(client.Client, true)
	dc, err := finder.DefaultDatacenter(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "1Error: %s\n", err)
		os.Exit(1)
	}
	refs := []types.ManagedObjectReference{dc.Reference()}

	m := view.NewManager(client.Client)

	v, err := m.CreateContainerView(ctx, client.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "3Error: %s\n", err)
		os.Exit(1)
	}

	var vms []mo.VirtualMachine
	err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{"summary"}, &vms)
	if err != nil {
		fmt.Fprintf(os.Stderr, "4Error: %s\n", err)
		os.Exit(1)
	}

	// Print summary per vm (see also: govc/vm/info.go)
	for _, vm := range vms {
		//fmt.Printf("%s: %s\n", vm.Summary.Config.Name, vm.Summary.Config.GuestFullName)
		//if vm.Summary.Config.Name == "adarsh-vm" {
		fmt.Println("VM ref :", vm.Summary.Config.Name)

		//}
	}

	// Setting up the event manager
	eventManager := event.NewManager(client.Client)
	err = eventManager.Events(ctx, refs, 10, false, false, handleEvent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "2Error: %s\n", err)
		os.Exit(1)
	}

	startService("apache2", "adarsh", "Datacenter", "adarsh-vm", "123")
	//fmt.Println("Connected to vCenter version \n", res)

}

func Shellout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func getService(serviceName, vmName, datacenter, vmUserName, vmPassWord string) {
	command := fmt.Sprintf("govc guest.run -vm=%s -dc=%s -l=%s:%s systemctl list-unit-files --type service | awk '{print $1}' | grep %s.service", vmName, datacenter, vmUserName, vmPassWord, serviceName)
	stdout, stderr, err := Shellout(command)

	if stderr != "" {
		fmt.Println("Stderr: ", stderr)
		return
	}

	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	stdout = stdout[:len(stdout)-1]
	fmt.Println(stdout + " is running")
}

func getServiceState(serviceName, vmName, datacenter, vmUserName, vmPassWord string) {
	command := fmt.Sprintf("govc guest.run -vm=%s -dc=%s -l=%s:%s systemctl show %s -p ActiveState --no-page | sed 's/ActiveState=//g'", vmName, datacenter, vmUserName, vmPassWord, serviceName)
	stdout, stderr, err := Shellout(command)

	if stderr != "" {
		fmt.Println("Stderr: ", stderr)
		return
	}

	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	stdout = stdout[:len(stdout)-1]
	fmt.Println(stdout)
}

func stopService(serviceName, vmName, datacenter, vmUserName, vmPassWord string) {
	command := fmt.Sprintf(`govc guest.run -vm=%s -dc=%s -l=%s:%s printf "%s" | sudo -S systemctl stop %s`, vmName, datacenter, vmUserName, vmPassWord, vmPassWord, serviceName)
	stdout, _, err := Shellout(command)

	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	fmt.Printf("%s", stdout)
}

func startService(serviceName, vmName, datacenter, vmUserName, vmPassWord string) {
	fmt.Println("startService ", vmUserName)
	//command := fmt.Sprintf("govc guest.run -l=%s:%s -vm %s printf chaos | sudo -S sh ~/go/src/vmware/sh1.sh", vmName, vmPassWord, vmUserName)
	command := fmt.Sprintf(`govc guest.run -l=%s:%s -vm %s "printf %s | sudo -S bash ~/go/src/vmware/sh.sh"`, vmName, vmPassWord, vmUserName, vmPassWord)
	// echo
	//command := fmt.Sprintf("govc guest.run -l=%s:%s -vm %s printf %s | sudo -S sh ~/home/adarsh/Desktop/vmware-http/sh1.sh", vmName, vmPassWord, vmUserName, vmPassWord)

	//command := fmt.Sprintf("sudo govc guest.run -vm=%s -l=%s:%s printf '%s' | sudo -S systemctl start %s", vmName, datacenter, vmUserName, vmPassWord, vmPassWord, serviceName)
	stdout, stderr, err := Shellout(command)
	fmt.Println("startService ", stderr)
	if err != nil {
		fmt.Println("err: 1", err)
		return
	}

	fmt.Printf("Output : %s", stdout)
}

//func main() {
// getService("apache2", "neel-vm-2", "Datacenter", "neel-vm", "123")
// getServiceState("apache2", "neel-vm-2", "Datacenter", "neel-vm", "123")
// stopService("apache2", "neel-vm-2", "Datacenter", "neel-vm", "123")
//	startService("apache2", "neel-vm-2", "Datacenter", "neel-vm", "123")
//}

// export GOVC_URL=http://106.51.78.18:3567/
// export GOVC_USERAME=Ajesh
// export GOVC_PASSWORD=mayadata@1658
// export GOVC_INSECURE=true
// source ~/.profile

// govc guest.upload -l=adarsh:123 -vm=adarsh-vm /home/oumkale/go/src/vmware/sh.sh  /home/adarsh/Desktop/vmware-http/sh.sh

// govc guest.run -l adarsh:123 -vm adarsh-vm-1  sh ~/a.sh
// Server@12
