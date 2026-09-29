package caas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &StorageCaasSearchEntityFunction{}

type StorageCaasSearchEntityFunction struct{}

func NewStorageCaasSearchEntityFunction() function.Function {
	return &StorageCaasSearchEntityFunction{}
}

func (f *StorageCaasSearchEntityFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "storage_caas_search_entity"
}

func (f *StorageCaasSearchEntityFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a Search Entity blueprint component",
		Description: "Builds a Search Entity component with the correct type for use in a fractal's components list. " +
			"If search is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Search Entity configuration",
				CustomType: components.NewConfigObjectType(map[string]attr.Type{
					"id":               types.StringType,
					"display_name":     types.StringType,
					"description":      types.StringType,
					"search":           components.ComponentObjectType,
					"extra_parameters": components.ParametersAttrType,
				}, "id"),
			},
		},
		Return: components.ComponentReturn(),
	}
}

type storageCaasSearchEntityConfig struct {
	Id              types.String `tfsdk:"id"`
	DisplayName     types.String `tfsdk:"display_name"`
	Description     types.String `tfsdk:"description"`
	Search          types.Object `tfsdk:"search"`
	ExtraParameters types.Map    `tfsdk:"extra_parameters"`
}

func (f *StorageCaasSearchEntityFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config storageCaasSearchEntityConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	var deps []string
	searchId, funcErr := components.ExtractDependency(config.Search, "Storage.CaaS.Search")
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}
	if searchId != "" {
		deps = append(deps, searchId)
	}

	parameters, funcErr := components.WithExtraParameters(nil, config.ExtraParameters)
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}

	result, funcErr := components.BuildComponent(
		config.Id.ValueString(),
		"Storage.CaaS.SearchEntity",
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
