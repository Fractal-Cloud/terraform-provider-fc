---
page_title: "custom_workloads_faas_workload Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a FaaS Workload blueprint component.
---

# function: custom_workloads_faas_workload

Creates a `CustomWorkloads.FaaS.Workload` component: a serverless function (AWS Lambda, Azure Functions, Google Cloud Functions). `runtime` and `handler` are required. `handler` is also written as `entryPoint`, the name Cloud Functions uses.

## Example Usage

```terraform
locals {
  resize = provider::fc::custom_workloads_faas_workload({
    id              = "image-resizer"
    display_name    = "Image Resizer"
    runtime         = "nodejs20.x"
    handler         = "index.handler"
    memory_mb       = 1024
    timeout_seconds = 60
    links = [
      { target = local.uploads, settings = { access = "read-write" } },
    ]
  })
}
```

## Signature

```text
custom_workloads_faas_workload(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `runtime` | String | Yes | Language runtime (e.g. `"nodejs20.x"`, `"python3.12"`, `"java21"`). |
| `handler` | String | Yes | Function entry point (e.g. `"index.handler"`). |
| `memory_mb` | Number | No | Memory allocation in MB. |
| `timeout_seconds` | Number | No | Maximum execution time in seconds. |
| `subnet` | Component Object | No | A Subnet component to add as a dependency. Must be a component returned by `network_and_compute_iaas_subnet`. |
| `links` | List of Object | No | Runtime relationships to other components. See [Link Object](#link-object). |
| `security_groups` | List of Component Object | No | SecurityGroup components the workload is a member of. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |

### Link Object

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `target` | Component Object | Yes | The component linked to. |
| `settings` | Object | No | Settings for the link; omit it when the link needs none. Values may be strings, numbers or bools (sent as strings); a list or object value (e.g. `redirectUris = ["https://..."]`) is sent as JSON. |

Common links: a `Storage.*.RelationalDatabase` with `settings = { access = "read-write" }` (or `"read-only"`), an ObjectStorage with `settings = { access = "read" }`, a messaging Entity with `settings = { access = "publish" }`, a `Security.PaaS.IdentityProvider` with `settings = { clientType = "web", redirectUris = [...] }`, another workload with `settings = { fromPort = 8080 }`.
