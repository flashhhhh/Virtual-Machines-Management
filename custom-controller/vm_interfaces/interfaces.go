package vminterfaces

type VirtualMachineConfig struct {
	Image 	   string
	Size        string
	SSHKeyIDs   []string
	SubnetID    string
	SecurityGroupIDs []string
	Name        string
}

type VirtualMachine struct {
	ID	string
	IP	string
}

type VMInterfaces interface {
	CreateVM(config VirtualMachineConfig) (*VirtualMachine, error)
	GetVM(id string) (*VirtualMachine, error)
	DeleteVM(id string) error
	GetVMStatus(instanceID string) (string, error)
}