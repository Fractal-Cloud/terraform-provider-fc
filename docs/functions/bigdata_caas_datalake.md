---
page_title: "bigdata_caas_datalake Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a containerized Datalake blueprint component.
---

# function: bigdata_caas_datalake

Creates a `BigData.CaaS.Datalake` component: a MinIO tenant on a container platform. Tenant sizing of the chosen offer (`servers`, `volumesPerServer`, `volumeSize`, `storageClass`) goes in `extra_parameters`. Jobs link to it with `settings = { purpose = "raw" | "curated" | "checkpoint" }`.

## Example Usage

```terraform
locals {
  lake = provider::fc::bigdata_caas_datalake({
    id                 = "lake"
    container_platform = local.k8s
    extra_parameters = {
      volumeSize = "500Gi"
    }
  })
}
```

## Signature

```text
bigdata_caas_datalake(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `container_platform` | Component Object | No | The ContainerPlatform to run on, added as a dependency. Must be a component returned by `network_and_compute_paas_container_platform`. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
