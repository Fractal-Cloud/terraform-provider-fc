package components

import (
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ObjectsFromDynamic returns the attributes of each object in a dynamic list
// attribute, for list-of-object attributes whose objects have optional
// fields (see LinksAttrType for why those are dynamic). what names the
// attribute in error messages.
func ObjectsFromDynamic(v types.Dynamic, what string) ([]map[string]attr.Value, *function.FuncError) {
	if v.IsNull() || v.IsUnknown() || v.IsUnderlyingValueNull() {
		return nil, nil
	}
	if v.IsUnderlyingValueUnknown() {
		return nil, function.NewFuncError(what + " must be known when the function is called")
	}

	var elements []attr.Value
	switch under := v.UnderlyingValue().(type) {
	case basetypes.TupleValue:
		elements = under.Elements()
	case basetypes.ListValue:
		elements = under.Elements()
	default:
		return nil, function.NewFuncError(fmt.Sprintf("%s must be a list of objects", what))
	}

	result := make([]map[string]attr.Value, len(elements))
	for i, element := range elements {
		obj, ok := element.(basetypes.ObjectValue)
		if !ok || obj.IsNull() {
			return nil, function.NewFuncError(fmt.Sprintf("%s[%d] must be an object", what, i))
		}
		result[i] = obj.Attributes()
	}
	return result, nil
}

// Int64Attr reads an optional whole-number attribute from an object decoded
// by ObjectsFromDynamic. It returns ok=false when the attribute is absent or
// null.
func Int64Attr(attrs map[string]attr.Value, name string) (value int64, ok bool, err *function.FuncError) {
	v, present := attrs[name]
	if !present || v.IsNull() {
		return 0, false, nil
	}
	if v.IsUnknown() {
		return 0, false, function.NewFuncError(name + " must be known when the function is called")
	}
	switch n := v.(type) {
	case basetypes.Int64Value:
		return n.ValueInt64(), true, nil
	case basetypes.NumberValue:
		f := n.ValueBigFloat()
		if f == nil {
			return 0, false, function.NewFuncError(name + " must be a number")
		}
		i, accuracy := f.Int64()
		if accuracy != big.Exact {
			return 0, false, function.NewFuncError(name + " must be a whole number")
		}
		return i, true, nil
	default:
		return 0, false, function.NewFuncError(name + " must be a number")
	}
}

// StringAttr reads an optional string attribute from an object decoded by
// ObjectsFromDynamic. It returns ok=false when the attribute is absent or
// null.
func StringAttr(attrs map[string]attr.Value, name string) (value string, ok bool, err *function.FuncError) {
	v, present := attrs[name]
	if !present || v.IsNull() {
		return "", false, nil
	}
	if v.IsUnknown() {
		return "", false, function.NewFuncError(name + " must be known when the function is called")
	}
	s, isString := v.(basetypes.StringValue)
	if !isString {
		return "", false, function.NewFuncError(name + " must be a string")
	}
	return s.ValueString(), true, nil
}
