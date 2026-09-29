package caas

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
	ft "fractal.cloud/terraform-provider-fc/internal/provider/functions/functiontest"
)

// Each function builds its catalog type, with the minimum the catalog
// requires, and depends on the component its dependency attribute names.
func TestBigdataCaasFunctions_BuildCatalogTypes(t *testing.T) {
	tests := []struct {
		new           func() function.Function
		name          string
		componentType string
		required      map[string]attr.Value
		dependency    string
		dependencyOn  string
	}{
		{NewBigdataCaasDistributedDataProcessingFunction, "bigdata_caas_distributed_data_processing", "BigData.CaaS.DistributedDataProcessing", nil, "container_platform", "NetworkAndCompute.PaaS.ContainerPlatform"},
		{NewBigdataCaasComputeClusterFunction, "bigdata_caas_compute_cluster", "BigData.CaaS.ComputeCluster", map[string]attr.Value{"spark_version": types.StringValue("3.5.1")}, "platform", "BigData.CaaS.DistributedDataProcessing"},
		{NewBigdataCaasDataProcessingJobFunction, "bigdata_caas_data_processing_job", "BigData.CaaS.DataProcessingJob", map[string]attr.Value{"job_type": types.StringValue("Python")}, "platform", "BigData.CaaS.DistributedDataProcessing"},
		{NewBigdataCaasMlExperimentFunction, "bigdata_caas_ml_experiment", "BigData.CaaS.MlExperiment", nil, "container_platform", "NetworkAndCompute.PaaS.ContainerPlatform"},
		{NewBigdataCaasDataCatalogFunction, "bigdata_caas_data_catalog", "BigData.CaaS.DataCatalog", nil, "container_platform", "NetworkAndCompute.PaaS.ContainerPlatform"},
		{NewBigdataCaasDatalakeFunction, "bigdata_caas_datalake", "BigData.CaaS.Datalake", nil, "container_platform", "NetworkAndCompute.PaaS.ContainerPlatform"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.new()
			var meta function.MetadataResponse
			f.Metadata(context.Background(), function.MetadataRequest{}, &meta)
			if meta.Name != tt.name {
				t.Errorf("name = %q, want %q", meta.Name, tt.name)
			}

			attrs := map[string]attr.Value{"id": types.StringValue("c")}
			for k, v := range tt.required {
				attrs[k] = v
			}
			if tt.dependency != "" {
				dep, err := components.BuildComponent("dep", tt.dependencyOn, types.StringNull(), types.StringNull(), types.StringNull(), nil, nil, nil)
				if err != nil {
					t.Fatalf("building dependency: %s", err.Text)
				}
				attrs[tt.dependency] = dep
			}

			c := ft.Component(t, ft.Run(t, f, ft.Object(t, attrs)))
			if got := c["type"].(types.String).ValueString(); got != tt.componentType {
				t.Errorf("type = %q, want %q", got, tt.componentType)
			}
			if tt.dependency != "" {
				if deps := ft.Strings(t, c, "dependencies_ids"); len(deps) != 1 || deps[0] != "dep" {
					t.Errorf("dependencies = %v, want [dep]", deps)
				}
			}
		})
	}
}
