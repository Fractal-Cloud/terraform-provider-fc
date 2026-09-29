---
page_title: "storage_caas_relational_database Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a containerized Relational Database blueprint component.
---

# function: storage_caas_relational_database

Creates a `Storage.CaaS.RelationalDatabase` component: a database on a `Storage.CaaS.RelationalDbms`. Workloads link to it with `settings = { access = "read-write" }` or `"read-only"`. Each linked workload gets its own role and receives `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USERNAME` and `DB_PASSWORD_REF`.

## Example Usage

```terraform
locals {
  orders = provider::fc::storage_caas_relational_database({
    id   = "orders"
    dbms = local.pg
  })
}
```

## Signature

```text
storage_caas_relational_database(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `dbms` | Component Object | No | The DBMS hosting the database, added as a dependency. Must be a component returned by `storage_caas_relational_dbms`. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
