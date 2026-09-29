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

// --- APIGateway ---

func TestAPIGatewayFunction_Metadata(t *testing.T) {
	f := NewCaaSAPIGatewayFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "api_management_caas_api_gateway" {
		t.Errorf("expected name %q, got %q", "api_management_caas_api_gateway", resp.Name)
	}
}

func TestAPIGatewayFunction_Definition(t *testing.T) {
	f := NewCaaSAPIGatewayFunction()
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

func TestAPIGatewayFunction_Run_NullDependency(t *testing.T) {
	f := NewCaaSAPIGatewayFunction()
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":                 types.StringType,
		"display_name":       types.StringType,
		"description":        types.StringType,
		"container_platform": components.ComponentObjectType,
		"links":              types.DynamicType,
	}, map[string]attr.Value{
		"id":                 types.StringValue("test-gw"),
		"display_name":       types.StringValue("Test Gateway"),
		"description":        types.StringNull(),
		"container_platform": types.ObjectNull(components.ComponentAttrTypes),
		"links":              types.DynamicNull(),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if id := attrs["id"].(types.String); id.ValueString() != "test-gw" {
		t.Errorf("expected id %q, got %q", "test-gw", id.ValueString())
	}
	if ct := attrs["type"].(types.String); ct.ValueString() != "APIManagement.CaaS.APIGateway" {
		t.Errorf("expected type %q, got %q", "APIManagement.CaaS.APIGateway", ct.ValueString())
	}
	if dn := attrs["display_name"].(types.String); dn.ValueString() != "Test Gateway" {
		t.Errorf("expected display_name %q, got %q", "Test Gateway", dn.ValueString())
	}
}

func TestAPIGatewayFunction_Run_WithDependency(t *testing.T) {
	f := NewCaaSAPIGatewayFunction()
	cp := buildTestComponent(t, "my-platform", "NetworkAndCompute.PaaS.ContainerPlatform")
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":                 types.StringType,
		"display_name":       types.StringType,
		"description":        types.StringType,
		"container_platform": components.ComponentObjectType,
		"links":              types.DynamicType,
	}, map[string]attr.Value{
		"id":                 types.StringValue("test-gw"),
		"display_name":       types.StringNull(),
		"description":        types.StringNull(),
		"container_platform": cp,
		"links":              types.DynamicNull(),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if id := attrs["id"].(types.String); id.ValueString() != "test-gw" {
		t.Errorf("expected id %q, got %q", "test-gw", id.ValueString())
	}
	if ct := attrs["type"].(types.String); ct.ValueString() != "APIManagement.CaaS.APIGateway" {
		t.Errorf("expected type %q, got %q", "APIManagement.CaaS.APIGateway", ct.ValueString())
	}
}
