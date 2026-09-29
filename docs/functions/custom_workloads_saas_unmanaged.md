---
page_title: "custom_workloads_saas_unmanaged Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates an unmanaged Custom Workloads component for an external service.
---

# function: custom_workloads_saas_unmanaged

Creates a `CustomWorkloads.SaaS.Unmanaged` component: a Custom Workloads service managed outside Fractal Cloud that other components link to. Its credentials stay in the environment secret store. `secret` references the environment secret by short name, and the agent grants each linked consumer read access to that secret and injects a reference to it (`*_REF`), never the raw value.

Consumers link to it with `links = [{ target = local.payments_api }]`. They can add `settings = { envPrefix = "..." }` to choose the prefix of the injected environment variables, which defaults to the component id.

## Example Usage

```terraform
locals {
  payments_api = provider::fc::custom_workloads_saas_unmanaged({
    id           = "payments-api"
    display_name = "Payments Provider API"
    secret       = provider::fc::secret_ref("payments-api-key")
  })
}
```

## Signature

```text
custom_workloads_saas_unmanaged(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `secret` | String | Yes | The environment secret holding the service's credentials, as `provider::fc::secret_ref("<short-name>")`. A raw secret value is rejected. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
