package saas

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"fractal.cloud/terraform-provider-fc/internal/provider/functions/unmanaged"
)

func NewStorageSaasUnmanagedFunction() function.Function {
	return unmanaged.New("storage_saas_unmanaged", "Storage.SaaS.Unmanaged", "Storage")
}
