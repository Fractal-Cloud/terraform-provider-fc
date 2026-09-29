---
page_title: "custom_workloads_caas_workload Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a CaaS Workload blueprint component.
---

# function: custom_workloads_caas_workload

Creates a `CustomWorkloads.CaaS.Workload` component: a containerized workload on a container platform. `platform` must be a ContainerPlatform and `subnet` a Subnet.

The blueprint does not choose the offer that runs the workload, and the CaaS offers spell the same setting differently. Each container setting is therefore written under every offer's key:

| Attribute | Parameters written | Read by |
|-----------|--------------------|---------|
| `container_image` | `containerImage`, `image` | Kubernetes and ECS; Cloud Run, Container Apps and Container Instances |
| `container_port` | `containerPort`, `port` | Kubernetes and ECS; Cloud Run, Container Apps and Container Instances |
| `replicas` | `replicas`, `desiredCount` | Kubernetes; ECS |
| `container_name` | `containerName` | ECS |
| `cpu`, `memory` | `cpu`, `memory` | all, in each offer's own units |

Settings only one offer reads go in `extra_parameters`. On Kubernetes that means `namespace`, `env`, `helmChart` and the like. On Container Apps it means `minReplicas` and `maxReplicas`.

## Example Usage

```terraform
locals {
  k8s = provider::fc::network_and_compute_paas_container_platform({
    id = "k8s-cluster"
  })

  api = provider::fc::custom_workloads_caas_workload({
    id              = "api-service"
    display_name    = "API Service"
    container_image = "myregistry/api:1.4.2"
    container_port  = 8080
    cpu             = "500m"
    memory          = "512Mi"
    replicas        = 3
    platform        = local.k8s
    links = [
      { target = local.orders_db, settings = { access = "read-write" } },
    ]
    extra_parameters = {
      namespace = "orders"
    }
  })
}
```

## Signature

```text
custom_workloads_caas_workload(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `container_image` | String | No | Container image to deploy. |
| `container_port` | Number | No | Port the container listens on. |
| `container_name` | String | No | Container name (ECS). |
| `cpu` | String | No | CPU allocation, in the chosen offer's units (e.g. `"500m"` on Kubernetes, `"256"` on ECS). |
| `memory` | String | No | Memory allocation, in the chosen offer's units (e.g. `"512Mi"` on Kubernetes, `"512"` on ECS). |
| `replicas` | Number | No | Number of running instances. |
| `platform` | Component Object | No | The ContainerPlatform to run on, added as a dependency. Must be a component returned by `network_and_compute_paas_container_platform`. |
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
