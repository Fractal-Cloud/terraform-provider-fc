---
page_title: "storage_caas_object_storage Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a containerized Object Storage blueprint component.
---

# function: storage_caas_object_storage

Creates a `Storage.CaaS.ObjectStorage` component: a MinIO bucket store on a container platform. Consumers link to it with `settings = { accessMode = "read" | "write" | "readWrite", scope = { path = "/", recursive = true } }`.

## Example Usage

```terraform
locals {
  assets = provider::fc::storage_caas_object_storage({
    id                 = "assets"
    container_platform = local.k8s
  })
}
```

## Signature

```text
storage_caas_object_storage(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `container_platform` | Component Object | No | The ContainerPlatform to run on, added as a dependency. Must be a component returned by `network_and_compute_paas_container_platform`. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
