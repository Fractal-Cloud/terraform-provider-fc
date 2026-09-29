package components

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// buildTestComponent is a helper that builds a minimal component object with just id and type.
func buildTestComponent(t *testing.T, id, componentType string) types.Object {
	t.Helper()
	obj, err := BuildComponent(id, componentType, types.StringNull(), types.StringNull(), types.StringNull(), nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to build test component: %s", err.Text)
	}
	return obj
}

// --- BuildComponent tests ---

func TestBuildComponent_MinimalFields(t *testing.T) {
	obj, funcErr := BuildComponent(
		"my-id", "Foo.Bar.Baz",
		types.StringNull(), types.StringNull(), types.StringNull(),
		nil, nil, nil,
	)
	if funcErr != nil {
		t.Fatalf("unexpected error: %s", funcErr.Text)
	}
	if obj.IsNull() || obj.IsUnknown() {
		t.Fatal("expected non-null, non-unknown object")
	}

	attrs := obj.Attributes()

	// id
	idVal, ok := attrs["id"].(types.String)
	if !ok {
		t.Fatal("id is not types.String")
	}
	if idVal.ValueString() != "my-id" {
		t.Errorf("expected id %q, got %q", "my-id", idVal.ValueString())
	}

	// type
	typeVal, ok := attrs["type"].(types.String)
	if !ok {
		t.Fatal("type is not types.String")
	}
	if typeVal.ValueString() != "Foo.Bar.Baz" {
		t.Errorf("expected type %q, got %q", "Foo.Bar.Baz", typeVal.ValueString())
	}

	// display_name should be null
	dnVal, ok := attrs["display_name"].(types.String)
	if !ok {
		t.Fatal("display_name is not types.String")
	}
	if !dnVal.IsNull() {
		t.Errorf("expected display_name to be null, got %q", dnVal.ValueString())
	}

	// description should be null
	descVal, ok := attrs["description"].(types.String)
	if !ok {
		t.Fatal("description is not types.String")
	}
	if !descVal.IsNull() {
		t.Errorf("expected description to be null, got %q", descVal.ValueString())
	}

	// version should default to "v1"
	verVal, ok := attrs["version"].(types.String)
	if !ok {
		t.Fatal("version is not types.String")
	}
	if verVal.IsNull() || verVal.ValueString() != "v1" {
		t.Errorf("expected version %q, got %q", "v1", verVal.ValueString())
	}

	// parameters should be null map
	paramsVal, ok := attrs["parameters"].(types.Map)
	if !ok {
		t.Fatal("parameters is not types.Map")
	}
	if !paramsVal.IsNull() {
		t.Error("expected parameters to be null")
	}

	// dependencies_ids should be null list
	depsVal, ok := attrs["dependencies_ids"].(types.List)
	if !ok {
		t.Fatal("dependencies_ids is not types.List")
	}
	if !depsVal.IsNull() {
		t.Error("expected dependencies_ids to be null")
	}

	// links should be null list
	linksVal, ok := attrs["links"].(types.List)
	if !ok {
		t.Fatal("links is not types.List")
	}
	if !linksVal.IsNull() {
		t.Error("expected links to be null")
	}
}

