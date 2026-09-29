package spec

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &Function{}

// Function is a component function built from a Spec.
type Function struct {
	spec       Spec
	configType components.ConfigObjectType
}

// New returns the function s describes. It panics on an inconsistent Spec,
// which is a programming error caught by the function's tests.
func New(s Spec) function.Function {
	attrTypes := map[string]attr.Type{
		"id":               types.StringType,
		"display_name":     types.StringType,
		"description":      types.StringType,
		"extra_parameters": components.ParametersAttrType,
	}
	required := []string{"id"}
	for _, a := range s.Attributes {
		addAttr(attrTypes, a.Name, attrType(a))
		if a.Required {
			required = append(required, a.Name)
		}
	}
	for _, d := range s.Dependencies {
		addAttr(attrTypes, d.Name, components.ComponentObjectType)
	}
	if s.Links {
		addAttr(attrTypes, "links", components.LinksAttrType)
	}
	return &Function{spec: s, configType: components.NewConfigObjectType(attrTypes, required...)}
}

func addAttr(attrTypes map[string]attr.Type, name string, typ attr.Type) {
	if _, taken := attrTypes[name]; taken {
		panic(fmt.Sprintf("attribute %q declared twice", name))
	}
	attrTypes[name] = typ
}

func attrType(a Attribute) attr.Type {
	switch a.Kind {
	case String:
		return types.StringType
	case Int64:
		return types.Int64Type
	case Bool:
		return types.BoolType
	case StringList:
		return types.ListType{ElemType: types.StringType}
	case StringMap:
		return types.MapType{ElemType: types.StringType}
	case Object:
		// Dynamic rather than an object type: an object type with optional
		// fields cannot be nested in a function argument (see
		// components.ConfigObjectType), so the fields are checked in Run.
		for _, f := range a.Fields {
			if f.Kind == Object {
				panic("nested Object fields are not supported")
			}
		}
		return types.DynamicType
	default:
		panic(fmt.Sprintf("attribute %q has unknown kind %d", a.Name, a.Kind))
	}
}

func (f *Function) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = f.spec.Name
}

func (f *Function) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     f.spec.Summary,
		Description: f.spec.Description,
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: f.spec.Summary + " configuration",
				CustomType:  f.configType,
			},
		},
		Return: components.ComponentReturn(),
	}
}

func (f *Function) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config components.ConfigObjectValue
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}
	attrs := config.Attributes()
	for _, name := range []string{"id", "display_name", "description"} {
		if v, ok := attrs[name]; ok && v.IsUnknown() {
			resp.Error = function.NewArgumentFuncError(0, name+" must be known when the function is called")
			return
		}
	}

	params := map[string]string{}
	for _, a := range f.spec.Attributes {
		value, ok, err := parameterValue(a, attrs[a.Name])
		if err != nil {
			resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("%s: %s", a.Name, err))
			return
		}
		if !ok {
			continue
		}
		for _, key := range append([]string{a.Key}, a.Aliases...) {
			params[key] = value
		}
	}

	var deps []string
	for _, d := range f.spec.Dependencies {
		obj, _ := attrs[d.Name].(types.Object)
		id, funcErr := components.ExtractDependency(obj, d.ComponentType)
		if funcErr != nil {
			resp.Error = function.ConcatFuncErrors(resp.Error, function.NewArgumentFuncError(0, d.Name+": "+funcErr.Text))
			return
		}
		if id != "" {
			deps = append(deps, id)
		}
	}

	var links []components.ComponentLink
	if f.spec.Links {
		dyn, _ := attrs["links"].(types.Dynamic)
		resolved, funcErr := components.LinksFromDynamic(dyn)
		if funcErr != nil {
			resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
			return
		}
		links = resolved
	}

	extra, _ := attrs["extra_parameters"].(types.Map)
	parameters, funcErr := components.WithExtraParameters(params, extra)
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}

	result, funcErr := components.BuildComponent(
		stringValue(attrs["id"]),
		f.spec.ComponentType,
		optionalString(attrs["display_name"]),
		optionalString(attrs["description"]),
		types.StringNull(),
		parameters,
		deps,
		links,
	)
	resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}

// parameterValue renders an attribute's value as the parameter string. It
// returns ok=false for a null value, which sets no parameter.
func parameterValue(a Attribute, v attr.Value) (string, bool, error) {
	if v == nil || v.IsNull() {
		return "", false, nil
	}
	if v.IsUnknown() {
		return "", false, fmt.Errorf("must be known when the function is called")
	}

	switch a.Kind {
	case String:
		s := v.(types.String).ValueString()
		if len(a.OneOf) > 0 && !contains(a.OneOf, s) {
			return "", false, fmt.Errorf("must be one of %s, got %q", strings.Join(quoted(a.OneOf), ", "), s)
		}
		return s, true, nil
	case Int64:
		return strconv.FormatInt(v.(types.Int64).ValueInt64(), 10), true, nil
	case Bool:
		return strconv.FormatBool(v.(types.Bool).ValueBool()), true, nil
	default:
		structured, err := jsonValue(a, v)
		if err != nil {
			return "", false, err
		}
		b, err := json.Marshal(structured)
		if err != nil {
			return "", false, err
		}
		return string(b), true, nil
	}
}

