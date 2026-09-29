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

// containerPlatformConfig builds a config object for functions that take id, display_name, description, container_platform.
func containerPlatformConfig(t *testing.T, id string, platform attr.Value) types.Object {
	t.Helper()
	obj, diags := types.ObjectValue(map[string]attr.Type{
		"id":                 types.StringType,
		"display_name":       types.StringType,
		"description":        types.StringType,
		"container_platform": components.ComponentObjectType,
	}, map[string]attr.Value{
		"id":                 types.StringValue(id),
		"display_name":       types.StringNull(),
		"description":        types.StringNull(),
		"container_platform": platform,
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}
	return obj
}

// --- Monitoring ---

func TestMonitoringFunction_Metadata(t *testing.T) {
	f := NewCaaSMonitoringFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "observability_caas_monitoring" {
		t.Errorf("expected name %q, got %q", "observability_caas_monitoring", resp.Name)
	}
}

func TestMonitoringFunction_Definition(t *testing.T) {
	f := NewCaaSMonitoringFunction()
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

func TestMonitoringFunction_Run_Minimal(t *testing.T) {
	f := NewCaaSMonitoringFunction()
	config := containerPlatformConfig(t, "mon-1", types.ObjectNull(components.ComponentAttrTypes))
	resp := runFunction(t, f, []attr.Value{config})
	attrs := getResultAttrs(t, resp)

	if attrs["id"].(types.String).ValueString() != "mon-1" {
		t.Errorf("expected id %q", "mon-1")
	}
	if attrs["type"].(types.String).ValueString() != "Observability.CaaS.Monitoring" {
		t.Errorf("expected type %q", "Observability.CaaS.Monitoring")
	}
}

func TestMonitoringFunction_Run_WithPlatform(t *testing.T) {
	f := NewCaaSMonitoringFunction()
	platform := buildTestComponent(t, "k8s-1", "NetworkAndCompute.PaaS.ContainerPlatform")
	config := containerPlatformConfig(t, "mon-1", platform)
	resp := runFunction(t, f, []attr.Value{config})
	attrs := getResultAttrs(t, resp)

	deps := attrs["dependencies_ids"].(types.List)
	if deps.IsNull() {
		t.Fatal("expected non-null dependencies")
	}
	if deps.Elements()[0].(types.String).ValueString() != "k8s-1" {
		t.Errorf("expected dependency %q", "k8s-1")
	}
}

// --- Tracing ---

func TestTracingFunction_Metadata(t *testing.T) {
	f := NewCaaSTracingFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "observability_caas_tracing" {
		t.Errorf("expected name %q, got %q", "observability_caas_tracing", resp.Name)
	}
}

func TestTracingFunction_Definition(t *testing.T) {
	f := NewCaaSTracingFunction()
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

func TestTracingFunction_Run(t *testing.T) {
	f := NewCaaSTracingFunction()
	platform := buildTestComponent(t, "k8s-1", "NetworkAndCompute.PaaS.ContainerPlatform")
	config := containerPlatformConfig(t, "trace-1", platform)
	resp := runFunction(t, f, []attr.Value{config})
	attrs := getResultAttrs(t, resp)

	if attrs["type"].(types.String).ValueString() != "Observability.CaaS.Tracing" {
		t.Errorf("expected type %q", "Observability.CaaS.Tracing")
	}
	deps := attrs["dependencies_ids"].(types.List)
	if deps.IsNull() || deps.Elements()[0].(types.String).ValueString() != "k8s-1" {
		t.Errorf("expected dependency %q", "k8s-1")
	}
}

// --- Logging ---

func TestLoggingFunction_Metadata(t *testing.T) {
	f := NewCaaSLoggingFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "observability_caas_logging" {
		t.Errorf("expected name %q, got %q", "observability_caas_logging", resp.Name)
	}
}

func TestLoggingFunction_Definition(t *testing.T) {
	f := NewCaaSLoggingFunction()
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

func TestLoggingFunction_Run(t *testing.T) {
	f := NewCaaSLoggingFunction()
	platform := buildTestComponent(t, "k8s-1", "NetworkAndCompute.PaaS.ContainerPlatform")
	config := containerPlatformConfig(t, "log-1", platform)
	resp := runFunction(t, f, []attr.Value{config})
	attrs := getResultAttrs(t, resp)

	if attrs["type"].(types.String).ValueString() != "Observability.CaaS.Logging" {
		t.Errorf("expected type %q", "Observability.CaaS.Logging")
	}
	deps := attrs["dependencies_ids"].(types.List)
	if deps.IsNull() || deps.Elements()[0].(types.String).ValueString() != "k8s-1" {
		t.Errorf("expected dependency %q", "k8s-1")
	}
}

func TestLoggingFunction_Run_WrongPlatformType(t *testing.T) {
	f := NewCaaSLoggingFunction()
	wrong := buildTestComponent(t, "vpc-1", "NetworkAndCompute.IaaS.VirtualNetwork")
	config := containerPlatformConfig(t, "log-1", wrong)
	resp := runFunction(t, f, []attr.Value{config})
	if resp.Error == nil {
		t.Fatal("expected error for wrong platform type")
	}
}
