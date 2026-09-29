package provider

import (
	"context"
	"testing"

	fractalCloud "fractal.cloud/terraform-provider-fc/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestStringValueOrNull(t *testing.T) {
	tests := []struct {
		name      string
		apiValue  string
		model     types.String
		wantValue string
		wantNull  bool
	}{
		{
			name:      "non-empty API value returns StringValue",
			apiValue:  "hello",
			model:     types.StringNull(),
			wantValue: "hello",
			wantNull:  false,
		},
		{
			name:      "empty API value with null model returns StringNull",
			apiValue:  "",
			model:     types.StringNull(),
			wantValue: "",
			wantNull:  true,
		},
		{
			name:      "empty API value with non-null model returns empty StringValue",
			apiValue:  "",
			model:     types.StringValue(""),
			wantValue: "",
			wantNull:  false,
		},
		{
			name:      "empty API value with populated model returns empty StringValue",
			apiValue:  "",
			model:     types.StringValue("previous"),
			wantValue: "",
			wantNull:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringValueOrNull(tt.apiValue, tt.model)
			if tt.wantNull {
				if !got.IsNull() {
					t.Errorf("expected null, got %q", got.ValueString())
				}
			} else {
				if got.IsNull() {
					t.Error("expected non-null, got null")
				}
				if got.ValueString() != tt.wantValue {
					t.Errorf("expected %q, got %q", tt.wantValue, got.ValueString())
				}
			}
		})
	}
}

func TestStringPointerToTFValue(t *testing.T) {
	strPtr := func(s string) *string { return &s }

	tests := []struct {
		name      string
		input     *string
		wantValue string
	}{
		{
			name:      "nil pointer returns empty string (not null)",
			input:     nil,
			wantValue: "",
		},
		{
			name:      "non-nil pointer with value",
			input:     strPtr("hello"),
			wantValue: "hello",
		},
		{
			name:      "non-nil pointer with empty string",
			input:     strPtr(""),
			wantValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringPointerToTFValue(tt.input)
			if got.IsNull() {
				t.Error("expected non-null value, got null")
			}
			if got.ValueString() != tt.wantValue {
				t.Errorf("expected %q, got %q", tt.wantValue, got.ValueString())
			}
		})
	}
}

func TestBoolPointerToTFValue(t *testing.T) {
	boolPtr := func(b bool) *bool { return &b }

	tests := []struct {
		name      string
		input     *bool
		wantValue bool
	}{
		{
			name:      "nil pointer returns false (not null)",
			input:     nil,
			wantValue: false,
		},
		{
			name:      "non-nil pointer with true",
			input:     boolPtr(true),
			wantValue: true,
		},
		{
			name:      "non-nil pointer with false",
			input:     boolPtr(false),
			wantValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := boolPointerToTFValue(tt.input)
			if got.IsNull() {
				t.Error("expected non-null value, got null")
			}
			if got.ValueBool() != tt.wantValue {
				t.Errorf("expected %v, got %v", tt.wantValue, got.ValueBool())
			}
		})
	}
}

