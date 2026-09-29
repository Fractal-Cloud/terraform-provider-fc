// Package functiontest holds helpers shared by the component function tests.
package functiontest

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

// Complete does to test arguments what Terraform does before a call: every
// attribute a config parameter declares but the argument omits is set to
// null, so a test states only the attributes it is about. An attribute the
// function does not declare fails the test.
func Complete(t *testing.T, f function.Function, args []attr.Value) []attr.Value {
	t.Helper()
	ctx := context.Background()

	var def function.DefinitionResponse
	f.Definition(ctx, function.DefinitionRequest{}, &def)

	completed := make([]attr.Value, len(args))
	copy(completed, args)
	for i, param := range def.Definition.Parameters {
		if i >= len(args) {
			break
		}
		objParam, ok := param.(function.ObjectParameter)
		if !ok {
			continue
		}
		configType, ok := objParam.CustomType.(components.ConfigObjectType)
		if !ok {
			continue
		}
		given, ok := args[i].(types.Object)
		if !ok || given.IsNull() || given.IsUnknown() {
			continue
		}

		values := make(map[string]attr.Value, len(configType.AttrTypes))
		for name, typ := range configType.AttrTypes {
			if v, ok := given.Attributes()[name]; ok {
				values[name] = v
				continue
			}
			null, err := typ.ValueFromTerraform(ctx, tftypes.NewValue(typ.TerraformType(ctx), nil))
			if err != nil {
				t.Fatalf("building null %s for %q: %v", typ, name, err)
			}
			values[name] = null
		}
		for name := range given.Attributes() {
			if _, ok := configType.AttrTypes[name]; !ok {
				t.Fatalf("argument sets %q, which the function does not declare", name)
			}
		}

		obj, diags := types.ObjectValue(configType.AttrTypes, values)
		if diags.HasError() {
			t.Fatalf("completing config argument: %v", diags)
		}
		completed[i] = obj
	}
	return completed
}
