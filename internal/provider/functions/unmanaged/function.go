// Package unmanaged implements the <Domain>.SaaS.Unmanaged component
// functions, which model an external service the platform does not manage.
package unmanaged

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
	"fractal.cloud/terraform-provider-fc/internal/provider/functions/secrets"
)

var _ function.Function = &Function{}

// Function builds an Unmanaged component of one infrastructure domain.
//
// An Unmanaged component carries the credentials of the external service as
// an environment-secret reference in its `secret` parameter; the agent grants
// linked consumers read access to that secret and injects a reference to it,
// never the raw value. The agent rejects an Unmanaged component with no
// parameters, so `secret` is required, and a raw value is refused here so it
// never reaches the blueprint.
type Function struct {
	name          string
	componentType string
	domain        string
}

// New returns the function called name that builds componentType; domain is
// the human-readable domain used in its documentation.
func New(name, componentType, domain string) function.Function {
	return &Function{name: name, componentType: componentType, domain: domain}
}

func (f *Function) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = f.name
}

func (f *Function) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: fmt.Sprintf("Creates an unmanaged %s blueprint component", f.domain),
		Description: fmt.Sprintf("Builds a %s component: an external %s service that consumers link to. "+
			"secret references the environment secret holding its credentials, built with secret_ref().", f.componentType, f.domain),
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: fmt.Sprintf("Unmanaged %s configuration", f.domain),
				CustomType: components.NewConfigObjectType(map[string]attr.Type{
					"id":               types.StringType,
					"display_name":     types.StringType,
					"description":      types.StringType,
					"secret":           types.StringType,
					"extra_parameters": components.ParametersAttrType,
				}, "id", "secret"),
			},
		},
		Return: components.ComponentReturn(),
	}
}

type config struct {
	Id              types.String `tfsdk:"id"`
	DisplayName     types.String `tfsdk:"display_name"`
	Description     types.String `tfsdk:"description"`
	Secret          types.String `tfsdk:"secret"`
	ExtraParameters types.Map    `tfsdk:"extra_parameters"`
}

func (f *Function) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var cfg config
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &cfg))
	if resp.Error != nil {
		return
	}

	if cfg.Secret.IsNull() {
		resp.Error = function.NewArgumentFuncError(0, "secret is required")
		return
	}
	if _, ok := secrets.ParseSecretRef(cfg.Secret.ValueString()); !ok {
		resp.Error = function.NewArgumentFuncError(0,
			`secret must reference an environment secret, e.g. secret = provider::fc::secret_ref("openai-key"); `+
				"raw secret values are not accepted")
		return
	}

	parameters, funcErr := components.WithExtraParameters(map[string]string{"secret": cfg.Secret.ValueString()}, cfg.ExtraParameters)
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}

	result, funcErr := components.BuildComponent(
		cfg.Id.ValueString(),
		f.componentType,
		components.OptionalString(cfg.DisplayName),
		components.OptionalString(cfg.Description),
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
