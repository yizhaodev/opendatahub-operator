package batchgateway

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	engineTypes "github.com/k8s-manifest-kit/engine/pkg/types"
	helmPkg "github.com/k8s-manifest-kit/renderer-helm/pkg"

	componentApi "github.com/opendatahub-io/opendatahub-operator/v2/api/components/v1alpha1"
	odhtypes "github.com/opendatahub-io/opendatahub-operator/v2/pkg/controller/types"
)

const (
	chartName         = "batch-gateway"
	defaultChartsPath = "opt/charts"
)

func initializeHelmCharts(_ context.Context, rr *odhtypes.ReconciliationRequest) error {
	instance, ok := rr.Instance.(*componentApi.LLMBatchGateway)
	if !ok {
		return fmt.Errorf("resource instance is not LLMBatchGateway (got %T)", rr.Instance)
	}

	rr.HelmCharts = []odhtypes.HelmChartInfo{
		{
			Source: helmPkg.Source{
				Chart:       filepath.Join(chartsBasePath(), chartName),
				ReleaseName: "batch-gateway",
				Values: func(_ context.Context) (engineTypes.Values, error) {
					return specToHelmValues(instance)
				},
			},
		},
	}

	return nil
}

func chartsBasePath() string {
	if p := os.Getenv("DEFAULT_CHARTS_PATH"); p != "" {
		return p
	}
	return defaultChartsPath
}
