package spec

import (
	"context"
	"math/big"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
	ft "fractal.cloud/terraform-provider-fc/internal/provider/functions/functiontest"
)

var testSpec = Spec{
	Name:          "test_caas_thing",
	ComponentType: "Test.CaaS.Thing",
	Summary:       "Creates a thing",
	Attributes: []Attribute{
		{Name: "version", Key: "version", Aliases: []string{"engineVersion"}, Kind: String, Required: true},
		{Name: "tier", Key: "tier", Kind: String, OneOf: []string{"small", "large"}},
		{Name: "replicas", Key: "replicas", Kind: Int64},
		{Name: "enabled", Key: "enabled", Kind: Bool},
		{Name: "arguments", Key: "arguments", Kind: StringList},
		{Name: "conf", Key: "conf", Kind: StringMap},
		{Name: "policy", Key: "policy", Kind: Object, Fields: []Attribute{
			{Name: "min_length", Key: "minLength", Kind: Int64},
			{Name: "require_symbols", Key: "requireSymbols", Kind: Bool},
		}},
	},
	Dependencies: []Dependency{{Name: "platform", ComponentType: "NetworkAndCompute.PaaS.ContainerPlatform"}},
	Links:        true,
}

func component(t *testing.T, id, componentType string) types.Object {
	t.Helper()
	obj, err := components.BuildComponent(id, componentType, types.StringNull(), types.StringNull(), types.StringNull(), nil, nil, nil)
	if err != nil {
		t.Fatalf("building %s: %s", id, err.Text)
	}
	return obj
}

func policy(t *testing.T, attrs map[string]attr.Value) attr.Value {
	t.Helper()
	return types.DynamicValue(ft.Object(t, attrs))
}

func TestSpec_DeclaresOnlyIdAndRequiredAttributesMandatory(t *testing.T) {
	f := New(testSpec)
	var def function.DefinitionResponse
	f.Definition(context.Background(), function.DefinitionRequest{}, &def)
	obj := def.Definition.Parameters[0].(function.ObjectParameter).CustomType.TerraformType(context.Background()).(tftypes.Object)

	for _, name := range []string{"id", "version"} {
		if _, optional := obj.OptionalAttributes[name]; optional {
			t.Errorf("%s must be required", name)
		}
	}
	for _, name := range []string{"tier", "replicas", "platform", "links", "extra_parameters", "display_name"} {
		if _, optional := obj.OptionalAttributes[name]; !optional {
			t.Errorf("%s must be optional", name)
		}
	}
}

func TestSpec_Run_WritesEveryKind(t *testing.T) {
	c := ft.Component(t, ft.Run(t, New(testSpec), ft.Object(t, map[string]attr.Value{
		"id":        types.StringValue("thing"),
		"version":   types.StringValue("17"),
		"tier":      types.StringValue("large"),
		"replicas":  types.Int64Value(3),
		"enabled":   types.BoolValue(false),
		"arguments": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("--a"), types.StringValue("b")}),
		"conf":      types.MapValueMust(types.StringType, map[string]attr.Value{"spark.x": types.StringValue("1")}),
		"policy":    policy(t, map[string]attr.Value{"min_length": types.NumberValue(big.NewFloat(12)), "require_symbols": types.BoolValue(true)}),
		"platform":  component(t, "k8s", "NetworkAndCompute.PaaS.ContainerPlatform"),
	})))

	ft.ExpectParameters(t, c, map[string]string{
		"version":       "17",
		"engineVersion": "17",
		"tier":          "large",
		"replicas":      "3",
		"enabled":       "false",
		"arguments":     `["--a","b"]`,
		"conf":          `{"spark.x":"1"}`,
		"policy":        `{"minLength":12,"requireSymbols":true}`,
	})
	if got := c["type"].(types.String).ValueString(); got != "Test.CaaS.Thing" {
		t.Errorf("type = %q", got)
	}
	if deps := ft.Strings(t, c, "dependencies_ids"); len(deps) != 1 || deps[0] != "k8s" {
		t.Errorf("dependencies = %v, want [k8s]", deps)
	}
}

func TestSpec_Run_OmittedAttributesSetNoParameters(t *testing.T) {
	c := ft.Component(t, ft.Run(t, New(testSpec), ft.Object(t, map[string]attr.Value{
		"id":      types.StringValue("thing"),
		"version": types.StringValue("17"),
	})))
	ft.ExpectParameters(t, c, map[string]string{"version": "17", "engineVersion": "17"})
	if deps := ft.Strings(t, c, "dependencies_ids"); len(deps) != 0 {
		t.Errorf("dependencies = %v, want none", deps)
	}
}

func TestSpec_Run_Links(t *testing.T) {
	c := ft.Component(t, ft.Run(t, New(testSpec), ft.Object(t, map[string]attr.Value{
		"id":      types.StringValue("thing"),
		"version": types.StringValue("17"),
		"links": ft.Tuple(t, ft.Object(t, map[string]attr.Value{
			"target": component(t, "idp", "Security.PaaS.IdentityProvider"),
			"settings": ft.Object(t, map[string]attr.Value{
				"clientType": types.StringValue("spa"),
			}),
		})),
	})))
	if ft.Links(t, c)["idp"]["clientType"] != "spa" {
		t.Errorf("links = %v", ft.Links(t, c))
	}
}

func TestSpec_Run_Errors(t *testing.T) {
	tests := []struct {
		name  string
		attrs map[string]attr.Value
		want  string
	}{
		{"value outside OneOf", map[string]attr.Value{"tier": types.StringValue("medium")}, `tier: must be one of "large", "small", got "medium"`},
		{"unknown object field", map[string]attr.Value{"policy": policy(t, map[string]attr.Value{"max_length": types.NumberValue(big.NewFloat(3))})}, `policy: unexpected attribute "max_length"`},
		{"fractional object field", map[string]attr.Value{"policy": policy(t, map[string]attr.Value{"min_length": types.NumberValue(big.NewFloat(1.5))})}, "policy: min_length: must be a whole number"},
		{"dependency of wrong type", map[string]attr.Value{"platform": component(t, "vpc", "NetworkAndCompute.IaaS.VirtualNetwork")}, "platform: expected component of type"},
		{"extra parameter set by attribute", map[string]attr.Value{
			"extra_parameters": types.MapValueMust(types.StringType, map[string]attr.Value{"engineVersion": types.StringValue("16")}),
		}, "extra_parameters.engineVersion is already set"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := map[string]attr.Value{"id": types.StringValue("thing"), "version": types.StringValue("17")}
			for k, v := range tt.attrs {
				attrs[k] = v
			}
			resp := ft.Run(t, New(testSpec), ft.Object(t, attrs))
			if resp.Error == nil || !strings.Contains(resp.Error.Text, tt.want) {
				t.Errorf("error = %v, want it to contain %q", resp.Error, tt.want)
			}
		})
	}
}

func TestNew_PanicsOnDuplicateAttribute(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected a panic for an attribute clashing with a built-in one")
		}
	}()
	New(Spec{Name: "x", ComponentType: "X", Attributes: []Attribute{{Name: "description", Key: "description", Kind: String}}})
}
