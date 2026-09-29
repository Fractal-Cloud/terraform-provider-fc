package caas

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"fractal.cloud/terraform-provider-fc/internal/provider/functions/spec"
)

func NewBigdataCaasDataCatalogFunction() function.Function {
	return spec.New(spec.Spec{
		Name:          "bigdata_caas_data_catalog",
		ComponentType: "BigData.CaaS.DataCatalog",
		Summary:       "Creates a containerized Data Catalog blueprint component",
		Description:   "Builds a BigData.CaaS.DataCatalog component (Unity Catalog) on a container platform.",
		Attributes: []spec.Attribute{
			{Name: "unity_catalog_version", Key: "unityCatalogVersion", Kind: spec.String},
			{Name: "replicas", Key: "replicas", Kind: spec.Int64},
		},
		Dependencies: []spec.Dependency{{Name: "container_platform", ComponentType: "NetworkAndCompute.PaaS.ContainerPlatform"}},
	})
}
