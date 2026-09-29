package components

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// LinksAttrType is the type of a function's `links` attribute: a list of
// `{ target = <component>, settings = { ... } }` objects in which `settings`
// may be omitted.
//
// It is dynamic because a typed list cannot express that. Terraform only
// omits an attribute that the type constraint declares optional, and
// terraform-plugin-go panics decoding a list whose element type declares
// optional attributes. A dynamic value accepts any list or tuple, and
// LinksFromDynamic enforces the shape instead.
var LinksAttrType = types.DynamicType

// LinksFromDynamic decodes a `links` attribute into ComponentLinks. A scalar
// setting value (`fromPort = 8080`, `protocol = "tcp"`) is kept as its string
// form. A list or object value (`redirectUris = ["https://..."]`,
// `scope = { path = "/", recursive = true }`) is kept as JSON, which the
// client sends as the JSON array or object the agent expects.
func LinksFromDynamic(v types.Dynamic) ([]ComponentLink, *function.FuncError) {
	if v.IsNull() || v.IsUnknown() || v.IsUnderlyingValueNull() {
		return nil, nil
	}
	if v.IsUnderlyingValueUnknown() {
		return nil, function.NewFuncError("links must be known when the function is called")
	}

	var elements []attr.Value
	switch under := v.UnderlyingValue().(type) {
	case basetypes.TupleValue:
		elements = under.Elements()
	case basetypes.ListValue:
		elements = under.Elements()
	case basetypes.SetValue:
		elements = under.Elements()
	default:
		return nil, function.NewFuncError(fmt.Sprintf("links must be a list of objects, got %s", v.UnderlyingValue().Type(context.Background())))
	}

	result := make([]ComponentLink, 0, len(elements))
	for i, element := range elements {
		link, err := linkFromValue(element)
		if err != nil {
			return nil, function.NewFuncError(fmt.Sprintf("links[%d]: %s", i, err.Text))
		}
		result = append(result, link)
	}
	return result, nil
}

func linkFromValue(v attr.Value) (ComponentLink, *function.FuncError) {
	obj, ok := v.(basetypes.ObjectValue)
	if !ok || obj.IsNull() {
		return ComponentLink{}, function.NewFuncError("each link must be an object with a target and optional settings")
	}

	attrs := obj.Attributes()
	for name := range attrs {
		if name != "target" && name != "settings" {
			return ComponentLink{}, function.NewFuncError(fmt.Sprintf("unexpected attribute %q; a link has only target and settings", name))
		}
	}

	targetValue, ok := attrs["target"]
	if !ok {
		return ComponentLink{}, function.NewFuncError("target is required")
	}
	target, ok := targetValue.(basetypes.ObjectValue)
	if !ok {
		return ComponentLink{}, function.NewFuncError("target must be a component object")
	}
	targetId, err := ExtractComponentId(target)
	if err != nil {
		return ComponentLink{}, err
	}

	settings, err := linkSettings(attrs["settings"])
	if err != nil {
		return ComponentLink{}, err
	}
	return ComponentLink{ComponentId: targetId, Settings: settings}, nil
}

func linkSettings(v attr.Value) (map[string]string, *function.FuncError) {
	if v == nil || v.IsNull() {
		return nil, nil
	}
	if v.IsUnknown() {
		return nil, function.NewFuncError("settings must be known when the function is called")
	}

	var entries map[string]attr.Value
	switch s := v.(type) {
	case basetypes.ObjectValue:
		entries = s.Attributes()
	case basetypes.MapValue:
		entries = s.Elements()
	default:
		return nil, function.NewFuncError("settings must be a map of values")
	}

	keys := make([]string, 0, len(entries))
	for k := range entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	settings := make(map[string]string, len(entries))
	for _, k := range keys {
		s, err := settingString(entries[k])
		if err != nil {
			return nil, function.NewFuncError(fmt.Sprintf("settings.%s: %s", k, err.Text))
		}
		settings[k] = s
	}
	return settings, nil
}

func settingString(v attr.Value) (string, *function.FuncError) {
	if v.IsNull() || v.IsUnknown() {
		return "", function.NewFuncError("must be a known, non-null value")
	}
	switch s := v.(type) {
	case basetypes.StringValue:
		return s.ValueString(), nil
	case basetypes.NumberValue:
		return s.ValueBigFloat().Text('f', -1), nil
	case basetypes.BoolValue:
		return strconv.FormatBool(s.ValueBool()), nil
	}

	structured, err := jsonValue(v)
	if err != nil {
		return "", err
	}
	b, marshalErr := json.Marshal(structured)
	if marshalErr != nil {
		return "", function.NewFuncError(marshalErr.Error())
	}
	return string(b), nil
}

// jsonValue converts a Terraform value into the Go value it encodes as JSON.
func jsonValue(v attr.Value) (any, *function.FuncError) {
	if v.IsUnknown() {
		return nil, function.NewFuncError("must be known when the function is called")
	}
	if v.IsNull() {
		return nil, nil
	}
	switch s := v.(type) {
	case basetypes.StringValue:
		return s.ValueString(), nil
	case basetypes.NumberValue:
		return json.Number(s.ValueBigFloat().Text('g', -1)), nil
	case basetypes.BoolValue:
		return s.ValueBool(), nil
	case basetypes.ListValue:
		return jsonArray(s.Elements())
	case basetypes.TupleValue:
		return jsonArray(s.Elements())
	case basetypes.SetValue:
		return jsonArray(s.Elements())
	case basetypes.MapValue:
		return jsonObject(s.Elements())
	case basetypes.ObjectValue:
		return jsonObject(s.Attributes())
	default:
		return nil, function.NewFuncError(fmt.Sprintf("unsupported value of type %s", v.Type(context.Background())))
	}
}

func jsonArray(elements []attr.Value) (any, *function.FuncError) {
	out := make([]any, len(elements))
	for i, e := range elements {
		v, err := jsonValue(e)
		if err != nil {
			return nil, err
		}
		out[i] = v
	}
	return out, nil
}

func jsonObject(entries map[string]attr.Value) (any, *function.FuncError) {
	out := make(map[string]any, len(entries))
	for k, e := range entries {
		v, err := jsonValue(e)
		if err != nil {
			return nil, err
		}
		out[k] = v
	}
	return out, nil
}
