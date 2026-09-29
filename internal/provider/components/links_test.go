package components

import (
	"math/big"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func linkTarget(t *testing.T, id string) types.Object {
	t.Helper()
	obj, err := BuildComponent(id, "NetworkAndCompute.IaaS.VirtualMachine", types.StringNull(), types.StringNull(), types.StringNull(), nil, nil, nil)
	if err != nil {
		t.Fatalf("building target: %s", err.Text)
	}
	return obj
}

func objectOf(t *testing.T, attrs map[string]attr.Value) types.Object {
	t.Helper()
	attrTypes := make(map[string]attr.Type, len(attrs))
	for k, v := range attrs {
		attrTypes[k] = v.Type(t.Context())
	}
	obj, diags := types.ObjectValue(attrTypes, attrs)
	if diags.HasError() {
		t.Fatalf("building object: %v", diags)
	}
	return obj
}

func tupleOf(t *testing.T, elems ...attr.Value) types.Dynamic {
	t.Helper()
	elemTypes := make([]attr.Type, len(elems))
	for i, e := range elems {
		elemTypes[i] = e.Type(t.Context())
	}
	tuple, diags := types.TupleValue(elemTypes, elems)
	if diags.HasError() {
		t.Fatalf("building tuple: %v", diags)
	}
	return types.DynamicValue(tuple)
}

func TestLinksFromDynamic_NullIsNoLinks(t *testing.T) {
	links, err := LinksFromDynamic(types.DynamicNull())
	if err != nil || links != nil {
		t.Errorf("got %v, %v; want no links and no error", links, err)
	}
}

// A tuple is what Terraform passes for `[ {...}, {...} ]` when the elements
// differ in shape, as they do when only some links carry settings.
func TestLinksFromDynamic_MixedShapes(t *testing.T) {
	withoutSettings := objectOf(t, map[string]attr.Value{"target": linkTarget(t, "sg")})
	withSettings := objectOf(t, map[string]attr.Value{
		"target": linkTarget(t, "api"),
		"settings": objectOf(t, map[string]attr.Value{
			"fromPort": types.NumberValue(big.NewFloat(8080)),
			"protocol": types.StringValue("tcp"),
			"enabled":  types.BoolValue(true),
		}),
	})

	links, err := LinksFromDynamic(tupleOf(t, withoutSettings, withSettings))
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Text)
	}
	if len(links) != 2 {
		t.Fatalf("len(links) = %d, want 2", len(links))
	}
	if links[0].ComponentId != "sg" || links[0].Settings != nil {
		t.Errorf("links[0] = %+v, want target sg without settings", links[0])
	}
	want := map[string]string{"fromPort": "8080", "protocol": "tcp", "enabled": "true"}
	for k, v := range want {
		if links[1].Settings[k] != v {
			t.Errorf("links[1].settings[%s] = %q, want %q", k, links[1].Settings[k], v)
		}
	}
}

func TestLinksFromDynamic_AcceptsTypedListWithMapSettings(t *testing.T) {
	settings := types.MapValueMust(types.StringType, map[string]attr.Value{"access": types.StringValue("read")})
	link := objectOf(t, map[string]attr.Value{"target": linkTarget(t, "bucket"), "settings": settings})
	list := types.ListValueMust(link.Type(t.Context()), []attr.Value{link})

	links, err := LinksFromDynamic(types.DynamicValue(list))
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Text)
	}
	if len(links) != 1 || links[0].Settings["access"] != "read" {
		t.Errorf("links = %+v, want one link with access=read", links)
	}
}

func TestLinksFromDynamic_Errors(t *testing.T) {
	tests := []struct {
		name  string
		links func(t *testing.T) types.Dynamic
		want  string
	}{
		{
			"not a list",
			func(t *testing.T) types.Dynamic { return types.DynamicValue(types.StringValue("x")) },
			"links must be a list of objects",
		},
		{
			"missing target",
			func(t *testing.T) types.Dynamic {
				return tupleOf(t, objectOf(t, map[string]attr.Value{"settings": types.MapNull(types.StringType)}))
			},
			"links[0]: target is required",
		},
		{
			"unexpected attribute",
			func(t *testing.T) types.Dynamic {
				return tupleOf(t, objectOf(t, map[string]attr.Value{"target": linkTarget(t, "a"), "port": types.StringValue("80")}))
			},
			`unexpected attribute "port"`,
		},
		{
			"target is not a component",
			func(t *testing.T) types.Dynamic {
				return tupleOf(t, objectOf(t, map[string]attr.Value{"target": types.StringValue("a")}))
			},
			"target must be a component object",
		},
		{
			"null setting value",
			func(t *testing.T) types.Dynamic {
				return tupleOf(t, objectOf(t, map[string]attr.Value{
					"target":   linkTarget(t, "a"),
					"settings": objectOf(t, map[string]attr.Value{"access": types.StringNull()}),
				}))
			},
			"settings.access: must be a known, non-null value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LinksFromDynamic(tt.links(t))
			if err == nil {
				t.Fatal("expected an error")
			}
			if !strings.Contains(err.Text, tt.want) {
				t.Errorf("error = %q, want it to contain %q", err.Text, tt.want)
			}
		})
	}
}

func TestLinksFromDynamic_StructuredSettingsBecomeJSON(t *testing.T) {
	redirectUris, _ := types.TupleValue(
		[]attr.Type{types.StringType, types.StringType},
		[]attr.Value{types.StringValue("https://a/cb"), types.StringValue("https://b/cb")},
	)
	scope := objectOf(t, map[string]attr.Value{"path": types.StringValue("/raw"), "recursive": types.BoolValue(true)})
	link := objectOf(t, map[string]attr.Value{
		"target": linkTarget(t, "idp"),
		"settings": objectOf(t, map[string]attr.Value{
			"clientType":   types.StringValue("web"),
			"redirectUris": redirectUris,
			"scope":        scope,
		}),
	})

	links, err := LinksFromDynamic(tupleOf(t, link))
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Text)
	}
	got := links[0].Settings
	if got["clientType"] != "web" {
		t.Errorf("clientType = %q", got["clientType"])
	}
	if want := `["https://a/cb","https://b/cb"]`; got["redirectUris"] != want {
		t.Errorf("redirectUris = %s, want %s", got["redirectUris"], want)
	}
	if want := `{"path":"/raw","recursive":true}`; got["scope"] != want {
		t.Errorf("scope = %s, want %s", got["scope"], want)
	}
}
