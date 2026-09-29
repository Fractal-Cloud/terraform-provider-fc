package caas

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"fractal.cloud/terraform-provider-fc/internal/provider/functions/spec"
)

func NewBigdataCaasDistributedDataProcessingFunction() function.Function {
	return spec.New(spec.Spec{
		Name:          "bigdata_caas_distributed_data_processing",
		ComponentType: "BigData.CaaS.DistributedDataProcessing",
		Summary:       "Creates a containerized Distributed Data Processing blueprint component",
		Description: "Builds a BigData.CaaS.DistributedDataProcessing component (Spark on Kubernetes) on a container platform. " +
			"Operator settings of the chosen offer (operatorVersion, sparkVersion, namespace) go in extra_parameters. " +
			"Link it to a BigData.CaaS.Datalake to mount the lake.",
		Dependencies: []spec.Dependency{{Name: "container_platform", ComponentType: "NetworkAndCompute.PaaS.ContainerPlatform"}},
		Links:        true,
	})
}
