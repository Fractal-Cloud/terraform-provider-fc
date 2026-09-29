package paas

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
	"fractal.cloud/terraform-provider-fc/internal/provider/functions/functiontest"
)

func buildTestComponent(t *testing.T, id, componentType string) types.Object {
	t.Helper()
	obj, err := components.BuildComponent(id, componentType, types.StringNull(), types.StringNull(), types.StringNull(), nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to build test component: %s", err.Text)
	}
	return obj
}

func runFunction(t *testing.T, f function.Function, args []attr.Value) *function.RunResponse {
	t.Helper()
	ctx := context.Background()
	req := function.RunRequest{
		Arguments: function.NewArgumentsData(functiontest.Complete(t, f, args)),
	}
	resp := &function.RunResponse{
		Result: function.NewResultData(types.ObjectNull(components.ComponentAttrTypes)),
	}
	f.Run(ctx, req, resp)
	return resp
}

func getResultAttrs(t *testing.T, resp *function.RunResponse) map[string]attr.Value {
	t.Helper()
	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error.Text)
	}
	result := resp.Result.Value()
	obj, ok := result.(types.Object)
	if !ok {
		t.Fatalf("expected types.Object result, got %T", result)
	}
	return obj.Attributes()
}

// simpleConfig builds a config with just id, display_name, description
func simpleConfig(t *testing.T, id string) types.Object {
	t.Helper()
	obj, diags := types.ObjectValue(map[string]attr.Type{
		"id": types.StringType, "display_name": types.StringType, "description": types.StringType,
	}, map[string]attr.Value{
		"id": types.StringValue(id), "display_name": types.StringNull(), "description": types.StringNull(),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}
	return obj
}

// dbmsChildConfig builds a config with id, display_name, description, dbms
func dbmsChildConfig(t *testing.T, id string, dbms attr.Value) types.Object {
	t.Helper()
	obj, diags := types.ObjectValue(map[string]attr.Type{
		"id": types.StringType, "display_name": types.StringType, "description": types.StringType,
		"dbms": components.ComponentObjectType,
	}, map[string]attr.Value{
		"id": types.StringValue(id), "display_name": types.StringNull(), "description": types.StringNull(),
		"dbms": dbms,
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}
	return obj
}

// --- FilesAndBlobs ---

func TestFilesAndBlobsFunction_Metadata(t *testing.T) {
	f := NewStoragePaasFilesAndBlobsFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_files_and_blobs" {
		t.Errorf("expected name %q, got %q", "storage_paas_files_and_blobs", resp.Name)
	}
}

func TestFilesAndBlobsFunction_Definition(t *testing.T) {
	f := NewStoragePaasFilesAndBlobsFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

func TestFilesAndBlobsFunction_Run(t *testing.T) {
	f := NewStoragePaasFilesAndBlobsFunction()
	resp := runFunction(t, f, []attr.Value{simpleConfig(t, "blob-1")})
	attrs := getResultAttrs(t, resp)
	if attrs["type"].(types.String).ValueString() != "Storage.PaaS.FilesAndBlobs" {
		t.Errorf("expected type %q", "Storage.PaaS.FilesAndBlobs")
	}
}

// --- RelationalDbms ---

func TestRelationalDbmsFunction_Metadata(t *testing.T) {
	f := NewStoragePaasRelationalDbmsFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_relational_dbms" {
		t.Errorf("expected name %q, got %q", "storage_paas_relational_dbms", resp.Name)
	}
}

func TestRelationalDbmsFunction_Definition(t *testing.T) {
	f := NewStoragePaasRelationalDbmsFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

func TestRelationalDbmsFunction_Run(t *testing.T) {
	f := NewStoragePaasRelationalDbmsFunction()
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id": types.StringType, "display_name": types.StringType, "description": types.StringType,
		"engine_version": types.StringType,
	}, map[string]attr.Value{
		"id": types.StringValue("rds-1"), "display_name": types.StringNull(), "description": types.StringNull(),
		"engine_version": types.StringValue("15.4"),
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)
	if attrs["type"].(types.String).ValueString() != "Storage.PaaS.RelationalDbms" {
		t.Errorf("expected type %q", "Storage.PaaS.RelationalDbms")
	}
	params := attrs["parameters"].(types.Map)
	if params.Elements()["version"].(types.String).ValueString() != "15.4" {
		t.Errorf("expected version param %q", "15.4")
	}
}

// --- RelationalDatabase ---

func TestRelationalDatabaseFunction_Metadata(t *testing.T) {
	f := NewStoragePaasRelationalDatabaseFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_relational_database" {
		t.Errorf("expected name %q, got %q", "storage_paas_relational_database", resp.Name)
	}
}

