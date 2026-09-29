---
page_title: "fc_fractal Resource - Fractal Cloud"
subcategory: ""
description: |-
  Manages a fractal (blueprint) definition in Fractal Cloud.
---

# fc_fractal (Resource)

Manages a fractal definition (blueprint). A fractal is a reusable, governed infrastructure pattern composed of blueprint components. Components are built using provider functions and express dependencies via direct object references.

## Example Usage

```terraform
locals {
  k8s_cluster = provider::fc::network_and_compute_paas_container_platform({
    id           = "k8s-cluster"
    display_name = "Kubernetes Cluster"
  })

  api_service = provider::fc::custom_workloads_caas_workload({
    id              = "api-service"
    display_name    = "API Service"
    container_image = "my-registry/api-service:latest"
    container_port  = 8080
    platform        = local.k8s_cluster
  })
}

resource "fc_fractal" "microservice" {
  bounded_context_id = fc_personal_bounded_context.production.id
  name        = "microservice-template"
  version     = "1.0"
  description = "Standard microservice architecture blueprint"

  components = [
    local.k8s_cluster,
    local.api_service,
  ]
}
```

## Schema

### Required

- `bounded_context_id` (Object) Composite bounded context identifier containing `type`, `owner_id`, and `short_name`.
- `name` (String) Name of the fractal.
- `version` (String) Version of the fractal.
- `components` (List of Object) List of blueprint components. Each component has:
  - `id` (String, Required) Component identifier.
  - `type` (String, Required) Component type string.
  - `display_name` (String, Optional) Human-readable name.
  - `description` (String, Optional) Component description.
  - `version` (String, Optional) Component version.
  - `parameters` (Map of String, Optional) Configuration parameters. See [Parameter and setting values](#parameter-and-setting-values).
  - `dependencies_ids` (List of String, Optional) IDs of components this depends on.
  - `links` (List of Object, Optional) Links to other components, each with `component_id` (String) and `settings` (Map of String).

### Optional

- `description` (String) Description of the fractal. Defaults to `""`.
- `is_private` (Boolean) Whether the fractal is private. Defaults to `false`.

### Read-Only

- `created_at` (String) Creation timestamp.


### Parameter and setting values

Parameter and link-setting values are strings. A value that is a JSON object or array, such as the output of `jsonencode()` or `secret_ref()`, is sent to the API as that JSON; every other value is sent as the string it is. So a parameter the catalog types as a list or an object is written with `jsonencode([...])`, and a plain-text value that happens to be valid JSON object or array text is sent as JSON too. A refresh keeps your text as long as the API returns the same JSON, whatever its key order or spacing.
