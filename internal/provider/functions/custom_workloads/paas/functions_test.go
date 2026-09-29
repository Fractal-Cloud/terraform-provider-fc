package paas

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
	ft "fractal.cloud/terraform-provider-fc/internal/provider/functions/functiontest"
)

func component(t *testing.T, id, componentType string) types.Object {
	t.Helper()
	obj, err := components.BuildComponent(id, componentType, types.StringNull(), types.StringNull(), types.StringNull(), nil, nil, nil)
	if err != nil {
		t.Fatalf("building %s: %s", id, err.Text)
	}
	return obj
}

func TestWorkloadFunction_Metadata(t *testing.T) {
	resp := &function.MetadataResponse{}
	NewWorkloadFunction().Metadata(context.Background(), function.MetadataRequest{}, resp)
	if resp.Name != "custom_workloads_paas_workload" {
		t.Errorf("name = %q, want %q", resp.Name, "custom_workloads_paas_workload")
	}
}

func TestWorkloadFunction_Run_Minimal(t *testing.T) {
	c := ft.Component(t, ft.Run(t, NewWorkloadFunction(), ft.Object(t, map[string]attr.Value{
		"id": types.StringValue("api"),
	})))
	if got := c["type"].(types.String).ValueString(); got != "CustomWorkloads.PaaS.Workload" {
		t.Errorf("type = %q, want %q", got, "CustomWorkloads.PaaS.Workload")
	}
	ft.ExpectParameters(t, c, map[string]string{})
}

func TestWorkloadFunction_Run_LinksAndSecurityGroups(t *testing.T) {
	db := component(t, "orders-db", "Storage.PaaS.RelationalDatabase")
	sg := component(t, "web-sg", "NetworkAndCompute.IaaS.SecurityGroup")
	c := ft.Component(t, ft.Run(t, NewWorkloadFunction(), ft.Object(t, map[string]attr.Value{
		"id": types.StringValue("api"),
		"links": ft.Tuple(t, ft.Object(t, map[string]attr.Value{
			"target":   db,
			"settings": ft.Object(t, map[string]attr.Value{"access": types.StringValue("read-write")}),
		})),
		"security_groups": types.ListValueMust(components.ComponentObjectType, []attr.Value{sg}),
	})))

	links := ft.Links(t, c)
	if links["orders-db"]["access"] != "read-write" {
		t.Errorf("links = %v, want orders-db with access=read-write", links)
	}
	if _, ok := links["web-sg"]; !ok {
		t.Errorf("links = %v, want a web-sg membership link", links)
	}
}

func TestWorkloadFunction_Run_RejectsWrongSubnetType(t *testing.T) {
	resp := ft.Run(t, NewWorkloadFunction(), ft.Object(t, map[string]attr.Value{
		"id":     types.StringValue("api"),
		"subnet": component(t, "vpc", "NetworkAndCompute.IaaS.VirtualNetwork"),
	}))
	if resp.Error == nil || !strings.Contains(resp.Error.Text, "NetworkAndCompute.IaaS.Subnet") {
		t.Errorf("error = %v, want a subnet type error", resp.Error)
	}
}

func TestWorkloadFunction_Run_WritesContainerOfferKeys(t *testing.T) {
	c := ft.Component(t, ft.Run(t, NewWorkloadFunction(), ft.Object(t, map[string]attr.Value{
		"id":              types.StringValue("api"),
		"container_image": types.StringValue("gcr.io/acme/api:1"),
		"container_port":  types.Int64Value(8080),
		"cpu":             types.StringValue("1"),
		"memory":          types.StringValue("512Mi"),
	})))
	ft.ExpectParameters(t, c, map[string]string{
		"image":  "gcr.io/acme/api:1",
		"port":   "8080",
		"cpu":    "1",
		"memory": "512Mi",
	})
}
