// Package ai holds the AI component functions that have no delivery model in
// their type.
package ai

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"fractal.cloud/terraform-provider-fc/internal/provider/functions/spec"
)

func NewAgenticPlatformFunction() function.Function {
	return spec.New(spec.Spec{
		Name:          "ai_agentic_platform",
		ComponentType: "AI.AgenticPlatform",
		Summary:       "Creates an Agentic Platform blueprint component",
		Description: "Builds an AI.AgenticPlatform component: a governed runtime for AI agents. " +
			"Link it to a Security.PaaS.IdentityProvider for authentication and to AI.SaaS.Unmanaged components " +
			"for the model services its agents call.",
		Attributes: []spec.Attribute{
			{Name: "bounded_context", Key: "boundedContext", Kind: spec.String, Required: true},
			{Name: "sizing_tier", Key: "sizingTier", Kind: spec.String, OneOf: []string{"small", "standard", "large"}},
			{Name: "default_acf_level", Key: "defaultAcfLevel", Kind: spec.String, OneOf: []string{"ACF-1", "ACF-2"}},
			{Name: "habitat_runtime", Key: "habitatRuntime", Kind: spec.String, OneOf: []string{"auto", "firecracker", "gvisor"}},
		},
		Links: true,
	})
}
