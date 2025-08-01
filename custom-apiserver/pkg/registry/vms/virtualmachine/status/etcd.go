package status

import (
	"github.com/flashhhhh/Virtual-Machines-Management/custom-apiserver/pkg/apis/vms"
	"github.com/flashhhhh/Virtual-Machines-Management/custom-apiserver/pkg/registry"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
)

// NewStatusREST returns a RESTStorage object that will work against API services for status updates.
func NewStatusREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*registry.StatusREST, error) {
	strategy := NewStatusStrategy(scheme)

	store := &genericregistry.Store{
		NewFunc:                   func() runtime.Object { return &vms.VirtualMachine{} },
		NewListFunc:               func() runtime.Object { return &vms.VirtualMachineList{} },
		PredicateFunc:             MatchVirtualMachine,
		DefaultQualifiedResource:  vms.Resource("virtualmachines"),
		SingularQualifiedResource: vms.Resource("virtualmachine"),

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,

		TableConvertor: rest.NewDefaultTableConvertor(vms.Resource("virtualmachines")),
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &registry.StatusREST{Store: store}, nil
}