func TestMapBlueprintToState(t *testing.T) {
	ctx := context.Background()
	strPtr := func(s string) *string { return &s }
	boolPtr := func(b bool) *bool { return &b }

	t.Run("maps full component correctly", func(t *testing.T) {
		blueprint := &fractalCloud.Blueprint{
			Description: "test blueprint",
			IsPrivate:   true,
			CreatedAt:   "2026-01-01",
			Components: []fractalCloud.Component{
				{
					Id:                "comp-1",
					Type:              "BigData.PaaS.DistributedDataProcessing",
					DisplayName:       strPtr("My Workspace"),
					Description:       strPtr("A workspace"),
					Version:           strPtr("v2"),
					IsLocked:          boolPtr(true),
					RecreateOnFailure: boolPtr(false),
					Parameters:        map[string]string{"region": "us-east-1", "sku": "premium"},
					DependenciesIds:   []string{"dep-1", "dep-2"},
					Links: []fractalCloud.ComponentLink{
						{ComponentId: "link-target", Settings: map[string]string{"fromPort": "8080"}},
					},
					OutputFields: []string{"workspace_url", "workspace_id"},
				},
			},
		}

		model := &BlueprintModel{
			Components: types.ListNull(basetypes.ObjectType{AttrTypes: componentAttrTypes}),
		}

		var diags Diagnostics
		mapBlueprintToState(ctx, blueprint, model, &diags)

		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags.Errors())
		}

		if model.Description.ValueString() != "test blueprint" {
			t.Errorf("description: expected %q, got %q", "test blueprint", model.Description.ValueString())
		}
		if model.IsPrivate.ValueBool() != true {
			t.Errorf("is_private: expected true, got false")
		}
		if model.CreatedAt.ValueString() != "2026-01-01" {
			t.Errorf("created_at: expected %q, got %q", "2026-01-01", model.CreatedAt.ValueString())
		}

		var components []ComponentModel
		d := model.Components.ElementsAs(ctx, &components, false)
		if d.HasError() {
			t.Fatalf("failed to extract components: %v", d.Errors())
		}
		if len(components) != 1 {
			t.Fatalf("expected 1 component, got %d", len(components))
		}

		comp := components[0]
		if comp.Id.ValueString() != "comp-1" {
			t.Errorf("id: expected %q, got %q", "comp-1", comp.Id.ValueString())
		}
		if comp.Type.ValueString() != "BigData.PaaS.DistributedDataProcessing" {
			t.Errorf("type: expected %q, got %q", "BigData.PaaS.DistributedDataProcessing", comp.Type.ValueString())
		}
		if comp.DisplayName.ValueString() != "My Workspace" {
			t.Errorf("display_name: expected %q, got %q", "My Workspace", comp.DisplayName.ValueString())
		}
		if comp.Description.ValueString() != "A workspace" {
			t.Errorf("description: expected %q, got %q", "A workspace", comp.Description.ValueString())
		}
		if comp.Version.ValueString() != "v2" {
			t.Errorf("version: expected %q, got %q", "v2", comp.Version.ValueString())
		}
		if comp.IsLocked.ValueBool() != true {
			t.Errorf("is_locked: expected true, got false")
		}
		if comp.RecreateOnFailure.ValueBool() != false {
			t.Errorf("recreate_on_failure: expected false, got true")
		}

		// Check parameters
		params := make(map[string]string)
		d = comp.Parameters.ElementsAs(ctx, &params, false)
		if d.HasError() {
			t.Fatalf("failed to extract parameters: %v", d.Errors())
		}
		if params["region"] != "us-east-1" {
			t.Errorf("parameters[region]: expected %q, got %q", "us-east-1", params["region"])
		}
		if params["sku"] != "premium" {
			t.Errorf("parameters[sku]: expected %q, got %q", "premium", params["sku"])
		}

		// Check dependencies
		var deps []string
		d = comp.DependenciesIds.ElementsAs(ctx, &deps, false)
		if d.HasError() {
			t.Fatalf("failed to extract dependencies: %v", d.Errors())
		}
		if len(deps) != 2 || deps[0] != "dep-1" || deps[1] != "dep-2" {
			t.Errorf("dependencies: expected [dep-1, dep-2], got %v", deps)
		}

		// Check links
		var links []LinkModel
		d = comp.Links.ElementsAs(ctx, &links, false)
		if d.HasError() {
			t.Fatalf("failed to extract links: %v", d.Errors())
		}
		if len(links) != 1 {
			t.Fatalf("expected 1 link, got %d", len(links))
		}
		if links[0].ComponentId.ValueString() != "link-target" {
			t.Errorf("link component_id: expected %q, got %q", "link-target", links[0].ComponentId.ValueString())
		}
		linkSettings := make(map[string]string)
		d = links[0].Settings.ElementsAs(ctx, &linkSettings, false)
		if d.HasError() {
			t.Fatalf("failed to extract link settings: %v", d.Errors())
		}
		if linkSettings["fromPort"] != "8080" {
			t.Errorf("link settings[fromPort]: expected %q, got %q", "8080", linkSettings["fromPort"])
		}

		// Check output fields
		var outFields []string
		d = comp.OutputFields.ElementsAs(ctx, &outFields, false)
		if d.HasError() {
			t.Fatalf("failed to extract output fields: %v", d.Errors())
		}
		if len(outFields) != 2 || outFields[0] != "workspace_url" || outFields[1] != "workspace_id" {
			t.Errorf("output_fields: expected [workspace_url, workspace_id], got %v", outFields)
		}
	})

	t.Run("nil slices and maps produce empty collections not null", func(t *testing.T) {
		blueprint := &fractalCloud.Blueprint{
			Description: "empty test",
			Components: []fractalCloud.Component{
				{
					Id:              "comp-empty",
					Type:            "NetworkAndCompute.IaaS.AwsVpc",
					Parameters:      nil,
					DependenciesIds: nil,
					Links:           nil,
					OutputFields:    nil,
				},
			},
		}

		model := &BlueprintModel{
			Components: types.ListNull(basetypes.ObjectType{AttrTypes: componentAttrTypes}),
		}

		var diags Diagnostics
		mapBlueprintToState(ctx, blueprint, model, &diags)

		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags.Errors())
		}

		var components []ComponentModel
		d := model.Components.ElementsAs(ctx, &components, false)
		if d.HasError() {
			t.Fatalf("failed to extract components: %v", d.Errors())
		}

		comp := components[0]

		// Parameters should be empty map, not null
		if comp.Parameters.IsNull() {
			t.Error("parameters should be empty map, not null")
		}
		params := make(map[string]string)
		d = comp.Parameters.ElementsAs(ctx, &params, false)
		if d.HasError() {
			t.Fatalf("failed to extract parameters: %v", d.Errors())
		}
		if len(params) != 0 {
			t.Errorf("expected 0 parameters, got %d", len(params))
		}

		// Dependencies should be empty list, not null
		if comp.DependenciesIds.IsNull() {
			t.Error("dependencies_ids should be empty list, not null")
		}
		var deps []string
		d = comp.DependenciesIds.ElementsAs(ctx, &deps, false)
		if d.HasError() {
			t.Fatalf("failed to extract dependencies: %v", d.Errors())
		}
		if len(deps) != 0 {
			t.Errorf("expected 0 dependencies, got %d", len(deps))
		}

		// Links should be empty list, not null
		if comp.Links.IsNull() {
			t.Error("links should be empty list, not null")
		}

		// Output fields should be empty list, not null
		if comp.OutputFields.IsNull() {
			t.Error("output_fields should be empty list, not null")
		}
	})

	t.Run("version preserved from prior state when API returns empty", func(t *testing.T) {
		// Build a prior state model with a component that has version "v1"
		priorComponent := ComponentModel{
			Id:                types.StringValue("comp-versioned"),
			Type:              types.StringValue("BigData.PaaS.ComputeCluster"),
			DisplayName:       types.StringValue(""),
			Description:       types.StringValue(""),
			Version:           types.StringValue("v1"),
			IsLocked:          types.BoolValue(false),
			RecreateOnFailure: types.BoolValue(false),
		}

		emptyParams, _ := types.MapValueFrom(ctx, types.StringType, map[string]string{})
		priorComponent.Parameters = emptyParams

		emptyList, _ := types.ListValueFrom(ctx, types.StringType, []string{})
		priorComponent.DependenciesIds = emptyList
		priorComponent.OutputFields = emptyList

		emptyLinks, _ := types.ListValueFrom(ctx, basetypes.ObjectType{AttrTypes: linkAttrTypes}, []LinkModel{})
		priorComponent.Links = emptyLinks

		priorComponents, _ := types.ListValueFrom(ctx, basetypes.ObjectType{AttrTypes: componentAttrTypes}, []ComponentModel{priorComponent})

		model := &BlueprintModel{
			Components: priorComponents,
		}

		blueprint := &fractalCloud.Blueprint{
			Description: "version test",
			Components: []fractalCloud.Component{
				{
					Id:      "comp-versioned",
					Type:    "BigData.PaaS.ComputeCluster",
					Version: nil, // API returns no version
				},
			},
		}

		var diags Diagnostics
		mapBlueprintToState(ctx, blueprint, model, &diags)

		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags.Errors())
		}

		var components []ComponentModel
		d := model.Components.ElementsAs(ctx, &components, false)
		if d.HasError() {
			t.Fatalf("failed to extract components: %v", d.Errors())
		}

		if components[0].Version.ValueString() != "v1" {
			t.Errorf("version: expected preserved %q, got %q", "v1", components[0].Version.ValueString())
		}
	})

	t.Run("version from API used when non-empty", func(t *testing.T) {
		// Prior state has "v1" but API returns "v3" — API wins
		priorComponent := ComponentModel{
			Id:                types.StringValue("comp-versioned"),
			Type:              types.StringValue("BigData.PaaS.ComputeCluster"),
			DisplayName:       types.StringValue(""),
			Description:       types.StringValue(""),
			Version:           types.StringValue("v1"),
			IsLocked:          types.BoolValue(false),
			RecreateOnFailure: types.BoolValue(false),
		}

		emptyParams, _ := types.MapValueFrom(ctx, types.StringType, map[string]string{})
		priorComponent.Parameters = emptyParams

		emptyList, _ := types.ListValueFrom(ctx, types.StringType, []string{})
		priorComponent.DependenciesIds = emptyList
		priorComponent.OutputFields = emptyList

		emptyLinks, _ := types.ListValueFrom(ctx, basetypes.ObjectType{AttrTypes: linkAttrTypes}, []LinkModel{})
		priorComponent.Links = emptyLinks

		priorComponents, _ := types.ListValueFrom(ctx, basetypes.ObjectType{AttrTypes: componentAttrTypes}, []ComponentModel{priorComponent})

		model := &BlueprintModel{
			Components: priorComponents,
		}

		blueprint := &fractalCloud.Blueprint{
			Description: "version from API test",
			Components: []fractalCloud.Component{
				{
					Id:      "comp-versioned",
					Type:    "BigData.PaaS.ComputeCluster",
					Version: strPtr("v3"),
				},
			},
		}

		var diags Diagnostics
		mapBlueprintToState(ctx, blueprint, model, &diags)

		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags.Errors())
		}

		var components []ComponentModel
		d := model.Components.ElementsAs(ctx, &components, false)
		if d.HasError() {
			t.Fatalf("failed to extract components: %v", d.Errors())
		}

		if components[0].Version.ValueString() != "v3" {
			t.Errorf("version: expected API value %q, got %q", "v3", components[0].Version.ValueString())
		}
	})
}
