package saas

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"fractal.cloud/terraform-provider-fc/internal/provider/functions/unmanaged"
)

func NewUnmanagedFunction() function.Function {
	return unmanaged.New("network_and_compute_saas_unmanaged", "NetworkAndCompute.SaaS.Unmanaged", "Network and Compute")
}
