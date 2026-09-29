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

var ingressRuleAttrTypes = map[string]attr.Type{
	"from_port":           types.Int64Type,
	"to_port":             types.Int64Type,
	"protocol":            types.StringType,
	"source_cidr":         types.StringType,
	"source_component_id": types.StringType,
}

func (f *SecurityGroupFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a SecurityGroup blueprint component",
		Description: "Builds a SecurityGroup component with the correct type and parameters for use in a fractal's components list. If vpc is provided, it is validated as a VirtualNetwork and automatically added as a dependency. Ingress rules are serialized into the component parameters.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "SecurityGroup configuration",
				CustomType: components.NewConfigObjectType(map[string]attr.Type{
					"id":           types.StringType,
					"display_name": types.StringType,
					"description":  types.StringType,
					"vpc":          components.ComponentObjectType,
					"ingress_rules": types.ListType{
						ElemType: types.ObjectType{AttrTypes: ingressRuleAttrTypes},
					},
					"extra_parameters": components.ParametersAttrType,
				}, "id"),
			},
		},
		Return: components.ComponentReturn(),
	}
}

type securityGroupConfig struct {
	Id              types.String `tfsdk:"id"`
	DisplayName     types.String `tfsdk:"display_name"`
	Description     types.String `tfsdk:"description"`
	Vpc             types.Object `tfsdk:"vpc"`
	IngressRules    types.List   `tfsdk:"ingress_rules"`
	ExtraParameters types.Map    `tfsdk:"extra_parameters"`
}

type ingressRuleConfig struct {
	FromPort          types.Int64  `tfsdk:"from_port"`
	ToPort            types.Int64  `tfsdk:"to_port"`
	Protocol          types.String `tfsdk:"protocol"`
	SourceCidr        types.String `tfsdk:"source_cidr"`
	SourceComponentId types.String `tfsdk:"source_component_id"`
}

type ingressRuleJSON struct {
	Protocol          string `json:"protocol"`
	FromPort          int64  `json:"fromPort"`
	ToPort            int64  `json:"toPort"`
	SourceCidr        string `json:"sourceCidr,omitempty"`
	SourceComponentId string `json:"sourceComponentId,omitempty"`
}

func (f *SecurityGroupFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config securityGroupConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}

	// SecurityGroup description is also set as a parameter
	if !config.Description.IsNull() && !config.Description.IsUnknown() {
		params["description"] = config.Description.ValueString()
	}

	// Serialize ingress rules to JSON
	if !config.IngressRules.IsNull() && !config.IngressRules.IsUnknown() {
		var rules []ingressRuleConfig
		diags := config.IngressRules.ElementsAs(ctx, &rules, false)
		if diags.HasError() {
			resp.Error = function.NewFuncError("failed to parse ingress_rules")
			return
		}

		if len(rules) > 0 {
			jsonRules := make([]ingressRuleJSON, len(rules))
			for i, rule := range rules {
				protocol := "tcp"
				if !rule.Protocol.IsNull() && !rule.Protocol.IsUnknown() {
					protocol = rule.Protocol.ValueString()
				}

				toPort := rule.FromPort.ValueInt64()
				if !rule.ToPort.IsNull() && !rule.ToPort.IsUnknown() {
					toPort = rule.ToPort.ValueInt64()
				}

				jsonRules[i] = ingressRuleJSON{
					Protocol: protocol,
					FromPort: rule.FromPort.ValueInt64(),
					ToPort:   toPort,
				}

				if !rule.SourceCidr.IsNull() && !rule.SourceCidr.IsUnknown() {
					jsonRules[i].SourceCidr = rule.SourceCidr.ValueString()
				}
				if !rule.SourceComponentId.IsNull() && !rule.SourceComponentId.IsUnknown() {
					jsonRules[i].SourceComponentId = rule.SourceComponentId.ValueString()
				}
			}

			b, err := json.Marshal(jsonRules)
			if err != nil {
				resp.Error = function.NewFuncError(fmt.Sprintf("failed to serialize ingress_rules: %s", err))
				return
			}
			params["ingressRules"] = string(b)
		}
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
