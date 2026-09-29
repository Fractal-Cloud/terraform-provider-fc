package provider

import (
	"context"
	"fmt"

	"fractal.cloud/terraform-provider-fc/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &FractalsDataSource{}
	_ datasource.DataSourceWithConfigure = &FractalsDataSource{}
)

// linkAttrTypes defines the attribute types for a component link object.
var linkAttrTypes = map[string]attr.Type{
	"component_id": basetypes.StringType{},
	"settings":     basetypes.MapType{ElemType: basetypes.StringType{}},
}

// componentAttrTypes defines the attribute types for a component object.
var componentAttrTypes = map[string]attr.Type{
	"id":                  basetypes.StringType{},
	"type":                basetypes.StringType{},
	"display_name":        basetypes.StringType{},
	"description":         basetypes.StringType{},
	"version":             basetypes.StringType{},
	"is_locked":           basetypes.BoolType{},
	"recreate_on_failure": basetypes.BoolType{},
	"parameters":          basetypes.MapType{ElemType: basetypes.StringType{}},
	"dependencies_ids":    basetypes.ListType{ElemType: basetypes.StringType{}},
	"links":               basetypes.ListType{ElemType: basetypes.ObjectType{AttrTypes: linkAttrTypes}},
	"output_fields":       basetypes.ListType{ElemType: basetypes.StringType{}},
}

// NewFractalDataSource is a helper function to simplify the provider implementation.
func NewFractalDataSource() datasource.DataSource {
	return &FractalsDataSource{}
}

// FractalsDataSource is the data source implementation.
type FractalsDataSource struct {
	client *fractalCloud.Client
}

