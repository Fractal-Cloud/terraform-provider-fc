package functiontest

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

// Object builds an object value whose type is inferred from its values, the
// way a Terraform object literal is typed.
func Object(t *testing.T, attrs map[string]attr.Value) types.Object {
	t.Helper()
	attrTypes := make(map[string]attr.Type, len(attrs))
	for k, v := range attrs {
		attrTypes[k] = v.Type(context.Background())
	}
	obj, diags := types.ObjectValue(attrTypes, attrs)
	if diags.HasError() {
		t.Fatalf("building object: %v", diags)
	}
	return obj
}

// Tuple builds a dynamic value holding a tuple, which is what Terraform
// passes for a list literal such as `links = [ {...}, {...} ]`.
func Tuple(t *testing.T, elems ...attr.Value) types.Dynamic {
	t.Helper()
	elemTypes := make([]attr.Type, len(elems))
	for i, e := range elems {
		elemTypes[i] = e.Type(context.Background())
	}
	tuple, diags := types.TupleValue(elemTypes, elems)
	if diags.HasError() {
		t.Fatalf("building tuple: %v", diags)
	}
	return types.DynamicValue(tuple)
}

// Run calls a component function with a single config argument, completed
// as Terraform would complete it.
func Run(t *testing.T, f function.Function, config types.Object) *function.RunResponse {
	t.Helper()
	resp := &function.RunResponse{Result: function.NewResultData(types.ObjectNull(components.ComponentAttrTypes))}
	f.Run(context.Background(), function.RunRequest{
		Arguments: function.NewArgumentsData(Complete(t, f, []attr.Value{config})),
	}, resp)
	return resp
}

// Component returns the attributes of the component a successful Run built.
func Component(t *testing.T, resp *function.RunResponse) map[string]attr.Value {
	t.Helper()
	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error.Text)
	}
	obj, ok := resp.Result.Value().(types.Object)
	if !ok {
		t.Fatalf("expected an object result, got %T", resp.Result.Value())
	}
	return obj.Attributes()
}

// Parameters returns a built component's parameters.
func Parameters(t *testing.T, component map[string]attr.Value) map[string]string {
	t.Helper()
	params := map[string]string{}
	m, ok := component["parameters"].(types.Map)
	if !ok || m.IsNull() {
		return params
	}
	for k, v := range m.Elements() {
		params[k] = v.(types.String).ValueString()
	}
	return params
}

// Strings returns the elements of a built component's string list attribute
// (dependencies_ids).
func Strings(t *testing.T, component map[string]attr.Value, name string) []string {
	t.Helper()
	l, ok := component[name].(types.List)
	if !ok || l.IsNull() {
		return nil
	}
	out := make([]string, len(l.Elements()))
	for i, v := range l.Elements() {
		out[i] = v.(types.String).ValueString()
	}
	return out
}

// Links returns a built component's links as target id → settings.
func Links(t *testing.T, component map[string]attr.Value) map[string]map[string]string {
	t.Helper()
	out := map[string]map[string]string{}
	l, ok := component["links"].(types.List)
	if !ok || l.IsNull() {
		return out
	}
	for _, v := range l.Elements() {
		attrs := v.(types.Object).Attributes()
		settings := map[string]string{}
		if m, ok := attrs["settings"].(types.Map); ok && !m.IsNull() {
			for k, s := range m.Elements() {
				settings[k] = s.(types.String).ValueString()
			}
		}
		out[attrs["component_id"].(types.String).ValueString()] = settings
	}
	return out
}

// ExpectParameters fails the test unless the component's parameters are
// exactly want.
func ExpectParameters(t *testing.T, component map[string]attr.Value, want map[string]string) {
	t.Helper()
	got := Parameters(t, component)
	if len(got) != len(want) {
		t.Errorf("parameters = %v, want %v", got, want)
		return
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("parameters[%s] = %q, want %q (all: %v)", k, got[k], v, got)
		}
	}
}
