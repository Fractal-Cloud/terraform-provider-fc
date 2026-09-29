package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"fractal.cloud/terraform-provider-fc/internal/client"
)

func TestSameJSON(t *testing.T) {
	tests := []struct {
		a, b string
		want bool
	}{
		{"x", "x", true},
		{"8080", "8080.0", false},
		{`{"b":1,"a":[1,2]}`, `{"a":[1,2],"b":1}`, true},
		{`[{"protocol":"tcp","fromPort":443}]`, `[ { "fromPort": 443, "protocol": "tcp" } ]`, true},
		{`{"a":1}`, `{"a":2}`, false},
		{`[1,2]`, `[2,1]`, false},
		{`{"a":1}`, `{"a":1} trailing`, false},
		{"{not json", "{not json ", false},
	}
	for _, tt := range tests {
		if got := sameJSON(tt.a, tt.b); got != tt.want {
			t.Errorf("sameJSON(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

// A refresh must not report a diff when the API returns a structured value
// re-encoded, and must when the value really changed.
func TestMapBlueprintToState_KeepsConfiguredJSONText(t *testing.T) {
	ctx := context.Background()
	configured := `[{"protocol":"tcp","fromPort":443,"toPort":443,"sourceCidr":"0.0.0.0/0"}]`
	reencoded := `[{"fromPort":443,"protocol":"tcp","sourceCidr":"0.0.0.0/0","toPort":443}]`
	secret := `{"$envSecret":"api-key"}`

	blueprint := func(rules, token, retention string) *fractalCloud.Blueprint {
		return &fractalCloud.Blueprint{Components: []fractalCloud.Component{{
			Id:         "sg",
			Type:       "NetworkAndCompute.IaaS.SecurityGroup",
			Parameters: map[string]string{"ingressRules": rules, "retention": retention},
			Links:      []fractalCloud.ComponentLink{{ComponentId: "api", Settings: map[string]string{"token": token}}},
		}}}
	}

	model := &BlueprintModel{Components: types.ListNull(basetypes.ObjectType{AttrTypes: componentAttrTypes})}
	var diags Diagnostics
	mapBlueprintToState(ctx, blueprint(configured, secret, "7"), model, &diags)
	mapBlueprintToState(ctx, blueprint(reencoded, `{ "$envSecret" : "api-key" }`, "14"), model, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags.Errors())
	}

	var comps []ComponentModel
	model.Components.ElementsAs(ctx, &comps, false)
	params := stringMap(ctx, comps[0].Parameters)
	if params["ingressRules"] != configured {
		t.Errorf("ingressRules = %s, want the configured text %s", params["ingressRules"], configured)
	}
	if params["retention"] != "14" {
		t.Errorf("retention = %q, want the changed value 14", params["retention"])
	}
	links := priorLinks(ctx, comps[0].Links)
	if links[0].settings["token"] != secret {
		t.Errorf("link token = %s, want the configured text %s", links[0].settings["token"], secret)
	}
}
