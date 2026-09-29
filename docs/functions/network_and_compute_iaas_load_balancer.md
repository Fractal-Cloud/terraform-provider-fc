---
page_title: "network_and_compute_iaas_load_balancer Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a LoadBalancer blueprint component.
---

# function: network_and_compute_iaas_load_balancer

Creates a LoadBalancer blueprint component. Use `links` to connect backend compute targets and `security_groups` for security group membership.

## Example Usage

```terraform
locals {
  sg = provider::fc::network_and_compute_iaas_security_group({
    id  = "lb-sg"
    vpc = local.vpc
  })

  lb = provider::fc::network_and_compute_iaas_load_balancer({
    id              = "app-lb"
    display_name    = "Application Load Balancer"
    security_groups = [local.sg]
    links = [
      {
        target   = local.web_server
        settings = {
          fromPort = "80"
          toPort   = "80"
          protocol = "tcp"
        }
      }
    ]
  })
}
```

## Signature

```text
network_and_compute_iaas_load_balancer(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `links` | List of Object | No | Links to backend compute components. Each link has a `target` (component object) and optional `settings` object. |
| `security_groups` | List of Component Object | No | SecurityGroup components for SG membership. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
