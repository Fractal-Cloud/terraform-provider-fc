package saas

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
	"fractal.cloud/terraform-provider-fc/internal/provider/functions/functiontest"
)

func TestUnmanagedFunction_Metadata(t *testing.T) {
	resp := &function.MetadataResponse{}
	NewMessagingSaasUnmanagedFunction().Metadata(context.Background(), function.MetadataRequest{}, resp)
	if resp.Name != "messaging_saas_unmanaged" {
		t.Errorf("expected name %q, got %q", "messaging_saas_unmanaged", resp.Name)
	}
}

func TestUnmanagedFunction_BuildsTypedComponentWithSecret(t *testing.T) {
	f := NewMessagingSaasUnmanagedFunction()
	config := types.ObjectValueMust(map[string]attr.Type{
		"id":     types.StringType,
		"secret": types.StringType,
	}, map[string]attr.Value{
		"id":     types.StringValue("external"),
		"secret": types.StringValue(`{"$envSecret":"external-key"}`),
	})
	resp := &function.RunResponse{Result: function.NewResultData(types.ObjectNull(components.ComponentAttrTypes))}
	f.Run(context.Background(), function.RunRequest{
		Arguments: function.NewArgumentsData(functiontest.Complete(t, f, []attr.Value{config})),
	}, resp)
	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error.Text)
	}

	attrs := resp.Result.Value().(types.Object).Attributes()
	if got := attrs["type"].(types.String).ValueString(); got != "Messaging.SaaS.Unmanaged" {
		t.Errorf("type = %q, want %q", got, "Messaging.SaaS.Unmanaged")
	}
	params := attrs["parameters"].(types.Map).Elements()
	if got := params["secret"].(types.String).ValueString(); got != `{"$envSecret":"external-key"}` {
		t.Errorf("parameters.secret = %s", got)
	}
}
