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

var _ function.Function = &BigdataPaasComputeClusterFunction{}

type BigdataPaasComputeClusterFunction struct{}

func NewBigdataPaasComputeClusterFunction() function.Function {
	return &BigdataPaasComputeClusterFunction{}
}

func (f *BigdataPaasComputeClusterFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "bigdata_paas_compute_cluster"
}

func (f *BigdataPaasComputeClusterFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a BigData PaaS Compute Cluster blueprint component",
		Description: "Builds a BigData PaaS Compute Cluster component with the correct type for use in a fractal's components list. " +
			"If platform is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Compute Cluster configuration",
				AttributeTypes: map[string]attr.Type{
					"id":                       types.StringType,
					"display_name":             types.StringType,
					"description":              types.StringType,
					"platform":                 components.ComponentObjectType,
					"cluster_name":             types.StringType,
					"spark_version":            types.StringType,
					"num_workers":              types.Int64Type,
					"min_workers":              types.Int64Type,
					"max_workers":              types.Int64Type,
					"auto_termination_minutes": types.Int64Type,
					"spark_conf":               types.MapType{ElemType: types.StringType},
					"pypi_libraries":           types.ListType{ElemType: types.StringType},
					"maven_libraries":          types.ListType{ElemType: types.StringType},
					"links": types.ListType{
						ElemType: types.ObjectType{AttrTypes: components.GenericLinkAttrTypes},
					},
				},
			},
		},
		Return: components.ComponentReturn(),
	}
}

type bigdataPaasComputeClusterConfig struct {
	Id                     types.String `tfsdk:"id"`
	DisplayName            types.String `tfsdk:"display_name"`
	Description            types.String `tfsdk:"description"`
	Platform               types.Object `tfsdk:"platform"`
	ClusterName            types.String `tfsdk:"cluster_name"`
	SparkVersion           types.String `tfsdk:"spark_version"`
	NumWorkers             types.Int64  `tfsdk:"num_workers"`
	MinWorkers             types.Int64  `tfsdk:"min_workers"`
	MaxWorkers             types.Int64  `tfsdk:"max_workers"`
	AutoTerminationMinutes types.Int64  `tfsdk:"auto_termination_minutes"`
	SparkConf              types.Map    `tfsdk:"spark_conf"`
	PypiLibraries          types.List   `tfsdk:"pypi_libraries"`
	MavenLibraries         types.List   `tfsdk:"maven_libraries"`
	Links                  types.List   `tfsdk:"links"`
}

func (f *BigdataPaasComputeClusterFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config bigdataPaasComputeClusterConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}

	if !config.ClusterName.IsNull() && !config.ClusterName.IsUnknown() {
		params["clusterName"] = config.ClusterName.ValueString()
	}
	if !config.SparkVersion.IsNull() && !config.SparkVersion.IsUnknown() {
		params["sparkVersion"] = config.SparkVersion.ValueString()
	}
	if !config.NumWorkers.IsNull() && !config.NumWorkers.IsUnknown() {
		params["numWorkers"] = fmt.Sprintf("%d", config.NumWorkers.ValueInt64())
	}
	if !config.MinWorkers.IsNull() && !config.MinWorkers.IsUnknown() {
		params["minWorkers"] = fmt.Sprintf("%d", config.MinWorkers.ValueInt64())
	}
	if !config.MaxWorkers.IsNull() && !config.MaxWorkers.IsUnknown() {
		params["maxWorkers"] = fmt.Sprintf("%d", config.MaxWorkers.ValueInt64())
	}
	if !config.AutoTerminationMinutes.IsNull() && !config.AutoTerminationMinutes.IsUnknown() {
		params["autoTerminationMinutes"] = fmt.Sprintf("%d", config.AutoTerminationMinutes.ValueInt64())
	}

	if !config.SparkConf.IsNull() && !config.SparkConf.IsUnknown() {
		sparkConfMap := make(map[string]string, len(config.SparkConf.Elements()))
		for k, v := range config.SparkConf.Elements() {
			sparkConfMap[k] = v.(types.String).ValueString()
		}
		sparkConfJSON, err := json.Marshal(sparkConfMap)
		if err != nil {
			resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("failed to marshal spark_conf: "+err.Error()))
			return
		}
		params["sparkConf"] = string(sparkConfJSON)
	}

	if !config.PypiLibraries.IsNull() && !config.PypiLibraries.IsUnknown() {
		pypiList := make([]string, 0, len(config.PypiLibraries.Elements()))
		for _, v := range config.PypiLibraries.Elements() {
			pypiList = append(pypiList, v.(types.String).ValueString())
		}
		pypiJSON, err := json.Marshal(pypiList)
		if err != nil {
			resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("failed to marshal pypi_libraries: "+err.Error()))
			return
		}
		params["pypiLibraries"] = string(pypiJSON)
	}

	if !config.MavenLibraries.IsNull() && !config.MavenLibraries.IsUnknown() {
		mavenList := make([]string, 0, len(config.MavenLibraries.Elements()))
		for _, v := range config.MavenLibraries.Elements() {
			mavenList = append(mavenList, v.(types.String).ValueString())
		}
		mavenJSON, err := json.Marshal(mavenList)
		if err != nil {
			resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("failed to marshal maven_libraries: "+err.Error()))
			return
		}
		params["mavenLibraries"] = string(mavenJSON)
	}

	var deps []string
	platformId, funcErr := components.ExtractDependency(config.Platform, "BigData.PaaS.DistributedDataProcessing")
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}
	if platformId != "" {
		deps = append(deps, platformId)
	}

	var links []components.ComponentLink
	if !config.Links.IsNull() && !config.Links.IsUnknown() {
		var genericLinks []components.GenericLinkConfig
		diags := config.Links.ElementsAs(ctx, &genericLinks, false)
		if diags.HasError() {
			resp.Error = function.NewFuncError("failed to parse links")
			return
		}
		resolved, funcErr := components.GenericLinksToComponentLinks(genericLinks)
		if funcErr != nil {
			resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
			return
		}
		links = append(links, resolved...)
	}

	result, funcErr := components.BuildComponent(
		config.Id.ValueString(),
		"BigData.PaaS.ComputeCluster",
		components.OptionalString(config.DisplayName),
		components.OptionalString(config.Description),
		types.StringNull(),
		params,
		deps,
		links,
	)
	resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
