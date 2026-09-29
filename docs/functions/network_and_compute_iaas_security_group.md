---
page_title: "network_and_compute_iaas_security_group Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a SecurityGroup blueprint component.
---

# function: network_and_compute_iaas_security_group

Creates a SecurityGroup blueprint component. If `vpc` is provided, it is validated to ensure it is a VirtualNetwork component and added as a dependency; without it the agent places the group in the environment's network. Compute components join the group through their `security_groups` attribute.

`ingress_rules` admit traffic from CIDR ranges. Traffic between two components is not a rule here: link the source component to the target with `settings = { fromPort = 8080 }`, and the agent derives the rules on both components' managed security groups.

## Example Usage

```terraform
locals {
  vpc = provider::fc::network_and_compute_iaas_virtual_network({
    id         = "main-vpc"
    cidr_block = "10.0.0.0/16"
  })

  sg = provider::fc::network_and_compute_iaas_security_group({
    id           = "web-sg"
    display_name = "Web Security Group"
    vpc          = local.vpc
    ingress_rules = [
      { from_port = 443, source_cidr = "0.0.0.0/0" },
      { from_port = 8000, to_port = 8100, protocol = "udp", source_cidr = "10.0.0.0/8" },
    ]
  })
}
```

## Signature

```text
network_and_compute_iaas_security_group(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `vpc` | Component Object | No | A VirtualNetwork component to add as a dependency. Must be a component returned by `network_and_compute_iaas_virtual_network`. |
| `ingress_rules` | List of Object | No | Ingress rules. Each rule supports the fields below; omitted optional fields take their defaults. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |

### Ingress Rule Object

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `from_port` | Number | Yes | Start of the port range. |
| `to_port` | Number | No | End of the port range. Defaults to `from_port`. |
| `protocol` | String | No | Protocol (`"tcp"`, `"udp"`, `"icmp"`). Defaults to `"tcp"`. |
| `source_cidr` | String | Yes | Source CIDR block (e.g. `"0.0.0.0/0"`). |

`source_component_id` is no longer accepted. Link the source component to the target instead.