func TestRelationalDatabaseFunction_Definition(t *testing.T) {
	f := NewStoragePaasRelationalDatabaseFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

func TestRelationalDatabaseFunction_Run(t *testing.T) {
	f := NewStoragePaasRelationalDatabaseFunction()
	dbms := buildTestComponent(t, "rds-1", "Storage.PaaS.RelationalDbms")
	configObj, diags := types.ObjectValue(map[string]attr.Type{
		"id": types.StringType, "display_name": types.StringType, "description": types.StringType,
		"collation": types.StringType, "charset": types.StringType, "dbms": components.ComponentObjectType,
	}, map[string]attr.Value{
		"id": types.StringValue("db-1"), "display_name": types.StringNull(), "description": types.StringNull(),
		"collation": types.StringValue("en_US.utf8"), "charset": types.StringValue("utf8"), "dbms": dbms,
	})
	if diags.HasError() {
		t.Fatalf("failed to build config: %s", diags.Errors())
	}

	resp := runFunction(t, f, []attr.Value{configObj})
	attrs := getResultAttrs(t, resp)
	if attrs["type"].(types.String).ValueString() != "Storage.PaaS.RelationalDatabase" {
		t.Errorf("expected type %q", "Storage.PaaS.RelationalDatabase")
	}
	deps := attrs["dependencies_ids"].(types.List)
	if deps.IsNull() || deps.Elements()[0].(types.String).ValueString() != "rds-1" {
		t.Errorf("expected dependency %q", "rds-1")
	}
	params := attrs["parameters"].(types.Map)
	if params.Elements()["collation"].(types.String).ValueString() != "en_US.utf8" {
		t.Errorf("expected collation param %q", "en_US.utf8")
	}
}

// --- DocumentDbms ---

func TestDocumentDbmsFunction_Metadata(t *testing.T) {
	f := NewStoragePaasDocumentDbmsFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_document_dbms" {
		t.Errorf("expected name %q, got %q", "storage_paas_document_dbms", resp.Name)
	}
}

func TestDocumentDbmsFunction_Definition(t *testing.T) {
	f := NewStoragePaasDocumentDbmsFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

func TestDocumentDbmsFunction_Run(t *testing.T) {
	f := NewStoragePaasDocumentDbmsFunction()
	resp := runFunction(t, f, []attr.Value{simpleConfig(t, "docdb-1")})
	attrs := getResultAttrs(t, resp)
	if attrs["type"].(types.String).ValueString() != "Storage.PaaS.DocumentDbms" {
		t.Errorf("expected type %q", "Storage.PaaS.DocumentDbms")
	}
}

// --- DocumentDatabase ---

func TestDocumentDatabaseFunction_Metadata(t *testing.T) {
	f := NewStoragePaasDocumentDatabaseFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_document_database" {
		t.Errorf("expected name %q, got %q", "storage_paas_document_database", resp.Name)
	}
}

func TestDocumentDatabaseFunction_Definition(t *testing.T) {
	f := NewStoragePaasDocumentDatabaseFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

func TestDocumentDatabaseFunction_Run(t *testing.T) {
	f := NewStoragePaasDocumentDatabaseFunction()
	dbms := buildTestComponent(t, "docdb-1", "Storage.PaaS.DocumentDbms")
	resp := runFunction(t, f, []attr.Value{dbmsChildConfig(t, "doc-1", dbms)})
	attrs := getResultAttrs(t, resp)
	if attrs["type"].(types.String).ValueString() != "Storage.PaaS.DocumentDatabase" {
		t.Errorf("expected type %q", "Storage.PaaS.DocumentDatabase")
	}
	deps := attrs["dependencies_ids"].(types.List)
	if deps.IsNull() || deps.Elements()[0].(types.String).ValueString() != "docdb-1" {
		t.Errorf("expected dependency %q", "docdb-1")
	}
}

// --- ColumnOrientedDbms ---

func TestColumnOrientedDbmsFunction_Metadata(t *testing.T) {
	f := NewStoragePaasColumnOrientedDbmsFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_column_oriented_dbms" {
		t.Errorf("expected name %q, got %q", "storage_paas_column_oriented_dbms", resp.Name)
	}
}

func TestColumnOrientedDbmsFunction_Definition(t *testing.T) {
	f := NewStoragePaasColumnOrientedDbmsFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

func TestColumnOrientedDbmsFunction_Run(t *testing.T) {
	f := NewStoragePaasColumnOrientedDbmsFunction()
	resp := runFunction(t, f, []attr.Value{simpleConfig(t, "col-1")})
	attrs := getResultAttrs(t, resp)
	if attrs["type"].(types.String).ValueString() != "Storage.PaaS.ColumnOrientedDbms" {
		t.Errorf("expected type %q", "Storage.PaaS.ColumnOrientedDbms")
	}
}

// --- ColumnOrientedEntity ---

func TestColumnOrientedEntityFunction_Metadata(t *testing.T) {
	f := NewStoragePaasColumnOrientedEntityFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_column_oriented_entity" {
		t.Errorf("expected name %q, got %q", "storage_paas_column_oriented_entity", resp.Name)
	}
}

