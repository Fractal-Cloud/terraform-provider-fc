package caas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &MessagingCaasBrokerFunction{}

type MessagingCaasBrokerFunction struct{}

func NewMessagingCaasBrokerFunction() function.Function {
	return &MessagingCaasBrokerFunction{}
}

func (f *MessagingCaasBrokerFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "messaging_caas_broker"
}

func (f *MessagingCaasBrokerFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a Messaging CaaS Broker blueprint component",
		Description: "Builds a Messaging CaaS Broker component with the correct type for use in a fractal's components list. " +
			"If container_platform is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Broker configuration",
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

type messagingCaasBrokerConfig struct {
	Id                types.String `tfsdk:"id"`
	DisplayName       types.String `tfsdk:"display_name"`
	Description       types.String `tfsdk:"description"`
	ContainerPlatform types.Object `tfsdk:"container_platform"`
	ExtraParameters   types.Map    `tfsdk:"extra_parameters"`
}

func (f *MessagingCaasBrokerFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config messagingCaasBrokerConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	var deps []string
	cpId, funcErr := components.ExtractDependency(config.ContainerPlatform, "NetworkAndCompute.PaaS.ContainerPlatform")
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
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
		"Messaging.CaaS.Broker",
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
