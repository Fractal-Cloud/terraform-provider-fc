package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// priorLink is a link as the prior state holds it, used to keep its setting
// strings when the API returns the same values re-encoded.
type priorLink struct {
	componentId string
	settings    map[string]string
}

func priorLinks(ctx context.Context, links types.List) []priorLink {
	if links.IsNull() || links.IsUnknown() {
		return nil
	}
	var models []LinkModel
	if d := links.ElementsAs(ctx, &models, false); d.HasError() {
		return nil
	}
	out := make([]priorLink, len(models))
	for i, m := range models {
		out[i] = priorLink{componentId: m.ComponentId.ValueString(), settings: stringMap(ctx, m.Settings)}
	}
	return out
}

func stringMap(ctx context.Context, m types.Map) map[string]string {
	if m.IsNull() || m.IsUnknown() {
		return nil
	}
	out := map[string]string{}
	if d := m.ElementsAs(ctx, &out, false); d.HasError() {
		return nil
	}
	return out
}