func TestColumnOrientedEntityFunction_Definition(t *testing.T) {
	f := NewStoragePaasColumnOrientedEntityFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

func TestColumnOrientedEntityFunction_Run(t *testing.T) {
	f := NewStoragePaasColumnOrientedEntityFunction()
	dbms := buildTestComponent(t, "col-1", "Storage.PaaS.ColumnOrientedDbms")
	resp := runFunction(t, f, []attr.Value{dbmsChildConfig(t, "entity-1", dbms)})
	attrs := getResultAttrs(t, resp)
	if attrs["type"].(types.String).ValueString() != "Storage.PaaS.ColumnOrientedEntity" {
		t.Errorf("expected type %q", "Storage.PaaS.ColumnOrientedEntity")
	}
	deps := attrs["dependencies_ids"].(types.List)
	if deps.IsNull() || deps.Elements()[0].(types.String).ValueString() != "col-1" {
		t.Errorf("expected dependency %q", "col-1")
	}
}

// --- KeyValueDbms ---

func TestKeyValueDbmsFunction_Metadata(t *testing.T) {
	f := NewStoragePaasKeyValueDbmsFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_key_value_dbms" {
		t.Errorf("expected name %q, got %q", "storage_paas_key_value_dbms", resp.Name)
	}
}

func TestKeyValueDbmsFunction_Definition(t *testing.T) {
	f := NewStoragePaasKeyValueDbmsFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

func TestKeyValueDbmsFunction_Run(t *testing.T) {
	f := NewStoragePaasKeyValueDbmsFunction()
	resp := runFunction(t, f, []attr.Value{simpleConfig(t, "kv-1")})
	attrs := getResultAttrs(t, resp)
	if attrs["type"].(types.String).ValueString() != "Storage.PaaS.KeyValueDbms" {
		t.Errorf("expected type %q", "Storage.PaaS.KeyValueDbms")
	}
}

// --- KeyValueEntity ---

func TestKeyValueEntityFunction_Metadata(t *testing.T) {
	f := NewStoragePaasKeyValueEntityFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_key_value_entity" {
		t.Errorf("expected name %q, got %q", "storage_paas_key_value_entity", resp.Name)
	}
}

func TestKeyValueEntityFunction_Definition(t *testing.T) {
	f := NewStoragePaasKeyValueEntityFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

func TestKeyValueEntityFunction_Run(t *testing.T) {
	f := NewStoragePaasKeyValueEntityFunction()
	dbms := buildTestComponent(t, "kv-1", "Storage.PaaS.KeyValueDbms")
	resp := runFunction(t, f, []attr.Value{dbmsChildConfig(t, "entity-1", dbms)})
	attrs := getResultAttrs(t, resp)
	if attrs["type"].(types.String).ValueString() != "Storage.PaaS.KeyValueEntity" {
		t.Errorf("expected type %q", "Storage.PaaS.KeyValueEntity")
	}
	deps := attrs["dependencies_ids"].(types.List)
	if deps.IsNull() || deps.Elements()[0].(types.String).ValueString() != "kv-1" {
		t.Errorf("expected dependency %q", "kv-1")
	}
}

// --- GraphDbms ---

func TestGraphDbmsFunction_Metadata(t *testing.T) {
	f := NewStoragePaasGraphDbmsFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_graph_dbms" {
		t.Errorf("expected name %q, got %q", "storage_paas_graph_dbms", resp.Name)
	}
}

func TestGraphDbmsFunction_Definition(t *testing.T) {
	f := NewStoragePaasGraphDbmsFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

func TestGraphDbmsFunction_Run(t *testing.T) {
	f := NewStoragePaasGraphDbmsFunction()
	resp := runFunction(t, f, []attr.Value{simpleConfig(t, "graph-1")})
	attrs := getResultAttrs(t, resp)
	if attrs["type"].(types.String).ValueString() != "Storage.PaaS.GraphDbms" {
		t.Errorf("expected type %q", "Storage.PaaS.GraphDbms")
	}
}

// --- GraphDatabase ---

func TestGraphDatabaseFunction_Metadata(t *testing.T) {
	f := NewStoragePaasGraphDatabaseFunction()
	req := function.MetadataRequest{}
	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), req, resp)
	if resp.Name != "storage_paas_graph_database" {
		t.Errorf("expected name %q, got %q", "storage_paas_graph_database", resp.Name)
	}
}

func TestGraphDatabaseFunction_Definition(t *testing.T) {
	f := NewStoragePaasGraphDatabaseFunction()
	req := function.DefinitionRequest{}
	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), req, resp)
	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}
	if resp.Definition.Return == nil {
		t.Error("expected non-nil return type")
	}
}

func TestGraphDatabaseFunction_Run(t *testing.T) {
	f := NewStoragePaasGraphDatabaseFunction()
	dbms := buildTestComponent(t, "graph-1", "Storage.PaaS.GraphDbms")
	resp := runFunction(t, f, []attr.Value{dbmsChildConfig(t, "db-1", dbms)})
	attrs := getResultAttrs(t, resp)
	if attrs["type"].(types.String).ValueString() != "Storage.PaaS.GraphDatabase" {
		t.Errorf("expected type %q", "Storage.PaaS.GraphDatabase")
	}
	deps := attrs["dependencies_ids"].(types.List)
	if deps.IsNull() || deps.Elements()[0].(types.String).ValueString() != "graph-1" {
		t.Errorf("expected dependency %q", "graph-1")
	}
}
