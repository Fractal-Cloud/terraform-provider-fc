---
page_title: "bigdata_paas_ml_experiment Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates an ML Experiment tracking component.
---

# function: bigdata_paas_ml_experiment

Creates an ML Experiment tracking component (e.g. MLflow). If `platform` is provided, it is validated to ensure it is a Distributed Data Processing component and added as a dependency. This ensures the experiment is not reconciled until the platform is active.

## Example Usage

```terraform
locals {
  databricks = provider::fc::bigdata_paas_distributed_data_processing({
    id           = "analytics-platform"
    display_name = "Analytics Platform"
  })

  experiment = provider::fc::bigdata_paas_ml_experiment({
    id              = "model-training"
    display_name    = "Model Training Experiment"
    platform        = local.databricks
    experiment_name = "/Shared/experiments/churn-prediction"
  })
}
```

## Signature

```text
bigdata_paas_ml_experiment(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `platform` | Component Object | No | A Distributed Data Processing component to add as a dependency. Must be a component returned by `bigdata_paas_distributed_data_processing`. |
| `experiment_name` | String | No | Name or path of the ML experiment. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
