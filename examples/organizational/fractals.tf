data "fc_fractal" "existing_fractal" {
  bounded_context_id = data.fc_organizational_bounded_context.existing_bounded_context.id
  name               = "existing-fractal"
  version            = "1.0"
}

# Multi-tier IaaS Fractal with full dependency and link wiring.
# Components are defined as locals so that dependencies and links use
# direct object references (type-checked) instead of string IDs.
locals {
  org_main_vpc = provider::fc::network_and_compute_iaas_virtual_network({
    id           = "main-vpc"
    display_name = "Main VPC"
    description  = "Primary VPC for the organization IaaS architecture"
    cidr_block   = "10.0.0.0/16"
  })

  org_web_subnet = provider::fc::network_and_compute_iaas_subnet({
    id           = "web-subnet"
    display_name = "Web Tier Subnet"
    description  = "Web tier subnet"
    cidr_block   = "10.0.1.0/24"
    vpc          = local.org_main_vpc
  })

  org_app_subnet = provider::fc::network_and_compute_iaas_subnet({
    id           = "app-subnet"
    display_name = "App Tier Subnet"
    description  = "Application tier subnet"
    cidr_block   = "10.0.2.0/24"
    vpc          = local.org_main_vpc
  })

  org_web_sg = provider::fc::network_and_compute_iaas_security_group({
    id           = "web-sg"
    display_name = "Web Security Group"
    description  = "Allow HTTPS from internet"
    vpc          = local.org_main_vpc
    ingress_rules = [
      {
        from_port   = 443
        source_cidr = "0.0.0.0/0"
      }
    ]
  })

  # No ingress rules: traffic from the web tier is admitted by the
  # web-server -> app-server link below, from which the agent derives the
  # rules on both servers' managed security groups.
  org_app_sg = provider::fc::network_and_compute_iaas_security_group({
    id           = "app-sg"
    display_name = "App Security Group"
    description  = "Application tier; reachable only through links"
    vpc          = local.org_main_vpc
  })

  # App server: depends on app-subnet, member of app-sg
  org_app_server = provider::fc::network_and_compute_iaas_virtual_machine({
    id              = "app-server"
    display_name    = "Application Server"
    description     = "Application tier server"
    subnet          = local.org_app_subnet
    security_groups = [local.org_app_sg]
    links           = []
  })

  # Web server: depends on web-subnet, member of web-sg, links to app-server on port 8080
  org_web_server = provider::fc::network_and_compute_iaas_virtual_machine({
    id              = "web-server"
    display_name    = "Web Server"
    description     = "Web tier server"
    subnet          = local.org_web_subnet
    security_groups = [local.org_web_sg]
    links = [
      {
        target = local.org_app_server
        settings = {
          fromPort = 8080
        }
      }
    ]
  })
}

resource "fc_fractal" "org_iaas" {
  bounded_context_id = data.fc_organizational_bounded_context.existing_bounded_context.id
  name               = "org-iaas"
  version            = "1.0"
  description        = "Organization multi-tier IaaS architecture"

  components = [
    local.org_main_vpc,
    local.org_web_subnet,
    local.org_app_subnet,
    local.org_web_sg,
    local.org_app_sg,
    local.org_web_server,
    local.org_app_server,
  ]
}

output "fractal" {
  value = data.fc_fractal.existing_fractal
}
