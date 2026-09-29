package components

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestWithExtraParameters_NullKeepsDerived(t *testing.T) {
	derived := map[string]string{"version": "16"}
	got, err := WithExtraParameters(derived, types.MapNull(types.StringType))
	if err != nil || got["version"] != "16" || len(got) != 1 {
		t.Errorf("got %v, %v; want derived parameters unchanged", got, err)
	}
}

func TestWithExtraParameters_Merges(t *testing.T) {
	extra := types.MapValueMust(types.StringType, map[string]attr.Value{"targetTags": types.StringValue(`["web"]`)})
	got, err := WithExtraParameters(map[string]string{"version": "16"}, extra)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Text)
	}
	if got["version"] != "16" || got["targetTags"] != `["web"]` {
		t.Errorf("got %v", got)
	}
}

func TestWithExtraParameters_MergesIntoNil(t *testing.T) {
	extra := types.MapValueMust(types.StringType, map[string]attr.Value{"namespace": types.StringValue("data")})
	got, err := WithExtraParameters(nil, extra)
	if err != nil || got["namespace"] != "data" {
		t.Errorf("got %v, %v", got, err)
	}
}

func TestWithExtraParameters_RejectsKeySetByAttribute(t *testing.T) {
	extra := types.MapValueMust(types.StringType, map[string]attr.Value{"version": types.StringValue("17")})
	_, err := WithExtraParameters(map[string]string{"version": "16"}, extra)
	if err == nil || !strings.Contains(err.Text, "extra_parameters.version is already set") {
		t.Errorf("err = %v, want a conflict error", err)
	}
}
