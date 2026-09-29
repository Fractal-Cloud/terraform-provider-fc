package paas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &ContainerPlatformFunction{}

type ContainerPlatformFunction struct{}

func NewContainerPlatformFunction() function.Function {
	return &ContainerPlatformFunction{}
}

func (f *ContainerPlatformFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "network_and_compute_paas_container_platform"
}

var nodePoolAttrTypes = map[string]attr.Type{
	"name":                types.StringType,
	"disk_size_gb":        types.Int64Type,
	"min_node_count":      types.Int64Type,
	"max_node_count":      types.Int64Type,
	"max_pods_per_node":   types.Int64Type,
	"autoscaling_enabled": types.BoolType,
	"initial_node_count":  types.Int64Type,
	"max_surge":           types.Int64Type,
}

func (f *ContainerPlatformFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a ContainerPlatform blueprint component",
		Description: "Builds a ContainerPlatform (managed Kubernetes) component with the correct type for use in a fractal's components list.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "ContainerPlatform configuration",
				CustomType: components.NewConfigObjectType(map[string]attr.Type{
					"id":           types.StringType,
					"display_name": types.StringType,
					"description":  types.StringType,
					"node_pools": types.ListType{
						ElemType: types.ObjectType{AttrTypes: nodePoolAttrTypes},
					},
					"extra_parameters": components.ParametersAttrType,
				}, "id"),
			},
		},
		Return: components.ComponentReturn(),
	}
}

type containerPlatformConfig struct {
	Id              types.String `tfsdk:"id"`
	DisplayName     types.String `tfsdk:"display_name"`
	Description     types.String `tfsdk:"description"`
	NodePools       types.List   `tfsdk:"node_pools"`
	ExtraParameters types.Map    `tfsdk:"extra_parameters"`
}

type nodePoolConfig struct {
	Name               types.String `tfsdk:"name"`
	DiskSizeGb         types.Int64  `tfsdk:"disk_size_gb"`
	MinNodeCount       types.Int64  `tfsdk:"min_node_count"`
	MaxNodeCount       types.Int64  `tfsdk:"max_node_count"`
	MaxPodsPerNode     types.Int64  `tfsdk:"max_pods_per_node"`
	AutoscalingEnabled types.Bool   `tfsdk:"autoscaling_enabled"`
	InitialNodeCount   types.Int64  `tfsdk:"initial_node_count"`
	MaxSurge           types.Int64  `tfsdk:"max_surge"`
}

type nodePoolJSON struct {
	Name               string `json:"name"`
	DiskSizeGb         int64  `json:"diskSizeGb,omitempty"`
	MinNodeCount       int64  `json:"minNodeCount,omitempty"`
	MaxNodeCount       int64  `json:"maxNodeCount,omitempty"`
	MaxPodsPerNode     int64  `json:"maxPodsPerNode,omitempty"`
	AutoscalingEnabled *bool  `json:"autoscalingEnabled,omitempty"`
	InitialNodeCount   int64  `json:"initialNodeCount,omitempty"`
	MaxSurge           int64  `json:"maxSurge,omitempty"`
}

func (f *ContainerPlatformFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config containerPlatformConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}

	if !config.NodePools.IsNull() && !config.NodePools.IsUnknown() {
		var pools []nodePoolConfig
		diags := config.NodePools.ElementsAs(ctx, &pools, false)
		if diags.HasError() {
			resp.Error = function.NewFuncError("failed to parse node_pools")
			return
		}

		if len(pools) > 0 {
			jsonPools := make([]nodePoolJSON, len(pools))
			for i, pool := range pools {
				jsonPools[i] = nodePoolJSON{
					Name: pool.Name.ValueString(),
				}
				if !pool.DiskSizeGb.IsNull() && !pool.DiskSizeGb.IsUnknown() {
					jsonPools[i].DiskSizeGb = pool.DiskSizeGb.ValueInt64()
				}
				if !pool.MinNodeCount.IsNull() && !pool.MinNodeCount.IsUnknown() {
					jsonPools[i].MinNodeCount = pool.MinNodeCount.ValueInt64()
				}
				if !pool.MaxNodeCount.IsNull() && !pool.MaxNodeCount.IsUnknown() {
					jsonPools[i].MaxNodeCount = pool.MaxNodeCount.ValueInt64()
				}
				if !pool.MaxPodsPerNode.IsNull() && !pool.MaxPodsPerNode.IsUnknown() {
					jsonPools[i].MaxPodsPerNode = pool.MaxPodsPerNode.ValueInt64()
				}
				if !pool.AutoscalingEnabled.IsNull() && !pool.AutoscalingEnabled.IsUnknown() {
					v := pool.AutoscalingEnabled.ValueBool()
					jsonPools[i].AutoscalingEnabled = &v
				}
				if !pool.InitialNodeCount.IsNull() && !pool.InitialNodeCount.IsUnknown() {
					jsonPools[i].InitialNodeCount = pool.InitialNodeCount.ValueInt64()
				}
				if !pool.MaxSurge.IsNull() && !pool.MaxSurge.IsUnknown() {
					jsonPools[i].MaxSurge = pool.MaxSurge.ValueInt64()
				}
			}

			b, err := json.Marshal(jsonPools)
			if err != nil {
				resp.Error = function.NewFuncError(fmt.Sprintf("failed to serialize node_pools: %s", err))
				return
			}
			params["nodePools"] = string(b)
		}
	}

	parameters, funcErr := components.WithExtraParameters(params, config.ExtraParameters)
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}

	result, funcErr := components.BuildComponent(
		config.Id.ValueString(),
		"NetworkAndCompute.PaaS.ContainerPlatform",
		components.OptionalString(config.DisplayName),
		components.OptionalString(config.Description),
		types.StringNull(),
		parameters,
		nil,
		nil,
	)
	resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
