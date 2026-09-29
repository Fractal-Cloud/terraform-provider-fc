package paas

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"fractal.cloud/terraform-provider-fc/internal/provider/functions/spec"
)

func NewIdentityProviderFunction() function.Function {
	return spec.New(spec.Spec{
		Name:          "security_paas_identity_provider",
		ComponentType: "Security.PaaS.IdentityProvider",
		Summary:       "Creates an Identity Provider blueprint component",
		Description: "Builds a Security.PaaS.IdentityProvider component: a user directory (Amazon Cognito, Microsoft Entra " +
			"External ID). Its attributes are guardrails that cap what any client may request. Each component that links " +
			"to it gets one OAuth app client, shaped by the link's settings: clientType (\"web\", \"spa\" or \"machine\"), " +
			"redirectUris, logoutUris and scopes.",
		Attributes: []spec.Attribute{
			{Name: "user_directory_name", Key: "userDirectoryName", Kind: spec.String},
			{Name: "mfa_configuration", Key: "mfaConfiguration", Kind: spec.String, OneOf: []string{"OFF", "ON", "OPTIONAL"}},
			{Name: "mfa_break_glass_user_ids", Key: "mfaBreakGlassUserIds", Kind: spec.StringList},
			{Name: "session_duration", Key: "sessionDuration", Kind: spec.Int64},
			{Name: "password_policy", Key: "passwordPolicy", Kind: spec.Object, Fields: []spec.Attribute{
				{Name: "min_length", Key: "minLength", Kind: spec.Int64},
				{Name: "require_uppercase", Key: "requireUppercase", Kind: spec.Bool},
				{Name: "require_lowercase", Key: "requireLowercase", Kind: spec.Bool},
				{Name: "require_numbers", Key: "requireNumbers", Kind: spec.Bool},
				{Name: "require_symbols", Key: "requireSymbols", Kind: spec.Bool},
			}},
		},
	})
}
