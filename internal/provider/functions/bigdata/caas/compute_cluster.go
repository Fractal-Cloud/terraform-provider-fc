package caas

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"fractal.cloud/terraform-provider-fc/internal/provider/functions/spec"
)

func NewBigdataCaasComputeClusterFunction() function.Function {
	return spec.New(spec.Spec{
		Name:          "bigdata_caas_compute_cluster",
		ComponentType: "BigData.CaaS.ComputeCluster",
		Summary:       "Creates a containerized Spark Compute Cluster blueprint component",
		Description:   "Builds a BigData.CaaS.ComputeCluster component on a BigData.CaaS.DistributedDataProcessing platform. spark_version is required.",
		Attributes: []spec.Attribute{
			{Name: "image", Key: "image", Kind: spec.String},
			{Name: "spark_version", Key: "sparkVersion", Kind: spec.String, Required: true},
			{Name: "driver_cores", Key: "driverCores", Kind: spec.String},
			{Name: "driver_memory", Key: "driverMemory", Kind: spec.String},
			{Name: "executor_cores", Key: "executorCores", Kind: spec.String},
			{Name: "executor_memory", Key: "executorMemory", Kind: spec.String},
			{Name: "executor_instances", Key: "executorInstances", Kind: spec.Int64},
		},
		Dependencies: []spec.Dependency{{Name: "platform", ComponentType: "BigData.CaaS.DistributedDataProcessing"}},
	})
}
