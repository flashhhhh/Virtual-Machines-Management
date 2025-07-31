package virtualmachine

import (
	"context"
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apiserver/pkg/admission"

	"github.com/flashhhhh/Virtual-Machines-Management/custom-apiserver/pkg/admission/initializer"
	"github.com/flashhhhh/Virtual-Machines-Management/custom-apiserver/pkg/apis/vms"
	informers "github.com/flashhhhh/Virtual-Machines-Management/custom-apiserver/pkg/generated/informers/externalversions"
	listers "github.com/flashhhhh/Virtual-Machines-Management/custom-apiserver/pkg/generated/listers/vms/v1alpha1"
)

// Register registers a plugin
func Register(plugins *admission.Plugins) {
	plugins.Register("VirtualMachine", func(config io.Reader) (admission.Interface, error) {
		return New()
	})
}

// The Plugin structure
type Plugin struct {
	*admission.Handler
	virtualMachineLister listers.VirtualMachineLister
}

var _ = initializer.WantsInternalVMSInformerFactory(&Plugin{})

// Admit ensures that the object in-flight is of kind VirtualMachine.
func (d *Plugin) Admit(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) error {
	// we are only interested in VirtualMachines
	if a.GetKind().GroupKind() != vms.Kind("VirtualMachine") {
		return nil
	}

	if !d.WaitForReady() {
		return admission.NewForbidden(a, fmt.Errorf("not yet ready to handle request"))
	}

	metaAccessor, err := meta.Accessor(a.GetObject())
	if err != nil {
		return err
	}
	virtualMachineName := metaAccessor.GetName()

	if len(virtualMachineName) > 10 {
		return errors.NewForbidden(
			a.GetResource().GroupResource(),
			a.GetName(),
			fmt.Errorf("the length of virtual machine's name mustn't be greater than 10"),
		)
	}

	return nil
}

func (d *Plugin) SetInternalVMSInformerFactory(f informers.SharedInformerFactory) {
	d.virtualMachineLister = f.Vms().V1alpha1().VirtualMachines().Lister()
	d.SetReadyFunc(f.Vms().V1alpha1().VirtualMachines().Informer().HasSynced)
}

// ValidateInitialization checks whether the plugin was correctly initialized.
func (d *Plugin) ValidateInitialization() error {
	if d.virtualMachineLister == nil {
		return fmt.Errorf("missing virtual machine lister")
	}
	return nil
}

// New creates a new VirtualMachine admission plugin.
func New() (*Plugin, error) {
	p := &Plugin{
		Handler: admission.NewHandler(admission.Create, admission.Update),
	}

	return p, nil
}
