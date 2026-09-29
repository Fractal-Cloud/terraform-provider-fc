package paas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &MessagingPaasEntityFunction{}

type MessagingPaasEntityFunction struct{}

func NewMessagingPaasEntityFunction() function.Function {
	return &MessagingPaasEntityFunction{}
}

func (f *MessagingPaasEntityFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "messaging_paas_entity"
}

func (f *MessagingPaasEntityFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a Messaging PaaS Entity blueprint component",
		Description: "Builds a Messaging PaaS Entity component with the correct type for use in a fractal's components list. " +
			"If broker is provided, it is automatically added as a dependency.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Entity configuration",
				CustomType: components.NewConfigObjectType(map[string]attr.Type{
					"id":                      types.StringType,
					"display_name":            types.StringType,
					"description":             types.StringType,
					"message_retention_hours": types.Int64Type,
					"broker":                  components.ComponentObjectType,
					"extra_parameters":        components.ParametersAttrType,
				}, "id"),
			},
		},
		Return: components.ComponentReturn(),
	}
}

type messagingPaasEntityConfig struct {
	Id                    types.String `tfsdk:"id"`
	DisplayName           types.String `tfsdk:"display_name"`
	Description           types.String `tfsdk:"description"`
	MessageRetentionHours types.Int64  `tfsdk:"message_retention_hours"`
	Broker                types.Object `tfsdk:"broker"`
	ExtraParameters       types.Map    `tfsdk:"extra_parameters"`
}

func (f *MessagingPaasEntityFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config messagingPaasEntityConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}

	if !config.MessageRetentionHours.IsNull() && !config.MessageRetentionHours.IsUnknown() {
		params["messageRetentionHours"] = fmt.Sprintf("%d", config.MessageRetentionHours.ValueInt64())
	}

	var deps []string
	brokerId, funcErr := components.ExtractDependency(config.Broker, "Messaging.PaaS.Broker")
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}
	if brokerId != "" {
		deps = append(deps, brokerId)
	}

	parameters, funcErr := components.WithExtraParameters(params, config.ExtraParameters)
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}

	result, funcErr := components.BuildComponent(
		config.Id.ValueString(),
		"Messaging.PaaS.Entity",
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
