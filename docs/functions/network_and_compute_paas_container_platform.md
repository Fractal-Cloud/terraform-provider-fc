---
page_title: "network_and_compute_paas_container_platform Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a ContainerPlatform (managed Kubernetes) blueprint component.
---

# function: network_and_compute_paas_container_platform

Creates a ContainerPlatform (managed Kubernetes) blueprint component. Node pools are serialized into parameters for the agent to reconcile against the target cloud provider (EKS, AKS, GKE, etc.).

## Example Usage

```terraform
locals {
  k8s = provider::fc::network_and_compute_paas_container_platform({
    id           = "k8s-cluster"
    display_name = "Production Kubernetes"
    node_pools = [
      {
        name               = "default"
        disk_size_gb       = 50
        min_node_count     = 2
        max_node_count     = 10
        max_pods_per_node  = 110
        autoscaling_enabled = true
        initial_node_count = 3
        max_surge          = 1
      }
    ]
  })
}
```

## Signature

```text
network_and_compute_paas_container_platform(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `node_pools` | List of Object | No | List of node pool configuration objects. See fields below. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |

### Node Pool Object

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `name` | String | Yes | Name of the node pool. |
| `disk_size_gb` | Number | No | Disk size in GB for each node. |
| `min_node_count` | Number | No | Minimum number of nodes (used with autoscaling). |
| `max_node_count` | Number | No | Maximum number of nodes (used with autoscaling). |
| `max_pods_per_node` | Number | No | Maximum number of pods per node. |
| `autoscaling_enabled` | Boolean | No | Whether cluster autoscaler is enabled for this pool. |
| `initial_node_count` | Number | No | Initial number of nodes when the pool is created. |
| `max_surge` | Number | No | Maximum number of extra nodes created during upgrades. |