func TestBuildComponent_AllFields(t *testing.T) {
	params := map[string]string{
		"cidr_block": "10.0.0.0/16",
		"region":     "us-east-1",
	}
	deps := []string{"dep-1", "dep-2"}
	links := []ComponentLink{
		{
			ComponentId: "link-target-1",
			Settings: map[string]string{
				"fromPort": "8080",
				"toPort":   "8080",
				"protocol": "tcp",
			},
		},
	}

	obj, funcErr := BuildComponent(
		"full-id", "NetworkAndCompute.IaaS.AwsVpc",
		types.StringValue("My VPC"),
		types.StringValue("A test VPC"),
		types.StringValue("1.0.0"),
		params, deps, links,
	)
	if funcErr != nil {
		t.Fatalf("unexpected error: %s", funcErr.Text)
	}

	attrs := obj.Attributes()

	// id
	idVal := attrs["id"].(types.String)
	if idVal.ValueString() != "full-id" {
		t.Errorf("expected id %q, got %q", "full-id", idVal.ValueString())
	}

	// type
	typeVal := attrs["type"].(types.String)
	if typeVal.ValueString() != "NetworkAndCompute.IaaS.AwsVpc" {
		t.Errorf("expected type %q, got %q", "NetworkAndCompute.IaaS.AwsVpc", typeVal.ValueString())
	}

	// display_name
	dnVal := attrs["display_name"].(types.String)
	if dnVal.IsNull() || dnVal.ValueString() != "My VPC" {
		t.Errorf("expected display_name %q, got %q", "My VPC", dnVal.ValueString())
	}

	// description
	descVal := attrs["description"].(types.String)
	if descVal.IsNull() || descVal.ValueString() != "A test VPC" {
		t.Errorf("expected description %q, got %q", "A test VPC", descVal.ValueString())
	}

	// version
	verVal := attrs["version"].(types.String)
	if verVal.IsNull() || verVal.ValueString() != "1.0.0" {
		t.Errorf("expected version %q, got %q", "1.0.0", verVal.ValueString())
	}

	// parameters
	paramsVal := attrs["parameters"].(types.Map)
	if paramsVal.IsNull() {
		t.Fatal("expected parameters to be non-null")
	}
	paramElems := paramsVal.Elements()
	if len(paramElems) != 2 {
		t.Fatalf("expected 2 parameters, got %d", len(paramElems))
	}
	cidr, ok := paramElems["cidr_block"].(types.String)
	if !ok || cidr.ValueString() != "10.0.0.0/16" {
		t.Errorf("expected cidr_block %q, got %v", "10.0.0.0/16", paramElems["cidr_block"])
	}
	region, ok := paramElems["region"].(types.String)
	if !ok || region.ValueString() != "us-east-1" {
		t.Errorf("expected region %q, got %v", "us-east-1", paramElems["region"])
	}

	// dependencies_ids
	depsVal := attrs["dependencies_ids"].(types.List)
	if depsVal.IsNull() {
		t.Fatal("expected dependencies_ids to be non-null")
	}
	depElems := depsVal.Elements()
	if len(depElems) != 2 {
		t.Fatalf("expected 2 dependencies, got %d", len(depElems))
	}
	dep0 := depElems[0].(types.String)
	if dep0.ValueString() != "dep-1" {
		t.Errorf("expected dep[0] %q, got %q", "dep-1", dep0.ValueString())
	}
	dep1 := depElems[1].(types.String)
	if dep1.ValueString() != "dep-2" {
		t.Errorf("expected dep[1] %q, got %q", "dep-2", dep1.ValueString())
	}

	// links
	linksVal := attrs["links"].(types.List)
	if linksVal.IsNull() {
		t.Fatal("expected links to be non-null")
	}
	linkElems := linksVal.Elements()
	if len(linkElems) != 1 {
		t.Fatalf("expected 1 link, got %d", len(linkElems))
	}
	linkObj := linkElems[0].(types.Object)
	linkAttrs := linkObj.Attributes()
	linkCid := linkAttrs["component_id"].(types.String)
	if linkCid.ValueString() != "link-target-1" {
		t.Errorf("expected link component_id %q, got %q", "link-target-1", linkCid.ValueString())
	}
	linkSettings := linkAttrs["settings"].(types.Map)
	if linkSettings.IsNull() {
		t.Fatal("expected link settings to be non-null")
	}
	settingsElems := linkSettings.Elements()
	if settingsElems["fromPort"].(types.String).ValueString() != "8080" {
		t.Errorf("expected fromPort %q, got %q", "8080", settingsElems["fromPort"].(types.String).ValueString())
	}
}

func TestBuildComponent_EmptyParams(t *testing.T) {
	obj, funcErr := BuildComponent(
		"id", "Type.X.Y",
		types.StringNull(), types.StringNull(), types.StringNull(),
		map[string]string{}, nil, nil,
	)
	if funcErr != nil {
		t.Fatalf("unexpected error: %s", funcErr.Text)
	}

	paramsVal := obj.Attributes()["parameters"].(types.Map)
	if !paramsVal.IsNull() {
		t.Error("expected empty params map to produce null parameters value")
	}
}

