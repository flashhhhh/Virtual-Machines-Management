package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:prerelease-lifecycle-gen:introduced=1.0
// +k8s:prerelease-lifecycle-gen:removed=1.10

type VirtualMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []VirtualMachine `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type VirtualMachineSpec struct {
	Image            string   `json:"image" protobuf:"bytes,1,opt,name=image"`
	Size             string   `json:"size" protobuf:"bytes,2,opt,name=size"`
	SSHKeyIDs        []string `json:"sshKeyIDs" protobuf:"bytes,3,rep,name=sshKeyIDs"`
	SubnetID         string   `json:"subnetID" protobuf:"bytes,4,opt,name=subnetID"`
	SecurityGroupIDs []string `json:"securityGroupIDs" protobuf:"bytes,5,rep,name=securityGroupIDs"`
}

type VirtualMachineStatus struct {
	Phase       string `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase"`
	ID          string `json:"id,omitempty" protobuf:"bytes,2,opt,name=id"`
	IP          string `json:"ip,omitempty" protobuf:"bytes,3,opt,name=ip"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:prerelease-lifecycle-gen:introduced=1.0
// +k8s:prerelease-lifecycle-gen:removed=1.10

type VirtualMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   VirtualMachineSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status VirtualMachineStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}