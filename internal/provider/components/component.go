package components

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// LinkAttrTypes defines the attribute types for a component link object.
var LinkAttrTypes = map[string]attr.Type{
	"component_id": types.StringType,
	"settings":     types.MapType{ElemType: types.StringType},
}

// LinkObjectType is the ObjectType for a single component link.
var LinkObjectType = types.ObjectType{AttrTypes: LinkAttrTypes}

// ComponentAttrTypes defines the attribute types for a component object.
// This must exactly match the nested object schema in fractal_resource.go.
var ComponentAttrTypes = map[string]attr.Type{
	"id":                  types.StringType,
	"type":                types.StringType,
	"display_name":        types.StringType,
	"description":         types.StringType,
	"version":             types.StringType,
	"is_locked":           types.BoolType,
	"recreate_on_failure": types.BoolType,
	"parameters":          types.MapType{ElemType: types.StringType},
	"dependencies_ids":    types.ListType{ElemType: types.StringType},
	"links":               types.ListType{ElemType: LinkObjectType},
	"output_fields":       types.ListType{ElemType: types.StringType},
}

// ComponentObjectType is the ObjectType for a component reference in function parameters.
var ComponentObjectType = types.ObjectType{AttrTypes: ComponentAttrTypes}

// ComponentLink represents a resolved link ready to be set on the component.
type ComponentLink struct {
	ComponentId string
	Settings    map[string]string
}

// ComponentReturn returns the standard function return type for all component functions.
func ComponentReturn() function.ObjectReturn {
	return function.ObjectReturn{
		AttributeTypes: ComponentAttrTypes,
	}
}

// BuildComponent constructs a types.Object representing a blueprint component.
func BuildComponent(
	id string,
	componentType string,
	displayName types.String,
	description types.String,
	version types.String,
	parameters map[string]string,
	dependenciesIds []string,
	links []ComponentLink,
) (types.Object, *function.FuncError) {
	// Build parameters map value
	var parametersValue attr.Value
	if len(parameters) > 0 {
		elems := make(map[string]attr.Value, len(parameters))
		for k, v := range parameters {
			elems[k] = types.StringValue(v)
		}
		mv, diags := types.MapValue(types.StringType, elems)
		if diags.HasError() {
			return types.ObjectNull(ComponentAttrTypes), function.NewFuncError("failed to build parameters map")
		}
		parametersValue = mv
	} else {
		parametersValue = types.MapNull(types.StringType)
	}

	// Build dependencies list value
	var depsValue attr.Value
	if len(dependenciesIds) > 0 {
		elems := make([]attr.Value, len(dependenciesIds))
		for i, dep := range dependenciesIds {
			elems[i] = types.StringValue(dep)
		}
		lv, diags := types.ListValue(types.StringType, elems)
		if diags.HasError() {
			return types.ObjectNull(ComponentAttrTypes), function.NewFuncError("failed to build dependencies list")
		}
		depsValue = lv
	} else {
		depsValue = types.ListNull(types.StringType)
	}

	// Build links list value
	var linksValue attr.Value
	if len(links) > 0 {
		linkElems := make([]attr.Value, len(links))
		for i, link := range links {
			// Build settings map for this link
			var settingsValue attr.Value
			if len(link.Settings) > 0 {
				settingsElems := make(map[string]attr.Value, len(link.Settings))
				for k, v := range link.Settings {
					settingsElems[k] = types.StringValue(v)
				}
				sv, diags := types.MapValue(types.StringType, settingsElems)
				if diags.HasError() {
					return types.ObjectNull(ComponentAttrTypes), function.NewFuncError("failed to build link settings map")
				}
				settingsValue = sv
			} else {
				settingsValue = types.MapNull(types.StringType)
			}

			linkObj, diags := types.ObjectValue(LinkAttrTypes, map[string]attr.Value{
				"component_id": types.StringValue(link.ComponentId),
				"settings":     settingsValue,
			})
			if diags.HasError() {
				return types.ObjectNull(ComponentAttrTypes), function.NewFuncError("failed to build link object")
			}
			linkElems[i] = linkObj
		}
		lv, diags := types.ListValue(LinkObjectType, linkElems)
		if diags.HasError() {
			return types.ObjectNull(ComponentAttrTypes), function.NewFuncError("failed to build links list")
		}
		linksValue = lv
	} else {
		linksValue = types.ListNull(LinkObjectType)
	}

	// Default version to "v1" if not explicitly provided
	if version.IsNull() || version.IsUnknown() {
		version = types.StringValue("v1")
	}

	attrs := map[string]attr.Value{
		"id":                  types.StringValue(id),
		"type":                types.StringValue(componentType),
		"display_name":        displayName,
		"description":         description,
		"version":             version,
		"is_locked":           types.BoolNull(),
		"recreate_on_failure": types.BoolNull(),
		"parameters":          parametersValue,
		"dependencies_ids":    depsValue,
		"links":               linksValue,
		"output_fields":       types.ListNull(types.StringType),
	}

	obj, diags := types.ObjectValue(ComponentAttrTypes, attrs)
	if diags.HasError() {
		return types.ObjectNull(ComponentAttrTypes), function.NewFuncError("failed to build component object")
	}
	return obj, nil
}

