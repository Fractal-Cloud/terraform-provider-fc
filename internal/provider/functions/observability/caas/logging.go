package caas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &CaaSLoggingFunction{}

type CaaSLoggingFunction struct{}

func NewCaaSLoggingFunction() function.Function {
	return &CaaSLoggingFunction{}
}

func (f *CaaSLoggingFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "observability_caas_logging"
}

func (f *CaaSLoggingFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a containerized Logging blueprint component",
		Description: "Builds a CaaS Logging component with the correct type for use in a fractal's components list. The container_platform dependency is validated as a NetworkAndCompute.PaaS.ContainerPlatform component.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "CaaS Logging configuration",
				CustomType: components.NewConfigObjectType(map[string]attr.Type{
					"id":                 types.StringType,
					"display_name":       types.StringType,
					"description":        types.StringType,
					"container_platform": components.ComponentObjectType,
					"extra_parameters":   components.ParametersAttrType,
				}, "id"),
			},
		},
		Return: components.ComponentReturn(),
	}
}

type caasLoggingConfig struct {
	Id                types.String `tfsdk:"id"`
	DisplayName       types.String `tfsdk:"display_name"`
	Description       types.String `tfsdk:"description"`
	ContainerPlatform types.Object `tfsdk:"container_platform"`
	ExtraParameters   types.Map    `tfsdk:"extra_parameters"`
}

func (f *CaaSLoggingFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config caasLoggingConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	var deps []string
	cpId, funcErr := components.ExtractDependency(config.ContainerPlatform, "NetworkAndCompute.PaaS.ContainerPlatform")
	resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	if resp.Error != nil {
		return
	}
	if cpId != "" {
		deps = append(deps, cpId)
	}

	parameters, funcErr := components.WithExtraParameters(nil, config.ExtraParameters)
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}

	result, funcErr := components.BuildComponent(
		config.Id.ValueString(),
		"Observability.CaaS.Logging",
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
