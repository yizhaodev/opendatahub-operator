package aigateway

import (
	"context"
	"fmt"

	operatorv1 "github.com/openshift/api/operator/v1"

	componentApi "github.com/opendatahub-io/opendatahub-operator/v2/api/components/v1alpha1"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/cluster"
	odhtypes "github.com/opendatahub-io/opendatahub-operator/v2/pkg/controller/types"
	odhdeploy "github.com/opendatahub-io/opendatahub-operator/v2/pkg/deploy"
)

func initialize(ctx context.Context, rr *odhtypes.ReconciliationRequest) error {
	ai, ok := rr.Instance.(*componentApi.AIGateway)
	if !ok {
		return fmt.Errorf("resource instance %v is not a componentApi.AIGateway", rr.Instance)
	}

	if ai.Spec.BatchGateway.ManagementState == operatorv1.Managed {
		rr.Manifests = append(rr.Manifests, batchGatewayManifestPath(rr.ManifestsBasePath))
	}

	appNamespace, err := cluster.ApplicationNamespace(ctx, rr.Client)
	if err != nil {
		return err
	}

	err = odhdeploy.ApplyParams(
		batchGatewayManifestPath(rr.ManifestsBasePath).String(),
		"params.env",
		nil,
		map[string]string{"namespace": appNamespace},
	)
	if err != nil {
		return fmt.Errorf("failed to update params.env: %w", err)
	}

	return nil
}
