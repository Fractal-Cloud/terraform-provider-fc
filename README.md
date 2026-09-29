# Terraform Provider for Fractal Cloud

[![build status][build-image]][build-url]
[![codecov][codecov-image]][codecov-url]
[![License: AGPL v3](https://img.shields.io/badge/License-AGPLv3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8.svg)](https://go.dev/)
[![Terraform](https://img.shields.io/badge/Terraform-1.1+-7B42BC.svg)](https://www.terraform.io/)

The Terraform Provider for [Fractal Cloud](https://fractal.cloud) enables platform and operations teams to manage Fractal Cloud governance resources through Terraform. It covers organizational structure, bounded contexts, environments, and fractal (blueprint) publication workflows.

Fractal Cloud is a platform engineering solution that delivers secure, compliant infrastructure across any cloud. It provides developers with ready-to-use building blocks and architecture templates while centralizing automation and governance for operations teams.

## Provider Scope

### Resources

| Resource | Description |
|---|---|
| `fc_personal_bounded_context` | Personal bounded contexts (scoped to a user account) |
| `fc_organizational_bounded_context` | Organizational bounded contexts (scoped to an organization) |
| `fc_fractal` | Fractal definitions (blueprints) -- publish, update, and manage component composition |
| `fc_management_environment` | Management environments for governance |
| `fc_operational_environment` | Operational environments for runtime governance |

### Data Sources

| Data Source | Description |
|---|---|
| `fc_personal_bounded_context` | Look up an existing personal bounded context |
| `fc_organizational_bounded_context` | Look up an existing organizational bounded context |
| `fc_organization` | Look up an existing organization |
| `fc_fractal` | Look up an existing fractal definition |

### Provider Functions

The provider includes 58 blueprint component builder functions organized by infrastructure domain and delivery model, plus `secret_ref` for referencing environment secrets. Function names follow the full component coordinate: `provider::fc::{domain}_{delivery_model}_{component}`. These functions create component objects for use in a fractal's `components` list. Dependencies between components are expressed as direct object references (type-checked at plan time) rather than string IDs.

Only `id` and the attributes a function's documentation marks as required need to be set; every other attribute may be omitted. Every function also accepts `extra_parameters`, a map of offer-specific parameters with no attribute of their own. A parameter value that is a JSON object or array string (from `jsonencode()` or `secret_ref()`) is sent as JSON; any other value is sent as a string.

<details>
<summary>NetworkAndCompute (7 functions)</summary>

| Function | Description |
|---|---|
| `provider::fc::network_and_compute_iaas_virtual_network` | Virtual network / VPC |
| `provider::fc::network_and_compute_iaas_subnet` | Subnet within a virtual network |
| `provider::fc::network_and_compute_iaas_load_balancer` | Load balancer |
| `provider::fc::network_and_compute_iaas_security_group` | Security group / firewall rules |
| `provider::fc::network_and_compute_iaas_virtual_machine` | Virtual machine / compute instance |
| `provider::fc::network_and_compute_paas_container_platform` | Container orchestration platform (e.g. Kubernetes) |
| `provider::fc::network_and_compute_saas_unmanaged` | External / unmanaged network resource |

</details>

<details>
<summary>CustomWorkloads (5 functions)</summary>

| Function | Description |
|---|---|
| `provider::fc::custom_workloads_caas_workload` | Generic container workload (CaaS) |
| `provider::fc::custom_workloads_iaas_workload` | IaaS workload |
| `provider::fc::custom_workloads_paas_workload` | PaaS workload |
| `provider::fc::custom_workloads_faas_workload` | FaaS / serverless workload |
| `provider::fc::custom_workloads_saas_unmanaged` | External / unmanaged workload resource |

</details>

<details>
<summary>Storage (17 functions)</summary>

| Function | Description |
|---|---|
| `provider::fc::storage_paas_files_and_blobs` | Object / blob / file storage |
| `provider::fc::storage_paas_relational_dbms` | Relational DBMS platform |
| `provider::fc::storage_paas_relational_database` | Relational database |
| `provider::fc::storage_paas_document_dbms` | Document DBMS platform |
| `provider::fc::storage_paas_document_database` | Document database |
| `provider::fc::storage_paas_column_oriented_dbms` | Column-oriented DBMS platform |
| `provider::fc::storage_paas_column_oriented_entity` | Column-oriented entity |
| `provider::fc::storage_paas_key_value_dbms` | Key-value DBMS platform |
| `provider::fc::storage_paas_key_value_entity` | Key-value entity |
| `provider::fc::storage_paas_graph_dbms` | Graph DBMS platform |
| `provider::fc::storage_paas_graph_database` | Graph database |
| `provider::fc::storage_caas_search` | Search platform |
| `provider::fc::storage_caas_search_entity` | Search entity / index |
| `provider::fc::storage_caas_relational_dbms` | Containerized relational DBMS (CloudNativePG) |
| `provider::fc::storage_caas_relational_database` | Database on a containerized DBMS |
| `provider::fc::storage_caas_object_storage` | Containerized object storage (MinIO) |
| `provider::fc::storage_saas_unmanaged` | External / unmanaged storage resource |

</details>

<details>
<summary>Messaging (5 functions)</summary>

| Function | Description |
|---|---|
| `provider::fc::messaging_paas_broker` | PaaS message broker |
| `provider::fc::messaging_paas_entity` | PaaS message broker entity (topic/queue) |
| `provider::fc::messaging_caas_broker` | CaaS message broker |
| `provider::fc::messaging_caas_entity` | CaaS message broker entity |
| `provider::fc::messaging_saas_unmanaged` | External / unmanaged messaging resource |

</details>

<details>
<summary>BigData (12 functions)</summary>

| Function | Description |
|---|---|
| `provider::fc::bigdata_paas_distributed_data_processing` | Distributed data processing platform |
| `provider::fc::bigdata_paas_compute_cluster` | Compute cluster |
| `provider::fc::bigdata_paas_data_processing_job` | Data processing job |
| `provider::fc::bigdata_paas_ml_experiment` | ML experiment |
| `provider::fc::bigdata_paas_datalake` | Data lake |
| `provider::fc::bigdata_caas_distributed_data_processing` | Spark on Kubernetes |
| `provider::fc::bigdata_caas_compute_cluster` | Spark cluster on Kubernetes |
| `provider::fc::bigdata_caas_data_processing_job` | Spark application on Kubernetes |
| `provider::fc::bigdata_caas_ml_experiment` | MLflow tracking server |
| `provider::fc::bigdata_caas_data_catalog` | Unity Catalog |
| `provider::fc::bigdata_caas_datalake` | Containerized data lake (MinIO) |
| `provider::fc::bigdata_saas_unmanaged` | External / unmanaged big data resource |

</details>

<details>
<summary>APIManagement (3 functions)</summary>

| Function | Description |
|---|---|
| `provider::fc::api_management_paas_api_gateway` | PaaS API gateway |
| `provider::fc::api_management_caas_api_gateway` | CaaS API gateway |
| `provider::fc::api_management_saas_unmanaged` | External / unmanaged API management resource |

</details>

<details>
<summary>Observability (4 functions)</summary>

| Function | Description |
|---|---|
| `provider::fc::observability_caas_monitoring` | Monitoring platform |
| `provider::fc::observability_caas_tracing` | Distributed tracing |
| `provider::fc::observability_caas_logging` | Logging platform |
| `provider::fc::observability_saas_unmanaged` | External / unmanaged observability resource |

</details>

<details>
<summary>Security (3 functions)</summary>

| Function | Description |
|---|---|
| `provider::fc::security_caas_service_mesh_security` | Service mesh security |
| `provider::fc::security_paas_identity_provider` | Identity provider (Cognito, Entra External ID) |
| `provider::fc::security_saas_unmanaged` | External / unmanaged security resource |

</details>

<details>
<summary>AI (2 functions)</summary>

| Function | Description |
|---|---|
| `provider::fc::ai_agentic_platform` | Governed runtime for AI agents |
| `provider::fc::ai_saas_unmanaged` | External AI / LLM service |

</details>

#### Secrets

A secret never goes into a blueprint. Define it on the environment, and reference it by short name with `provider::fc::secret_ref("<short-name>")`. The reference works in any parameter or link setting, and the agent resolves it from the environment secret store at reconciliation time. Unmanaged components require one:

```hcl
locals {
  openai = provider::fc::ai_saas_unmanaged({
    id     = "openai"
    secret = provider::fc::secret_ref("openai-api-key")
  })
}
```

### Not in Scope

This provider does **not** manage:

- **Live Systems** -- created by development teams via the Fractal SDK or UI
- Application-level deployments or runtime operations
- CI/CD pipeline-triggered infrastructure instantiation

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.8 (provider functions)
- [Go](https://golang.org/doc/install) >= 1.26 (to build the provider)
- A [Fractal Cloud](https://fractal.cloud) account with a service account

## Getting Started

### Provider Configuration

```hcl
terraform {
  required_providers {
    fc = {
      source  = "registry.terraform.io/fractal-cloud/fc"
      version = "~> 2.0"
    }
  }
  required_version = ">= 1.1.0"
}

provider "fc" {
  service_account_id     = var.fc_service_account_id
  service_account_secret = var.fc_service_account_secret
}
```

### Authentication

The provider authenticates using a Fractal Cloud service account. Credentials can be set via:

- **HCL attributes**: `service_account_id` and `service_account_secret`
- **Environment variables**: `FRACTAL_CLOUD_SERVICE_ACCOUNT_ID` and `FRACTAL_CLOUD_SERVICE_ACCOUNT_SECRET`

Environment variables are used as defaults and can be overridden by HCL attributes.

| Attribute | Environment Variable | Required | Description |
|---|---|---|---|
| `service_account_id` | `FRACTAL_CLOUD_SERVICE_ACCOUNT_ID` | Yes | Service account identifier |
| `service_account_secret` | `FRACTAL_CLOUD_SERVICE_ACCOUNT_SECRET` | Yes | Service account secret |
| `host` | -- | No | API endpoint (defaults to `https://api.fractal.cloud`) |

## Usage Examples

### Personal Bounded Context

```hcl
resource "fc_personal_bounded_context" "production" {
  short_name   = "production"
  display_name = "Production"
  description  = "Production bounded context"
}
```

### Organizational Bounded Context

```hcl
data "fc_organization" "my_org" {
  id = "your-organization-id"
}

resource "fc_organizational_bounded_context" "platform" {
  short_name      = "platform"
  organization_id = data.fc_organization.my_org.id
  display_name    = "Platform"
  description     = "Shared platform bounded context"
}
```

### Fractal with Blueprint Components

Components are defined as locals so that dependencies are expressed via direct object references (type-checked at plan time) instead of copy-pasted string IDs.

```hcl
locals {
  k8s_cluster = provider::fc::network_and_compute_paas_container_platform({
    id           = "k8s-cluster"
    display_name = "Kubernetes Cluster"
  })

  db_platform = provider::fc::storage_paas_relational_dbms({
    id             = "database-platform"
    display_name   = "Database Platform"
    engine_version = "16"
  })

  app_database = provider::fc::storage_paas_relational_database({
    id           = "app-database"
    display_name = "Application Database"
    dbms         = local.db_platform     # type-checked reference
  })

  api_service = provider::fc::custom_workloads_caas_workload({
    id              = "api-service"
    display_name    = "API Service"
    container_image = "my-registry/api-service:latest"
    container_port  = 8080
    cpu             = "500m"
    memory          = "1Gi"
    replicas        = 2
    platform        = local.k8s_cluster  # type-checked reference
    links = [
      { target = local.app_database, settings = { access = "read-write" } },
    ]
  })
}

resource "fc_fractal" "microservice" {
  bounded_context_id = fc_personal_bounded_context.production.id
  name        = "microservice-template"
  version     = "1.0"
  description = "Standard microservice architecture blueprint"

  components = [
    local.k8s_cluster,
    local.api_service,
    local.db_platform,
    local.app_database,
  ]
}
```

### IaaS Architecture with Dependencies and Links

```hcl
locals {
  main_vpc = provider::fc::network_and_compute_iaas_virtual_network({
    id           = "main-vpc"
    display_name = "Main VPC"
    cidr_block   = "10.0.0.0/16"
  })

  public_subnet = provider::fc::network_and_compute_iaas_subnet({
    id                = "public-subnet"
    display_name      = "Public Subnet"
    cidr_block        = "10.0.1.0/24"
    vpc               = local.main_vpc   # type-checked dependency
  })

  web_sg = provider::fc::network_and_compute_iaas_security_group({
    id          = "web-sg"
    description = "Allow HTTPS from the internet"
    vpc         = local.main_vpc         # type-checked dependency
    ingress_rules = [
      {
        from_port   = 443
        source_cidr = "0.0.0.0/0"
      }
    ]
  })

  web_server = provider::fc::network_and_compute_iaas_virtual_machine({
    id              = "web-server"
    display_name    = "Web Server"
    subnet          = local.public_subnet   # type-checked dependency
    security_groups = [local.web_sg]        # type-checked SG membership
  })
}

resource "fc_fractal" "iaas" {
  bounded_context_id = fc_personal_bounded_context.production.id
  name        = "iaas-architecture"
  version     = "1.0"
  description = "IaaS with network, security, and compute"

  components = [
    local.main_vpc,
    local.public_subnet,
    local.web_sg,
    local.web_server,
  ]
}
```

### Data Source Lookup

```hcl
data "fc_personal_bounded_context" "existing" {
  short_name = "production"
}

data "fc_fractal" "existing" {
  bounded_context_id = data.fc_personal_bounded_context.existing.id
  name    = "microservice-template"
  version = "1.0"
}
```

See the [`examples/`](examples/) directory for complete working configurations.

## Building the Provider

```sh
git clone https://github.com/Fractal-Cloud/terraform-provider-fractalcloud.git
cd terraform-provider-fractalcloud
make build
```

### Local Development

To use a locally built provider, add a dev override to your `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/fractal-cloud/fc" = "/path/to/your/go/bin"
  }
  direct {}
}
```

Then build and install:

```sh
make install
```

### Running Tests

```sh
# Unit tests
make test

# Acceptance tests (requires a Fractal Cloud account)
make testacc
```

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to get started.

## Security

To report a security vulnerability, please see [SECURITY.md](SECURITY.md).

## License

This project is licensed under the GNU Affero General Public License v3.0 -- see the [LICENSE](LICENSE) file for details.

[build-image]: https://github.com/Fractal-Cloud/terraform-provider-fractalcloud/actions/workflows/pr.yml/badge.svg
[build-url]: https://github.com/Fractal-Cloud/terraform-provider-fractalcloud/actions/workflows/pr.yml
[codecov-image]: https://codecov.io/gh/Fractal-Cloud/terraform-provider-fractalcloud/branch/main/graph/badge.svg
[codecov-url]: https://codecov.io/gh/Fractal-Cloud/terraform-provider-fractalcloud
