package caas

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"fractal.cloud/terraform-provider-fc/internal/provider/functions/spec"
)

func NewStorageCaasRelationalDbmsFunction() function.Function {
	return spec.New(spec.Spec{
		Name:          "storage_caas_relational_dbms",
		ComponentType: "Storage.CaaS.RelationalDbms",
		Summary:       "Creates a containerized Relational DBMS blueprint component",
		Description: "Builds a Storage.CaaS.RelationalDbms component (PostgreSQL via CloudNativePG) on a container platform. " +
			"engine_version is the version the cluster is created with; a running cluster changes version only through " +
			"extra_parameters.postgresqlVersion, within its major. age enables the Apache AGE graph extension.",
		Attributes: []spec.Attribute{
			{Name: "engine_version", Key: "version", Kind: spec.String},
			{Name: "age", Key: "age", Kind: spec.Bool},
		},
		Dependencies: []spec.Dependency{{Name: "container_platform", ComponentType: "NetworkAndCompute.PaaS.ContainerPlatform"}},
	})
}
