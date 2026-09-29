---
page_title: "bigdata_caas_data_catalog Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a containerized Data Catalog blueprint component.
---

# function: bigdata_caas_data_catalog

Creates a `BigData.CaaS.DataCatalog` component: a Unity Catalog server on a container platform.

## Example Usage

```terraform
locals {
  catalog = provider::fc::bigdata_caas_data_catalog({
    id                 = "catalog"
    container_platform = local.k8s
  })
}
```

## Signature

```text
bigdata_caas_data_catalog(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `unity_catalog_version` | String | No | Unity Catalog version. |
| `replicas` | Number | No | Server replicas (platform default `2`). |
| `container_platform` | Component Object | No | The ContainerPlatform to run on, added as a dependency. Must be a component returned by `network_and_compute_paas_container_platform`. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
