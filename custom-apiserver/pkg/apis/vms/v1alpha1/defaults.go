package v1alpha1

import "k8s.io/apimachinery/pkg/runtime"

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

// SetDefaults_VirtualMachineSpec sets defaults for VirtualMachine spec
func SetDefaults_VirtualMachineSpec(obj *VirtualMachineSpec) {
	// if (obj.ReferenceType == nil || len(*obj.ReferenceType) == 0) && len(obj.Reference) != 0 {
	// 	t := VirtualMachineReferenceType
	// 	obj.ReferenceType = &t
	// }

	if obj.Image == "" {
		obj.Image = "ami-020cba7c55df1f615"
	}

	if (obj.Size == "") {
		obj.Size = "t2.micro"
	}
}
