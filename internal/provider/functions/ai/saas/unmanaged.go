package saas

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"fractal.cloud/terraform-provider-fc/internal/provider/functions/unmanaged"
)

func NewAiSaasUnmanagedFunction() function.Function {
	return unmanaged.New("ai_saas_unmanaged", "AI.SaaS.Unmanaged", "AI")
}
