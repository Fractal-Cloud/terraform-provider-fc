package caas

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
	if resp.Name != "custom_workloads_caas_workload" {
		t.Errorf("name = %q, want %q", resp.Name, "custom_workloads_caas_workload")
	}
}

func TestWorkloadFunction_Run_Minimal(t *testing.T) {
	c := ft.Component(t, ft.Run(t, NewWorkloadFunction(), ft.Object(t, map[string]attr.Value{
		"id": types.StringValue("api"),
	})))
	if got := c["type"].(types.String).ValueString(); got != "CustomWorkloads.CaaS.Workload" {
		t.Errorf("type = %q, want %q", got, "CustomWorkloads.CaaS.Workload")
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

func TestWorkloadFunction_Run_WritesEveryCaaSOfferKey(t *testing.T) {
	platform := component(t, "k8s", "NetworkAndCompute.PaaS.ContainerPlatform")
	subnet := component(t, "private", "NetworkAndCompute.IaaS.Subnet")
	c := ft.Component(t, ft.Run(t, NewWorkloadFunction(), ft.Object(t, map[string]attr.Value{
		"id":              types.StringValue("api"),
		"container_image": types.StringValue("ghcr.io/acme/api:1.2"),
		"container_port":  types.Int64Value(8080),
		"container_name":  types.StringValue("api"),
		"cpu":             types.StringValue("500m"),
		"memory":          types.StringValue("512Mi"),
		"replicas":        types.Int64Value(3),
		"platform":        platform,
		"subnet":          subnet,
		"extra_parameters": types.MapValueMust(types.StringType, map[string]attr.Value{
			"namespace": types.StringValue("orders"),
		}),
	})))

	ft.ExpectParameters(t, c, map[string]string{
		"containerImage": "ghcr.io/acme/api:1.2",
		"image":          "ghcr.io/acme/api:1.2",
		"containerPort":  "8080",
		"port":           "8080",
		"containerName":  "api",
		"cpu":            "500m",
		"memory":         "512Mi",
		"replicas":       "3",
		"desiredCount":   "3",
		"namespace":      "orders",
	})
	if deps := ft.Strings(t, c, "dependencies_ids"); len(deps) != 2 || deps[0] != "k8s" || deps[1] != "private" {
		t.Errorf("dependencies = %v, want [k8s private]", deps)
	}
}

func TestWorkloadFunction_Run_RejectsExtraParameterSetByAttribute(t *testing.T) {
	resp := ft.Run(t, NewWorkloadFunction(), ft.Object(t, map[string]attr.Value{
		"id":       types.StringValue("api"),
		"replicas": types.Int64Value(2),
		"extra_parameters": types.MapValueMust(types.StringType, map[string]attr.Value{
			"replicas": types.StringValue("5"),
		}),
	}))
	if resp.Error == nil || !strings.Contains(resp.Error.Text, "extra_parameters.replicas") {
		t.Errorf("error = %v, want a conflict on replicas", resp.Error)
	}
}
