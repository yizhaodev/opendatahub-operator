package batchgateway

import (
	"context"
	"errors"
	"fmt"

	operatorv1 "github.com/openshift/api/operator/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/opendatahub-io/opendatahub-operator/v2/api/common"
	componentApi "github.com/opendatahub-io/opendatahub-operator/v2/api/components/v1alpha1"
	dscv2 "github.com/opendatahub-io/opendatahub-operator/v2/api/datasciencecluster/v2"
	"github.com/opendatahub-io/opendatahub-operator/v2/internal/controller/components"
	"github.com/opendatahub-io/opendatahub-operator/v2/internal/controller/status"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/controller/conditions"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/controller/types"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/operatorconfig"
)

const (
	ComponentName      = componentApi.BatchGatewayComponentName
	ReadyConditionType = componentApi.BatchGatewayKind + status.ReadySuffix
)

type componentHandler struct{}

func NewHandler() *componentHandler { return &componentHandler{} }

func (s *componentHandler) GetName() string {
	return ComponentName
}

func (s *componentHandler) Init(_ common.Platform, _ operatorconfig.OperatorSettings) error {
	return nil
}

// NewCRObject returns nil — the user creates the LLMBatchGateway CR manually
// with the full spec (apiServer, processor, GC images, inference gateway config, etc.).
// The DSC only controls enablement via ManagementState.
func (s *componentHandler) NewCRObject(_ context.Context, _ client.Client, _ *dscv2.DataScienceCluster) (common.PlatformObject, error) {
	return nil, nil
}

func (s *componentHandler) IsEnabled(dsc *dscv2.DataScienceCluster) bool {
	if dsc.Spec.Components.Kserve.ManagementState != operatorv1.Managed {
		return false
	}
	return dsc.Spec.Components.Kserve.BatchGateway.ManagementState == operatorv1.Managed
}

func (s *componentHandler) UpdateDSCStatus(ctx context.Context, rr *types.ReconciliationRequest) (metav1.ConditionStatus, error) {
	cs := metav1.ConditionUnknown

	dsc, ok := rr.Instance.(*dscv2.DataScienceCluster)
	if !ok {
		return cs, errors.New("failed to convert to DataScienceCluster")
	}

	rr.Conditions.MarkFalse(ReadyConditionType)

	if !s.IsEnabled(dsc) {
		ms := dsc.Spec.Components.Kserve.BatchGateway.ManagementState
		if ms == "" {
			ms = operatorv1.Removed
		}
		rr.Conditions.MarkFalse(
			ReadyConditionType,
			conditions.WithReason(string(ms)),
			conditions.WithMessage("Component ManagementState is set to %s", string(ms)),
			conditions.WithSeverity(common.ConditionSeverityInfo),
		)
		return cs, nil
	}

	c := componentApi.LLMBatchGateway{}
	c.Name = componentApi.BatchGatewayInstanceName

	if err := rr.Client.Get(ctx, client.ObjectKeyFromObject(&c), &c); err != nil {
		if k8serr.IsNotFound(err) {
			rr.Conditions.MarkFalse(
				ReadyConditionType,
				conditions.WithReason(status.NotReadyReason),
				conditions.WithMessage("LLMBatchGateway CR not available yet"),
			)
			return metav1.ConditionFalse, nil
		}
		return cs, fmt.Errorf("failed to get LLMBatchGateway: %w", err)
	}

	if !c.GetDeletionTimestamp().IsZero() {
		rr.Conditions.MarkFalse(
			ReadyConditionType,
			conditions.WithReason(status.DeletingReason),
			conditions.WithMessage(status.DeletingMessage),
		)
		return metav1.ConditionFalse, nil
	}

	_ = components.NormalizeManagementState(dsc.Spec.Components.Kserve.BatchGateway.ManagementState)

	if rc := conditions.FindStatusCondition(c.GetStatus(), status.ConditionTypeReady); rc != nil {
		rr.Conditions.MarkFrom(ReadyConditionType, *rc)
		cs = rc.Status
	} else {
		rr.Conditions.MarkFalse(
			ReadyConditionType,
			conditions.WithReason(status.NotReadyReason),
			conditions.WithMessage("LLMBatchGateway CR exists but has no ready condition yet"),
		)
		cs = metav1.ConditionFalse
	}

	return cs, nil
}

func (s *componentHandler) NewComponentReconciler(ctx context.Context, mgr ctrl.Manager) error {
	return newReconciler(ctx, mgr)
}
