package iaas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &SecurityGroupFunction{}

type SecurityGroupFunction struct{}

func NewSecurityGroupFunction() function.Function {
	return &SecurityGroupFunction{}
}

func (f *SecurityGroupFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "network_and_compute_iaas_security_group"
}

func (f *SecurityGroupFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a SecurityGroup blueprint component",
		Description: "Builds a SecurityGroup component with the correct type and parameters for use in a fractal's components list. " +
			"If vpc is provided, it is validated as a VirtualNetwork and added as a dependency; without it the agent places the " +
			"group in the environment's network. ingress_rules admit traffic from CIDR ranges. Traffic between components is " +
			"declared with links between them, not with rules here.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "SecurityGroup configuration",
				CustomType: components.NewConfigObjectType(map[string]attr.Type{
					"id":               types.StringType,
					"display_name":     types.StringType,
					"description":      types.StringType,
					"vpc":              components.ComponentObjectType,
					"ingress_rules":    types.DynamicType,
					"extra_parameters": components.ParametersAttrType,
				}, "id"),
			},
		},
		Return: components.ComponentReturn(),
	}
}

type securityGroupConfig struct {
	Id              types.String  `tfsdk:"id"`
	DisplayName     types.String  `tfsdk:"display_name"`
	Description     types.String  `tfsdk:"description"`
	Vpc             types.Object  `tfsdk:"vpc"`
	IngressRules    types.Dynamic `tfsdk:"ingress_rules"`
	ExtraParameters types.Map     `tfsdk:"extra_parameters"`
}

// ingressRule is one entry of the ingressRules parameter, in the shape the
// agents' NetworkSecurityRule parser reads.
type ingressRule struct {
	Protocol   string `json:"protocol"`
	FromPort   int64  `json:"fromPort"`
	ToPort     int64  `json:"toPort"`
	SourceCidr string `json:"sourceCidr"`
}

// ingressRules decodes the ingress_rules attribute. from_port and
// source_cidr are required; to_port defaults to from_port and protocol to tcp.
func ingressRules(v types.Dynamic) ([]ingressRule, *function.FuncError) {
	objects, err := components.ObjectsFromDynamic(v, "ingress_rules")
	if err != nil {
		return nil, err
	}

	rules := make([]ingressRule, len(objects))
	for i, attrs := range objects {
		for name := range attrs {
			switch name {
			case "from_port", "to_port", "protocol", "source_cidr":
			case "source_component_id":
				return nil, function.NewFuncError(fmt.Sprintf(
					"ingress_rules[%d]: source_component_id is no longer supported; admit traffic from a component "+
						"by linking that component to the target with settings = { fromPort = ... }", i))
			default:
				return nil, function.NewFuncError(fmt.Sprintf("ingress_rules[%d]: unexpected attribute %q", i, name))
			}
		}

		fromPort, ok, err := components.Int64Attr(attrs, "from_port")
		if err != nil {
			return nil, function.NewFuncError(fmt.Sprintf("ingress_rules[%d]: %s", i, err.Text))
		}
		if !ok {
			return nil, function.NewFuncError(fmt.Sprintf("ingress_rules[%d]: from_port is required", i))
		}
		toPort, ok, err := components.Int64Attr(attrs, "to_port")
		if err != nil {
			return nil, function.NewFuncError(fmt.Sprintf("ingress_rules[%d]: %s", i, err.Text))
		}
		if !ok {
			toPort = fromPort
		}
		protocol, ok, err := components.StringAttr(attrs, "protocol")
		if err != nil {
			return nil, function.NewFuncError(fmt.Sprintf("ingress_rules[%d]: %s", i, err.Text))
		}
		if !ok {
			protocol = "tcp"
		}
		sourceCidr, ok, err := components.StringAttr(attrs, "source_cidr")
		if err != nil {
			return nil, function.NewFuncError(fmt.Sprintf("ingress_rules[%d]: %s", i, err.Text))
		}
		if !ok {
			return nil, function.NewFuncError(fmt.Sprintf("ingress_rules[%d]: source_cidr is required", i))
		}

		rules[i] = ingressRule{Protocol: protocol, FromPort: fromPort, ToPort: toPort, SourceCidr: sourceCidr}
	}
	return rules, nil
}

func (f *SecurityGroupFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config securityGroupConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}

	// The agents read the group's own description from its parameters.
	if !config.Description.IsNull() && !config.Description.IsUnknown() {
		params["description"] = config.Description.ValueString()
	}

	rules, funcErr := ingressRules(config.IngressRules)
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}
	if len(rules) > 0 {
		b, err := json.Marshal(rules)
		if err != nil {
			resp.Error = function.NewFuncError(fmt.Sprintf("failed to serialize ingress_rules: %s", err))
			return
		}
		params["ingressRules"] = string(b)
	}

	var deps []string
	vpcId, funcErr := components.ExtractDependency(config.Vpc, "NetworkAndCompute.IaaS.VirtualNetwork")
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}
	if vpcId != "" {
		deps = append(deps, vpcId)
	}

	parameters, funcErr := components.WithExtraParameters(params, config.ExtraParameters)
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}

	result, funcErr := components.BuildComponent(
		config.Id.ValueString(),
		"NetworkAndCompute.IaaS.SecurityGroup",
		components.OptionalString(config.DisplayName),
		components.OptionalString(config.Description),
		types.StringNull(),
		parameters,
		deps,
		nil,
	)
	resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
