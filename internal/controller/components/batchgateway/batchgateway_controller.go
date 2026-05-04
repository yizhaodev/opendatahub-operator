package batchgateway

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	componentApi "github.com/opendatahub-io/opendatahub-operator/v2/api/components/v1alpha1"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/controller/actions/deploy"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/controller/actions/gc"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/controller/actions/render/helm"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/controller/actions/status/deployments"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/controller/predicates"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/controller/reconciler"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/metadata/labels"
)

var conditionTypes = []string{
	"DeploymentsAvailable",
}

func newReconciler(ctx context.Context, mgr ctrl.Manager) error {
	_, err := reconciler.ReconcilerFor(mgr, &componentApi.LLMBatchGateway{}).
		Owns(&appsv1.Deployment{}, reconciler.WithPredicates(predicates.DefaultDeploymentPredicate)).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&rbacv1.Role{}).
		Owns(&rbacv1.RoleBinding{}).
		WithAction(initializeHelmCharts).
		WithAction(helm.NewAction(
			helm.WithLabel(labels.ODH.Component(ComponentName), labels.True),
			helm.WithLabel(labels.K8SCommon.PartOf, ComponentName),
		)).
		WithAction(deploy.NewAction(
			deploy.WithCache(),
		)).
		WithAction(deployments.NewAction()).
		WithAction(gc.NewAction()).
		WithConditions(conditionTypes...).
		Build(ctx)

	return err
}
