---
page_title: "messaging_paas_entity Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a PaaS Message Entity (topic/queue) component.
---

# function: messaging_paas_entity

Creates a PaaS Message Entity (topic/queue) component. If `broker` is provided, it is added as a dependency to ensure the broker is provisioned before the entity.

## Example Usage

```terraform
locals {
  broker = provider::fc::messaging_paas_broker({
    id = "msg-broker"
  })

  topic = provider::fc::messaging_paas_entity({
    id                      = "order-events"
    display_name            = "Order Events Topic"
    message_retention_hours = 72
    broker                  = local.broker
  })
}
```

## Signature

```text
messaging_paas_entity(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `message_retention_hours` | Number | No | Number of hours to retain messages. |
| `broker` | Component Object | No | The Message Broker component to depend on. If provided, added as a dependency. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
