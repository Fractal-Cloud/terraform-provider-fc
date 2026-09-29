## 2.0.1 (Unreleased)

BUG FIXES:

* The provider address is `registry.terraform.io/fractal-cloud/fc`. The Terraform Registry namespace is the GitHub organization, `Fractal-Cloud`; the `fractalcloud/fc` address used before could not be published. Update `required_providers` to `source = "fractal-cloud/fc"`.
* Releases include `terraform-provider-fc_<version>_manifest.json`, listed in the signed checksums. The Terraform Registry reads the protocol version (6.0) from it; without it the Registry assumes protocol 5, which this provider does not speak.

## 2.0.0 (September 29, 2026)

BREAKING CHANGES:

* Unmanaged functions (`*_saas_unmanaged`) require `secret`, an environment-secret reference built with `secret_ref()`. Raw secret values are rejected. Without it the agent could not reconcile the component.
* `custom_workloads_caas_workload`: `desired_count` is replaced by `replicas`. The value is written as both `replicas` (Kubernetes) and `desiredCount` (ECS).
* `custom_workloads_paas_workload`: `container_name` and `desired_count` are removed, and `container_image`/`container_port` are written as `image`/`port`, the keys the PaaS offers read.
* `custom_workloads_iaas_workload`: container attributes are removed. The offer deploys from git; pass its settings in `extra_parameters`.
* `custom_workloads_faas_workload`: container attributes are removed; `runtime` and `handler` are required.
* `network_and_compute_iaas_security_group`: `ingress_rules[*].source_component_id` is removed; link the source component to the target instead. `source_cidr` is required on every rule.
* `storage_paas_relational_dbms`: `engine_version` is required.
* `links` is a dynamic list: `settings` may be omitted, and setting values may be numbers, bools, lists or objects.
* A parameter or link-setting string that is a JSON object or array (it starts with `{` or `[` and parses) is sent to the API as that JSON instead of as a string. This is what makes list and object parameters and `secret_ref()` work. A value meant as literal text that happens to parse, such as `"[1]"`, is affected too.

FEATURES:

* Every attribute a function does not mark as required may now be omitted. Before, Terraform rejected any call that left out an attribute, including the documented optional ones, so the IaaS examples did not validate.
* **New Provider Function:** `secret_ref(short_name)` references an environment secret from any parameter or link setting.
* **New Provider Functions:** `ai_agentic_platform`, `ai_saas_unmanaged`, `security_paas_identity_provider`, `bigdata_caas_distributed_data_processing`, `bigdata_caas_compute_cluster`, `bigdata_caas_data_processing_job`, `bigdata_caas_ml_experiment`, `bigdata_caas_data_catalog`, `bigdata_caas_datalake`, `storage_caas_relational_dbms`, `storage_caas_relational_database`, `storage_caas_object_storage`.
* Every function accepts `extra_parameters`, for offer-specific parameters that have no attribute.
* `bigdata_paas_data_processing_job` accepts `links` (Datalake `purpose`/`path`, messaging Entity `access`).
* `storage_paas_relational_dbms` accepts `age` (Apache AGE graph extension).

BUG FIXES:

* Array and object parameter values (`sparkConf`, library lists, `nodePools`, `ingressRules`, secret references) are sent as JSON instead of as strings containing JSON, which agents ignored.
* A refresh no longer reports a diff when the API returns a structured value with different key order or spacing.
* Numeric values read from the API keep their literal text.

NOTES:

* Built with Go 1.26. The dependencies are updated, including terraform-plugin-framework 1.19.0 and grpc 1.84.0.

## 1.0.0 (March 13, 2026)

FEATURES:

* **New Resource:** `fc_personal_bounded_context` - Manage personal bounded contexts
* **New Resource:** `fc_organizational_bounded_context` - Manage organizational bounded contexts
* **New Resource:** `fc_fractal` - Manage fractal definitions (blueprints) with component composition
* **New Resource:** `fc_management_environment` - Manage governance environments
* **New Resource:** `fc_operational_environment` - Manage runtime environments

* **New Data Source:** `fc_personal_bounded_context` - Look up personal bounded contexts
* **New Data Source:** `fc_organizational_bounded_context` - Look up organizational bounded contexts
* **New Data Source:** `fc_organization` - Look up organizations
* **New Data Source:** `fc_fractal` - Look up fractal definitions

* **New Provider Functions:** 46 blueprint component builder functions across 8 infrastructure domains:
  * NetworkAndCompute (7): virtual network, subnet, load balancer, security group, virtual machine, container platform, unmanaged
  * CustomWorkloads (5): CaaS, IaaS, PaaS, FaaS workloads, unmanaged
  * Storage (14): files/blobs, relational/document/column-oriented/key-value/graph databases, search, unmanaged
  * Messaging (5): PaaS and CaaS brokers and entities, unmanaged
  * BigData (6): distributed data processing, compute cluster, data processing job, ML experiment, datalake, unmanaged
  * APIManagement (3): PaaS and CaaS API gateways, unmanaged
  * Observability (4): monitoring, tracing, logging, unmanaged
  * Security (2): service mesh security, unmanaged
