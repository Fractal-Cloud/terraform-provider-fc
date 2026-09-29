package unmanaged

import "github.com/hashicorp/terraform-plugin-framework/types"

// config is the argument of an Unmanaged component function.
type config struct {
	Id              types.String `tfsdk:"id"`
	DisplayName     types.String `tfsdk:"display_name"`
	Description     types.String `tfsdk:"description"`
	Secret          types.String `tfsdk:"secret"`
	ExtraParameters types.Map    `tfsdk:"extra_parameters"`
}
