---
page_title: "storage_paas_relational_database Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a Relational Database blueprint component.
---

# function: storage_paas_relational_database

Creates a Relational Database blueprint component. If `dbms` is provided, it is validated to ensure it is a Relational DBMS component and added as a dependency, ensuring the database server is active before the database is created.

## Example Usage

```terraform
locals {
  dbms = provider::fc::storage_paas_relational_dbms({
    id             = "postgres-server"
    display_name   = "PostgreSQL Server"
    engine_version = "16"
  })

  database = provider::fc::storage_paas_relational_database({
    id        = "app-db"
    display_name = "Application Database"
    collation = "en_US.UTF-8"
    charset   = "UTF8"
    dbms      = local.dbms
  })
}
```

## Signature

```text
storage_paas_relational_database(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `collation` | String | No | Database collation (e.g. `"en_US.UTF-8"`). |
| `charset` | String | No | Database character set (e.g. `"UTF8"`). |
| `dbms` | Component Object | No | A Relational DBMS component to add as a dependency. Must be a component returned by `storage_paas_relational_dbms`. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
