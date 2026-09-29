---
page_title: "bigdata_paas_distributed_data_processing Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a Distributed Data Processing platform component.
---

# function: bigdata_paas_distributed_data_processing

Creates a Distributed Data Processing platform component (e.g. Databricks workspace). This component represents a managed data processing platform that other BigData components such as compute clusters, jobs, and ML experiments can depend on.

## Example Usage

```terraform
locals {
  data_lake = provider::fc::bigdata_paas_datalake({
    id           = "data-lake"
    display_name = "Data Lake"
  })

  databricks = provider::fc::bigdata_paas_distributed_data_processing({
    id           = "analytics-platform"
    display_name = "Analytics Platform"
    links = [
      {
        target   = local.data_lake
        settings = {
          mountName = "datalake"
          path      = "/"
        }
      }
    ]
  })
}
```

## Signature

```text
bigdata_paas_distributed_data_processing(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `links` | List of Link Object | No | Links to other components with optional settings. Each link has a `target` (Component Object) and optional `settings` (Map of String). |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
