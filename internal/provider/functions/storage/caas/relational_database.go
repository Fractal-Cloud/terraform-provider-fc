package caas

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"fractal.cloud/terraform-provider-fc/internal/provider/functions/spec"
)

func NewStorageCaasRelationalDatabaseFunction() function.Function {
	return spec.New(spec.Spec{
		Name:          "storage_caas_relational_database",
		ComponentType: "Storage.CaaS.RelationalDatabase",
		Summary:       "Creates a containerized Relational Database blueprint component",
		Description: "Builds a Storage.CaaS.RelationalDatabase component: a database on a Storage.CaaS.RelationalDbms. " +
			"Workloads link to it with settings = { access = \"read-write\" | \"read-only\" }.",
		Dependencies: []spec.Dependency{{Name: "dbms", ComponentType: "Storage.CaaS.RelationalDbms"}},
	})
}
