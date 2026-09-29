---
page_title: "api_management_caas_api_gateway Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a CaaS API Gateway component.
---

# function: api_management_caas_api_gateway

Creates a CaaS API Gateway component. If `container_platform` is provided, it is validated to ensure it is a Container Orchestrator component and added as a dependency. This ensures the API gateway is not reconciled until the container platform is active.

## Example Usage

```terraform
locals {
  k8s = provider::fc::network_and_compute_paas_container_platform({
    id           = "main-k8s"
    display_name = "Main Kubernetes Cluster"
  })

  api_service = provider::fc::custom_workloads_caas_workload({
    id              = "api-service"
    display_name    = "API Service"
    container_image = "my-registry/api:latest"
    container_port  = 8080
    platform        = local.k8s
  })

  api_gw = provider::fc::api_management_caas_api_gateway({
    id                 = "ingress-gateway"
    display_name       = "Ingress API Gateway"
    container_platform = local.k8s
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
api_management_caas_api_gateway(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `container_platform` | Component Object | No | A Container Orchestrator component to add as a dependency. Must be a container platform component. |
| `links` | List of Link Object | No | Links to other components (e.g. workloads) with optional settings. Each link has a `target` (Component Object) and optional `settings` (Map of String). |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
