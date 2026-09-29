package paas

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
	resp.Name = "custom_workloads_paas_workload"
}

func (f *WorkloadFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a PaaS Workload blueprint component",
		Description: "Builds a PaaS Workload component with the correct type and parameters for use in a fractal's components list. " +
			"Subnet is a component object reference with type validation. " +
			"container_image and container_port are written as image and port, the keys the container-based PaaS offers read. " +
			"Use links to define runtime relationships to other components, and security_groups for SG membership.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "PaaS Workload configuration",
				CustomType: components.NewConfigObjectType(map[string]attr.Type{
					"id":               types.StringType,
					"display_name":     types.StringType,
					"description":      types.StringType,
					"container_image":  types.StringType,
					"container_port":   types.Int64Type,
					"cpu":              types.StringType,
					"memory":           types.StringType,
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
	Cpu             types.String  `tfsdk:"cpu"`
	Memory          types.String  `tfsdk:"memory"`
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

	// The container-based PaaS offers (Cloud Run, Container Instances) read
	// image and port; the Web App offer deploys from git and takes its
	// settings through extra_parameters.
	if !config.ContainerImage.IsNull() && !config.ContainerImage.IsUnknown() {
		params["image"] = config.ContainerImage.ValueString()
	}
	if !config.ContainerPort.IsNull() && !config.ContainerPort.IsUnknown() {
		params["port"] = fmt.Sprintf("%d", config.ContainerPort.ValueInt64())
	}
	if !config.Cpu.IsNull() && !config.Cpu.IsUnknown() {
		params["cpu"] = config.Cpu.ValueString()
	}
	if !config.Memory.IsNull() && !config.Memory.IsUnknown() {
		params["memory"] = config.Memory.ValueString()
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
		"CustomWorkloads.PaaS.Workload",
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
