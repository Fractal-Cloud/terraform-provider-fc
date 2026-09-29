---
page_title: "storage_paas_document_database Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a Document Database component.
---

# function: storage_paas_document_database

Creates a Document Database component. If `dbms` is provided, it is added as a dependency to ensure the DBMS is provisioned before the database.

## Example Usage

```terraform
locals {
  doc_dbms = provider::fc::storage_paas_document_dbms({
    id = "doc-dbms"
  })

  doc_db = provider::fc::storage_paas_document_database({
    id           = "doc-db"
    display_name = "Orders Database"
    dbms         = local.doc_dbms
  })
}
```

## Signature

```text
storage_paas_document_database(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `dbms` | Component Object | No | The Document DBMS component to depend on. If provided, added as a dependency. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
