package secrets

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func runSecretRef(t *testing.T, shortName string) *function.RunResponse {
	t.Helper()
	resp := &function.RunResponse{Result: function.NewResultData(types.StringNull())}
	NewSecretRefFunction().Run(context.Background(), function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{types.StringValue(shortName)}),
	}, resp)
	return resp
}

func TestSecretRef_BuildsTaggedReference(t *testing.T) {
	resp := runSecretRef(t, "openai-key")
	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error.Text)
	}
	got := resp.Result.Value().(types.String).ValueString()
	if want := `{"$envSecret":"openai-key"}`; got != want {
		t.Errorf("secret_ref = %s, want %s", got, want)
	}
}

func TestSecretRef_EscapesShortName(t *testing.T) {
	got, err := SecretRef(`a"b`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := `{"$envSecret":"a\"b"}`; got != want {
		t.Errorf("SecretRef = %s, want %s", got, want)
	}
}

func TestSecretRef_RejectsEmptyName(t *testing.T) {
	if resp := runSecretRef(t, "  "); resp.Error == nil {
		t.Error("expected an error for an empty short name")
	}
}
