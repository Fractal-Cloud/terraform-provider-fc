---
page_title: "storage_paas_key_value_entity Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a Key-Value Entity component.
---

# function: storage_paas_key_value_entity

Creates a Key-Value Entity component. If `dbms` is provided, it is added as a dependency to ensure the DBMS is provisioned before the entity.

## Example Usage

```terraform
locals {
  kv_dbms = provider::fc::storage_paas_key_value_dbms({
    id = "kv-dbms"
  })

  kv_entity = provider::fc::storage_paas_key_value_entity({
    id           = "kv-entity"
    display_name = "Session Cache"
    dbms         = local.kv_dbms
  })
}
```

## Signature

```text
storage_paas_key_value_entity(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `dbms` | Component Object | No | The Key-Value DBMS component to depend on. If provided, added as a dependency. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
