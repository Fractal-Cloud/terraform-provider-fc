---
page_title: "storage_paas_column_oriented_dbms Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a Column-Oriented DBMS component.
---

# function: storage_paas_column_oriented_dbms

Creates a Column-Oriented DBMS component. This represents a managed column-oriented database management system such as Cassandra or HBase.

## Example Usage

```terraform
locals {
  col_dbms = provider::fc::storage_paas_column_oriented_dbms({
    id           = "col-dbms"
    display_name = "Column Store"
    description  = "Column-oriented database for analytics"
  })
}
```

## Signature

```text
storage_paas_column_oriented_dbms(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
