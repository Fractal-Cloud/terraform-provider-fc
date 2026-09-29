---
page_title: "api_management_paas_api_gateway Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a PaaS API Gateway component.
---

# function: api_management_paas_api_gateway

Creates a PaaS API Gateway component. This represents a managed API gateway service for routing, throttling, and securing API traffic.

## Example Usage

```terraform
locals {
  api_service = provider::fc::custom_workloads_paas_workload({
    id              = "api-service"
    display_name    = "API Service"
    container_image = "my-registry/api:latest"
    container_port  = 8080
  })

  api_gw = provider::fc::api_management_paas_api_gateway({
    id           = "main-api-gateway"
    display_name = "Main API Gateway"
    description  = "Public-facing API gateway"
    links = [
      {
        target   = local.api_service
        settings = {
          protocol = "http"
        }
      }
    ]
  })
}
```

## Signature

```text
api_management_paas_api_gateway(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `links` | List of Link Object | No | Links to other components (e.g. workloads) with optional settings. Each link has a `target` (Component Object) and optional `settings` (Map of String). |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
