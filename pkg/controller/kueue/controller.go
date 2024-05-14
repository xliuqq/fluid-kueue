package kueue

import (
	"context"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	kueue "sigs.k8s.io/kueue/apis/kueue/v1beta1"
	"sigs.k8s.io/kueue/pkg/util/admissioncheck"
	"sigs.k8s.io/kueue/pkg/workload"
)

type Controller struct {
	client client.Client
	helper *provisioningConfigHelper
	record record.EventRecorder
}

// used for check *kueue.AdmissionCheckParametersReference
type provisioningConfigHelper = admissioncheck.ConfigHelper[*kueue.ProvisioningRequestConfig, kueue.ProvisioningRequestConfig]

func newProvisioningConfigHelper(c client.Client) (*provisioningConfigHelper, error) {
	return admissioncheck.NewConfigHelper[*kueue.ProvisioningRequestConfig](c)
}

func NewController(client client.Client, record record.EventRecorder) (*Controller, error) {
	helper, err := newProvisioningConfigHelper(client)
	if err != nil {
		return nil, err
	}
	return &Controller{
		client: client,
		record: record,
		helper: helper,
	}, nil
}

// Reconcile performs a full reconciliation for the object referred to by the Request.
// The Controller will requeue the Request to be processed again if an error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (c *Controller) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	wl := &kueue.Workload{}
	log := ctrl.LoggerFrom(ctx)
	log.V(2).Info("Reconcile workload")

	err := c.client.Get(ctx, req.NamespacedName, wl)
	if err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// if admitted, no need to reconcile.
	if workload.IsAdmitted(wl) {
		return reconcile.Result{}, nil
	}

	// get the lists of relevant checks
	relevantChecks, err := admissioncheck.FilterForController(ctx, c.client, wl.Status.AdmissionChecks, ControllerName)
	if err != nil {
		return reconcile.Result{}, err
	}

	// get the waiting data operation

	// check the data operation status

	// if the data operation is not ready, requeue

	// if the data operation is ready, update the corresponding admission check State

	// add Contition `Admitted` Status `True`

	return reconcile.Result{}, nil
}