func TestBuildComponent_LinksWithAndWithoutSettings(t *testing.T) {
	links := []ComponentLink{
		{
			ComponentId: "with-settings",
			Settings: map[string]string{
				"fromPort": "443",
			},
		},
		{
			ComponentId: "without-settings",
			Settings:    nil,
		},
	}

	obj, funcErr := BuildComponent(
		"id", "Type.X.Y",
		types.StringNull(), types.StringNull(), types.StringNull(),
		nil, nil, links,
	)
	if funcErr != nil {
		t.Fatalf("unexpected error: %s", funcErr.Text)
	}

	linksVal := obj.Attributes()["links"].(types.List)
	if linksVal.IsNull() {
		t.Fatal("expected links to be non-null")
	}
	linkElems := linksVal.Elements()
	if len(linkElems) != 2 {
		t.Fatalf("expected 2 links, got %d", len(linkElems))
	}

	// First link: has settings
	link0 := linkElems[0].(types.Object)
	settings0 := link0.Attributes()["settings"].(types.Map)
	if settings0.IsNull() {
		t.Error("expected first link settings to be non-null")
	}
	elems0 := settings0.Elements()
	if elems0["fromPort"].(types.String).ValueString() != "443" {
		t.Errorf("expected fromPort %q, got %q", "443", elems0["fromPort"].(types.String).ValueString())
	}

	// Second link: no settings (null)
	link1 := linkElems[1].(types.Object)
	settings1 := link1.Attributes()["settings"].(types.Map)
	if !settings1.IsNull() {
		t.Error("expected second link settings to be null")
	}
}

// --- ExtractComponentId tests ---

func TestExtractComponentId_Valid(t *testing.T) {
	obj := buildTestComponent(t, "extract-me", "Foo.Bar.Baz")
	id, funcErr := ExtractComponentId(obj)
	if funcErr != nil {
		t.Fatalf("unexpected error: %s", funcErr.Text)
	}
	if id != "extract-me" {
		t.Errorf("expected id %q, got %q", "extract-me", id)
	}
}

func TestExtractComponentId_NullObject(t *testing.T) {
	obj := types.ObjectNull(ComponentAttrTypes)
	_, funcErr := ExtractComponentId(obj)
	if funcErr == nil {
		t.Fatal("expected error for null object, got nil")
	}
}

func TestExtractComponentId_UnknownObject(t *testing.T) {
	obj := types.ObjectUnknown(ComponentAttrTypes)
	_, funcErr := ExtractComponentId(obj)
	if funcErr == nil {
		t.Fatal("expected error for unknown object, got nil")
	}
}

// --- ValidateComponentType tests ---

func TestValidateComponentType_Match(t *testing.T) {
	obj := buildTestComponent(t, "id", "Foo.Bar.Baz")
	funcErr := ValidateComponentType(obj, "Foo.Bar.Baz")
	if funcErr != nil {
		t.Fatalf("unexpected error: %s", funcErr.Text)
	}
}

func TestValidateComponentType_Mismatch(t *testing.T) {
	obj := buildTestComponent(t, "id", "Foo.Bar.Baz")
	funcErr := ValidateComponentType(obj, "Wrong.Type")
	if funcErr == nil {
		t.Fatal("expected error for type mismatch, got nil")
	}
	if !strings.Contains(funcErr.Text, "Foo.Bar.Baz") {
		t.Errorf("expected error to contain actual type %q, got %q", "Foo.Bar.Baz", funcErr.Text)
	}
	if !strings.Contains(funcErr.Text, "Wrong.Type") {
		t.Errorf("expected error to contain expected type %q, got %q", "Wrong.Type", funcErr.Text)
	}
}

func TestValidateComponentType_NullType(t *testing.T) {
	// Build an object with a null type field manually
	attrs := map[string]attr.Value{
		"id":                  types.StringValue("id"),
		"type":                types.StringNull(),
		"display_name":        types.StringNull(),
		"description":         types.StringNull(),
		"version":             types.StringNull(),
		"is_locked":           types.BoolNull(),
		"recreate_on_failure": types.BoolNull(),
		"parameters":          types.MapNull(types.StringType),
		"dependencies_ids":    types.ListNull(types.StringType),
		"links":               types.ListNull(LinkObjectType),
		"output_fields":       types.ListNull(types.StringType),
	}
	obj, diags := types.ObjectValue(ComponentAttrTypes, attrs)
	if diags.HasError() {
		t.Fatalf("failed to build test object: %s", diags.Errors())
	}

	funcErr := ValidateComponentType(obj, "Foo.Bar.Baz")
	if funcErr == nil {
		t.Fatal("expected error for null type field, got nil")
	}
}

// --- ExtractDependency tests ---

func TestExtractDependency_Valid(t *testing.T) {
	obj := buildTestComponent(t, "dep-id", "NetworkAndCompute.IaaS.AwsVpc")
	id, funcErr := ExtractDependency(obj, "NetworkAndCompute.IaaS.AwsVpc")
	if funcErr != nil {
		t.Fatalf("unexpected error: %s", funcErr.Text)
	}
	if id != "dep-id" {
		t.Errorf("expected id %q, got %q", "dep-id", id)
	}
}

