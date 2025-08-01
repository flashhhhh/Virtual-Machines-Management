package registry

import (
	"fmt"

	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
)

// REST implements a RESTStorage for API services against etcd
type REST struct {
	*genericregistry.Store
}

// RESTInPeace is just a simple function that panics on error.
// Otherwise returns the given storage object. It is meant to be
// a wrapper for vms registries.
func RESTInPeace(storage *REST, err error) *REST {
	if err != nil {
		err = fmt.Errorf("unable to create REST storage for a resource due to %v, will die", err)
		panic(err)
	}
	return storage
}


type StatusREST struct {
	*genericregistry.Store
}

// StatusRESTInPeace is just a simple function that panics on error.
// Otherwise returns the given status storage object. It is meant to be
// a wrapper for vms registries.
func StatusRESTInPeace(storage *StatusREST, err error) *StatusREST {
	if err != nil {
		err = fmt.Errorf("unable to create REST storage for a resource status due to %v, will die", err)
		panic(err)
	}
	return storage
}