package saas

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
	"fractal.cloud/terraform-provider-fc/internal/provider/functions/functiontest"
)

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

// --- Unmanaged ---

func TestUnmanagedFunction_Metadata(t *testing.T) {
	f := NewSaaSUnmanagedFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "observability_saas_unmanaged" {
		t.Errorf("expected name %q, got %q", "observability_saas_unmanaged", resp.Name)
	}
}

func TestUnmanagedFunction_Definition(t *testing.T) {
	f := NewSaaSUnmanagedFunction()
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

func TestUnmanagedFunction_Run(t *testing.T) {
	f := NewSaaSUnmanagedFunction()
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":           types.StringType,
		"display_name": types.StringType,
		"description":  types.StringType,
	}, map[string]attr.Value{
		"id":           types.StringValue("test-id"),
		"display_name": types.StringValue("Test Name"),
		"description":  types.StringNull(),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if id := attrs["id"].(types.String); id.ValueString() != "test-id" {
		t.Errorf("expected id %q, got %q", "test-id", id.ValueString())
	}
	if ct := attrs["type"].(types.String); ct.ValueString() != "Observability.SaaS.Unmanaged" {
		t.Errorf("expected type %q, got %q", "Observability.SaaS.Unmanaged", ct.ValueString())
	}
	if dn := attrs["display_name"].(types.String); dn.ValueString() != "Test Name" {
		t.Errorf("expected display_name %q, got %q", "Test Name", dn.ValueString())
	}
}
