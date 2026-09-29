---
page_title: "storage_paas_document_dbms Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a Document DBMS blueprint component.
---

# function: storage_paas_document_dbms

Creates a Document DBMS blueprint component for managed document/NoSQL database services such as DynamoDB, Cosmos DB, or Firestore.

## Example Usage

```terraform
locals {
  docdb = provider::fc::storage_paas_document_dbms({
    id           = "doc-store"
    display_name = "Document Database"
    description  = "NoSQL document store for user profiles"
  })
}
```

## Signature

```text
storage_paas_document_dbms(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
