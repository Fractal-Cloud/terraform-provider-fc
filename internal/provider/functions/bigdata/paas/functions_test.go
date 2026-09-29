package paas

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
	"fractal.cloud/terraform-provider-fc/internal/provider/functions/functiontest"
)

// testLinkAttrTypes is the shape of one element of a function's links list.
var testLinkAttrTypes = map[string]attr.Type{
	"target":   components.ComponentObjectType,
	"settings": types.MapType{ElemType: types.StringType},
}

func buildTestComponent(t *testing.T, id, componentType string) types.Object {
	t.Helper()
	obj, err := components.BuildComponent(id, componentType, types.StringNull(), types.StringNull(), types.StringNull(), nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to build test component: %s", err.Text)
	}
	return obj
}

func runFunction(t *testing.T, f function.Function, args []attr.Value) *function.RunResponse {
	t.Helper()
	ctx := context.Background()
	req := function.RunRequest{
		Arguments: function.NewArgumentsData(functiontest.Complete(t, f, args)),
	}
	resp := &function.RunResponse{
		Result: function.NewResultData(types.ObjectNull(components.ComponentAttrTypes)),
	}
	f.Run(ctx, req, resp)
	return resp
}

func getResultAttrs(t *testing.T, resp *function.RunResponse) map[string]attr.Value {
	t.Helper()
	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error.Text)
	}
	result := resp.Result.Value()
	obj, ok := result.(types.Object)
	if !ok {
		t.Fatalf("expected types.Object result, got %T", result)
	}
	return obj.Attributes()
}

// --- DistributedDataProcessing ---

func TestDistributedDataProcessingFunction_Metadata(t *testing.T) {
	f := NewBigdataPaasDistributedDataProcessingFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "bigdata_paas_distributed_data_processing" {
		t.Errorf("expected name %q, got %q", "bigdata_paas_distributed_data_processing", resp.Name)
	}
}

