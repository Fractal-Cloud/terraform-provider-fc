---
page_title: "ai_agentic_platform Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates an Agentic Platform blueprint component.
---

# function: ai_agentic_platform

Creates an `AI.AgenticPlatform` component: a governed runtime for AI agents. Link it to a `Security.PaaS.IdentityProvider` for authentication and to `AI.SaaS.Unmanaged` components for the model services its agents call.

## Example Usage

```terraform
locals {
  agents = provider::fc::ai_agentic_platform({
    id                = "agents"
    bounded_context   = "payments"
    sizing_tier       = "standard"
    default_acf_level = "ACF-2"
    links = [
      { target = local.idp, settings = { clientType = "web", redirectUris = ["https://agents.example.com/callback"] } },
      { target = local.openai },
    ]
  })
}
```

## Signature

```text
ai_agentic_platform(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `bounded_context` | String | Yes | The bounded context the platform serves. |
| `sizing_tier` | String | No | `"small"`, `"standard"` or `"large"`. The platform defaults to `"standard"`. |
| `default_acf_level` | String | No | Default agent containment level, `"ACF-1"` or `"ACF-2"`. The platform defaults to `"ACF-2"`. |
| `habitat_runtime` | String | No | Sandbox runtime for agents: `"auto"`, `"firecracker"` or `"gvisor"`. The platform defaults to `"auto"`. |
| `links` | List of Object | No | Runtime relationships to other components: each `{ target = <component>, settings = { ... } }`, with `settings` optional. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