// Configure adds the provider configured client to the data source.
func (d *FractalsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*fractalCloud.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *fractalCloud.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Metadata returns the data source type name.
func (d *FractalsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fractal"
}

// Schema defines the schema for the data source.
func (d *FractalsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bounded_context_id": schema.ObjectAttribute{
				Required: true,
				AttributeTypes: map[string]attr.Type{
					"type":       basetypes.StringType{},
					"owner_id":   basetypes.StringType{},
					"short_name": basetypes.StringType{},
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"version": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"is_private": schema.BoolAttribute{
				Computed: true,
			},
			"components": schema.ListAttribute{
				Computed:    true,
				ElementType: basetypes.ObjectType{AttrTypes: componentAttrTypes},
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

type LinkModel struct {
	ComponentId types.String `tfsdk:"component_id"`
	Settings    types.Map    `tfsdk:"settings"`
}

type ComponentModel struct {
	Id                types.String `tfsdk:"id"`
	Type              types.String `tfsdk:"type"`
	DisplayName       types.String `tfsdk:"display_name"`
	Description       types.String `tfsdk:"description"`
	Version           types.String `tfsdk:"version"`
	IsLocked          types.Bool   `tfsdk:"is_locked"`
	RecreateOnFailure types.Bool   `tfsdk:"recreate_on_failure"`
	Parameters        types.Map    `tfsdk:"parameters"`
	DependenciesIds   types.List   `tfsdk:"dependencies_ids"`
	Links             types.List   `tfsdk:"links"`
	OutputFields      types.List   `tfsdk:"output_fields"`
}

type BoundedContextIdModel struct {
	Type      types.String `tfsdk:"type"`
	OwnerId   types.String `tfsdk:"owner_id"`
	ShortName types.String `tfsdk:"short_name"`
}

// BlueprintModel maps fractal schema data.
type BlueprintModel struct {
	BoundedContextId BoundedContextIdModel `tfsdk:"bounded_context_id"`
	Name             types.String          `tfsdk:"name"`
	Version          types.String          `tfsdk:"version"`
	Description      types.String          `tfsdk:"description"`
	IsPrivate        types.Bool            `tfsdk:"is_private"`
	Components       types.List            `tfsdk:"components"`
	CreatedAt        types.String          `tfsdk:"created_at"`
}

// Read refreshes the Terraform state with the latest data.
func (d *FractalsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config BlueprintModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fractalId := fractalIdFromModel(config)
	tflog.Debug(ctx, "reading fractal data source", map[string]any{
		"fractal_id": fractalId.ToString(),
	})

	blueprint, err := d.client.GetBlueprint(ctx, fractalId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Fractal",
			fmt.Sprintf("Could not read fractal %q: %s", fractalId.ToString(), err),
		)
		return
	}

	if blueprint == nil {
		resp.Diagnostics.AddError(
			"Fractal Not Found",
			fmt.Sprintf("No fractal found with id %q.", fractalId.ToString()),
		)
		return
	}

	state := BlueprintModel{
		BoundedContextId: config.BoundedContextId,
		Name:             config.Name,
		Version:          config.Version,
	}

	mapBlueprintToState(ctx, blueprint, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// mapBlueprintToState maps an API Blueprint response into the Terraform BlueprintModel.
func mapBlueprintToState(
	ctx context.Context,
	blueprint *fractalCloud.Blueprint,
	model *BlueprintModel,
	diags *Diagnostics,
) {
	// Extract prior state component values so we can preserve fields the API doesn't round-trip.
	priorVersions := make(map[string]types.String)
	priorParameters := make(map[string]map[string]string)
	priorLinkSettings := make(map[string][]priorLink)
	if !model.Components.IsNull() && !model.Components.IsUnknown() {
		var priorComponents []ComponentModel
		d := model.Components.ElementsAs(ctx, &priorComponents, false)
		if !d.HasError() {
			for _, pc := range priorComponents {
				id := pc.Id.ValueString()
				priorVersions[id] = pc.Version
				priorParameters[id] = stringMap(ctx, pc.Parameters)
				priorLinkSettings[id] = priorLinks(ctx, pc.Links)
			}
		}
	}

	components := make([]ComponentModel, len(blueprint.Components))
	for i, component := range blueprint.Components {
		params := preferPriorValues(priorParameters[component.Id], component.Parameters)
		if params == nil {
			params = map[string]string{}
		}
		parameters, d := types.MapValueFrom(ctx, types.StringType, params)
		diags.Append(d...)

		depIds := component.DependenciesIds
		if depIds == nil {
			depIds = []string{}
		}
		dependenciesIds, d := types.ListValueFrom(ctx, types.StringType, depIds)
		diags.Append(d...)

		links := make([]LinkModel, len(component.Links))
		priorByTarget := priorLinksByTarget(priorLinkSettings[component.Id])
		for j, link := range component.Links {
			linkSettings := link.Settings
			if prior := priorByTarget[link.ComponentId]; len(prior) > 0 {
				linkSettings = preferPriorValues(prior[0], linkSettings)
				priorByTarget[link.ComponentId] = prior[1:]
			}
			if linkSettings == nil {
				linkSettings = map[string]string{}
			}
			settings, d := types.MapValueFrom(ctx, types.StringType, linkSettings)
			diags.Append(d...)

			links[j] = LinkModel{
				ComponentId: types.StringValue(link.ComponentId),
				Settings:    settings,
			}
		}

		linksToMap, d := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: linkAttrTypes,
		}, links)
		diags.Append(d...)

		outFields := component.OutputFields
		if outFields == nil {
			outFields = []string{}
		}
		outputFields, d := types.ListValueFrom(ctx, types.StringType, outFields)
		diags.Append(d...)

		// The API doesn't return version — preserve from prior state if available.
		version := stringPointerToTFValue(component.Version)
		if version.ValueString() == "" {
			if v, ok := priorVersions[component.Id]; ok && !v.IsNull() && !v.IsUnknown() {
				version = v
			}
		}

		components[i] = ComponentModel{
			Id:                types.StringValue(component.Id),
			Type:              types.StringValue(component.Type),
			DisplayName:       stringPointerToTFValue(component.DisplayName),
			Description:       stringPointerToTFValue(component.Description),
			Version:           version,
			IsLocked:          boolPointerToTFValue(component.IsLocked),
			RecreateOnFailure: boolPointerToTFValue(component.RecreateOnFailure),
			Parameters:        parameters,
			DependenciesIds:   dependenciesIds,
			Links:             linksToMap,
			OutputFields:      outputFields,
		}
	}

	if diags.HasError() {
		return
	}

	componentsToMap, d := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: componentAttrTypes,
	}, components)
	diags.Append(d...)

	model.Description = types.StringValue(blueprint.Description)
	model.IsPrivate = types.BoolValue(blueprint.IsPrivate)
	model.Components = componentsToMap
	model.CreatedAt = types.StringValue(blueprint.CreatedAt)
}
