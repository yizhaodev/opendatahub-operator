package aigateway

import (
	componentApi "github.com/opendatahub-io/opendatahub-operator/v2/api/components/v1alpha1"
	"github.com/opendatahub-io/opendatahub-operator/v2/internal/controller/status"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/controller/types"
)

const (
	ComponentName = componentApi.AIGatewayComponentName

	ReadyConditionType = componentApi.AIGatewayKind + status.ReadySuffix

	// LegacyComponentName is the name of the component that is assigned to deployments
	// via Kustomize. Since a deployment selector is immutable, we can't upgrade existing
	// deployment to the new component name, so keep it around till we figure out a solution.
	LegacyComponentName = "llm-d-batch-gateway"
)

var (
	imageParamMap = map[string]string{
		"LLM_D_BATCH_GATEWAY_OPERATOR_IMAGE":  "RELATED_IMAGE_ODH_BATCH_GATEWAY_OPERATOR_IMAGE",
		"LLM_D_BATCH_GATEWAY_APISERVER_IMAGE": "RELATED_IMAGE_ODH_LLM_D_BATCH_GATEWAY_APISERVER_IMAGE",
		"LLM_D_BATCH_GATEWAY_PROCESSOR_IMAGE": "RELATED_IMAGE_ODH_LLM_D_BATCH_GATEWAY_PROCESSOR_IMAGE",
		"LLM_D_BATCH_GATEWAY_GC_IMAGE":        "RELATED_IMAGE_ODH_LLM_D_BATCH_GATEWAY_GC_IMAGE",
	}

	conditionTypes = []string{
		status.ConditionDeploymentsAvailable,
	}
)

func batchGatewayManifestPath(basePath string) types.ManifestInfo {
	return types.ManifestInfo{
		Path:       basePath,
		ContextDir: "batchgateway",
		SourcePath: "overlays/odh",
	}
}
