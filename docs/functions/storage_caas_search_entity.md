---
page_title: "storage_caas_search_entity Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a Search Entity component.
---

# function: storage_caas_search_entity

Creates a Search Entity component. If `search` is provided, it is added as a dependency to ensure the search platform is provisioned before the entity.

## Example Usage

```terraform
locals {
  search = provider::fc::storage_caas_search({
    id = "search-platform"
  })

  search_entity = provider::fc::storage_caas_search_entity({
    id           = "products-index"
    display_name = "Products Index"
    search       = local.search
  })
}
```

## Signature

```text
storage_caas_search_entity(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `search` | Component Object | No | The Search Platform component to depend on. If provided, added as a dependency. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
