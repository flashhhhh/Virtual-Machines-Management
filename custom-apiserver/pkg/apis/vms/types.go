package vms

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type VirtualMachineList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items           []VirtualMachine
}

type VirtualMachineSpec struct {
	Image		string
	Size		string
	SSHKeyIDs	[]string
	SubnetID	string
	SecurityGroupIDs []string
}

type VirtualMachineStatus struct {
	Phase	string
	ID		string
	IP		string
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type VirtualMachine struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   VirtualMachineSpec
	Status VirtualMachineStatus
}