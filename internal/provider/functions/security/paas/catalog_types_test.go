package paas

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
	ft "fractal.cloud/terraform-provider-fc/internal/provider/functions/functiontest"
)

// Each function builds its catalog type, with the minimum the catalog
// requires, and depends on the component its dependency attribute names.
func TestSecurityPaasFunctions_BuildCatalogTypes(t *testing.T) {
	tests := []struct {
		new           func() function.Function
		name          string
		componentType string
		required      map[string]attr.Value
		dependency    string
		dependencyOn  string
	}{
		{NewIdentityProviderFunction, "security_paas_identity_provider", "Security.PaaS.IdentityProvider", nil, "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.new()
			var meta function.MetadataResponse
			f.Metadata(context.Background(), function.MetadataRequest{}, &meta)
			if meta.Name != tt.name {
				t.Errorf("name = %q, want %q", meta.Name, tt.name)
			}

			attrs := map[string]attr.Value{"id": types.StringValue("c")}
			for k, v := range tt.required {
				attrs[k] = v
			}
			if tt.dependency != "" {
				dep, err := components.BuildComponent("dep", tt.dependencyOn, types.StringNull(), types.StringNull(), types.StringNull(), nil, nil, nil)
				if err != nil {
					t.Fatalf("building dependency: %s", err.Text)
				}
				attrs[tt.dependency] = dep
			}

			c := ft.Component(t, ft.Run(t, f, ft.Object(t, attrs)))
			if got := c["type"].(types.String).ValueString(); got != tt.componentType {
				t.Errorf("type = %q, want %q", got, tt.componentType)
			}
			if tt.dependency != "" {
				if deps := ft.Strings(t, c, "dependencies_ids"); len(deps) != 1 || deps[0] != "dep" {
					t.Errorf("dependencies = %v, want [dep]", deps)
				}
			}
		})
	}
}
