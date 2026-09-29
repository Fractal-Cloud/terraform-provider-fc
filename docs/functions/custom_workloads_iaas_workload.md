---
page_title: "custom_workloads_iaas_workload Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates an IaaS Workload blueprint component.
---

# function: custom_workloads_iaas_workload

Creates a `CustomWorkloads.IaaS.Workload` component: a workload deployed onto a virtual machine. `vm` must be a VirtualMachine and is added as a dependency. The offer deploys from a git repository, and its settings (`sshRepositoryURI`, `repoId`, `branchName`, `roles`, and the SSH key secret names) go in `extra_parameters`.

## Example Usage

```terraform
locals {
  app = provider::fc::custom_workloads_iaas_workload({
    id           = "legacy-app"
    display_name = "Legacy App"
    vm           = local.app_server
    extra_parameters = {
      sshRepositoryURI = "git@github.com:acme/legacy-app.git"
      branchName       = "main"
    }
  })
}
```

## Signature

```text
custom_workloads_iaas_workload(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `vm` | Component Object | No | The VirtualMachine to deploy onto, added as a dependency. Must be a component returned by `network_and_compute_iaas_virtual_machine`. |
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