func TestExtractDependency_WrongType(t *testing.T) {
	obj := buildTestComponent(t, "dep-id", "NetworkAndCompute.IaaS.AwsVpc")
	_, funcErr := ExtractDependency(obj, "Storage.PaaS.AwsRds")
	if funcErr == nil {
		t.Fatal("expected error for wrong type, got nil")
	}
}

func TestExtractDependency_NullObject(t *testing.T) {
	obj := types.ObjectNull(ComponentAttrTypes)
	id, funcErr := ExtractDependency(obj, "NetworkAndCompute.IaaS.AwsVpc")
	if funcErr != nil {
		t.Fatalf("unexpected error: %s", funcErr.Text)
	}
	if id != "" {
		t.Errorf("expected empty string for null object, got %q", id)
	}
}

// An unknown dependency is an error, not an absent one: dropping it would
// silently lose the ordering constraint.
func TestExtractDependency_UnknownObject(t *testing.T) {
	obj := types.ObjectUnknown(ComponentAttrTypes)
	if _, funcErr := ExtractDependency(obj, "NetworkAndCompute.IaaS.AwsVpc"); funcErr == nil {
		t.Fatal("expected an error for an unknown dependency")
	}
}

// --- SgMembershipLinks tests ---

func TestSgMembershipLinks_Valid(t *testing.T) {
	sg1 := buildTestComponent(t, "sg-1", "NetworkAndCompute.IaaS.SecurityGroup")
	sg2 := buildTestComponent(t, "sg-2", "NetworkAndCompute.IaaS.SecurityGroup")

	result, funcErr := SgMembershipLinks([]types.Object{sg1, sg2})
	if funcErr != nil {
		t.Fatalf("unexpected error: %s", funcErr.Text)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 links, got %d", len(result))
	}

	if result[0].ComponentId != "sg-1" {
		t.Errorf("expected first link component id %q, got %q", "sg-1", result[0].ComponentId)
	}
	if result[0].Settings != nil {
		t.Errorf("expected first link settings to be nil, got %v", result[0].Settings)
	}

	if result[1].ComponentId != "sg-2" {
		t.Errorf("expected second link component id %q, got %q", "sg-2", result[1].ComponentId)
	}
	if result[1].Settings != nil {
		t.Errorf("expected second link settings to be nil, got %v", result[1].Settings)
	}
}

func TestSgMembershipLinks_WrongType(t *testing.T) {
	vnet := buildTestComponent(t, "vnet-1", "NetworkAndCompute.IaaS.VirtualNetwork")

	_, funcErr := SgMembershipLinks([]types.Object{vnet})
	if funcErr == nil {
		t.Fatal("expected error for wrong component type, got nil")
	}
}

func TestSgMembershipLinks_Empty(t *testing.T) {
	result, funcErr := SgMembershipLinks([]types.Object{})
	if funcErr != nil {
		t.Fatalf("unexpected error: %s", funcErr.Text)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 links, got %d", len(result))
	}
}

// --- OptionalString tests ---

func TestOptionalString_Value(t *testing.T) {
	v := types.StringValue("hello")
	result := OptionalString(v)
	if result.IsNull() {
		t.Fatal("expected non-null result")
	}
	if result.ValueString() != "hello" {
		t.Errorf("expected %q, got %q", "hello", result.ValueString())
	}
}

func TestOptionalString_Null(t *testing.T) {
	v := types.StringNull()
	result := OptionalString(v)
	if !result.IsNull() {
		t.Errorf("expected null result, got %q", result.ValueString())
	}
}

func TestOptionalString_Unknown(t *testing.T) {
	v := types.StringUnknown()
	result := OptionalString(v)
	if !result.IsNull() {
		t.Error("expected null result for unknown input")
	}
}

// --- ComponentReturn tests ---

func TestComponentReturn_HasCorrectAttrTypes(t *testing.T) {
	ret := ComponentReturn()
	if len(ret.AttributeTypes) != len(ComponentAttrTypes) {
		t.Fatalf("expected %d attribute types, got %d", len(ComponentAttrTypes), len(ret.AttributeTypes))
	}
	for key, expectedType := range ComponentAttrTypes {
		gotType, ok := ret.AttributeTypes[key]
		if !ok {
			t.Errorf("missing attribute type for key %q", key)
			continue
		}
		if !gotType.Equal(expectedType) {
			t.Errorf("attribute type mismatch for key %q: expected %v, got %v", key, expectedType, gotType)
		}
	}
}
