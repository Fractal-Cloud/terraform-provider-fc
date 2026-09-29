package caas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &WorkloadFunction{}

type WorkloadFunction struct{}

func NewWorkloadFunction() function.Function {
	return &WorkloadFunction{}
}

func (f *WorkloadFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "custom_workloads_caas_workload"
}

func (f *WorkloadFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a CaaS Workload blueprint component",
		Description: "Builds a CaaS Workload component with the correct type and parameters for use in a fractal's components list. " +
			"Platform and subnet are component object references with type validation. " +
			"Use links to define runtime relationships to other components, and security_groups for SG membership.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "CaaS Workload configuration",
				CustomType: components.NewConfigObjectType(map[string]attr.Type{
					"id":               types.StringType,
					"display_name":     types.StringType,
					"description":      types.StringType,
					"container_image":  types.StringType,
					"container_port":   types.Int64Type,
					"container_name":   types.StringType,
					"cpu":              types.StringType,
					"memory":           types.StringType,
					"desired_count":    types.Int64Type,
					"platform":         components.ComponentObjectType,
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

type workloadConfig struct {
	Id              types.String  `tfsdk:"id"`
	DisplayName     types.String  `tfsdk:"display_name"`
	Description     types.String  `tfsdk:"description"`
	ContainerImage  types.String  `tfsdk:"container_image"`
	ContainerPort   types.Int64   `tfsdk:"container_port"`
	ContainerName   types.String  `tfsdk:"container_name"`
	Cpu             types.String  `tfsdk:"cpu"`
	Memory          types.String  `tfsdk:"memory"`
	DesiredCount    types.Int64   `tfsdk:"desired_count"`
	Platform        types.Object  `tfsdk:"platform"`
	Subnet          types.Object  `tfsdk:"subnet"`
	Links           types.Dynamic `tfsdk:"links"`
	SecurityGroups  types.List    `tfsdk:"security_groups"`
	ExtraParameters types.Map     `tfsdk:"extra_parameters"`
}

func (f *WorkloadFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config workloadConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}

	if !config.ContainerImage.IsNull() && !config.ContainerImage.IsUnknown() {
		params["containerImage"] = config.ContainerImage.ValueString()
	}
	if !config.ContainerPort.IsNull() && !config.ContainerPort.IsUnknown() {
		params["containerPort"] = fmt.Sprintf("%d", config.ContainerPort.ValueInt64())
	}
	if !config.ContainerName.IsNull() && !config.ContainerName.IsUnknown() {
		params["containerName"] = config.ContainerName.ValueString()
	}
	if !config.Cpu.IsNull() && !config.Cpu.IsUnknown() {
		params["cpu"] = config.Cpu.ValueString()
	}
	if !config.Memory.IsNull() && !config.Memory.IsUnknown() {
		params["memory"] = config.Memory.ValueString()
	}
	if !config.DesiredCount.IsNull() && !config.DesiredCount.IsUnknown() {
		params["desiredCount"] = fmt.Sprintf("%d", config.DesiredCount.ValueInt64())
	}

	var deps []string

	platformId, funcErr := components.ExtractDependency(config.Platform, "NetworkAndCompute.PaaS.ContainerPlatform")
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}
	if platformId != "" {
		deps = append(deps, platformId)
	}

	subnetId, funcErr := components.ExtractDependency(config.Subnet, "NetworkAndCompute.IaaS.Subnet")
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}
	if subnetId != "" {
		deps = append(deps, subnetId)
	}

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

	parameters, funcErr := components.WithExtraParameters(params, config.ExtraParameters)
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}

	result, funcErr := components.BuildComponent(
		config.Id.ValueString(),
		"CustomWorkloads.CaaS.Workload",
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
