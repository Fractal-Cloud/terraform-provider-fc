package paas

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
	"fractal.cloud/terraform-provider-fc/internal/provider/functions/functiontest"
)

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

// --- ContainerPlatform ---

func TestContainerPlatformFunction_Metadata(t *testing.T) {
	f := NewContainerPlatformFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "network_and_compute_paas_container_platform" {
		t.Errorf("expected name %q, got %q", "network_and_compute_paas_container_platform", resp.Name)
	}
}

func TestContainerPlatformFunction_Definition(t *testing.T) {
	f := NewContainerPlatformFunction()
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

func TestContainerPlatformFunction_Run_Minimal(t *testing.T) {
	f := NewContainerPlatformFunction()
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":           types.StringType,
		"display_name": types.StringType,
		"description":  types.StringType,
		"node_pools": types.ListType{
			ElemType: types.ObjectType{AttrTypes: nodePoolAttrTypes},
		},
	}, map[string]attr.Value{
		"id":           types.StringValue("test-k8s"),
		"display_name": types.StringValue("Test K8s"),
		"description":  types.StringNull(),
		"node_pools":   types.ListNull(types.ObjectType{AttrTypes: nodePoolAttrTypes}),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if id := attrs["id"].(types.String); id.ValueString() != "test-k8s" {
		t.Errorf("expected id %q, got %q", "test-k8s", id.ValueString())
	}
	if ct := attrs["type"].(types.String); ct.ValueString() != "NetworkAndCompute.PaaS.ContainerPlatform" {
		t.Errorf("expected type %q, got %q", "NetworkAndCompute.PaaS.ContainerPlatform", ct.ValueString())
	}
}

func TestContainerPlatformFunction_Run_WithNodePools(t *testing.T) {
	f := NewContainerPlatformFunction()

	pool, poolDiags := types.ObjectValue(nodePoolAttrTypes, map[string]attr.Value{
		"name":                types.StringValue("default-pool"),
		"disk_size_gb":        types.Int64Value(100),
		"min_node_count":      types.Int64Value(1),
		"max_node_count":      types.Int64Value(5),
		"max_pods_per_node":   types.Int64Null(),
		"autoscaling_enabled": types.BoolValue(true),
		"initial_node_count":  types.Int64Value(2),
		"max_surge":           types.Int64Null(),
	})
	if poolDiags.HasError() {
		t.Fatalf("failed to build pool object: %s", poolDiags.Errors())
	}

	poolList, listDiags := types.ListValue(types.ObjectType{AttrTypes: nodePoolAttrTypes}, []attr.Value{pool})
	if listDiags.HasError() {
		t.Fatalf("failed to build pool list: %s", listDiags.Errors())
	}

	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":           types.StringType,
		"display_name": types.StringType,
		"description":  types.StringType,
		"node_pools": types.ListType{
			ElemType: types.ObjectType{AttrTypes: nodePoolAttrTypes},
		},
	}, map[string]attr.Value{
		"id":           types.StringValue("test-k8s"),
		"display_name": types.StringValue("Test K8s"),
		"description":  types.StringNull(),
		"node_pools":   poolList,
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if ct := attrs["type"].(types.String); ct.ValueString() != "NetworkAndCompute.PaaS.ContainerPlatform" {
		t.Errorf("expected type %q, got %q", "NetworkAndCompute.PaaS.ContainerPlatform", ct.ValueString())
	}

	// Verify node pools parameter
	params := attrs["parameters"].(types.Map)
	if params.IsNull() {
		t.Fatal("expected non-null parameters")
	}
	elems := params.Elements()
	nodePoolsJSON, ok := elems["nodePools"].(types.String)
	if !ok {
		t.Fatal("expected nodePools parameter")
	}

	var pools []map[string]interface{}
	if err := json.Unmarshal([]byte(nodePoolsJSON.ValueString()), &pools); err != nil {
		t.Fatalf("failed to unmarshal nodePools: %s", err)
	}
	if len(pools) != 1 {
		t.Fatalf("expected 1 pool, got %d", len(pools))
	}
	if pools[0]["name"] != "default-pool" {
		t.Errorf("expected pool name %q, got %q", "default-pool", pools[0]["name"])
	}
}
