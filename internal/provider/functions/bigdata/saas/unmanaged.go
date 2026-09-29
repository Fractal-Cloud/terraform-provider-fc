package saas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &BigdataSaasUnmanagedFunction{}

type BigdataSaasUnmanagedFunction struct{}

func NewBigdataSaasUnmanagedFunction() function.Function {
	return &BigdataSaasUnmanagedFunction{}
}

func (f *BigdataSaasUnmanagedFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "bigdata_saas_unmanaged"
}

func (f *BigdataSaasUnmanagedFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a BigData SaaS Unmanaged blueprint component",
		Description: "Builds a BigData SaaS Unmanaged component with the correct type for use in a fractal's components list.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Unmanaged configuration",
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

type bigdataSaasUnmanagedConfig struct {
	Id              types.String `tfsdk:"id"`
	DisplayName     types.String `tfsdk:"display_name"`
	Description     types.String `tfsdk:"description"`
	ExtraParameters types.Map    `tfsdk:"extra_parameters"`
}

func (f *BigdataSaasUnmanagedFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config bigdataSaasUnmanagedConfig
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
		"BigData.SaaS.Unmanaged",
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
