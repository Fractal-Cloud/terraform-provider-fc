package caas

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

var workloadAttrTypes = map[string]attr.Type{
	"id":              types.StringType,
	"display_name":    types.StringType,
	"description":     types.StringType,
	"container_image": types.StringType,
	"container_port":  types.Int64Type,
	"container_name":  types.StringType,
	"cpu":             types.StringType,
	"memory":          types.StringType,
	"desired_count":   types.Int64Type,
	"platform":        components.ComponentObjectType,
	"subnet":          components.ComponentObjectType,
	"links":           types.DynamicType,
	"security_groups": types.ListType{ElemType: components.ComponentObjectType},
}

func TestWorkloadFunction_Metadata(t *testing.T) {
	f := NewWorkloadFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "custom_workloads_caas_workload" {
		t.Errorf("expected name %q, got %q", "custom_workloads_caas_workload", resp.Name)
	}
}

func TestWorkloadFunction_Definition(t *testing.T) {
	f := NewWorkloadFunction()
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

func TestWorkloadFunction_Run_Minimal(t *testing.T) {
	f := NewWorkloadFunction()
	configObj, diags := types.ObjectValue(workloadAttrTypes, map[string]attr.Value{
		"id":              types.StringValue("workload-1"),
		"display_name":    types.StringNull(),
		"description":     types.StringNull(),
		"container_image": types.StringNull(),
		"container_port":  types.Int64Null(),
		"container_name":  types.StringNull(),
		"cpu":             types.StringNull(),
		"memory":          types.StringNull(),
		"desired_count":   types.Int64Null(),
		"platform":        types.ObjectNull(components.ComponentAttrTypes),
		"subnet":          types.ObjectNull(components.ComponentAttrTypes),
		"links":           types.DynamicNull(),
		"security_groups": types.ListNull(components.ComponentObjectType),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if attrs["id"].(types.String).ValueString() != "workload-1" {
		t.Errorf("expected id %q", "workload-1")
	}
	if attrs["type"].(types.String).ValueString() != "CustomWorkloads.CaaS.Workload" {
		t.Errorf("expected type %q", "CustomWorkloads.CaaS.Workload")
	}
}

func TestWorkloadFunction_Run_WithDepsAndParams(t *testing.T) {
	f := NewWorkloadFunction()
	platform := buildTestComponent(t, "k8s-1", "NetworkAndCompute.PaaS.ContainerPlatform")
	subnet := buildTestComponent(t, "subnet-1", "NetworkAndCompute.IaaS.Subnet")
	sg := buildTestComponent(t, "sg-1", "NetworkAndCompute.IaaS.SecurityGroup")
	target := buildTestComponent(t, "workload-2", "CustomWorkloads.CaaS.Workload")

	linkSettings, _ := types.MapValue(types.StringType, map[string]attr.Value{
		"fromPort": types.StringValue("8080"),
	})
	genericLink, diags := types.ObjectValue(testLinkAttrTypes, map[string]attr.Value{
		"target":   target,
		"settings": linkSettings,
	})
	if diags.HasError() {
		t.Fatalf("failed to build generic link: %s", diags.Errors())
	}
	linkList, diags := types.ListValue(types.ObjectType{AttrTypes: testLinkAttrTypes}, []attr.Value{genericLink})
	if diags.HasError() {
		t.Fatalf("failed to build link list: %s", diags.Errors())
	}
	sgList, diags := types.ListValue(components.ComponentObjectType, []attr.Value{sg})
	if diags.HasError() {
		t.Fatalf("failed to build sg list: %s", diags.Errors())
	}

	configObj, diags := types.ObjectValue(workloadAttrTypes, map[string]attr.Value{
		"id":              types.StringValue("workload-1"),
		"display_name":    types.StringValue("My Workload"),
		"description":     types.StringNull(),
		"container_image": types.StringValue("nginx:latest"),
		"container_port":  types.Int64Value(80),
		"container_name":  types.StringValue("web"),
		"cpu":             types.StringValue("256"),
		"memory":          types.StringValue("512"),
		"desired_count":   types.Int64Value(3),
		"platform":        platform,
		"subnet":          subnet,
		"links":           types.DynamicValue(linkList),
		"security_groups": sgList,
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	// Check dependencies
	deps := attrs["dependencies_ids"].(types.List)
	if deps.IsNull() {
		t.Fatal("expected non-null dependencies")
	}
	depElems := deps.Elements()
	if len(depElems) != 2 {
		t.Fatalf("expected 2 dependencies, got %d", len(depElems))
	}
	if depElems[0].(types.String).ValueString() != "k8s-1" {
		t.Errorf("expected first dep %q", "k8s-1")
	}
	if depElems[1].(types.String).ValueString() != "subnet-1" {
		t.Errorf("expected second dep %q", "subnet-1")
	}

	// Check parameters
	params := attrs["parameters"].(types.Map)
	elems := params.Elements()
	if elems["containerImage"].(types.String).ValueString() != "nginx:latest" {
		t.Errorf("expected containerImage %q", "nginx:latest")
	}
	if elems["containerPort"].(types.String).ValueString() != "80" {
		t.Errorf("expected containerPort %q", "80")
	}

	// Check links
	linksVal := attrs["links"].(types.List)
	linkElems := linksVal.Elements()
	if len(linkElems) != 2 {
		t.Fatalf("expected 2 links, got %d", len(linkElems))
	}
}
