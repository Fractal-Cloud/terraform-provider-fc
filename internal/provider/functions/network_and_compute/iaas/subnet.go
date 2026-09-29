package iaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &SubnetFunction{}

type SubnetFunction struct{}

func NewSubnetFunction() function.Function {
	return &SubnetFunction{}
}

func (f *SubnetFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "network_and_compute_iaas_subnet"
}

func (f *SubnetFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a Subnet blueprint component",
		Description: "Builds a Subnet component with the correct type and parameters for use in a fractal's components list. If vpc is provided, it is validated as a VirtualNetwork and automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Subnet configuration",
				CustomType: components.NewConfigObjectType(map[string]attr.Type{
					"id":               types.StringType,
					"display_name":     types.StringType,
					"description":      types.StringType,
					"cidr_block":       types.StringType,
					"vpc":              components.ComponentObjectType,
					"extra_parameters": components.ParametersAttrType,
				}, "id"),
			},
		},
		Return: components.ComponentReturn(),
	}
}

type subnetConfig struct {
	Id              types.String `tfsdk:"id"`
	DisplayName     types.String `tfsdk:"display_name"`
	Description     types.String `tfsdk:"description"`
	CidrBlock       types.String `tfsdk:"cidr_block"`
	Vpc             types.Object `tfsdk:"vpc"`
	ExtraParameters types.Map    `tfsdk:"extra_parameters"`
}

func (f *SubnetFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config subnetConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}
	if !config.CidrBlock.IsNull() && !config.CidrBlock.IsUnknown() {
		params["cidrBlock"] = config.CidrBlock.ValueString()
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
		"NetworkAndCompute.IaaS.Subnet",
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
