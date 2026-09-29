package components

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.ObjectValuable = ConfigObjectValue{}

// ConfigObjectValue is the value of a ConfigObjectType. It behaves exactly as
// an ObjectValue; it exists so that the value reports its ConfigObjectType,
// which lists and nested objects check element types against.
type ConfigObjectValue struct {
	basetypes.ObjectValue

	typ ConfigObjectType
}

func (v ConfigObjectValue) Type(_ context.Context) attr.Type {
	return v.typ
}

func (v ConfigObjectValue) Equal(other attr.Value) bool {
	o, ok := other.(ConfigObjectValue)
	if !ok {
		return false
	}
	return v.typ.Equal(o.typ) && v.ObjectValue.Equal(o.ObjectValue)
}

func (v ConfigObjectValue) ToObjectValue(_ context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	return v.ObjectValue, nil
}
