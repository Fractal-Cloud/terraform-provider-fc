package caas

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"fractal.cloud/terraform-provider-fc/internal/provider/functions/spec"
)

func NewBigdataCaasMlExperimentFunction() function.Function {
	return spec.New(spec.Spec{
		Name:          "bigdata_caas_ml_experiment",
		ComponentType: "BigData.CaaS.MlExperiment",
		Summary:       "Creates a containerized ML Experiment tracking blueprint component",
		Description:   "Builds a BigData.CaaS.MlExperiment component (an MLflow tracking server) on a container platform.",
		Attributes: []spec.Attribute{
			{Name: "mlflow_version", Key: "mlflowVersion", Kind: spec.String},
			{Name: "backend_store_uri", Key: "backendStoreUri", Kind: spec.String},
			{Name: "artifact_root", Key: "artifactRoot", Kind: spec.String},
			{Name: "replicas", Key: "replicas", Kind: spec.Int64},
			{Name: "service_port", Key: "servicePort", Kind: spec.Int64},
		},
		Dependencies: []spec.Dependency{{Name: "container_platform", ComponentType: "NetworkAndCompute.PaaS.ContainerPlatform"}},
	})
}
