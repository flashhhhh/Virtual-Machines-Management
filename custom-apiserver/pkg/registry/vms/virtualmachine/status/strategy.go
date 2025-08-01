/*
Copyright 2017 The Kubernetes Authors.

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

package status

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"

	"github.com/flashhhhh/Virtual-Machines-Management/custom-apiserver/pkg/apis/vms"
	"github.com/flashhhhh/Virtual-Machines-Management/custom-apiserver/pkg/apis/vms/validation"
)

// NewStatusStrategy creates and returns a virtualMachineStatusStrategy instance
func NewStatusStrategy(typer runtime.ObjectTyper) virtualMachineStatusStrategy {
	return virtualMachineStatusStrategy{typer, names.SimpleNameGenerator}
}

// GetAttrs returns labels.Set, fields.Set, and error in case the given runtime.Object is not a VirtualMachine
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	apiserver, ok := obj.(*vms.VirtualMachine)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a Virtual Machine")
	}
	return labels.Set(apiserver.ObjectMeta.Labels), SelectableFields(apiserver), nil
}

// SelectableFields returns a field set that represents the object.
func SelectableFields(obj *vms.VirtualMachine) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}

// MatchVirtualMachine is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiserver only interested in specific labels/fields.
func MatchVirtualMachine(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

type virtualMachineStatusStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (virtualMachineStatusStrategy) NamespaceScoped() bool {
	return true
}

func (virtualMachineStatusStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
}

func (virtualMachineStatusStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	// Only allow updating status — preserve spec
	newVM := obj.(*vms.VirtualMachine)
	oldVM := old.(*vms.VirtualMachine)

	newVM.Spec = oldVM.Spec
}

func (virtualMachineStatusStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	newVM, ok := obj.(*vms.VirtualMachine)
	if !ok {
		return field.ErrorList{field.Invalid(field.NewPath("obj"), obj, "not a VirtualMachine")}
	}
	return validation.ValidateVirtualMachine(newVM)
}

func (virtualMachineStatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (virtualMachineStatusStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

func (virtualMachineStatusStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (virtualMachineStatusStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (s virtualMachineStatusStrategy) Canonicalize(obj runtime.Object) {}

func (s virtualMachineStatusStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}
