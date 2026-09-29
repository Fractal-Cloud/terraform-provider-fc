package components

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ basetypes.ObjectTypable = ConfigObjectType{}

// ConfigObjectType is an object type for function arguments in which every
// attribute not listed as required may be omitted by the caller.
//
// The framework's own ObjectType declares every attribute as required, so a
// function taking `{ id = "db" }` against a five-attribute object fails in
// Terraform with "attributes ... are required". Terraform supports optional
// object attributes in type constraints; this type declares them, and
// Terraform then fills omitted attributes with null before the call.
//
// Use it only as the type of a function's top-level object parameter.
// terraform-plugin-go panics decoding a value in which a type with optional
// attributes is nested, whether as a list element or as a null attribute,
// so nested objects whose fields are optional are declared dynamic and
// checked in Run instead (see LinksAttrType).
type ConfigObjectType struct {
	basetypes.ObjectType

	// required lists the attributes the caller must set. Kept as a sorted
	// slice rather than a map so the type stays cheap to compare.
	required []string
}

// NewConfigObjectType builds a ConfigObjectType over attrTypes in which only
// the attributes named in required are mandatory.
func NewConfigObjectType(attrTypes map[string]attr.Type, required ...string) ConfigObjectType {
	for _, name := range required {
		if _, ok := attrTypes[name]; !ok {
			panic(fmt.Sprintf("required attribute %q is not declared", name))
		}
	}
	sorted := append([]string(nil), required...)
	sort.Strings(sorted)
	return ConfigObjectType{
		ObjectType: basetypes.ObjectType{AttrTypes: attrTypes},
		required:   sorted,
	}
}

func (t ConfigObjectType) isRequired(name string) bool {
	i := sort.SearchStrings(t.required, name)
	return i < len(t.required) && t.required[i] == name
}

// TerraformType returns the object type constraint with every non-required
// attribute declared optional.
func (t ConfigObjectType) TerraformType(ctx context.Context) tftypes.Type {
	attributeTypes := make(map[string]tftypes.Type, len(t.AttrTypes))
	optional := map[string]struct{}{}
	for name, typ := range t.AttrTypes {
		attributeTypes[name] = typ.TerraformType(ctx)
		if !t.isRequired(name) {
			optional[name] = struct{}{}
		}
	}
	if len(optional) == 0 {
		optional = nil
	}
	return tftypes.Object{
		AttributeTypes:     attributeTypes,
		OptionalAttributes: optional,
	}
}

// ValueFromTerraform decodes an argument value. Terraform sends every
// attribute (null when omitted), and the value's own type never carries the
// optional markers, so the value is matched against the attribute set rather
// than against TerraformType.
func (t ConfigObjectType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return t.nullValue(), nil
	}
	inType, ok := in.Type().(tftypes.Object)
	if !ok {
		return nil, fmt.Errorf("expected %s, got %s", t.String(), in.Type())
	}
	if len(inType.AttributeTypes) != len(t.AttrTypes) {
		return nil, fmt.Errorf("expected %s, got %s", t.String(), in.Type())
	}
	for name := range inType.AttributeTypes {
		if _, ok := t.AttrTypes[name]; !ok {
			return nil, fmt.Errorf("unexpected attribute %q for %s", name, t.String())
		}
	}
	if !in.IsKnown() {
		return ConfigObjectValue{ObjectValue: basetypes.NewObjectUnknown(t.AttrTypes), typ: t}, nil
	}
	if in.IsNull() {
		return t.nullValue(), nil
	}

	raw := map[string]tftypes.Value{}
	if err := in.As(&raw); err != nil {
		return nil, err
	}
	attributes := make(map[string]attr.Value, len(raw))
	for name, v := range raw {
		a, err := t.AttrTypes[name].ValueFromTerraform(ctx, v)
		if err != nil {
			return nil, err
		}
		attributes[name] = a
	}
	obj, diags := basetypes.NewObjectValue(t.AttrTypes, attributes)
	if diags.HasError() {
		return nil, fmt.Errorf("building %s: %v", t.String(), diags)
	}
	return ConfigObjectValue{ObjectValue: obj, typ: t}, nil
}

func (t ConfigObjectType) nullValue() ConfigObjectValue {
	return ConfigObjectValue{ObjectValue: basetypes.NewObjectNull(t.AttrTypes), typ: t}
}

// Equal reports whether candidate is a ConfigObjectType with the same
// attributes and the same required set.
func (t ConfigObjectType) Equal(candidate attr.Type) bool {
	other, ok := candidate.(ConfigObjectType)
	if !ok {
		return false
	}
	if !t.ObjectType.Equal(other.ObjectType) {
		return false
	}
	if len(t.required) != len(other.required) {
		return false
	}
	for i := range t.required {
		if t.required[i] != other.required[i] {
			return false
		}
	}
	return true
}

func (t ConfigObjectType) String() string {
	return fmt.Sprintf("ConfigObjectType[%s; required: %s]", t.ObjectType.String(), strings.Join(t.required, ","))
}

func (t ConfigObjectType) ValueType(_ context.Context) attr.Value {
	return ConfigObjectValue{typ: t}
}

func (t ConfigObjectType) ValueFromObject(_ context.Context, obj basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	return ConfigObjectValue{ObjectValue: obj, typ: t}, nil
}
