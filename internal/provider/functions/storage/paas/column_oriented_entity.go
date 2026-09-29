package paas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &StoragePaasColumnOrientedEntityFunction{}

type StoragePaasColumnOrientedEntityFunction struct{}

func NewStoragePaasColumnOrientedEntityFunction() function.Function {
	return &StoragePaasColumnOrientedEntityFunction{}
}

func (f *StoragePaasColumnOrientedEntityFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "storage_paas_column_oriented_entity"
}

func (f *StoragePaasColumnOrientedEntityFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a Column-Oriented Entity blueprint component",
		Description: "Builds a Column-Oriented Entity component with the correct type for use in a fractal's components list. " +
			"If dbms is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Column-Oriented Entity configuration",
				CustomType: components.NewConfigObjectType(map[string]attr.Type{
					"id":               types.StringType,
					"display_name":     types.StringType,
					"description":      types.StringType,
					"dbms":             components.ComponentObjectType,
					"extra_parameters": components.ParametersAttrType,
				}, "id"),
			},
		},
		Return: components.ComponentReturn(),
	}
}

type storagePaasColumnOrientedEntityConfig struct {
	Id              types.String `tfsdk:"id"`
	DisplayName     types.String `tfsdk:"display_name"`
	Description     types.String `tfsdk:"description"`
	Dbms            types.Object `tfsdk:"dbms"`
	ExtraParameters types.Map    `tfsdk:"extra_parameters"`
}

func (f *StoragePaasColumnOrientedEntityFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config storagePaasColumnOrientedEntityConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	var deps []string
	dbmsId, funcErr := components.ExtractDependency(config.Dbms, "Storage.PaaS.ColumnOrientedDbms")
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}
	if dbmsId != "" {
		deps = append(deps, dbmsId)
	}

	parameters, funcErr := components.WithExtraParameters(nil, config.ExtraParameters)
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}

	result, funcErr := components.BuildComponent(
		config.Id.ValueString(),
		"Storage.PaaS.ColumnOrientedEntity",
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
