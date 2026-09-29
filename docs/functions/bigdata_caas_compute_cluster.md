---
page_title: "bigdata_caas_compute_cluster Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a containerized Spark Compute Cluster blueprint component.
---

# function: bigdata_caas_compute_cluster

Creates a `BigData.CaaS.ComputeCluster` component: a Spark cluster on a `BigData.CaaS.DistributedDataProcessing` platform.

## Example Usage

```terraform
locals {
  cluster = provider::fc::bigdata_caas_compute_cluster({
    id                 = "analytics"
    platform           = local.spark
    spark_version      = "3.5.1"
    executor_instances = 4
    executor_memory    = "4g"
  })
}
```

## Signature

```text
bigdata_caas_compute_cluster(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `image` | String | No | Spark container image. |
| `spark_version` | String | Yes | Spark version. |
| `driver_cores` | String | No | Driver cores (platform default `"1"`). |
| `driver_memory` | String | No | Driver memory (platform default `"1g"`). |
| `executor_cores` | String | No | Cores per executor (platform default `"2"`). |
| `executor_memory` | String | No | Memory per executor (platform default `"2g"`). |
| `executor_instances` | Number | No | Number of executors (platform default `2`). |
| `platform` | Component Object | No | The Spark platform to run on, added as a dependency. Must be a component returned by `bigdata_caas_distributed_data_processing`. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
