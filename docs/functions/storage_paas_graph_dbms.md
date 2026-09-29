---
page_title: "storage_paas_graph_dbms Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a Graph DBMS component.
---

# function: storage_paas_graph_dbms

Creates a Graph DBMS component. This represents a managed graph database management system such as Neo4j or Neptune.

## Example Usage

```terraform
locals {
  graph_dbms = provider::fc::storage_paas_graph_dbms({
    id           = "graph-dbms"
    display_name = "Graph Store"
    description  = "Graph database for relationship data"
  })
}
```

## Signature

```text
storage_paas_graph_dbms(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
