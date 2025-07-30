package initializer

import (
	"k8s.io/apiserver/pkg/admission"
	informers "custom-apiserver/pkg/generated/informers/externalversions"
)

// WantsInternalVMSInformerFactory defines a function which sets InformerFactory for admission plugins that need it
type WantsInternalVMSInformerFactory interface {
	SetInternalVMSInformerFactory(informers.SharedInformerFactory)
	admission.InitializationValidator
}