// ExtractComponentId extracts the "id" field from a component object reference.
func ExtractComponentId(obj types.Object) (string, *function.FuncError) {
	if obj.IsNull() || obj.IsUnknown() {
		return "", function.NewFuncError("component reference is null")
	}
	attrs := obj.Attributes()
	idAttr, ok := attrs["id"]
	if !ok {
		return "", function.NewFuncError("component reference has no id field")
	}
	idStr, ok := idAttr.(types.String)
	if !ok || idStr.IsNull() || idStr.IsUnknown() {
		return "", function.NewFuncError("component reference id is not a valid string")
	}
	return idStr.ValueString(), nil
}

// ValidateComponentType checks that a component object has the expected type string.
func ValidateComponentType(obj types.Object, expectedType string) *function.FuncError {
	attrs := obj.Attributes()
	typeAttr, ok := attrs["type"]
	if !ok {
		return function.NewFuncError("component reference has no type field")
	}
	typeStr, ok := typeAttr.(types.String)
	if !ok || typeStr.IsNull() || typeStr.IsUnknown() {
		return function.NewFuncError("component reference type is not a valid string")
	}
	if typeStr.ValueString() != expectedType {
		return function.NewFuncError(fmt.Sprintf(
			"expected component of type %q but got %q",
			expectedType, typeStr.ValueString(),
		))
	}
	return nil
}

// ExtractDependency extracts a component's ID after validating its type.
// Returns the ID and any error. If the object is null/unknown, returns empty string and no error.
func ExtractDependency(obj types.Object, expectedType string) (string, *function.FuncError) {
	if obj.IsNull() || obj.IsUnknown() {
		return "", nil
	}
	if err := ValidateComponentType(obj, expectedType); err != nil {
		return "", err
	}
	return ExtractComponentId(obj)
}

// SgMembershipLinks converts a list of security group component objects to membership ComponentLinks.
func SgMembershipLinks(sgObjects []types.Object) ([]ComponentLink, *function.FuncError) {
	result := make([]ComponentLink, len(sgObjects))
	for i, sgObj := range sgObjects {
		id, err := ExtractDependency(sgObj, "NetworkAndCompute.IaaS.SecurityGroup")
		if err != nil {
			return nil, err
		}
		result[i] = ComponentLink{
			ComponentId: id,
			Settings:    nil,
		}
	}
	return result, nil
}

// OptionalString returns the types.String value if non-null, otherwise types.StringNull().
func OptionalString(v types.String) types.String {
	if v.IsNull() || v.IsUnknown() {
		return types.StringNull()
	}
	return v
}

// ParametersAttrType is the type of a function's `extra_parameters` attribute:
// extra parameters, for keys an offer reads that the function has no
// attribute for. Values follow the fc_fractal convention (a JSON object or
// array string is sent as JSON; secret_ref() output is one).
var ParametersAttrType = types.MapType{ElemType: types.StringType}

// WithExtraParameters returns the parameters a function derived from its
// attributes plus the caller's `extra_parameters`. A key set both ways is an error
// rather than a silent override, because the attribute usually validates or
// translates the value.
func WithExtraParameters(derived map[string]string, extra types.Map) (map[string]string, *function.FuncError) {
	if extra.IsNull() || extra.IsUnknown() || len(extra.Elements()) == 0 {
		return derived, nil
	}
	merged := make(map[string]string, len(derived)+len(extra.Elements()))
	for k, v := range derived {
		merged[k] = v
	}
	for k, v := range extra.Elements() {
		s, ok := v.(types.String)
		if !ok || s.IsNull() || s.IsUnknown() {
			return nil, function.NewFuncError(fmt.Sprintf("extra_parameters.%s must be a known string", k))
		}
		if _, taken := merged[k]; taken {
			return nil, function.NewFuncError(fmt.Sprintf(
				"extra_parameters.%s is already set by one of this function's attributes; set it there instead", k))
		}
		merged[k] = s.ValueString()
	}
	return merged, nil
}
