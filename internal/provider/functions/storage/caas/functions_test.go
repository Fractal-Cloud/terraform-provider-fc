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

func TestSearchFunction_Metadata(t *testing.T) {
	f := NewStorageCaasSearchFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_caas_search" {
		t.Errorf("expected name %q, got %q", "storage_caas_search", resp.Name)
	}
}

func TestSearchFunction_Definition(t *testing.T) {
	f := NewStorageCaasSearchFunction()
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

func TestSearchFunction_Run_WithPlatform(t *testing.T) {
	f := NewStorageCaasSearchFunction()
	platform := buildTestComponent(t, "k8s-1", "NetworkAndCompute.PaaS.ContainerPlatform")
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":                 types.StringType,
		"display_name":       types.StringType,
		"description":        types.StringType,
		"container_platform": components.ComponentObjectType,
	}, map[string]attr.Value{
		"id":                 types.StringValue("search-1"),
		"display_name":       types.StringNull(),
		"description":        types.StringNull(),
		"container_platform": platform,
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if attrs["type"].(types.String).ValueString() != "Storage.CaaS.Search" {
		t.Errorf("expected type %q", "Storage.CaaS.Search")
	}
	deps := attrs["dependencies_ids"].(types.List)
	if deps.IsNull() || deps.Elements()[0].(types.String).ValueString() != "k8s-1" {
		t.Errorf("expected dependency %q", "k8s-1")
	}
}

func TestSearchEntityFunction_Metadata(t *testing.T) {
	f := NewStorageCaasSearchEntityFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_caas_search_entity" {
		t.Errorf("expected name %q, got %q", "storage_caas_search_entity", resp.Name)
	}
}

func TestSearchEntityFunction_Definition(t *testing.T) {
	f := NewStorageCaasSearchEntityFunction()
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

func TestSearchEntityFunction_Run_WithSearch(t *testing.T) {
	f := NewStorageCaasSearchEntityFunction()
	search := buildTestComponent(t, "search-1", "Storage.CaaS.Search")
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":           types.StringType,
		"display_name": types.StringType,
		"description":  types.StringType,
		"search":       components.ComponentObjectType,
	}, map[string]attr.Value{
		"id":           types.StringValue("entity-1"),
		"display_name": types.StringNull(),
		"description":  types.StringNull(),
		"search":       search,
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if attrs["type"].(types.String).ValueString() != "Storage.CaaS.SearchEntity" {
		t.Errorf("expected type %q", "Storage.CaaS.SearchEntity")
	}
	deps := attrs["dependencies_ids"].(types.List)
	if deps.IsNull() || deps.Elements()[0].(types.String).ValueString() != "search-1" {
		t.Errorf("expected dependency %q", "search-1")
	}
}
