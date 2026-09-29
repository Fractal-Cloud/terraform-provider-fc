package iaas

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
	"fractal.cloud/terraform-provider-fc/internal/provider/functions/functiontest"
)

// buildTestComponent is a helper that builds a minimal component object with just id and type.

// testLinkAttrTypes is the shape of one element of a function's links list.
var testLinkAttrTypes = map[string]attr.Type{
	"target":   components.ComponentObjectType,
	"settings": types.MapType{ElemType: types.StringType},
}

func buildTestComponent(t *testing.T, id, componentType string) types.Object {
	t.Helper()
	obj, err := components.BuildComponent(id, componentType, types.StringNull(), types.StringNull(), types.StringNull(), nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to build test component: %s", err.Text)
	}
	return obj
}

// runFunction is a helper that calls a function's Run method with the given arguments and returns the response.
func runFunction(t *testing.T, f function.Function, args []attr.Value) *function.RunResponse {
	t.Helper()
	ctx := context.Background()
	req := function.RunRequest{
		Arguments: function.NewArgumentsData(functiontest.Complete(t, f, args)),
	}
	resp := &function.RunResponse{
		Result: function.NewResultData(types.ObjectNull(components.ComponentAttrTypes)),
	}
	f.Run(ctx, req, resp)
	return resp
}

// getResultAttrs extracts the result object's attributes from a successful response.
func getResultAttrs(t *testing.T, resp *function.RunResponse) map[string]attr.Value {
	t.Helper()
	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error.Text)
	}
	result := resp.Result.Value()
	obj, ok := result.(types.Object)
	if !ok {
		t.Fatalf("expected types.Object result, got %T", result)
	}
	return obj.Attributes()
}

// --- VirtualNetwork ---

func TestVirtualNetworkFunction_Metadata(t *testing.T) {
	f := NewVirtualNetworkFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "network_and_compute_iaas_virtual_network" {
		t.Errorf("expected name %q, got %q", "network_and_compute_iaas_virtual_network", resp.Name)
	}
}

