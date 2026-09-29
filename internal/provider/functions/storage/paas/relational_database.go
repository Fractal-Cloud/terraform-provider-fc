package paas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &StoragePaasRelationalDatabaseFunction{}

type StoragePaasRelationalDatabaseFunction struct{}

func NewStoragePaasRelationalDatabaseFunction() function.Function {
	return &StoragePaasRelationalDatabaseFunction{}
}

func (f *StoragePaasRelationalDatabaseFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "storage_paas_relational_database"
}

func (f *StoragePaasRelationalDatabaseFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a Relational Database blueprint component",
		Description: "Builds a Relational Database component with the correct type for use in a fractal's components list. " +
			"If dbms is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Relational Database configuration",
				CustomType: components.NewConfigObjectType(map[string]attr.Type{
					"id":               types.StringType,
					"display_name":     types.StringType,
					"description":      types.StringType,
					"collation":        types.StringType,
					"charset":          types.StringType,
					"dbms":             components.ComponentObjectType,
					"extra_parameters": components.ParametersAttrType,
				}, "id"),
			},
		},
		Return: components.ComponentReturn(),
	}
}

type storagePaasRelationalDatabaseConfig struct {
	Id              types.String `tfsdk:"id"`
	DisplayName     types.String `tfsdk:"display_name"`
	Description     types.String `tfsdk:"description"`
	Collation       types.String `tfsdk:"collation"`
	Charset         types.String `tfsdk:"charset"`
	Dbms            types.Object `tfsdk:"dbms"`
	ExtraParameters types.Map    `tfsdk:"extra_parameters"`
}

func (f *StoragePaasRelationalDatabaseFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config storagePaasRelationalDatabaseConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}

	if !config.Collation.IsNull() && !config.Collation.IsUnknown() {
		params["collation"] = config.Collation.ValueString()
	}
	if !config.Charset.IsNull() && !config.Charset.IsUnknown() {
		params["charset"] = config.Charset.ValueString()
	}

	var deps []string
	dbmsId, funcErr := components.ExtractDependency(config.Dbms, "Storage.PaaS.RelationalDbms")
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}
	if dbmsId != "" {
		deps = append(deps, dbmsId)
	}

	parameters, funcErr := components.WithExtraParameters(params, config.ExtraParameters)
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}

	result, funcErr := components.BuildComponent(
		config.Id.ValueString(),
		"Storage.PaaS.RelationalDatabase",
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
