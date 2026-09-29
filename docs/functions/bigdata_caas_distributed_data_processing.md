---
page_title: "bigdata_caas_distributed_data_processing Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a containerized Distributed Data Processing blueprint component.
---

# function: bigdata_caas_distributed_data_processing

Creates a `BigData.CaaS.DistributedDataProcessing` component: Spark on Kubernetes, run by the Spark operator on a container platform. Operator settings of the chosen offer (`operatorVersion`, `sparkVersion`, `namespace`, `enableWebhook`, `enableMetrics`) go in `extra_parameters`. Link it to a `BigData.CaaS.Datalake` to mount the lake.

## Example Usage

```terraform
locals {
  spark = provider::fc::bigdata_caas_distributed_data_processing({
    id                 = "spark"
    container_platform = local.k8s
    links              = [{ target = local.lake }]
    extra_parameters = {
      namespace = "spark"
    }
  })
}
```

## Signature

```text
bigdata_caas_distributed_data_processing(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `container_platform` | Component Object | No | The ContainerPlatform to run on, added as a dependency. Must be a component returned by `network_and_compute_paas_container_platform`. |
| `links` | List of Object | No | Runtime relationships to other components: each `{ target = <component>, settings = { ... } }`, with `settings` optional. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
