package main

import (
	"fmt"

	vminterfaces "github.com/flashhhhh/Virtual-Machines-Management/custom-controller/vm_interfaces"
	"github.com/flashhhhh/Virtual-Machines-Management/custom-controller/vm_interfaces/aws"
)

func main() {
	var provider vminterfaces.VMInterfaces
	var err error

	provider, err = aws.NewAWSProvider()
	if err != nil {
		panic(fmt.Sprintf("failed to create AWS provider: %v", err))
	}

	vmCfg := vminterfaces.VirtualMachineConfig{
		Name:       "my-virtual-machine",
		Size:       "t2.micro",
		Image:      "ami-0c94855ba95c71c99", // For AWS
		SSHKeyIDs:  []string{"kp-linux"},
		SecurityGroupIDs: []string{"sg-03577d4124c5d53e7"},
		SubnetID:   "subnet-080dc7385f0044c14",
	}

	vm, err := provider.CreateVM(vmCfg)
	if err != nil {
		println("Error creating VM:", err.Error())
	}

	fmt.Printf("Created VM: %+v\n", vm)

	fetchedVM, err := provider.GetVM(vm.ID)
	if err != nil {
		println("Error fetching VM:", err.Error())
	} else {
		fmt.Printf("Fetched VM: %+v\n", fetchedVM)
	}

	err = provider.DeleteVM(vm.ID)
	if err != nil {
		println("Error deleting VM:", err.Error())
	} else {
		fmt.Println("VM deleted successfully")
	}
}
