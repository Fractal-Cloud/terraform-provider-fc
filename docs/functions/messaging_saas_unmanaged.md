---
page_title: "messaging_saas_unmanaged Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates an unmanaged Messaging component for an external service.
---

# function: messaging_saas_unmanaged

Creates a `Messaging.SaaS.Unmanaged` component: a Messaging service managed outside Fractal Cloud that other components link to. Its credentials stay in the environment secret store. `secret` references the environment secret by short name, and the agent grants each linked consumer read access to that secret and injects a reference to it (`*_REF`), never the raw value.

Consumers link to it with `links = [{ target = local.external_broker }]`. They can add `settings = { envPrefix = "..." }` to choose the prefix of the injected environment variables, which defaults to the component id.

## Example Usage

```terraform
locals {
  external_broker = provider::fc::messaging_saas_unmanaged({
    id           = "external-broker"
    display_name = "External Kafka"
    secret       = provider::fc::secret_ref("kafka-credentials")
  })
}
```

## Signature

```text
messaging_saas_unmanaged(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `secret` | String | Yes | The environment secret holding the service's credentials, as `provider::fc::secret_ref("<short-name>")`. A raw secret value is rejected. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
