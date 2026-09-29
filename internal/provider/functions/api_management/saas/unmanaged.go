package saas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &SaaSUnmanagedFunction{}

type SaaSUnmanagedFunction struct{}

func NewSaaSUnmanagedFunction() function.Function {
	return &SaaSUnmanagedFunction{}
}

func (f *SaaSUnmanagedFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "api_management_saas_unmanaged"
}

func (f *SaaSUnmanagedFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates an unmanaged API Management blueprint component",
		Description: "Builds an unmanaged API Management component with the correct type for use in a fractal's components list.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Unmanaged API Management configuration",
				CustomType: components.NewConfigObjectType(map[string]attr.Type{
					"id":               types.StringType,
					"display_name":     types.StringType,
					"description":      types.StringType,
					"extra_parameters": components.ParametersAttrType,
				}, "id"),
			},
		},
		Return: components.ComponentReturn(),
	}
}

type saasUnmanagedConfig struct {
	Id              types.String `tfsdk:"id"`
	DisplayName     types.String `tfsdk:"display_name"`
	Description     types.String `tfsdk:"description"`
	ExtraParameters types.Map    `tfsdk:"extra_parameters"`
}

func (f *SaaSUnmanagedFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config saasUnmanagedConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	parameters, funcErr := components.WithExtraParameters(nil, config.ExtraParameters)
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}

	result, funcErr := components.BuildComponent(
		config.Id.ValueString(),
		"APIManagement.SaaS.Unmanaged",
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
