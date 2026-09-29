---
page_title: "bigdata_caas_ml_experiment Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a containerized ML Experiment tracking blueprint component.
---

# function: bigdata_caas_ml_experiment

Creates a `BigData.CaaS.MlExperiment` component: an MLflow tracking server on a container platform.

## Example Usage

```terraform
locals {
  mlflow = provider::fc::bigdata_caas_ml_experiment({
    id                 = "mlflow"
    container_platform = local.k8s
    replicas           = 2
  })
}
```

## Signature

```text
bigdata_caas_ml_experiment(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `mlflow_version` | String | No | MLflow version. |
| `backend_store_uri` | String | No | Tracking backend store URI (platform default `sqlite:///mlflow/mlflow.db`). |
| `artifact_root` | String | No | Artifact root (platform default `/mlflow/artifacts`). |
| `replicas` | Number | No | Server replicas (platform default `2`). |
| `service_port` | Number | No | Service port (platform default `5000`). |
| `container_platform` | Component Object | No | The ContainerPlatform to run on, added as a dependency. Must be a component returned by `network_and_compute_paas_container_platform`. |
| `extra_parameters` | Map of String | No | Additional parameters for keys the chosen offer reads that have no attribute here. A key an attribute already sets is rejected. JSON object or array strings (e.g. from `jsonencode()` or `secret_ref()`) are sent as JSON. |
