---
page_title: "network_and_compute_iaas_virtual_machine Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a VirtualMachine blueprint component.
---

# function: network_and_compute_iaas_virtual_machine

Creates a VirtualMachine blueprint component. The `subnet` dependency is type-validated to ensure it is a Subnet component. Use `links` to define runtime relationships to other components, and `security_groups` for security group membership.

## Example Usage

```terraform
locals {
  vpc = provider::fc::network_and_compute_iaas_virtual_network({
    id         = "main-vpc"
    cidr_block = "10.0.0.0/16"
  })

  subnet = provider::fc::network_and_compute_iaas_subnet({
    id         = "app-subnet"
    cidr_block = "10.0.1.0/24"
    vpc        = local.vpc
  })

  sg = provider::fc::network_and_compute_iaas_security_group({
    id  = "web-sg"
    vpc = local.vpc
  })

  vm = provider::fc::network_and_compute_iaas_virtual_machine({
    id              = "web-server"
    display_name    = "Web Server"
    subnet          = local.subnet
    security_groups = [local.sg]
    links = [
      {
        target   = local.backend_server
        settings = {
          fromPort = "8080"
          toPort   = "8080"
          protocol = "tcp"
        }
      }
    ]
  })
}
```

## Signature

```text
network_and_compute_iaas_virtual_machine(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `subnet` | Component Object | No | A Subnet component to add as a dependency. Must be a component returned by `network_and_compute_iaas_subnet`. |
| `links` | List of Object | No | Runtime relationship links to other components. Each link has a `target` (component object) and optional `settings` object. |
| `security_groups` | List of Component Object | No | SecurityGroup components for SG membership. Each must be a component returned by `network_and_compute_iaas_security_group`. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |

### Link Object

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `target` | Component Object | Yes | The target component object. |
| `settings` | Object | No | Settings for the link; omit it when the link needs none. Values may be strings, numbers or bools (sent as strings); a list or object value (e.g. `redirectUris = ["https://..."]`) is sent as JSON. |
