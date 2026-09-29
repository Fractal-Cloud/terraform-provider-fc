---
page_title: "messaging_paas_broker Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a PaaS Message Broker component.
---

# function: messaging_paas_broker

Creates a PaaS Message Broker component. This represents a managed message broker service such as Amazon SNS/SQS, Azure Service Bus, or Google Pub/Sub.

## Example Usage

```terraform
locals {
  broker = provider::fc::messaging_paas_broker({
    id           = "msg-broker"
    display_name = "Event Broker"
    description  = "Central message broker for event-driven architecture"
  })
}
```

## Signature

```text
messaging_paas_broker(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
