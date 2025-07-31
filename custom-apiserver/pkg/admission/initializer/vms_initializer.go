package initializer

import (
	informers "github.com/flashhhhh/Virtual-Machines-Management/custom-apiserver/pkg/generated/informers/externalversions"

	"k8s.io/apiserver/pkg/admission"
)

type pluginInitializer struct {
	informers informers.SharedInformerFactory
}

var _ admission.PluginInitializer = pluginInitializer{}

// NewPluginInitializer returns a new plugin initializer for the vms admission plugin.
func New(informers informers.SharedInformerFactory) admission.PluginInitializer {
	return pluginInitializer{
		informers: informers,
	}
}

// Initialize checks the initialization interfaces implemented by a plugin
// and provide the appropriate initialization data
func (i pluginInitializer) Initialize(plugin admission.Interface) {
	if wants, ok := plugin.(WantsInternalVMSInformerFactory); ok {
		wants.SetInternalVMSInformerFactory(i.informers)
	}
}
