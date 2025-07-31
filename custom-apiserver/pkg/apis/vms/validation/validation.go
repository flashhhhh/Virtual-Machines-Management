/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package validation

import (
	"github.com/flashhhhh/Virtual-Machines-Management/custom-apiserver/pkg/apis/vms"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateVirtualMachine validates a Virtual Machine.
func ValidateVirtualMachine(f *vms.VirtualMachine) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, ValidateVirtualMachineSpec(&f.Spec, field.NewPath("spec"))...)

	return allErrs
}

// ValidateVirtualMachineSpec validates a VirtualMachineSpec.
func ValidateVirtualMachineSpec(s *vms.VirtualMachineSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if s.Image == "" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("image"), s.Image, "virtual machine's image is empty"))
	}

	if s.SubnetID == "" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("subnetID"), s.SubnetID, "virtual machine's subnet ID is empty"))
	}

	return allErrs
}
