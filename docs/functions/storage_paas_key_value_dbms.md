---
page_title: "storage_paas_key_value_dbms Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a Key-Value DBMS component.
---

# function: storage_paas_key_value_dbms

Creates a Key-Value DBMS component. This represents a managed key-value database management system such as Redis or DynamoDB.

## Example Usage

```terraform
locals {
  kv_dbms = provider::fc::storage_paas_key_value_dbms({
    id           = "kv-dbms"
    display_name = "Cache Store"
    description  = "Key-value store for caching"
  })
}
```

## Signature

```text
storage_paas_key_value_dbms(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
