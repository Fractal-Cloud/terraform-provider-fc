package paas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &BigdataPaasMlExperimentFunction{}

type BigdataPaasMlExperimentFunction struct{}

func NewBigdataPaasMlExperimentFunction() function.Function {
	return &BigdataPaasMlExperimentFunction{}
}

func (f *BigdataPaasMlExperimentFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "bigdata_paas_ml_experiment"
}

func (f *BigdataPaasMlExperimentFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a BigData PaaS ML Experiment blueprint component",
		Description: "Builds a BigData PaaS ML Experiment component with the correct type for use in a fractal's components list. " +
			"If platform is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "ML Experiment configuration",
				AttributeTypes: map[string]attr.Type{
					"id":              types.StringType,
					"display_name":    types.StringType,
					"description":     types.StringType,
					"platform":        components.ComponentObjectType,
					"experiment_name": types.StringType,
				},
			},
		},
		Return: components.ComponentReturn(),
	}
}

type bigdataPaasMlExperimentConfig struct {
	Id             types.String `tfsdk:"id"`
	DisplayName    types.String `tfsdk:"display_name"`
	Description    types.String `tfsdk:"description"`
	Platform       types.Object `tfsdk:"platform"`
	ExperimentName types.String `tfsdk:"experiment_name"`
}

func (f *BigdataPaasMlExperimentFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config bigdataPaasMlExperimentConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}

	if !config.ExperimentName.IsNull() && !config.ExperimentName.IsUnknown() {
		params["experimentName"] = config.ExperimentName.ValueString()
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

	result, funcErr := components.BuildComponent(
		config.Id.ValueString(),
		"BigData.PaaS.MlExperiment",
		components.OptionalString(config.DisplayName),
		components.OptionalString(config.Description),
		types.StringNull(),
		params,
		deps,
		nil,
	)
	resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
