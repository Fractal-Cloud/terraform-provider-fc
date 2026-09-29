---
page_title: "custom_workloads_paas_workload Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a PaaS Workload blueprint component.
---

# function: custom_workloads_paas_workload

Creates a `CustomWorkloads.PaaS.Workload` component: a workload on a platform service (Cloud Run, Azure Container Instances, Azure Web App). `container_image` and `container_port` are written as `image` and `port`, the keys the container-based offers read. The Web App offer deploys from git and takes its settings (`sshRepositoryURI`, `branchName`, ...) through `extra_parameters`.

## Example Usage

```terraform
locals {
  web = provider::fc::custom_workloads_paas_workload({
    id              = "web"
    display_name    = "Web Frontend"
    container_image = "gcr.io/acme/web:2.0"
    container_port  = 8080
    cpu             = "1"
    memory          = "512Mi"
    extra_parameters = {
      maxInstances = "10"
    }
  })
}
```

## Signature

```text
custom_workloads_paas_workload(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `container_image` | String | No | Container image to deploy, written as `image`. |
| `container_port` | Number | No | Port the container listens on, written as `port`. |
| `cpu` | String | No | CPU allocation, in the chosen offer's units. |
| `memory` | String | No | Memory allocation, in the chosen offer's units. |
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
