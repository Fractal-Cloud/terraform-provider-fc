terraform {
  required_providers {
    fc = {
      source = "registry.terraform.io/fractal-cloud/fc"
    }
  }
}

provider "fc" {
  service_account_id     = "xxx"
  service_account_secret = "xxx"
}
