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

// --- Broker ---

func TestBrokerFunction_Metadata(t *testing.T) {
	f := NewMessagingPaasBrokerFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "messaging_paas_broker" {
		t.Errorf("expected name %q, got %q", "messaging_paas_broker", resp.Name)
	}
}

func TestBrokerFunction_Definition(t *testing.T) {
	f := NewMessagingPaasBrokerFunction()
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

func TestBrokerFunction_Run(t *testing.T) {
	f := NewMessagingPaasBrokerFunction()
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":           types.StringType,
		"display_name": types.StringType,
		"description":  types.StringType,
	}, map[string]attr.Value{
		"id":           types.StringValue("test-broker"),
		"display_name": types.StringValue("Test Broker"),
		"description":  types.StringNull(),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if id := attrs["id"].(types.String); id.ValueString() != "test-broker" {
		t.Errorf("expected id %q, got %q", "test-broker", id.ValueString())
	}
	if ct := attrs["type"].(types.String); ct.ValueString() != "Messaging.PaaS.Broker" {
		t.Errorf("expected type %q, got %q", "Messaging.PaaS.Broker", ct.ValueString())
	}
}

// --- Entity ---

func TestEntityFunction_Metadata(t *testing.T) {
	f := NewMessagingPaasEntityFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "messaging_paas_entity" {
		t.Errorf("expected name %q, got %q", "messaging_paas_entity", resp.Name)
	}
}

func TestEntityFunction_Definition(t *testing.T) {
	f := NewMessagingPaasEntityFunction()
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

func TestEntityFunction_Run_Minimal(t *testing.T) {
	f := NewMessagingPaasEntityFunction()
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":                      types.StringType,
		"display_name":            types.StringType,
		"description":             types.StringType,
		"message_retention_hours": types.Int64Type,
		"broker":                  components.ComponentObjectType,
	}, map[string]attr.Value{
		"id":                      types.StringValue("test-entity"),
		"display_name":            types.StringValue("Test Entity"),
		"description":             types.StringNull(),
		"message_retention_hours": types.Int64Null(),
		"broker":                  types.ObjectNull(components.ComponentAttrTypes),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if id := attrs["id"].(types.String); id.ValueString() != "test-entity" {
		t.Errorf("expected id %q, got %q", "test-entity", id.ValueString())
	}
	if ct := attrs["type"].(types.String); ct.ValueString() != "Messaging.PaaS.Entity" {
		t.Errorf("expected type %q, got %q", "Messaging.PaaS.Entity", ct.ValueString())
	}
}

func TestEntityFunction_Run_WithBrokerAndParams(t *testing.T) {
	f := NewMessagingPaasEntityFunction()
	broker := buildTestComponent(t, "my-broker", "Messaging.PaaS.Broker")

	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":                      types.StringType,
		"display_name":            types.StringType,
		"description":             types.StringType,
		"message_retention_hours": types.Int64Type,
		"broker":                  components.ComponentObjectType,
	}, map[string]attr.Value{
		"id":                      types.StringValue("test-entity"),
		"display_name":            types.StringValue("Test Entity"),
		"description":             types.StringNull(),
		"message_retention_hours": types.Int64Value(48),
		"broker":                  broker,
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if ct := attrs["type"].(types.String); ct.ValueString() != "Messaging.PaaS.Entity" {
		t.Errorf("expected type %q, got %q", "Messaging.PaaS.Entity", ct.ValueString())
	}

	// Check dependency on broker
	deps := attrs["dependencies_ids"].(types.List)
	if deps.IsNull() {
		t.Fatal("expected non-null dependencies")
	}
	depElems := deps.Elements()
	if len(depElems) != 1 || depElems[0].(types.String).ValueString() != "my-broker" {
		t.Errorf("expected dependency [my-broker], got %v", depElems)
	}

	// Check parameter
	params := attrs["parameters"].(types.Map)
	if params.IsNull() {
		t.Fatal("expected non-null parameters")
	}
	elems := params.Elements()
	if elems["messageRetentionHours"].(types.String).ValueString() != "48" {
		t.Errorf("expected messageRetentionHours %q, got %q", "48", elems["messageRetentionHours"].(types.String).ValueString())
	}
}
