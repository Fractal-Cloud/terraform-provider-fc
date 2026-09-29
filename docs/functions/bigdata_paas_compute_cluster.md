---
page_title: "bigdata_paas_compute_cluster Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a Compute Cluster component.
---

# function: bigdata_paas_compute_cluster

Creates a Compute Cluster component. If `platform` is provided, it is validated to ensure it is a Distributed Data Processing component and added as a dependency. This ensures the cluster is not reconciled until the platform is active.

## Example Usage

```terraform
locals {
  databricks = provider::fc::bigdata_paas_distributed_data_processing({
    id           = "analytics-platform"
    display_name = "Analytics Platform"
  })

  cluster = provider::fc::bigdata_paas_compute_cluster({
    id                       = "etl-cluster"
    display_name             = "ETL Cluster"
    platform                 = local.databricks
    cluster_name             = "etl-processing"
    spark_version            = "13.3.x-scala2.12"
    min_workers              = 1
    max_workers              = 4
    auto_termination_minutes = 30
    spark_conf = {
      "spark.sql.adaptive.enabled" = "true"
    }
    pypi_libraries  = ["pandas==2.0.0", "numpy"]
    maven_libraries = ["com.databricks:spark-xml_2.12:0.16.0"]
  })
}
```

## Signature

```text
bigdata_paas_compute_cluster(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `platform` | Component Object | No | A Distributed Data Processing component to add as a dependency. Must be a component returned by `bigdata_paas_distributed_data_processing`. |
| `cluster_name` | String | No | Name of the compute cluster. |
| `spark_version` | String | No | The Spark runtime version for the cluster. |
| `num_workers` | Number | No | Fixed number of workers. Use this for a non-autoscaling cluster. |
| `min_workers` | Number | No | Minimum number of workers for autoscaling. |
| `max_workers` | Number | No | Maximum number of workers for autoscaling. |
| `auto_termination_minutes` | Number | No | Number of idle minutes before the cluster is automatically terminated. |
| `spark_conf` | Map of String | No | A map of Spark configuration key-value pairs. |
| `pypi_libraries` | List of String | No | List of PyPI packages to install on the cluster. |
| `maven_libraries` | List of String | No | List of Maven coordinates to install on the cluster. |
| `links` | List of Link Object | No | Links to other components with optional settings. Each link has a `target` (Component Object) and optional `settings` (Map of String). |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
