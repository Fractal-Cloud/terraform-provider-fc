---
page_title: "storage_caas_relational_dbms Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a containerized Relational DBMS blueprint component.
---

# function: storage_caas_relational_dbms

Creates a `Storage.CaaS.RelationalDbms` component: a PostgreSQL cluster run by CloudNativePG on a container platform.

## Example Usage

```terraform
locals {
  pg = provider::fc::storage_caas_relational_dbms({
    id                 = "pg"
    container_platform = local.k8s
    engine_version     = "17.2"
    age                = true
  })
}
```

## Signature

```text
storage_caas_relational_dbms(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `engine_version` | String | No | PostgreSQL version. |
| `age` | Bool | No | Enables the Apache AGE graph extension. |
| `container_platform` | Component Object | No | The ContainerPlatform to run on, added as a dependency. Must be a component returned by `network_and_compute_paas_container_platform`. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
