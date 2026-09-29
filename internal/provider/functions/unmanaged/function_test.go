package unmanaged

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
	"fractal.cloud/terraform-provider-fc/internal/provider/functions/functiontest"
)

func run(t *testing.T, attrs map[string]attr.Value) *function.RunResponse {
	t.Helper()
	f := New("ai_saas_unmanaged", "AI.SaaS.Unmanaged", "AI")
	attrTypes := make(map[string]attr.Type, len(attrs))
	for k, v := range attrs {
		attrTypes[k] = v.Type(context.Background())
	}
	resp := &function.RunResponse{Result: function.NewResultData(types.ObjectNull(components.ComponentAttrTypes))}
	f.Run(context.Background(), function.RunRequest{
		Arguments: function.NewArgumentsData(functiontest.Complete(t, f, []attr.Value{types.ObjectValueMust(attrTypes, attrs)})),
	}, resp)
	return resp
}

func TestUnmanaged_RejectsRawSecret(t *testing.T) {
	resp := run(t, map[string]attr.Value{"id": types.StringValue("openai"), "secret": types.StringValue("sk-live-123")})
	if resp.Error == nil || !strings.Contains(resp.Error.Text, "raw secret values are not accepted") {
		t.Errorf("error = %v, want a raw-secret rejection", resp.Error)
	}
}

func TestUnmanaged_RequiresSecret(t *testing.T) {
	resp := run(t, map[string]attr.Value{"id": types.StringValue("openai")})
	if resp.Error == nil || !strings.Contains(resp.Error.Text, "secret is required") {
		t.Errorf("error = %v, want a missing-secret error", resp.Error)
	}
}

func TestUnmanaged_DeclaresSecretRequired(t *testing.T) {
	var def function.DefinitionResponse
	New("x", "AI.SaaS.Unmanaged", "AI").Definition(context.Background(), function.DefinitionRequest{}, &def)
	cfg := def.Definition.Parameters[0].(function.ObjectParameter).CustomType.(components.ConfigObjectType)
	obj := cfg.TerraformType(context.Background()).(tftypes.Object)
	if _, ok := obj.AttributeTypes["secret"]; !ok {
		t.Fatal("secret attribute missing")
	}
	if _, optional := obj.OptionalAttributes["secret"]; optional {
		t.Error("secret must be required")
	}
}
