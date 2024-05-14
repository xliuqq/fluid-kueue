package kueue

import (
	"context"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	kueue "sigs.k8s.io/kueue/apis/kueue/v1beta1"
)

type acReconciler struct {
	client client.Client
	helper *provisioningConfigHelper
}

var _ reconcile.Reconciler = (*acReconciler)(nil)

// Reconcile is used to set the AdmissionCheck Status.
func (a *acReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	ac := &kueue.AdmissionCheck{}
	// only reconcile current Admission Check
	if err := a.client.Get(ctx, req.NamespacedName, ac); err != nil || ac.Spec.ControllerName != ControllerName {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	currentCondition := ptr.Deref(apimeta.FindStatusCondition(ac.Status.Conditions, kueue.AdmissionCheckActive), metav1.Condition{})
	newCondition := metav1.Condition{
		Type:               kueue.AdmissionCheckActive,
		Status:             metav1.ConditionTrue,
		Reason:             "Active",
		Message:            "The admission check is active",
		ObservedGeneration: ac.Generation,
	}

	// check the parameter
	if _, err := a.helper.ConfigFromRef(ctx, ac.Spec.Parameters); err != nil {
		newCondition.Status = metav1.ConditionFalse
		newCondition.Reason = "BadParametersRef"
		newCondition.Message = err.Error()
	}

	// update status
	if currentCondition.Status != newCondition.Status {
		apimeta.SetStatusCondition(&ac.Status.Conditions, newCondition)
		return reconcile.Result{}, a.client.Status().Update(ctx, ac)
	}
	return reconcile.Result{}, nil
}
