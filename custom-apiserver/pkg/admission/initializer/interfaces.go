package initializer

import (
	informers "github.com/flashhhhh/Virtual-Machines-Management/custom-apiserver/pkg/generated/informers/externalversions"
	"k8s.io/apiserver/pkg/admission"
)

// WantsInternalVMSInformerFactory defines a function which sets InformerFactory for admission plugins that need it
type WantsInternalVMSInformerFactory interface {
	SetInternalVMSInformerFactory(informers.SharedInformerFactory)
	admission.InitializationValidator
}
