---
page_title: "bigdata_caas_data_processing_job Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a containerized Spark Data Processing Job blueprint component.
---

# function: bigdata_caas_data_processing_job

Creates a `BigData.CaaS.DataProcessingJob` component: a Spark application on a `BigData.CaaS.DistributedDataProcessing` platform. `main_class` is needed for Java and Scala jobs; `main_application_file` is needed to run.

## Example Usage

```terraform
locals {
  etl = provider::fc::bigdata_caas_data_processing_job({
    id                    = "nightly-etl"
    platform              = local.spark
    job_type              = "Python"
    main_application_file = "s3a://jobs/etl.py"
    arguments             = ["--date", "yesterday"]
    spark_conf            = { "spark.sql.shuffle.partitions" = "64" }
    schedule              = "0 2 * * *"
    links = [
      { target = local.lake, settings = { purpose = "curated", path = "orders" } },
    ]
  })
}
```

## Signature

```text
bigdata_caas_data_processing_job(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `job_type` | String | Yes | `"Java"`, `"Python"` or `"Scala"`. Written as the `type` parameter. |
| `mode` | String | No | `"cluster"` or `"client"`. |
| `main_class` | String | No | Main class, for Java and Scala jobs. |
| `main_application_file` | String | No | URI of the application file. |
| `arguments` | List of String | No | Arguments passed to the application. |
| `spark_conf` | Map of String | No | Spark configuration. |
| `schedule` | String | No | Cron expression; without it the job runs once. |
| `restart_policy` | String | No | `"Never"`, `"Always"` or `"OnFailure"`. |
| `max_retries` | Number | No | Retries on failure. |
| `concurrency_policy` | String | No | For scheduled jobs: `"Allow"`, `"Forbid"` or `"Replace"`. |
| `platform` | Component Object | No | The Spark platform to run on, added as a dependency. Must be a component returned by `bigdata_caas_distributed_data_processing`. |
| `links` | List of Object | No | Runtime relationships to other components: each `{ target = <component>, settings = { ... } }`, with `settings` optional. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |

A link to a `BigData.*.Datalake` takes `purpose` (`"raw"`, `"curated"` or `"checkpoint"`, required) and an optional `path`. A link to a `Messaging.*.Entity` takes `access` (`"publish"`, `"subscribe"` or `"publish-subscribe"`, required).
