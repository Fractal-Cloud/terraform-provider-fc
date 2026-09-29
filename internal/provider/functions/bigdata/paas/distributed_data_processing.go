package paas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &BigdataPaasDistributedDataProcessingFunction{}

type BigdataPaasDistributedDataProcessingFunction struct{}

func NewBigdataPaasDistributedDataProcessingFunction() function.Function {
	return &BigdataPaasDistributedDataProcessingFunction{}
}

func (f *BigdataPaasDistributedDataProcessingFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "bigdata_paas_distributed_data_processing"
}

func (f *BigdataPaasDistributedDataProcessingFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Creates a BigData PaaS Distributed Data Processing blueprint component",
		Description: "Builds a BigData PaaS Distributed Data Processing component with the correct type for use in a fractal's components list.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Distributed Data Processing configuration",
				CustomType: components.NewConfigObjectType(map[string]attr.Type{
					"id":               types.StringType,
					"display_name":     types.StringType,
					"description":      types.StringType,
					"links":            components.LinksAttrType,
					"extra_parameters": components.ParametersAttrType,
				}, "id"),
			},
		},
		Return: components.ComponentReturn(),
	}
}

type bigdataPaasDistributedDataProcessingConfig struct {
	Id              types.String  `tfsdk:"id"`
	DisplayName     types.String  `tfsdk:"display_name"`
	Description     types.String  `tfsdk:"description"`
	Links           types.Dynamic `tfsdk:"links"`
	ExtraParameters types.Map     `tfsdk:"extra_parameters"`
}

func (f *BigdataPaasDistributedDataProcessingFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config bigdataPaasDistributedDataProcessingConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	var links []components.ComponentLink
	resolved, funcErr := components.LinksFromDynamic(config.Links)
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}
	links = append(links, resolved...)

	parameters, funcErr := components.WithExtraParameters(nil, config.ExtraParameters)
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}

	result, funcErr := components.BuildComponent(
		config.Id.ValueString(),
		"BigData.PaaS.DistributedDataProcessing",
		components.OptionalString(config.DisplayName),
		components.OptionalString(config.Description),
		types.StringNull(),
		parameters,
		nil,
		links,
	)
	resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
