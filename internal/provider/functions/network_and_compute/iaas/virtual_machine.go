package iaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &VirtualMachineFunction{}

type VirtualMachineFunction struct{}

func NewVirtualMachineFunction() function.Function {
	return &VirtualMachineFunction{}
}

func (f *VirtualMachineFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "network_and_compute_iaas_virtual_machine"
}

func (f *VirtualMachineFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a VirtualMachine blueprint component",
		Description: "Builds a VirtualMachine (VM/EC2) component with the correct type for use in a fractal's components list. " +
			"If subnet is provided, it is validated as a Subnet and automatically added as a dependency. " +
			"Use links to define runtime relationships to other components, and security_groups for SG membership.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "VirtualMachine configuration",
				CustomType: components.NewConfigObjectType(map[string]attr.Type{
					"id":               types.StringType,
					"display_name":     types.StringType,
					"description":      types.StringType,
					"subnet":           components.ComponentObjectType,
					"links":            components.LinksAttrType,
					"security_groups":  types.ListType{ElemType: components.ComponentObjectType},
					"extra_parameters": components.ParametersAttrType,
				}, "id"),
			},
		},
		Return: components.ComponentReturn(),
	}
}

type virtualMachineConfig struct {
	Id              types.String  `tfsdk:"id"`
	DisplayName     types.String  `tfsdk:"display_name"`
	Description     types.String  `tfsdk:"description"`
	Subnet          types.Object  `tfsdk:"subnet"`
	Links           types.Dynamic `tfsdk:"links"`
	SecurityGroups  types.List    `tfsdk:"security_groups"`
	ExtraParameters types.Map     `tfsdk:"extra_parameters"`
}

func (f *VirtualMachineFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config virtualMachineConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	var deps []string
	subnetId, funcErr := components.ExtractDependency(config.Subnet, "NetworkAndCompute.IaaS.Subnet")
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}
	if subnetId != "" {
		deps = append(deps, subnetId)
	}

	// Build links from generic links and SG memberships
	var links []components.ComponentLink

	resolved, funcErr := components.LinksFromDynamic(config.Links)
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}
	links = append(links, resolved...)

	if !config.SecurityGroups.IsNull() && !config.SecurityGroups.IsUnknown() {
		var sgObjects []types.Object
		diags := config.SecurityGroups.ElementsAs(ctx, &sgObjects, false)
		if diags.HasError() {
			resp.Error = function.NewFuncError("failed to parse security_groups")
			return
		}
		sgLinks, funcErr := components.SgMembershipLinks(sgObjects)
		if funcErr != nil {
			resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
			return
		}
		links = append(links, sgLinks...)
	}

	parameters, funcErr := components.WithExtraParameters(nil, config.ExtraParameters)
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}

	result, funcErr := components.BuildComponent(
		config.Id.ValueString(),
		"NetworkAndCompute.IaaS.VirtualMachine",
		components.OptionalString(config.DisplayName),
		components.OptionalString(config.Description),
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