func jsonValue(a Attribute, v attr.Value) (any, error) {
	switch a.Kind {
	case String:
		return v.(types.String).ValueString(), nil
	case Int64:
		return v.(types.Int64).ValueInt64(), nil
	case Bool:
		return v.(types.Bool).ValueBool(), nil
	case StringList:
		out := []string{}
		for _, e := range v.(types.List).Elements() {
			s, ok := e.(types.String)
			if !ok || s.IsNull() || s.IsUnknown() {
				return nil, fmt.Errorf("elements must be known strings")
			}
			out = append(out, s.ValueString())
		}
		return out, nil
	case StringMap:
		out := map[string]string{}
		for k, e := range v.(types.Map).Elements() {
			s, ok := e.(types.String)
			if !ok || s.IsNull() || s.IsUnknown() {
				return nil, fmt.Errorf("value of %q must be a known string", k)
			}
			out[k] = s.ValueString()
		}
		return out, nil
	case Object:
		return objectJSON(a, v)
	default:
		return nil, fmt.Errorf("unknown kind %d", a.Kind)
	}
}

// objectJSON renders an Object attribute, a dynamic value holding the object
// literal the caller wrote, as a JSON object keyed by the fields' parameter
// keys.
func objectJSON(a Attribute, v attr.Value) (any, error) {
	dyn, ok := v.(types.Dynamic)
	if !ok || dyn.IsUnderlyingValueNull() {
		return map[string]any{}, nil
	}
	if dyn.IsUnderlyingValueUnknown() {
		return nil, fmt.Errorf("must be known when the function is called")
	}

	var entries map[string]attr.Value
	switch under := dyn.UnderlyingValue().(type) {
	case basetypes.ObjectValue:
		entries = under.Attributes()
	case basetypes.MapValue:
		entries = under.Elements()
	default:
		return nil, fmt.Errorf("must be an object")
	}

	fields := make(map[string]Attribute, len(a.Fields))
	for _, f := range a.Fields {
		fields[f.Name] = f
	}
	out := map[string]any{}
	for name, fv := range entries {
		field, known := fields[name]
		if !known {
			return nil, fmt.Errorf("unexpected attribute %q", name)
		}
		if fv.IsNull() {
			continue
		}
		if fv.IsUnknown() {
			return nil, fmt.Errorf("%s must be known when the function is called", name)
		}
		j, err := fieldJSON(field, fv)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		out[field.Key] = j
	}
	return out, nil
}

// fieldJSON converts one field of an Object attribute. Terraform passes the
// literal through a dynamic value, so numbers arrive as Number, not Int64.
func fieldJSON(f Attribute, v attr.Value) (any, error) {
	switch f.Kind {
	case String:
		s, ok := v.(types.String)
		if !ok {
			return nil, fmt.Errorf("must be a string")
		}
		if len(f.OneOf) > 0 && !contains(f.OneOf, s.ValueString()) {
			return nil, fmt.Errorf("must be one of %s", strings.Join(quoted(f.OneOf), ", "))
		}
		return s.ValueString(), nil
	case Int64:
		switch n := v.(type) {
		case types.Int64:
			return n.ValueInt64(), nil
		case types.Number:
			i, accuracy := n.ValueBigFloat().Int64()
			if accuracy != big.Exact {
				return nil, fmt.Errorf("must be a whole number")
			}
			return i, nil
		}
		return nil, fmt.Errorf("must be a number")
	case Bool:
		b, ok := v.(types.Bool)
		if !ok {
			return nil, fmt.Errorf("must be a bool")
		}
		return b.ValueBool(), nil
	default:
		return nil, fmt.Errorf("unsupported field kind %d", f.Kind)
	}
}

func contains(values []string, v string) bool {
	for _, candidate := range values {
		if candidate == v {
			return true
		}
	}
	return false
}

func quoted(values []string) []string {
	out := make([]string, len(values))
	for i, v := range values {
		out[i] = strconv.Quote(v)
	}
	sort.Strings(out)
	return out
}

func stringValue(v attr.Value) string {
	s, _ := v.(types.String)
	return s.ValueString()
}

func optionalString(v attr.Value) types.String {
	s, ok := v.(types.String)
	if !ok {
		return types.StringNull()
	}
	return components.OptionalString(s)
}
