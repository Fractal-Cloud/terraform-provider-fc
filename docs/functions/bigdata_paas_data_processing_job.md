---
page_title: "bigdata_paas_data_processing_job Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a Data Processing Job component.
---

# function: bigdata_paas_data_processing_job

Creates a Data Processing Job component. If `platform` is provided, it is validated to ensure it is a Distributed Data Processing component and added as a dependency. This ensures the job is not reconciled until the platform is active.

## Example Usage

```terraform
locals {
  databricks = provider::fc::bigdata_paas_distributed_data_processing({
    id           = "analytics-platform"
    display_name = "Analytics Platform"
  })

  etl_job = provider::fc::bigdata_paas_data_processing_job({
    id             = "daily-etl"
    display_name   = "Daily ETL Job"
    platform       = local.databricks
    job_name       = "daily-etl-pipeline"
    task_type      = "notebook"
    notebook_path  = "/Shared/etl/daily_pipeline"
    cron_schedule  = "0 0 6 * * ?"
    max_retries    = 2
    parameters     = ["--env=production", "--date=today"]
  })
}
```

## Signature

```text
bigdata_paas_data_processing_job(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `platform` | Component Object | No | A Distributed Data Processing component to add as a dependency. Must be a component returned by `bigdata_paas_distributed_data_processing`. |
| `job_name` | String | No | Name of the job. |
| `task_type` | String | No | The type of task to run (e.g. `"notebook"`, `"python"`, `"jar"`). |
| `notebook_path` | String | No | Path to the notebook to execute (used when `task_type` is `"notebook"`). |
| `python_file` | String | No | Path to the Python file to execute (used when `task_type` is `"python"`). |
| `main_class_name` | String | No | Main class name for JAR tasks (used when `task_type` is `"jar"`). |
| `jar_uri` | String | No | URI of the JAR file to execute (used when `task_type` is `"jar"`). |
| `cron_schedule` | String | No | Cron expression for scheduling the job. |
| `max_retries` | Number | No | Maximum number of retries on failure. |
| `existing_cluster` | Boolean | No | Whether to use an existing cluster instead of creating a new one. |
| `parameters` | List of String | No | List of parameters to pass to the job. |
| `links` | List of Object | No | Links to the data the job reads and writes. See [Link Object](#link-object). |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |

### Link Object

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `target` | Component Object | Yes | The component linked to. |
| `settings` | Object | No | Settings for the link. |

A link to a `BigData.*.Datalake` takes `purpose` (`"raw"`, `"curated"` or `"checkpoint"`, required) and an optional `path` inside the lake. A link to a `Messaging.*.Entity` takes `access` (`"publish"`, `"subscribe"` or `"publish-subscribe"`, required), plus the optional `consumerGroup` and `startingPosition`.

```terraform
links = [
  { target = local.lake,   settings = { purpose = "raw", path = "orders" } },
  { target = local.events, settings = { access = "subscribe", consumerGroup = "etl" } },
]
```
