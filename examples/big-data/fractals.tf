# Data Platform Fractal
#
# A complete data pipeline: a containerized producer pushes events into a
# message stream, which is consumed by a Spark cluster that lands data in a
# data lake.  A legacy Hadoop cluster is mounted alongside for archive access.
# ETL runs on a scheduled job and model training is tracked via MLflow.
#
# All components are defined as locals so that dependencies and links use
# direct object references (type-checked) rather than copy-pasted string IDs.

locals {
  # ── Container platform & producer workload ───────────────────────────
  k8s = provider::fc::network_and_compute_paas_container_platform({
    id           = "k8s"
    display_name = "Container Platform"
    description  = "Managed Container Platform for the data-producer workload"
    node_pools   = []
  })

  data_producer = provider::fc::custom_workloads_caas_workload({
    id              = "data-producer"
    display_name    = "Data Producer"
    description     = "Microservice that ingests external events and publishes them to the stream"
    container_name  = "producer"
    container_image = "my-registry/data-producer:latest"
    container_port  = 8080
    cpu             = "512"
    memory          = "1024"
    replicas        = 2
    platform        = local.k8s
    subnet          = null
    links = [
      {
        target = local.ingest_stream
        settings = {
          access = "publish"
        }
      }
    ]
    security_groups = []
  })

  # ── Messaging (event stream) ─────────────────────────────────────────
  event_broker = provider::fc::messaging_paas_broker({
    id           = "event-broker"
    display_name = "Event Broker"
    description  = "Managed message broker for the ingest stream"
  })

  ingest_stream = provider::fc::messaging_paas_entity({
    id                      = "ingest-stream"
    display_name            = "Ingest Stream"
    description             = "Topic that carries raw events from the producer to the Spark cluster"
    message_retention_hours = 72
    broker                  = local.event_broker
  })

  # ── Storage ──────────────────────────────────────────────────────────
  data_lake = provider::fc::bigdata_paas_datalake({
    id           = "data-lake"
    display_name = "Data Lake"
    description  = "Cloud object storage for raw and curated datasets"
  })

  # ── External (unmanaged) resource ────────────────────────────────────
  legacy_hadoop = provider::fc::bigdata_saas_unmanaged({
    id           = "legacy-hadoop"
    display_name = "Legacy Hadoop"
    description  = "On-premises Hadoop cluster with historical archive data"
    secret       = provider::fc::secret_ref("legacy-hadoop-credentials")
  })

  # ── Spark platform ──────────────────────────────────────────────────
  spark_platform = provider::fc::bigdata_paas_distributed_data_processing({
    id           = "spark-platform"
    display_name = "Spark Platform"
    description  = "Databricks workspace for all data workloads"
    links = [
      {
        target = local.data_lake
        settings = {
          mountName = "datalake"
          path      = "/"
        }
      },
      {
        target = local.legacy_hadoop
        settings = {
          mountName = "legacy-archive"
          path      = "/archive"
        }
      }
    ]
  })

  # ── Compute cluster (consumes from the ingest stream) ────────────────
  etl_cluster = provider::fc::bigdata_paas_compute_cluster({
    id                       = "etl-cluster"
    display_name             = "ETL Cluster"
    description              = "Autoscaling Spark cluster that consumes events from the ingest stream"
    platform                 = local.spark_platform
    cluster_name             = "etl-processing"
    spark_version            = "14.3.x-scala2.12"
    num_workers              = 2
    min_workers              = 1
    max_workers              = 8
    auto_termination_minutes = 20
    spark_conf = {
      "spark.sql.adaptive.enabled"                   = "true"
      "spark.databricks.delta.optimizeWrite.enabled" = "true"
    }
    pypi_libraries  = ["pandas>=2.0", "pyarrow>=14.0"]
    maven_libraries = ["org.apache.hadoop:hadoop-aws:3.3.4"]
    links = [
      {
        target = local.ingest_stream
        settings = {
          consumerGroup    = "$Default"
          startingPosition = "end"
        }
      }
    ]
  })

  # ── Scheduled ETL job ────────────────────────────────────────────────
  etl_job = provider::fc::bigdata_paas_data_processing_job({
    id               = "etl-job"
    display_name     = "ETL Job"
    description      = "Hourly job that transforms raw events into curated datasets"
    platform         = local.spark_platform
    job_name         = "hourly-etl"
    task_type        = "notebook"
    notebook_path    = "/Repos/data-team/pipelines/etl/transform"
    python_file      = null
    main_class_name  = null
    jar_uri          = null
    cron_schedule    = "0 * * * *"
    max_retries      = 2
    existing_cluster = true
    parameters       = ["--env=prod", "--layers=silver,gold"]
    links = [
      { target = local.data_lake, settings = { purpose = "curated", path = "silver" } },
    ]
  })

  # ── ML experiment tracking ───────────────────────────────────────────
  ml_experiment = provider::fc::bigdata_paas_ml_experiment({
    id              = "churn-experiment"
    display_name    = "Churn Prediction"
    description     = "MLflow experiment for customer churn prediction models"
    platform        = local.spark_platform
    experiment_name = "/Shared/ml-experiments/churn-prediction"
  })
}

resource "fc_fractal" "data_platform" {
  bounded_context_id = data.fc_personal_bounded_context.existing_bounded_context.id
  name               = "data-platform"
  version            = "1.0"
  description        = "Data platform: producer → stream → Spark → data lake, with legacy archive and ML tracking"

  components = [
    local.k8s,
    local.data_producer,
    local.event_broker,
    local.ingest_stream,
    local.data_lake,
    local.legacy_hadoop,
    local.spark_platform,
    local.etl_cluster,
    local.etl_job,
    local.ml_experiment,
  ]
}

output "data_platform" {
  value = fc_fractal.data_platform
}
