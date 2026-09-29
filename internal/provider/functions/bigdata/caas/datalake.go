package caas

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"fractal.cloud/terraform-provider-fc/internal/provider/functions/spec"
)

func NewBigdataCaasDatalakeFunction() function.Function {
	return spec.New(spec.Spec{
		Name:          "bigdata_caas_datalake",
		ComponentType: "BigData.CaaS.Datalake",
		Summary:       "Creates a containerized Datalake blueprint component",
		Description: "Builds a BigData.CaaS.Datalake component (MinIO) on a container platform. Tenant sizing of the chosen " +
			"offer (servers, volumesPerServer, volumeSize, storageClass) goes in extra_parameters. Jobs link to it with " +
			"settings = { purpose = \"raw\" | \"curated\" | \"checkpoint\" }.",
		Dependencies: []spec.Dependency{{Name: "container_platform", ComponentType: "NetworkAndCompute.PaaS.ContainerPlatform"}},
	})
}