func TestVirtualNetworkFunction_Definition(t *testing.T) {
	f := NewVirtualNetworkFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

// --- Subnet ---

func TestSubnetFunction_Metadata(t *testing.T) {
	f := NewSubnetFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "network_and_compute_iaas_subnet" {
		t.Errorf("expected name %q, got %q", "network_and_compute_iaas_subnet", resp.Name)
	}
}

func TestSubnetFunction_Definition(t *testing.T) {
	f := NewSubnetFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

// --- SecurityGroup ---

func TestSecurityGroupFunction_Metadata(t *testing.T) {
	f := NewSecurityGroupFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "network_and_compute_iaas_security_group" {
		t.Errorf("expected name %q, got %q", "network_and_compute_iaas_security_group", resp.Name)
	}
}

func TestSecurityGroupFunction_Definition(t *testing.T) {
	f := NewSecurityGroupFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

// --- VirtualMachine ---

func TestVirtualMachineFunction_Metadata(t *testing.T) {
	f := NewVirtualMachineFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "network_and_compute_iaas_virtual_machine" {
		t.Errorf("expected name %q, got %q", "network_and_compute_iaas_virtual_machine", resp.Name)
	}
}

func TestVirtualMachineFunction_Definition(t *testing.T) {
	f := NewVirtualMachineFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

// --- LoadBalancer ---

func TestLoadBalancerFunction_Metadata(t *testing.T) {
	f := NewLoadBalancerFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "network_and_compute_iaas_load_balancer" {
		t.Errorf("expected name %q, got %q", "network_and_compute_iaas_load_balancer", resp.Name)
	}
}

func TestLoadBalancerFunction_Definition(t *testing.T) {
	f := NewLoadBalancerFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

// --- Run() tests ---

func TestVirtualNetworkFunction_Run_Minimal(t *testing.T) {
	f := NewVirtualNetworkFunction()
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":           types.StringType,
		"display_name": types.StringType,
		"description":  types.StringType,
		"cidr_block":   types.StringType,
		"links":        types.DynamicType,
	}, map[string]attr.Value{
		"id":           types.StringValue("my-vpc"),
		"display_name": types.StringNull(),
		"description":  types.StringNull(),
		"cidr_block":   types.StringNull(),
		"links":        types.DynamicNull(),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if id := attrs["id"].(types.String); id.ValueString() != "my-vpc" {
		t.Errorf("expected id %q, got %q", "my-vpc", id.ValueString())
	}
	if ct := attrs["type"].(types.String); ct.ValueString() != "NetworkAndCompute.IaaS.VirtualNetwork" {
		t.Errorf("expected type %q, got %q", "NetworkAndCompute.IaaS.VirtualNetwork", ct.ValueString())
	}
	if params := attrs["parameters"].(types.Map); !params.IsNull() {
		t.Error("expected null parameters when no cidr_block set")
	}
}

func TestVirtualNetworkFunction_Run_WithParams(t *testing.T) {
	f := NewVirtualNetworkFunction()
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":           types.StringType,
		"display_name": types.StringType,
		"description":  types.StringType,
		"cidr_block":   types.StringType,
		"links":        types.DynamicType,
	}, map[string]attr.Value{
		"id":           types.StringValue("vpc-1"),
		"display_name": types.StringValue("My VPC"),
		"description":  types.StringValue("Test VPC"),
		"cidr_block":   types.StringValue("10.0.0.0/16"),
		"links":        types.DynamicNull(),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if dn := attrs["display_name"].(types.String); dn.ValueString() != "My VPC" {
		t.Errorf("expected display_name %q, got %q", "My VPC", dn.ValueString())
	}
	if desc := attrs["description"].(types.String); desc.ValueString() != "Test VPC" {
		t.Errorf("expected description %q, got %q", "Test VPC", desc.ValueString())
	}
	params := attrs["parameters"].(types.Map)
	if params.IsNull() {
		t.Fatal("expected non-null parameters")
	}
	elems := params.Elements()
	if cidr := elems["cidrBlock"].(types.String); cidr.ValueString() != "10.0.0.0/16" {
		t.Errorf("expected cidrBlock %q, got %q", "10.0.0.0/16", cidr.ValueString())
	}
}

func TestSubnetFunction_Run_Minimal(t *testing.T) {
	f := NewSubnetFunction()
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":           types.StringType,
		"display_name": types.StringType,
		"description":  types.StringType,
		"cidr_block":   types.StringType,
		"vpc":          components.ComponentObjectType,
	}, map[string]attr.Value{
		"id":           types.StringValue("subnet-1"),
		"display_name": types.StringNull(),
		"description":  types.StringNull(),
		"cidr_block":   types.StringNull(),
		"vpc":          types.ObjectNull(components.ComponentAttrTypes),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if id := attrs["id"].(types.String); id.ValueString() != "subnet-1" {
		t.Errorf("expected id %q, got %q", "subnet-1", id.ValueString())
	}
	if ct := attrs["type"].(types.String); ct.ValueString() != "NetworkAndCompute.IaaS.Subnet" {
		t.Errorf("expected type %q, got %q", "NetworkAndCompute.IaaS.Subnet", ct.ValueString())
	}
	if deps := attrs["dependencies_ids"].(types.List); !deps.IsNull() {
		t.Error("expected null dependencies when no vpc set")
	}
}

func TestSubnetFunction_Run_WithVpc(t *testing.T) {
	f := NewSubnetFunction()
	vpc := buildTestComponent(t, "my-vpc", "NetworkAndCompute.IaaS.VirtualNetwork")
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":           types.StringType,
		"display_name": types.StringType,
		"description":  types.StringType,
		"cidr_block":   types.StringType,
		"vpc":          components.ComponentObjectType,
	}, map[string]attr.Value{
		"id":           types.StringValue("subnet-1"),
		"display_name": types.StringNull(),
		"description":  types.StringNull(),
		"cidr_block":   types.StringValue("10.0.1.0/24"),
		"vpc":          vpc,
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	// Check dependencies include vpc
	deps := attrs["dependencies_ids"].(types.List)
	if deps.IsNull() {
		t.Fatal("expected non-null dependencies")
	}
	depElems := deps.Elements()
	if len(depElems) != 1 {
		t.Fatalf("expected 1 dependency, got %d", len(depElems))
	}
	if depElems[0].(types.String).ValueString() != "my-vpc" {
		t.Errorf("expected dependency %q, got %q", "my-vpc", depElems[0].(types.String).ValueString())
	}

	// Check params
	params := attrs["parameters"].(types.Map)
	elems := params.Elements()
	if elems["cidrBlock"].(types.String).ValueString() != "10.0.1.0/24" {
		t.Errorf("expected cidrBlock %q, got %q", "10.0.1.0/24", elems["cidrBlock"].(types.String).ValueString())
	}
}

func TestSubnetFunction_Run_WrongVpcType(t *testing.T) {
	f := NewSubnetFunction()
	wrongVpc := buildTestComponent(t, "sg-1", "NetworkAndCompute.IaaS.SecurityGroup")
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":           types.StringType,
		"display_name": types.StringType,
		"description":  types.StringType,
		"cidr_block":   types.StringType,
		"vpc":          components.ComponentObjectType,
	}, map[string]attr.Value{
		"id":           types.StringValue("subnet-1"),
		"display_name": types.StringNull(),
		"description":  types.StringNull(),
		"cidr_block":   types.StringNull(),
		"vpc":          wrongVpc,
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	if resp.Error == nil {
		t.Fatal("expected error for wrong vpc type, got nil")
	}
}

func TestSecurityGroupFunction_Run_Minimal(t *testing.T) {
	f := NewSecurityGroupFunction()
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":           types.StringType,
		"display_name": types.StringType,
		"description":  types.StringType,
		"vpc":          components.ComponentObjectType,
		"ingress_rules": types.ListType{
			ElemType: types.ObjectType{AttrTypes: ingressRuleAttrTypes},
		},
	}, map[string]attr.Value{
		"id":           types.StringValue("sg-1"),
		"display_name": types.StringNull(),
		"description":  types.StringNull(),
		"vpc":          types.ObjectNull(components.ComponentAttrTypes),
		"ingress_rules": types.ListNull(
			types.ObjectType{AttrTypes: ingressRuleAttrTypes},
		),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if id := attrs["id"].(types.String); id.ValueString() != "sg-1" {
		t.Errorf("expected id %q, got %q", "sg-1", id.ValueString())
	}
	if ct := attrs["type"].(types.String); ct.ValueString() != "NetworkAndCompute.IaaS.SecurityGroup" {
		t.Errorf("expected type %q, got %q", "NetworkAndCompute.IaaS.SecurityGroup", ct.ValueString())
	}
}

func TestSecurityGroupFunction_Run_WithIngressRules(t *testing.T) {
	f := NewSecurityGroupFunction()
	vpc := buildTestComponent(t, "my-vpc", "NetworkAndCompute.IaaS.VirtualNetwork")

	rule1, diags := types.ObjectValue(ingressRuleAttrTypes, map[string]attr.Value{
		"from_port":           types.Int64Value(80),
		"to_port":             types.Int64Value(443),
		"protocol":            types.StringValue("tcp"),
		"source_cidr":         types.StringValue("10.0.0.0/8"),
		"source_component_id": types.StringNull(),
	})
	if diags.HasError() {
		t.Fatalf("failed to build rule: %s", diags.Errors())
	}

	ruleList, diags := types.ListValue(types.ObjectType{AttrTypes: ingressRuleAttrTypes}, []attr.Value{rule1})
	if diags.HasError() {
		t.Fatalf("failed to build rule list: %s", diags.Errors())
	}

	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":           types.StringType,
		"display_name": types.StringType,
		"description":  types.StringType,
		"vpc":          components.ComponentObjectType,
		"ingress_rules": types.ListType{
			ElemType: types.ObjectType{AttrTypes: ingressRuleAttrTypes},
		},
	}, map[string]attr.Value{
		"id":            types.StringValue("sg-1"),
		"display_name":  types.StringNull(),
		"description":   types.StringValue("My SG"),
		"vpc":           vpc,
		"ingress_rules": ruleList,
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	// Check vpc dependency
	deps := attrs["dependencies_ids"].(types.List)
	if deps.IsNull() {
		t.Fatal("expected non-null dependencies")
	}
	depElems := deps.Elements()
	if len(depElems) != 1 || depElems[0].(types.String).ValueString() != "my-vpc" {
		t.Errorf("expected dependency [my-vpc], got %v", depElems)
	}

	// Check ingress rules serialized in parameters
	params := attrs["parameters"].(types.Map)
	elems := params.Elements()
	ingressRulesJSON := elems["ingressRules"].(types.String).ValueString()
	var rules []ingressRuleJSON
	if err := json.Unmarshal([]byte(ingressRulesJSON), &rules); err != nil {
		t.Fatalf("failed to parse ingressRules JSON: %s", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].FromPort != 80 {
		t.Errorf("expected fromPort 80, got %d", rules[0].FromPort)
	}
	if rules[0].ToPort != 443 {
		t.Errorf("expected toPort 443, got %d", rules[0].ToPort)
	}
	if rules[0].Protocol != "tcp" {
		t.Errorf("expected protocol %q, got %q", "tcp", rules[0].Protocol)
	}
	if rules[0].SourceCidr != "10.0.0.0/8" {
		t.Errorf("expected sourceCidr %q, got %q", "10.0.0.0/8", rules[0].SourceCidr)
	}

	// Check description param
	if elems["description"].(types.String).ValueString() != "My SG" {
		t.Errorf("expected description param %q, got %q", "My SG", elems["description"].(types.String).ValueString())
	}
}

func TestVirtualMachineFunction_Run_Minimal(t *testing.T) {
	f := NewVirtualMachineFunction()
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":              types.StringType,
		"display_name":    types.StringType,
		"description":     types.StringType,
		"subnet":          components.ComponentObjectType,
		"links":           types.DynamicType,
		"security_groups": types.ListType{ElemType: components.ComponentObjectType},
	}, map[string]attr.Value{
		"id":              types.StringValue("vm-1"),
		"display_name":    types.StringNull(),
		"description":     types.StringNull(),
		"subnet":          types.ObjectNull(components.ComponentAttrTypes),
		"links":           types.DynamicNull(),
		"security_groups": types.ListNull(components.ComponentObjectType),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if id := attrs["id"].(types.String); id.ValueString() != "vm-1" {
		t.Errorf("expected id %q, got %q", "vm-1", id.ValueString())
	}
	if ct := attrs["type"].(types.String); ct.ValueString() != "NetworkAndCompute.IaaS.VirtualMachine" {
		t.Errorf("expected type %q, got %q", "NetworkAndCompute.IaaS.VirtualMachine", ct.ValueString())
	}
}

func TestVirtualMachineFunction_Run_WithDepsAndLinks(t *testing.T) {
	f := NewVirtualMachineFunction()
	subnet := buildTestComponent(t, "subnet-1", "NetworkAndCompute.IaaS.Subnet")
	sg := buildTestComponent(t, "sg-1", "NetworkAndCompute.IaaS.SecurityGroup")
	target := buildTestComponent(t, "vm-2", "NetworkAndCompute.IaaS.VirtualMachine")

	// Build a generic link
	linkSettings, _ := types.MapValue(types.StringType, map[string]attr.Value{
		"fromPort": types.StringValue("8080"),
	})
	genericLink, diags := types.ObjectValue(testLinkAttrTypes, map[string]attr.Value{
		"target":   target,
		"settings": linkSettings,
	})
	if diags.HasError() {
		t.Fatalf("failed to build generic link: %s", diags.Errors())
	}

	linkList, diags := types.ListValue(types.ObjectType{AttrTypes: testLinkAttrTypes}, []attr.Value{genericLink})
	if diags.HasError() {
		t.Fatalf("failed to build link list: %s", diags.Errors())
	}

	sgList, diags := types.ListValue(components.ComponentObjectType, []attr.Value{sg})
	if diags.HasError() {
		t.Fatalf("failed to build sg list: %s", diags.Errors())
	}

	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":              types.StringType,
		"display_name":    types.StringType,
		"description":     types.StringType,
		"subnet":          components.ComponentObjectType,
		"links":           types.DynamicType,
		"security_groups": types.ListType{ElemType: components.ComponentObjectType},
	}, map[string]attr.Value{
		"id":              types.StringValue("vm-1"),
		"display_name":    types.StringValue("My VM"),
		"description":     types.StringNull(),
		"subnet":          subnet,
		"links":           types.DynamicValue(linkList),
		"security_groups": sgList,
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	// Check subnet dependency
	deps := attrs["dependencies_ids"].(types.List)
	if deps.IsNull() {
		t.Fatal("expected non-null dependencies")
	}
	depElems := deps.Elements()
	if len(depElems) != 1 || depElems[0].(types.String).ValueString() != "subnet-1" {
		t.Errorf("expected dependency [subnet-1], got %v", depElems)
	}

	// Check links (port link + SG membership)
	linksVal := attrs["links"].(types.List)
	if linksVal.IsNull() {
		t.Fatal("expected non-null links")
	}
	linkElems := linksVal.Elements()
	if len(linkElems) != 2 {
		t.Fatalf("expected 2 links (1 port + 1 SG), got %d", len(linkElems))
	}

	// First link: port link to vm-2
	link0 := linkElems[0].(types.Object)
	if link0.Attributes()["component_id"].(types.String).ValueString() != "vm-2" {
		t.Errorf("expected first link target %q, got %q", "vm-2", link0.Attributes()["component_id"].(types.String).ValueString())
	}
	settings0 := link0.Attributes()["settings"].(types.Map)
	if settings0.Elements()["fromPort"].(types.String).ValueString() != "8080" {
		t.Errorf("expected fromPort %q, got %q", "8080", settings0.Elements()["fromPort"].(types.String).ValueString())
	}

	// Second link: SG membership
	link1 := linkElems[1].(types.Object)
	if link1.Attributes()["component_id"].(types.String).ValueString() != "sg-1" {
		t.Errorf("expected second link target %q, got %q", "sg-1", link1.Attributes()["component_id"].(types.String).ValueString())
	}
}

func TestLoadBalancerFunction_Run_Minimal(t *testing.T) {
	f := NewLoadBalancerFunction()
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":              types.StringType,
		"display_name":    types.StringType,
		"description":     types.StringType,
		"links":           types.DynamicType,
		"security_groups": types.ListType{ElemType: components.ComponentObjectType},
	}, map[string]attr.Value{
		"id":              types.StringValue("lb-1"),
		"display_name":    types.StringNull(),
		"description":     types.StringNull(),
		"links":           types.DynamicNull(),
		"security_groups": types.ListNull(components.ComponentObjectType),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	if id := attrs["id"].(types.String); id.ValueString() != "lb-1" {
		t.Errorf("expected id %q, got %q", "lb-1", id.ValueString())
	}
	if ct := attrs["type"].(types.String); ct.ValueString() != "NetworkAndCompute.IaaS.LoadBalancer" {
		t.Errorf("expected type %q, got %q", "NetworkAndCompute.IaaS.LoadBalancer", ct.ValueString())
	}
}

func TestLoadBalancerFunction_Run_WithLinks(t *testing.T) {
	f := NewLoadBalancerFunction()
	target := buildTestComponent(t, "backend-1", "CustomWorkloads.CaaS.Workload")
	sg := buildTestComponent(t, "sg-1", "NetworkAndCompute.IaaS.SecurityGroup")

	lbLinkSettings, _ := types.MapValue(types.StringType, map[string]attr.Value{
		"fromPort": types.StringValue("80"),
		"toPort":   types.StringValue("8080"),
		"protocol": types.StringValue("tcp"),
	})
	genericLink, diags := types.ObjectValue(testLinkAttrTypes, map[string]attr.Value{
		"target":   target,
		"settings": lbLinkSettings,
	})
	if diags.HasError() {
		t.Fatalf("failed to build generic link: %s", diags.Errors())
	}

	linkList, diags := types.ListValue(types.ObjectType{AttrTypes: testLinkAttrTypes}, []attr.Value{genericLink})
	if diags.HasError() {
		t.Fatalf("failed to build link list: %s", diags.Errors())
	}

	sgList, diags := types.ListValue(components.ComponentObjectType, []attr.Value{sg})
	if diags.HasError() {
		t.Fatalf("failed to build sg list: %s", diags.Errors())
	}

	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id":              types.StringType,
		"display_name":    types.StringType,
		"description":     types.StringType,
		"links":           types.DynamicType,
		"security_groups": types.ListType{ElemType: components.ComponentObjectType},
	}, map[string]attr.Value{
		"id":              types.StringValue("lb-1"),
		"display_name":    types.StringValue("My LB"),
		"description":     types.StringNull(),
		"links":           types.DynamicValue(linkList),
		"security_groups": sgList,
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)

	linksVal := attrs["links"].(types.List)
	if linksVal.IsNull() {
		t.Fatal("expected non-null links")
	}
	linkElems := linksVal.Elements()
	if len(linkElems) != 2 {
		t.Fatalf("expected 2 links, got %d", len(linkElems))
	}

	// Port link
	link0 := linkElems[0].(types.Object)
	if link0.Attributes()["component_id"].(types.String).ValueString() != "backend-1" {
		t.Errorf("expected port link target %q", "backend-1")
	}
	settings := link0.Attributes()["settings"].(types.Map).Elements()
	if settings["fromPort"].(types.String).ValueString() != "80" {
		t.Errorf("expected fromPort %q", "80")
	}
	if settings["toPort"].(types.String).ValueString() != "8080" {
		t.Errorf("expected toPort %q", "8080")
	}

	// SG membership link
	link1 := linkElems[1].(types.Object)
	if link1.Attributes()["component_id"].(types.String).ValueString() != "sg-1" {
		t.Errorf("expected SG link target %q", "sg-1")
	}
}
