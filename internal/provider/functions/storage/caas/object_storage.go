package caas

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"fractal.cloud/terraform-provider-fc/internal/provider/functions/spec"
)

func NewStorageCaasObjectStorageFunction() function.Function {
	return spec.New(spec.Spec{
		Name:          "storage_caas_object_storage",
		ComponentType: "Storage.CaaS.ObjectStorage",
		Summary:       "Creates a containerized Object Storage blueprint component",
		Description: "Builds a Storage.CaaS.ObjectStorage component (MinIO) on a container platform. Consumers link to it with " +
			"settings = { accessMode = \"read\" | \"write\" | \"readWrite\", scope = { path = ..., recursive = ... } }.",
		Dependencies: []spec.Dependency{{Name: "container_platform", ComponentType: "NetworkAndCompute.PaaS.ContainerPlatform"}},
	})
}
