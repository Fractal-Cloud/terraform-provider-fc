---
page_title: "security_paas_identity_provider Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates an Identity Provider blueprint component.
---

# function: security_paas_identity_provider

Creates a `Security.PaaS.IdentityProvider` component: a user directory (Amazon Cognito on AWS, Microsoft Entra External ID on Azure). Its attributes are guardrails that cap what any client may request.

Each component that links to the identity provider gets exactly one OAuth app client. The client is shaped by the link's settings:

| Setting | Required | Description |
|---------|----------|-------------|
| `clientType` | Yes | `"web"` (confidential, with a client secret), `"spa"` (public) or `"machine"` (client credentials) |
| `redirectUris` | For `web` and `spa` | OAuth callback URLs |
| `logoutUris` | No | Post-logout redirect URLs |
| `scopes` | No | Requested OAuth scopes |

The linked component receives `OIDC_ISSUER_URI`, `OIDC_CLIENT_ID`, `OIDC_JWKS_URI` and `OIDC_SCOPES`. Confidential clients also receive `OIDC_CLIENT_SECRET_REF`, a reference to the secret store, never the secret.

## Example Usage

```terraform
locals {
  idp = provider::fc::security_paas_identity_provider({
    id                = "customers"
    mfa_configuration = "OPTIONAL"
    session_duration  = 3600
    password_policy = {
      min_length      = 12
      require_symbols = true
    }
  })

  web = provider::fc::custom_workloads_caas_workload({
    id       = "web"
    platform = local.k8s
    links = [
      { target = local.idp, settings = { clientType = "web", redirectUris = ["https://shop.example.com/callback"] } },
    ]
  })
}
```

## Signature

```text
security_paas_identity_provider(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `user_directory_name` | String | No | Name of the user directory. Defaults to the component id. |
| `mfa_configuration` | String | No | `"OFF"`, `"ON"` or `"OPTIONAL"`. On Entra, `"ON"` requires `mfa_break_glass_user_ids`. |
| `mfa_break_glass_user_ids` | List of String | No | Entra object ids excluded from the MFA policy (emergency-access accounts). |
| `session_duration` | Number | No | Session duration in seconds. |
| `password_policy` | Object | No | Password rules: `min_length` (number) and `require_uppercase`, `require_lowercase`, `require_numbers`, `require_symbols` (bools). All optional; Entra applies only `min_length`. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
