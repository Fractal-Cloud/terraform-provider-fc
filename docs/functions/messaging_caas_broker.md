---
page_title: "messaging_caas_broker Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a CaaS Message Broker component.
---

# function: messaging_caas_broker

Creates a CaaS Message Broker component. If `container_platform` is provided, it is added as a dependency to ensure the container platform is provisioned before the broker.

## Example Usage

```terraform
locals {
  k8s = provider::fc::network_and_compute_paas_container_platform({
    id = "k8s-cluster"
  })

  broker = provider::fc::messaging_caas_broker({
    id                 = "caas-broker"
    display_name       = "Kafka Broker"
    container_platform = local.k8s
  })
}
```

## Signature

```text
messaging_caas_broker(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `container_platform` | Component Object | No | The container platform component to depend on. If provided, added as a dependency. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
