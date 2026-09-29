---
page_title: "network_and_compute_iaas_virtual_network Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a VirtualNetwork (VPC) blueprint component.
---

# function: network_and_compute_iaas_virtual_network

Creates a VirtualNetwork (VPC) blueprint component. This is the foundational networking component that defines an isolated virtual network in which other components (subnets, security groups, VMs, etc.) are deployed.

## Example Usage

```terraform
locals {
  vpc = provider::fc::network_and_compute_iaas_virtual_network({
    id           = "main-vpc"
    display_name = "Main VPC"
    cidr_block   = "10.0.0.0/16"
  })
}
```

## Signature

```text
network_and_compute_iaas_virtual_network(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `cidr_block` | String | No | The CIDR block for the virtual network (e.g. `"10.0.0.0/16"`). |
| `links` | List of Link Object | No | Links to other components (e.g. another VirtualNetwork for peering). Each link has a `target` (Component Object) and optional `settings` (Map of String). |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
