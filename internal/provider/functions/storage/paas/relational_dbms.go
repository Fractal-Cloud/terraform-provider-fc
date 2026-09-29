package paas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &StoragePaasRelationalDbmsFunction{}

type StoragePaasRelationalDbmsFunction struct{}

func NewStoragePaasRelationalDbmsFunction() function.Function {
	return &StoragePaasRelationalDbmsFunction{}
}

func (f *StoragePaasRelationalDbmsFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "storage_paas_relational_dbms"
}

func (f *StoragePaasRelationalDbmsFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a Relational DBMS Platform blueprint component",
		Description: "Builds a Relational DBMS Platform component with the correct type for use in a fractal's components list.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Relational DBMS configuration",
				CustomType: components.NewConfigObjectType(map[string]attr.Type{
					"id":               types.StringType,
					"display_name":     types.StringType,
					"description":      types.StringType,
					"engine_version":   types.StringType,
					"extra_parameters": components.ParametersAttrType,
				}, "id"),
			},
		},
		Return: components.ComponentReturn(),
	}
}

type storagePaasRelationalDbmsConfig struct {
	Id              types.String `tfsdk:"id"`
	DisplayName     types.String `tfsdk:"display_name"`
	Description     types.String `tfsdk:"description"`
	EngineVersion   types.String `tfsdk:"engine_version"`
	ExtraParameters types.Map    `tfsdk:"extra_parameters"`
}

func (f *StoragePaasRelationalDbmsFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config storagePaasRelationalDbmsConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}

	if !config.EngineVersion.IsNull() && !config.EngineVersion.IsUnknown() {
		params["version"] = config.EngineVersion.ValueString()
	}

	parameters, funcErr := components.WithExtraParameters(params, config.ExtraParameters)
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}

	result, funcErr := components.BuildComponent(
		config.Id.ValueString(),
		"Storage.PaaS.RelationalDbms",
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
