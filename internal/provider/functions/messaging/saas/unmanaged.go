package saas

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"fractal.cloud/terraform-provider-fc/internal/provider/functions/unmanaged"
)

func NewMessagingSaasUnmanagedFunction() function.Function {
	return unmanaged.New("messaging_saas_unmanaged", "Messaging.SaaS.Unmanaged", "Messaging")
}
