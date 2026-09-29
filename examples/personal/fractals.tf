data "fc_fractal" "existing_fractal" {
  bounded_context_id = data.fc_personal_bounded_context.existing_bounded_context.id
  name               = "existing-fractal"
  version            = "1.0"
}

# IaaS Fractal with VPC, subnets, security groups, and VMs linked together.
# Components are defined as locals so that dependencies and links use
# direct object references (type-checked) instead of string IDs.
locals {
  main_vpc = provider::fc::network_and_compute_iaas_virtual_network({
    id           = "main-vpc"
    display_name = "Main VPC"
    description  = "Primary VPC for the IaaS fractal"
    cidr_block   = "10.0.0.0/16"
  })

  public_subnet = provider::fc::network_and_compute_iaas_subnet({
    id           = "public-subnet"
    display_name = "Public Subnet"
    description  = "Public-facing subnet"
    cidr_block   = "10.0.1.0/24"
    vpc          = local.main_vpc # type-checked dependency on VPC
  })

  web_sg = provider::fc::network_and_compute_iaas_security_group({
    id           = "web-sg"
    display_name = "Web Security Group"
    description  = "Allow HTTPS from the internet"
    vpc          = local.main_vpc # type-checked dependency on VPC
    ingress_rules = [
      {
        from_port   = 443
        source_cidr = "0.0.0.0/0"
      }
    ]
  })

  api_server = provider::fc::network_and_compute_iaas_virtual_machine({
    id              = "api-server"
    display_name    = "API Server"
    description     = "Backend API server"
    subnet          = local.public_subnet
    security_groups = [local.web_sg]
    links           = []
  })

  web_server = provider::fc::network_and_compute_iaas_virtual_machine({
    id              = "web-server"
    display_name    = "Web Server"
    description     = "Frontend web server"
    subnet          = local.public_subnet # type-checked dependency on subnet
    security_groups = [local.web_sg]      # type-checked SG membership link
    links = [
      {
        target = local.api_server
        settings = {
          fromPort = "8080"
        }
      }
    ]
  })
}

resource "fc_fractal" "iaas_fractal" {
  bounded_context_id = data.fc_personal_bounded_context.existing_bounded_context.id
  name               = "basic-iaas"
  version            = "1.0"
  description        = "IaaS Fractal with network, security, and compute"

  components = [
    local.main_vpc,
    local.public_subnet,
    local.web_sg,
    local.web_server,
    local.api_server,
  ]
}

# Container Fractal with workloads linked for inter-service traffic
locals {
  k8s_cluster = provider::fc::network_and_compute_paas_container_platform({
    id           = "k8s-cluster"
    display_name = "Kubernetes Cluster"
    description  = "Managed Kubernetes cluster for microservices"
    node_pools   = []
  })

  db_service = provider::fc::custom_workloads_caas_workload({
    id              = "db-service"
    display_name    = "Database Service"
    description     = "PostgreSQL database service"
    container_name  = "postgres"
    container_image = "postgres:16"
    container_port  = 5432
    cpu             = "1024"
    memory          = "2048"
    replicas        = 1
    platform        = local.k8s_cluster
    subnet          = null
    links           = []
    security_groups = []
  })

  api_service = provider::fc::custom_workloads_caas_workload({
    id              = "api-service"
    display_name    = "API Service"
    description     = "REST API backend service"
    container_name  = "api"
    container_image = "my-registry/api-service:latest"
    container_port  = 8080
    cpu             = "512"
    memory          = "1024"
    replicas        = 2
    platform        = local.k8s_cluster # type-checked dependency on platform
    subnet          = null
    security_groups = []
    links = [
      {
        target = local.db_service
        settings = {
          fromPort = "5432"
          protocol = "tcp"
        }
      }
    ]
  })
}

resource "fc_fractal" "container_fractal" {
  bounded_context_id = data.fc_personal_bounded_context.existing_bounded_context.id
  name               = "microservice"
  version            = "1.0"
  description        = "Containerized Microservice Fractal"

  components = [
    local.k8s_cluster,
    local.api_service,
    local.db_service,
  ]
}

output "fractal" {
  value = data.fc_fractal.existing_fractal
}
