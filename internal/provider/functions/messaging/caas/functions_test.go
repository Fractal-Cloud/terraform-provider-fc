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

func TestBrokerFunction_Metadata(t *testing.T) {
	f := NewMessagingCaasBrokerFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "messaging_caas_broker" {
		t.Errorf("expected name %q, got %q", "messaging_caas_broker", resp.Name)
	}
}

func TestBrokerFunction_Definition(t *testing.T) {
	f := NewMessagingCaasBrokerFunction()
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

func TestBrokerFunction_Run_WithPlatform(t *testing.T) {
	f := NewMessagingCaasBrokerFunction()
	platform := buildTestComponent(t, "k8s-1", "NetworkAndCompute.PaaS.ContainerPlatform")
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":                 types.StringType,
		"display_name":       types.StringType,
		"description":        types.StringType,
		"container_platform": components.ComponentObjectType,
	}, map[string]attr.Value{
		"id":                 types.StringValue("broker-1"),
		"display_name":       types.StringValue("My Broker"),
		"description":        types.StringNull(),
		"container_platform": platform,
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if attrs["type"].(types.String).ValueString() != "Messaging.CaaS.Broker" {
		t.Errorf("expected type %q", "Messaging.CaaS.Broker")
	}
	deps := attrs["dependencies_ids"].(types.List)
	if deps.IsNull() || deps.Elements()[0].(types.String).ValueString() != "k8s-1" {
		t.Errorf("expected dependency %q", "k8s-1")
	}
}

func TestEntityFunction_Metadata(t *testing.T) {
	f := NewMessagingCaasEntityFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "messaging_caas_entity" {
		t.Errorf("expected name %q, got %q", "messaging_caas_entity", resp.Name)
	}
}

func TestEntityFunction_Definition(t *testing.T) {
	f := NewMessagingCaasEntityFunction()
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

func TestEntityFunction_Run_WithBroker(t *testing.T) {
	f := NewMessagingCaasEntityFunction()
	broker := buildTestComponent(t, "broker-1", "Messaging.CaaS.Broker")
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":           types.StringType,
		"display_name": types.StringType,
		"description":  types.StringType,
		"broker":       components.ComponentObjectType,
	}, map[string]attr.Value{
		"id":           types.StringValue("entity-1"),
		"display_name": types.StringNull(),
		"description":  types.StringNull(),
		"broker":       broker,
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if attrs["type"].(types.String).ValueString() != "Messaging.CaaS.Entity" {
		t.Errorf("expected type %q", "Messaging.CaaS.Entity")
	}
	deps := attrs["dependencies_ids"].(types.List)
	if deps.IsNull() || deps.Elements()[0].(types.String).ValueString() != "broker-1" {
		t.Errorf("expected dependency %q", "broker-1")
	}
}
