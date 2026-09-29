package saas

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"fractal.cloud/terraform-provider-fc/internal/provider/functions/unmanaged"
)

func NewSaaSUnmanagedFunction() function.Function {
	return unmanaged.New("api_management_saas_unmanaged", "APIManagement.SaaS.Unmanaged", "API Management")
}
