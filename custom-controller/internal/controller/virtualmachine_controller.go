/*
Copyright 2025.

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

package controller

import (
	"context"
	"time"

	vmsv1alpha1 "github.com/flashhhhh/Virtual-Machines-Management/custom-apiserver/pkg/apis/vms/v1alpha1"
	vminterfaces "github.com/flashhhhh/Virtual-Machines-Management/custom-controller/vm_interfaces"
	"github.com/flashhhhh/Virtual-Machines-Management/custom-controller/vm_interfaces/aws"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// VirtualMachineReconciler reconciles a VirtualMachine object
type VirtualMachineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var AWSProvider *aws.AWSProvider

var (
	VMPending      = "pending"
	VMRunning      = "running"
	VMShuttingDown = "shutting-down"
	VMTerminated   = "terminated"
	VMStopped      = "stopped"
)

func init() {
	provider, err := aws.NewAWSProvider()
	if err != nil {
		panic("failed to initialize AWSProvider: " + err.Error())
	}
	AWSProvider = provider
}

// +kubebuilder:rbac:groups=vms.example.com,resources=virtualmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vms.example.com,resources=virtualmachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=vms.example.com,resources=virtualmachines/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the VirtualMachine object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *VirtualMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	virtualMachine := &vmsv1alpha1.VirtualMachine{}
	if err := r.Get(ctx, req.NamespacedName, virtualMachine); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.Info("VirtualMachine resource not found, ignoring since object must be deleted", "name", req.NamespacedName)
			return ctrl.Result{}, nil
		}

		logger.Error(err, "unable to fetch VirtualMachine")
		return ctrl.Result{}, err
	}

	// If the VirtualMachine is marked for deletion, we should handle it
	if !virtualMachine.ObjectMeta.DeletionTimestamp.IsZero() {
		logger.Info("Finalizing VirtualMachine", "id", virtualMachine.Status.ID)

		err := AWSProvider.DeleteVM(virtualMachine.Status.ID)
		if err != nil {
			logger.Error(err, "unable to delete VirtualMachine in AWS", "ID", virtualMachine.Status.ID)
			return ctrl.Result{}, err
		}

		// Remove finalizer if present
		finalizerName := "virtualmachine.finalizers.vms.example.com"
		for i, f := range virtualMachine.ObjectMeta.Finalizers {
			if f == finalizerName {
				virtualMachine.ObjectMeta.Finalizers = append(virtualMachine.ObjectMeta.Finalizers[:i], virtualMachine.ObjectMeta.Finalizers[i+1:]...)
				break
			}
		}
		if err := r.Update(ctx, virtualMachine); err != nil {
			logger.Error(err, "unable to remove finalizer from VirtualMachine")
			return ctrl.Result{}, err
		}

		logger.Info("VirtualMachine deleted in AWS", "name", virtualMachine.Name, "ID", virtualMachine.Status.ID)
		return ctrl.Result{}, nil
	}

	// Create VirtualMachine
	if virtualMachine.Status.Phase == "" {
		// Add finalizer if not present
		finalizerName := "virtualmachine.finalizers.vms.example.com"
		hasFinalizer := false
		for _, f := range virtualMachine.ObjectMeta.Finalizers {
			if f == finalizerName {
				hasFinalizer = true
				break
			}
		}
		if !hasFinalizer {
			virtualMachine.ObjectMeta.Finalizers = append(virtualMachine.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(ctx, virtualMachine); err != nil {
				logger.Error(err, "unable to add finalizer to VirtualMachine")
				return ctrl.Result{}, err
			}
		}

		// This is a new VirtualMachine, we can create the VirtualMachine in AWS
		virtualMachineConfig := vminterfaces.VirtualMachineConfig{
			Name:             virtualMachine.Name,
			Image:            virtualMachine.Spec.Image,
			Size:             virtualMachine.Spec.Size,
			SSHKeyIDs:        virtualMachine.Spec.SSHKeyIDs,
			SubnetID:         virtualMachine.Spec.SubnetID,
			SecurityGroupIDs: virtualMachine.Spec.SecurityGroupIDs,
		}
		vm, err := AWSProvider.CreateVM(virtualMachineConfig)
		if err != nil {
			logger.Error(err, "unable to create VirtualMachine in AWS")
			return ctrl.Result{}, err
		}

		// We can set its initial phase
		virtualMachine.Status.Phase = VMPending
		virtualMachine.Status.ID = vm.ID

		if err := r.Status().Update(ctx, virtualMachine); err != nil {
			logger.Error(err, "unable to update VirtualMachine status")
			return ctrl.Result{}, err
		}
		logger.Info("VirtualMachine created", "name", virtualMachine.Name, "ID", virtualMachine.Status.ID)
	}

	// Update the VirtualMachine status based on its current state in AWS
	vmStatus, err := AWSProvider.GetVMStatus(virtualMachine.Status.ID)
	if err != nil {
		logger.Error(err, "unable to get VirtualMachine status from AWS", "ID", virtualMachine.Status.ID)
		return ctrl.Result{}, err
	}

	if vmStatus == virtualMachine.Status.Phase {
		logger.Info("VirtualMachine status is unchanged", "ID", virtualMachine.Status.ID, "status", vmStatus)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Update the VirtualMachine status
	virtualMachine.Status.Phase = vmStatus

	if (vmStatus == VMRunning) {
		// If the VM is running, we can fetch its public IP
		vm, err := AWSProvider.GetVM(virtualMachine.Status.ID)
		if err != nil {
			logger.Error(err, "unable to get VirtualMachine details from AWS", "ID", virtualMachine.Status.ID)
			return ctrl.Result{}, err
		}

		virtualMachine.Status.IP = vm.IP
	}

	if err := r.Status().Update(ctx, virtualMachine); err != nil {
		logger.Error(err, "unable to update VirtualMachine status")
		return ctrl.Result{}, err
	}

	logger.Info("VirtualMachine status updated", "ID", virtualMachine.Status.ID, "status", virtualMachine.Status.Phase)
	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *VirtualMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&vmsv1alpha1.VirtualMachine{}).
		Named("virtualmachine").
		Complete(r)
}
