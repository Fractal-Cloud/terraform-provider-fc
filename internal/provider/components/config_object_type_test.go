package components

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var testConfigType = NewConfigObjectType(map[string]attr.Type{
	"id":    types.StringType,
	"label": types.StringType,
	"links": types.ListType{ElemType: types.StringType},
}, "id")

func TestConfigObjectType_DeclaresNonRequiredAttributesOptional(t *testing.T) {
	obj, ok := testConfigType.TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatalf("TerraformType is %T, want tftypes.Object", testConfigType.TerraformType(context.Background()))
	}
	if _, ok := obj.OptionalAttributes["id"]; ok {
		t.Error("id must stay required")
	}
	for _, name := range []string{"label", "links"} {
		if _, ok := obj.OptionalAttributes[name]; !ok {
			t.Errorf("%s must be optional", name)
		}
	}
}

func TestConfigObjectType_AllRequiredDeclaresNoOptionalAttributes(t *testing.T) {
	typ := NewConfigObjectType(map[string]attr.Type{"id": types.StringType}, "id")
	obj := typ.TerraformType(context.Background()).(tftypes.Object)
	if obj.OptionalAttributes != nil {
		t.Errorf("OptionalAttributes = %v, want nil", obj.OptionalAttributes)
	}
}

func TestNewConfigObjectType_PanicsOnUnknownRequiredAttribute(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected a panic for an undeclared required attribute")
		}
	}()
	NewConfigObjectType(map[string]attr.Type{"id": types.StringType}, "name")
}

// Terraform sends omitted optional attributes as nulls, in a value whose own
// type carries no optional markers.
func TestConfigObjectType_ValueFromTerraform_DecodesArgument(t *testing.T) {
	ctx := context.Background()
	linkType := tftypes.String
	plain := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"id":    tftypes.String,
		"label": tftypes.String,
		"links": tftypes.List{ElementType: linkType},
	}}
	in := tftypes.NewValue(plain, map[string]tftypes.Value{
		"id":    tftypes.NewValue(tftypes.String, "db"),
		"label": tftypes.NewValue(tftypes.String, nil),
		"links": tftypes.NewValue(tftypes.List{ElementType: linkType}, nil),
	})

	v, err := testConfigType.ValueFromTerraform(ctx, in)
	if err != nil {
		t.Fatalf("ValueFromTerraform: %v", err)
	}
	cv, ok := v.(ConfigObjectValue)
	if !ok {
		t.Fatalf("value is %T, want ConfigObjectValue", v)
	}
	if !cv.Type(ctx).Equal(testConfigType) {
		t.Error("value does not report its ConfigObjectType")
	}
	if got := cv.Attributes()["id"].(types.String).ValueString(); got != "db" {
		t.Errorf("id = %q, want db", got)
	}
	if !cv.Attributes()["label"].IsNull() {
		t.Error("omitted label must decode as null")
	}
}

func TestConfigObjectType_ValueFromTerraform_RejectsOtherAttributeSets(t *testing.T) {
	in := tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String}},
		map[string]tftypes.Value{"id": tftypes.NewValue(tftypes.String, "x")})
	if _, err := testConfigType.ValueFromTerraform(context.Background(), in); err == nil {
		t.Error("expected an error for a value missing attributes")
	}
}

func TestConfigObjectType_ValueFromTerraform_NullAndUnknown(t *testing.T) {
	ctx := context.Background()
	plain := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"id":    tftypes.String,
		"label": tftypes.String,
		"links": tftypes.List{ElementType: tftypes.String},
	}}

	null, err := testConfigType.ValueFromTerraform(ctx, tftypes.NewValue(plain, nil))
	if err != nil || !null.IsNull() {
		t.Errorf("null value = %v, %v; want a null value", null, err)
	}
	unknown, err := testConfigType.ValueFromTerraform(ctx, tftypes.NewValue(plain, tftypes.UnknownValue))
	if err != nil || !unknown.IsUnknown() {
		t.Errorf("unknown value = %v, %v; want an unknown value", unknown, err)
	}
}

func TestConfigObjectType_Equal(t *testing.T) {
	same := NewConfigObjectType(map[string]attr.Type{
		"id":    types.StringType,
		"label": types.StringType,
		"links": types.ListType{ElemType: types.StringType},
	}, "id")
	otherRequired := NewConfigObjectType(map[string]attr.Type{
		"id":    types.StringType,
		"label": types.StringType,
		"links": types.ListType{ElemType: types.StringType},
	}, "id", "label")

	if !testConfigType.Equal(same) {
		t.Error("identical types must be equal")
	}
	if testConfigType.Equal(otherRequired) {
		t.Error("types with different required sets must differ")
	}
	if testConfigType.Equal(basetypes.ObjectType{AttrTypes: testConfigType.AttrTypes}) {
		t.Error("a plain ObjectType must not equal a ConfigObjectType")
	}
}
