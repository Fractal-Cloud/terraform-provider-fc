---
page_title: "observability_caas_logging Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a CaaS Logging component.
---

# function: observability_caas_logging

Creates a CaaS Logging component. If `container_platform` is provided, it is validated to ensure it is a Container Orchestrator component and added as a dependency. This ensures the logging stack is not reconciled until the container platform is active.

## Example Usage

```terraform
locals {
  k8s = provider::fc::network_and_compute_paas_container_platform({
    id           = "main-k8s"
    display_name = "Main Kubernetes Cluster"
  })

  logging = provider::fc::observability_caas_logging({
    id                 = "cluster-logging"
    display_name       = "Centralized Logging"
    container_platform = local.k8s
  })
}
```

## Signature

```text
observability_caas_logging(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `container_platform` | Component Object | No | A Container Orchestrator component to add as a dependency. Must be a container platform component. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
