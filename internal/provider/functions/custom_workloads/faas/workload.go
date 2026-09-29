package faas

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
	resp.Name = "custom_workloads_faas_workload"
}

func (f *WorkloadFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a FaaS Workload blueprint component",
		Description: "Builds a FaaS Workload (serverless function) component with the correct type and parameters for use in a fractal's components list. " +
			"Subnet is a component object reference with type validation. " +
			"runtime and handler are required; handler is also written as entryPoint for Cloud Functions. " +
			"Use links to define runtime relationships to other components, and security_groups for SG membership.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "FaaS Workload configuration",
				CustomType: components.NewConfigObjectType(map[string]attr.Type{
					"id":               types.StringType,
					"display_name":     types.StringType,
					"description":      types.StringType,
					"runtime":          types.StringType,
					"memory_mb":        types.Int64Type,
					"timeout_seconds":  types.Int64Type,
					"handler":          types.StringType,
					"subnet":           components.ComponentObjectType,
					"links":            components.LinksAttrType,
					"security_groups":  types.ListType{ElemType: components.ComponentObjectType},
					"extra_parameters": components.ParametersAttrType,
				}, "id", "runtime", "handler"),
			},
		},
		Return: components.ComponentReturn(),
	}
}

type workloadConfig struct {
	Id              types.String  `tfsdk:"id"`
	DisplayName     types.String  `tfsdk:"display_name"`
	Description     types.String  `tfsdk:"description"`
	Runtime         types.String  `tfsdk:"runtime"`
	MemoryMb        types.Int64   `tfsdk:"memory_mb"`
	TimeoutSeconds  types.Int64   `tfsdk:"timeout_seconds"`
	Handler         types.String  `tfsdk:"handler"`
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

	if !config.Runtime.IsNull() && !config.Runtime.IsUnknown() {
		params["runtime"] = config.Runtime.ValueString()
	}
	if !config.MemoryMb.IsNull() && !config.MemoryMb.IsUnknown() {
		params["memoryMb"] = fmt.Sprintf("%d", config.MemoryMb.ValueInt64())
	}
	if !config.TimeoutSeconds.IsNull() && !config.TimeoutSeconds.IsUnknown() {
		params["timeoutSeconds"] = fmt.Sprintf("%d", config.TimeoutSeconds.ValueInt64())
	}
	if !config.Handler.IsNull() && !config.Handler.IsUnknown() {
		// GCP Cloud Functions call the entry point entryPoint.
		params["handler"] = config.Handler.ValueString()
		params["entryPoint"] = config.Handler.ValueString()
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
		"CustomWorkloads.FaaS.Workload",
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
