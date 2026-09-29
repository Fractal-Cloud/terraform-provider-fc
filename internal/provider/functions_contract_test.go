package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

// Every registered function must be callable from Terraform: a unique name,
// a config argument whose only mandatory attributes are deliberate, and no
// type terraform-plugin-go cannot decode.
func TestFunctions_Contract(t *testing.T) {
	ctx := context.Background()
	p := &fractalCloudProvider{}
	seen := map[string]bool{}

	for _, newFunction := range p.Functions(ctx) {
		f := newFunction()
		var meta function.MetadataResponse
		f.Metadata(ctx, function.MetadataRequest{}, &meta)
		name := meta.Name
		if seen[name] {
			t.Errorf("function %q is registered twice", name)
		}
		seen[name] = true

		var def function.DefinitionResponse
		f.Definition(ctx, function.DefinitionRequest{}, &def)
		if def.Definition.Summary == "" {
			t.Errorf("%s: missing summary", name)
		}

		for _, param := range def.Definition.Parameters {
			objParam, ok := param.(function.ObjectParameter)
			if !ok {
				continue
			}
			if len(objParam.AttributeTypes) > 0 {
				t.Errorf("%s: config declares AttributeTypes, which makes every attribute required; use components.NewConfigObjectType", name)
				continue
			}
			configType, ok := objParam.CustomType.(components.ConfigObjectType)
			if !ok {
				t.Errorf("%s: config type is %T, want components.ConfigObjectType", name, objParam.CustomType)
				continue
			}

			tfType := configType.TerraformType(ctx).(tftypes.Object)
			if _, optional := tfType.OptionalAttributes["id"]; optional {
				t.Errorf("%s: id must be required", name)
			}
			for attrName, attrType := range configType.AttrTypes {
				if hasOptionalAttributes(ctx, attrType) {
					t.Errorf("%s: attribute %q nests a type with optional attributes, which terraform-plugin-go panics decoding", name, attrName)
				}
			}
		}
	}

	if len(seen) == 0 {
		t.Fatal("no functions registered")
	}
}

func hasOptionalAttributes(ctx context.Context, typ attr.Type) bool {
	return tfTypeHasOptionalAttributes(typ.TerraformType(ctx))
}

func tfTypeHasOptionalAttributes(typ tftypes.Type) bool {
	switch tt := typ.(type) {
	case tftypes.Object:
		if len(tt.OptionalAttributes) > 0 {
			return true
		}
		for _, at := range tt.AttributeTypes {
			if tfTypeHasOptionalAttributes(at) {
				return true
			}
		}
	case tftypes.List:
		return tfTypeHasOptionalAttributes(tt.ElementType)
	case tftypes.Set:
		return tfTypeHasOptionalAttributes(tt.ElementType)
	case tftypes.Map:
		return tfTypeHasOptionalAttributes(tt.ElementType)
	case tftypes.Tuple:
		for _, et := range tt.ElementTypes {
			if tfTypeHasOptionalAttributes(et) {
				return true
			}
		}
	}
	return false
}

// A nested ConfigObjectType is exactly the shape the contract test rejects.
func TestFunctions_ContractDetectsNestedOptionalAttributes(t *testing.T) {
	nested := components.NewConfigObjectType(map[string]attr.Type{"x": basetypes.StringType{}})
	if !hasOptionalAttributes(context.Background(), basetypes.ListType{ElemType: nested}) {
		t.Error("a list of ConfigObjectType must be detected")
	}
	if hasOptionalAttributes(context.Background(), components.ComponentObjectType) {
		t.Error("a plain object type must not be flagged")
	}
}
