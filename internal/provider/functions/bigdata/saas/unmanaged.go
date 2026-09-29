package saas

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"fractal.cloud/terraform-provider-fc/internal/provider/functions/unmanaged"
)

func NewBigdataSaasUnmanagedFunction() function.Function {
	return unmanaged.New("bigdata_saas_unmanaged", "BigData.SaaS.Unmanaged", "Big Data")
}
