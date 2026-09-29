---
page_title: "storage_paas_column_oriented_entity Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a Column-Oriented Entity component.
---

# function: storage_paas_column_oriented_entity

Creates a Column-Oriented Entity component. If `dbms` is provided, it is added as a dependency to ensure the DBMS is provisioned before the entity.

## Example Usage

```terraform
locals {
  col_dbms = provider::fc::storage_paas_column_oriented_dbms({
    id = "col-dbms"
  })

  col_entity = provider::fc::storage_paas_column_oriented_entity({
    id           = "col-entity"
    display_name = "Events Table"
    dbms         = local.col_dbms
  })
}
```

## Signature

```text
storage_paas_column_oriented_entity(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `dbms` | Component Object | No | The Column-Oriented DBMS component to depend on. If provided, added as a dependency. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