func TestDistributedDataProcessingFunction_Definition(t *testing.T) {
	f := NewBigdataPaasDistributedDataProcessingFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

func TestDistributedDataProcessingFunction_Run(t *testing.T) {
	f := NewBigdataPaasDistributedDataProcessingFunction()
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":           types.StringType,
		"display_name": types.StringType,
		"description":  types.StringType,
		"links":        types.DynamicType,
	}, map[string]attr.Value{
		"id":           types.StringValue("databricks-1"),
		"display_name": types.StringValue("My Databricks"),
		"description":  types.StringNull(),
		"links":        types.DynamicNull(),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if attrs["type"].(types.String).ValueString() != "BigData.PaaS.DistributedDataProcessing" {
		t.Errorf("expected type %q", "BigData.PaaS.DistributedDataProcessing")
	}
}

// --- Datalake ---

func TestDatalakeFunction_Metadata(t *testing.T) {
	f := NewBigdataPaasDatalakeFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "bigdata_paas_datalake" {
		t.Errorf("expected name %q, got %q", "bigdata_paas_datalake", resp.Name)
	}
}

func TestDatalakeFunction_Definition(t *testing.T) {
	f := NewBigdataPaasDatalakeFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

func TestDatalakeFunction_Run(t *testing.T) {
	f := NewBigdataPaasDatalakeFunction()
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":           types.StringType,
		"display_name": types.StringType,
		"description":  types.StringType,
	}, map[string]attr.Value{
		"id":           types.StringValue("lake-1"),
		"display_name": types.StringValue("My Lake"),
		"description":  types.StringNull(),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)
	if attrs["type"].(types.String).ValueString() != "BigData.PaaS.Datalake" {
		t.Errorf("expected type %q", "BigData.PaaS.Datalake")
	}
}

// --- ComputeCluster ---

func TestComputeClusterFunction_Metadata(t *testing.T) {
	f := NewBigdataPaasComputeClusterFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "bigdata_paas_compute_cluster" {
		t.Errorf("expected name %q, got %q", "bigdata_paas_compute_cluster", resp.Name)
	}
}

func TestComputeClusterFunction_Definition(t *testing.T) {
	f := NewBigdataPaasComputeClusterFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

func TestComputeClusterFunction_Run(t *testing.T) {
	f := NewBigdataPaasComputeClusterFunction()
	platform := buildTestComponent(t, "databricks-1", "BigData.PaaS.DistributedDataProcessing")
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":                       types.StringType,
		"display_name":             types.StringType,
		"description":              types.StringType,
		"platform":                 components.ComponentObjectType,
		"cluster_name":             types.StringType,
		"spark_version":            types.StringType,
		"num_workers":              types.Int64Type,
		"min_workers":              types.Int64Type,
		"max_workers":              types.Int64Type,
		"auto_termination_minutes": types.Int64Type,
		"spark_conf":               types.MapType{ElemType: types.StringType},
		"pypi_libraries":           types.ListType{ElemType: types.StringType},
		"maven_libraries":          types.ListType{ElemType: types.StringType},
		"links":                    types.DynamicType,
	}, map[string]attr.Value{
		"id":                       types.StringValue("cluster-1"),
		"display_name":             types.StringNull(),
		"description":              types.StringNull(),
		"platform":                 platform,
		"cluster_name":             types.StringValue("my-cluster"),
		"spark_version":            types.StringValue("13.3"),
		"num_workers":              types.Int64Value(4),
		"min_workers":              types.Int64Null(),
		"max_workers":              types.Int64Null(),
		"auto_termination_minutes": types.Int64Value(120),
		"spark_conf":               types.MapNull(types.StringType),
		"pypi_libraries":           types.ListNull(types.StringType),
		"maven_libraries":          types.ListNull(types.StringType),
		"links":                    types.DynamicNull(),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if attrs["type"].(types.String).ValueString() != "BigData.PaaS.ComputeCluster" {
		t.Errorf("expected type %q", "BigData.PaaS.ComputeCluster")
	}

	deps := attrs["dependencies_ids"].(types.List)
	if deps.IsNull() || deps.Elements()[0].(types.String).ValueString() != "databricks-1" {
		t.Errorf("expected dependency %q", "databricks-1")
	}

	params := attrs["parameters"].(types.Map)
	elems := params.Elements()
	if elems["clusterName"].(types.String).ValueString() != "my-cluster" {
		t.Errorf("expected clusterName %q", "my-cluster")
	}
	if elems["sparkVersion"].(types.String).ValueString() != "13.3" {
		t.Errorf("expected sparkVersion %q", "13.3")
	}
	if elems["numWorkers"].(types.String).ValueString() != "4" {
		t.Errorf("expected numWorkers %q", "4")
	}
	if elems["autoTerminationMinutes"].(types.String).ValueString() != "120" {
		t.Errorf("expected autoTerminationMinutes %q", "120")
	}
}

func TestComputeClusterFunction_Run_AllParams(t *testing.T) {
	f := NewBigdataPaasComputeClusterFunction()
	platform := buildTestComponent(t, "databricks-1", "BigData.PaaS.DistributedDataProcessing")

	sparkConf, diags := types.MapValue(types.StringType, map[string]attr.Value{
		"spark.executor.memory": types.StringValue("4g"),
	})
	if diags.HasError() {
		t.Fatalf("failed to build spark_conf: %s", diags.Errors())
	}
	pypiLibs, diags := types.ListValue(types.StringType, []attr.Value{types.StringValue("pandas==1.5.0")})
	if diags.HasError() {
		t.Fatalf("failed to build pypi_libraries: %s", diags.Errors())
	}
	mavenLibs, diags := types.ListValue(types.StringType, []attr.Value{types.StringValue("com.example:lib:1.0")})
	if diags.HasError() {
		t.Fatalf("failed to build maven_libraries: %s", diags.Errors())
	}

	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":                       types.StringType,
		"display_name":             types.StringType,
		"description":              types.StringType,
		"platform":                 components.ComponentObjectType,
		"cluster_name":             types.StringType,
		"spark_version":            types.StringType,
		"num_workers":              types.Int64Type,
		"min_workers":              types.Int64Type,
		"max_workers":              types.Int64Type,
		"auto_termination_minutes": types.Int64Type,
		"spark_conf":               types.MapType{ElemType: types.StringType},
		"pypi_libraries":           types.ListType{ElemType: types.StringType},
		"maven_libraries":          types.ListType{ElemType: types.StringType},
		"links":                    types.DynamicType,
	}, map[string]attr.Value{
		"id":                       types.StringValue("cluster-2"),
		"display_name":             types.StringValue("Full Cluster"),
		"description":              types.StringNull(),
		"platform":                 platform,
		"cluster_name":             types.StringValue("full-cluster"),
		"spark_version":            types.StringValue("14.0"),
		"num_workers":              types.Int64Value(8),
		"min_workers":              types.Int64Value(2),
		"max_workers":              types.Int64Value(16),
		"auto_termination_minutes": types.Int64Value(60),
		"spark_conf":               sparkConf,
		"pypi_libraries":           pypiLibs,
		"maven_libraries":          mavenLibs,
		"links":                    types.DynamicNull(),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	params := attrs["parameters"].(types.Map)
	elems := params.Elements()
	if elems["minWorkers"].(types.String).ValueString() != "2" {
		t.Errorf("expected minWorkers %q", "2")
	}
	if elems["maxWorkers"].(types.String).ValueString() != "16" {
		t.Errorf("expected maxWorkers %q", "16")
	}
	if _, ok := elems["sparkConf"]; !ok {
		t.Error("expected sparkConf parameter")
	}
	if _, ok := elems["pypiLibraries"]; !ok {
		t.Error("expected pypiLibraries parameter")
	}
	if _, ok := elems["mavenLibraries"]; !ok {
		t.Error("expected mavenLibraries parameter")
	}
}

// --- DataProcessingJob ---

func TestDataProcessingJobFunction_Metadata(t *testing.T) {
	f := NewBigdataPaasDataProcessingJobFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "bigdata_paas_data_processing_job" {
		t.Errorf("expected name %q, got %q", "bigdata_paas_data_processing_job", resp.Name)
	}
}

func TestDataProcessingJobFunction_Definition(t *testing.T) {
	f := NewBigdataPaasDataProcessingJobFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

func TestDataProcessingJobFunction_Run(t *testing.T) {
	f := NewBigdataPaasDataProcessingJobFunction()
	platform := buildTestComponent(t, "databricks-1", "BigData.PaaS.DistributedDataProcessing")
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":               types.StringType,
		"display_name":     types.StringType,
		"description":      types.StringType,
		"platform":         components.ComponentObjectType,
		"job_name":         types.StringType,
		"task_type":        types.StringType,
		"notebook_path":    types.StringType,
		"python_file":      types.StringType,
		"main_class_name":  types.StringType,
		"jar_uri":          types.StringType,
		"cron_schedule":    types.StringType,
		"max_retries":      types.Int64Type,
		"existing_cluster": types.BoolType,
		"parameters":       types.ListType{ElemType: types.StringType},
	}, map[string]attr.Value{
		"id":               types.StringValue("job-1"),
		"display_name":     types.StringNull(),
		"description":      types.StringNull(),
		"platform":         platform,
		"job_name":         types.StringValue("etl-job"),
		"task_type":        types.StringValue("notebook"),
		"notebook_path":    types.StringValue("/jobs/etl"),
		"python_file":      types.StringNull(),
		"main_class_name":  types.StringNull(),
		"jar_uri":          types.StringNull(),
		"cron_schedule":    types.StringNull(),
		"max_retries":      types.Int64Value(3),
		"existing_cluster": types.BoolValue(true),
		"parameters":       types.ListNull(types.StringType),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if attrs["type"].(types.String).ValueString() != "BigData.PaaS.DataProcessingJob" {
		t.Errorf("expected type %q", "BigData.PaaS.DataProcessingJob")
	}
	params := attrs["parameters"].(types.Map)
	elems := params.Elements()
	if elems["jobName"].(types.String).ValueString() != "etl-job" {
		t.Errorf("expected jobName %q", "etl-job")
	}
	if elems["existingCluster"].(types.String).ValueString() != "true" {
		t.Errorf("expected existingCluster %q", "true")
	}
}

func TestDataProcessingJobFunction_Run_AllParams(t *testing.T) {
	f := NewBigdataPaasDataProcessingJobFunction()
	platform := buildTestComponent(t, "databricks-1", "BigData.PaaS.DistributedDataProcessing")
	paramsList, diags := types.ListValue(types.StringType, []attr.Value{
		types.StringValue("--input"),
		types.StringValue("s3://bucket/data"),
	})
	if diags.HasError() {
		t.Fatalf("failed to build parameters list: %s", diags.Errors())
	}

	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":               types.StringType,
		"display_name":     types.StringType,
		"description":      types.StringType,
		"platform":         components.ComponentObjectType,
		"job_name":         types.StringType,
		"task_type":        types.StringType,
		"notebook_path":    types.StringType,
		"python_file":      types.StringType,
		"main_class_name":  types.StringType,
		"jar_uri":          types.StringType,
		"cron_schedule":    types.StringType,
		"max_retries":      types.Int64Type,
		"existing_cluster": types.BoolType,
		"parameters":       types.ListType{ElemType: types.StringType},
	}, map[string]attr.Value{
		"id":               types.StringValue("job-2"),
		"display_name":     types.StringValue("Full Job"),
		"description":      types.StringNull(),
		"platform":         platform,
		"job_name":         types.StringValue("spark-job"),
		"task_type":        types.StringValue("spark_jar"),
		"notebook_path":    types.StringNull(),
		"python_file":      types.StringValue("s3://bucket/main.py"),
		"main_class_name":  types.StringValue("com.example.Main"),
		"jar_uri":          types.StringValue("s3://bucket/app.jar"),
		"cron_schedule":    types.StringValue("0 0 * * *"),
		"max_retries":      types.Int64Value(5),
		"existing_cluster": types.BoolValue(false),
		"parameters":       paramsList,
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	params := attrs["parameters"].(types.Map)
	elems := params.Elements()
	if elems["pythonFile"].(types.String).ValueString() != "s3://bucket/main.py" {
		t.Errorf("expected pythonFile %q", "s3://bucket/main.py")
	}
	if elems["mainClassName"].(types.String).ValueString() != "com.example.Main" {
		t.Errorf("expected mainClassName %q", "com.example.Main")
	}
	if elems["jarUri"].(types.String).ValueString() != "s3://bucket/app.jar" {
		t.Errorf("expected jarUri %q", "s3://bucket/app.jar")
	}
	if elems["cronSchedule"].(types.String).ValueString() != "0 0 * * *" {
		t.Errorf("expected cronSchedule %q", "0 0 * * *")
	}
	if _, ok := elems["parameters"]; !ok {
		t.Error("expected parameters parameter")
	}
}

// --- MlExperiment ---

func TestMlExperimentFunction_Metadata(t *testing.T) {
	f := NewBigdataPaasMlExperimentFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "bigdata_paas_ml_experiment" {
		t.Errorf("expected name %q, got %q", "bigdata_paas_ml_experiment", resp.Name)
	}
}

func TestMlExperimentFunction_Definition(t *testing.T) {
	f := NewBigdataPaasMlExperimentFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

func TestMlExperimentFunction_Run(t *testing.T) {
	f := NewBigdataPaasMlExperimentFunction()
	platform := buildTestComponent(t, "databricks-1", "BigData.PaaS.DistributedDataProcessing")
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":              types.StringType,
		"display_name":    types.StringType,
		"description":     types.StringType,
		"platform":        components.ComponentObjectType,
		"experiment_name": types.StringType,
	}, map[string]attr.Value{
		"id":              types.StringValue("exp-1"),
		"display_name":    types.StringNull(),
		"description":     types.StringNull(),
		"platform":        platform,
		"experiment_name": types.StringValue("my-experiment"),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if attrs["type"].(types.String).ValueString() != "BigData.PaaS.MlExperiment" {
		t.Errorf("expected type %q", "BigData.PaaS.MlExperiment")
	}
	deps := attrs["dependencies_ids"].(types.List)
	if deps.IsNull() || deps.Elements()[0].(types.String).ValueString() != "databricks-1" {
		t.Errorf("expected dependency %q", "databricks-1")
	}
	params := attrs["parameters"].(types.Map)
	elems := params.Elements()
	if elems["experimentName"].(types.String).ValueString() != "my-experiment" {
		t.Errorf("expected experimentName %q", "my-experiment")
	}
}

func TestDistributedDataProcessingFunction_Run_WithLinks(t *testing.T) {
	f := NewBigdataPaasDistributedDataProcessingFunction()
	datalake := buildTestComponent(t, "data-lake", "BigData.PaaS.Datalake")

	settingsMap, diags := types.MapValue(types.StringType, map[string]attr.Value{
		"mountName": types.StringValue("datalake"),
		"path":      types.StringValue("/"),
	})
	if diags.HasError() {
		t.Fatalf("failed to build settings: %s", diags.Errors())
	}

	linkObj, diags := types.ObjectValue(testLinkAttrTypes, map[string]attr.Value{
		"target":   datalake,
		"settings": settingsMap,
	})
	if diags.HasError() {
		t.Fatalf("failed to build link: %s", diags.Errors())
	}

	linkList, diags := types.ListValue(types.ObjectType{AttrTypes: testLinkAttrTypes}, []attr.Value{linkObj})
	if diags.HasError() {
		t.Fatalf("failed to build link list: %s", diags.Errors())
	}

	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":           types.StringType,
		"display_name": types.StringType,
		"description":  types.StringType,
		"links":        types.DynamicType,
	}, map[string]attr.Value{
		"id":           types.StringValue("spark-platform"),
		"display_name": types.StringValue("Spark Platform"),
		"description":  types.StringNull(),
		"links":        types.DynamicValue(linkList),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	linksVal := attrs["links"].(types.List)
	if linksVal.IsNull() {
		t.Fatal("expected non-null links")
	}
	linkElems := linksVal.Elements()
	if len(linkElems) != 1 {
		t.Fatalf("expected 1 link, got %d", len(linkElems))
	}
	linkAttrs := linkElems[0].(types.Object).Attributes()
	if linkAttrs["component_id"].(types.String).ValueString() != "data-lake" {
		t.Errorf("expected link component_id %q", "data-lake")
	}
	linkSettings := linkAttrs["settings"].(types.Map)
	if linkSettings.IsNull() {
		t.Fatal("expected non-null settings")
	}
	if linkSettings.Elements()["mountName"].(types.String).ValueString() != "datalake" {
		t.Errorf("expected mountName %q", "datalake")
	}
}
