package faas

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
	if resp.Name != "custom_workloads_faas_workload" {
		t.Errorf("name = %q, want %q", resp.Name, "custom_workloads_faas_workload")
	}
}

func TestWorkloadFunction_Run_Minimal(t *testing.T) {
	c := ft.Component(t, ft.Run(t, NewWorkloadFunction(), ft.Object(t, map[string]attr.Value{
		"id":      types.StringValue("api"),
		"runtime": types.StringValue("python3.12"),
		"handler": types.StringValue("app.handler"),
	})))
	if got := c["type"].(types.String).ValueString(); got != "CustomWorkloads.FaaS.Workload" {
		t.Errorf("type = %q, want %q", got, "CustomWorkloads.FaaS.Workload")
	}
	ft.ExpectParameters(t, c, map[string]string{"runtime": "python3.12", "handler": "app.handler", "entryPoint": "app.handler"})
}

func TestWorkloadFunction_Run_LinksAndSecurityGroups(t *testing.T) {
	db := component(t, "orders-db", "Storage.PaaS.RelationalDatabase")
	sg := component(t, "web-sg", "NetworkAndCompute.IaaS.SecurityGroup")
	c := ft.Component(t, ft.Run(t, NewWorkloadFunction(), ft.Object(t, map[string]attr.Value{
		"id":      types.StringValue("api"),
		"runtime": types.StringValue("python3.12"),
		"handler": types.StringValue("app.handler"),
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
		"id":      types.StringValue("api"),
		"runtime": types.StringValue("python3.12"),
		"handler": types.StringValue("app.handler"),
		"subnet":  component(t, "vpc", "NetworkAndCompute.IaaS.VirtualNetwork"),
	}))
	if resp.Error == nil || !strings.Contains(resp.Error.Text, "NetworkAndCompute.IaaS.Subnet") {
		t.Errorf("error = %v, want a subnet type error", resp.Error)
	}
}

func TestWorkloadFunction_Run_WritesFunctionParameters(t *testing.T) {
	c := ft.Component(t, ft.Run(t, NewWorkloadFunction(), ft.Object(t, map[string]attr.Value{
		"id":              types.StringValue("resize"),
		"runtime":         types.StringValue("nodejs20.x"),
		"handler":         types.StringValue("index.handler"),
		"memory_mb":       types.Int64Value(1024),
		"timeout_seconds": types.Int64Value(60),
	})))
	ft.ExpectParameters(t, c, map[string]string{
		"runtime":        "nodejs20.x",
		"handler":        "index.handler",
		"entryPoint":     "index.handler",
		"memoryMb":       "1024",
		"timeoutSeconds": "60",
	})
}
