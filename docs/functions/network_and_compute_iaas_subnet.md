---
page_title: "network_and_compute_iaas_subnet Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a Subnet blueprint component.
---

# function: network_and_compute_iaas_subnet

Creates a Subnet blueprint component. If `vpc` is provided, it is validated to ensure it is a VirtualNetwork component and added as a dependency. This ensures the subnet is not reconciled until the VPC is active.

## Example Usage

```terraform
locals {
  vpc = provider::fc::network_and_compute_iaas_virtual_network({
    id           = "main-vpc"
    display_name = "Main VPC"
    cidr_block   = "10.0.0.0/16"
  })

  subnet = provider::fc::network_and_compute_iaas_subnet({
    id                = "app-subnet"
    display_name      = "Application Subnet"
    cidr_block        = "10.0.1.0/24"
    vpc               = local.vpc
  })
}
```

## Signature

```text
network_and_compute_iaas_subnet(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `cidr_block` | String | No | The CIDR block for the subnet (e.g. `"10.0.1.0/24"`). |
| `vpc` | Component Object | No | A VirtualNetwork component to add as a dependency. Must be a component returned by `network_and_compute_iaas_virtual_network`. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
