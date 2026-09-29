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
			"age enables the Apache AGE graph extension.",
		Attributes: []spec.Attribute{
			// The catalog declares version; the CloudNativePG agent reads
			// postgresqlVersion.
			{Name: "engine_version", Key: "version", Aliases: []string{"postgresqlVersion"}, Kind: spec.String},
			{Name: "age", Key: "age", Kind: spec.Bool},
		},
		Dependencies: []spec.Dependency{{Name: "container_platform", ComponentType: "NetworkAndCompute.PaaS.ContainerPlatform"}},
	})
}
