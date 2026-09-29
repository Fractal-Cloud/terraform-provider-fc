package saas

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"fractal.cloud/terraform-provider-fc/internal/provider/functions/unmanaged"
)

func NewUnmanagedFunction() function.Function {
	return unmanaged.New("custom_workloads_saas_unmanaged", "CustomWorkloads.SaaS.Unmanaged", "Custom Workloads")
}
