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

func TestServiceMeshSecurityFunction_Metadata(t *testing.T) {
	f := NewCaaSServiceMeshSecurityFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "security_caas_service_mesh_security" {
		t.Errorf("expected name %q, got %q", "security_caas_service_mesh_security", resp.Name)
	}
}

func TestServiceMeshSecurityFunction_Definition(t *testing.T) {
	f := NewCaaSServiceMeshSecurityFunction()
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

func TestServiceMeshSecurityFunction_Run_Minimal(t *testing.T) {
	f := NewCaaSServiceMeshSecurityFunction()
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":                 types.StringType,
		"display_name":       types.StringType,
		"description":        types.StringType,
		"container_platform": components.ComponentObjectType,
	}, map[string]attr.Value{
		"id":                 types.StringValue("sms-1"),
		"display_name":       types.StringNull(),
		"description":        types.StringNull(),
		"container_platform": types.ObjectNull(components.ComponentAttrTypes),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if attrs["id"].(types.String).ValueString() != "sms-1" {
		t.Errorf("expected id %q", "sms-1")
	}
	if attrs["type"].(types.String).ValueString() != "Security.CaaS.ServiceMeshSecurity" {
		t.Errorf("expected type %q", "Security.CaaS.ServiceMeshSecurity")
	}
}

func TestServiceMeshSecurityFunction_Run_WithPlatform(t *testing.T) {
	f := NewCaaSServiceMeshSecurityFunction()
	platform := buildTestComponent(t, "k8s-1", "NetworkAndCompute.PaaS.ContainerPlatform")
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":                 types.StringType,
		"display_name":       types.StringType,
		"description":        types.StringType,
		"container_platform": components.ComponentObjectType,
	}, map[string]attr.Value{
		"id":                 types.StringValue("sms-1"),
		"display_name":       types.StringNull(),
		"description":        types.StringNull(),
		"container_platform": platform,
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	deps := attrs["dependencies_ids"].(types.List)
	if deps.IsNull() {
		t.Fatal("expected non-null dependencies")
	}
	if deps.Elements()[0].(types.String).ValueString() != "k8s-1" {
		t.Errorf("expected dependency %q", "k8s-1")
	}
}
