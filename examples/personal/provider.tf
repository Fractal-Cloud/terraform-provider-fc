terraform {
  required_providers {
    fc = {
      source = "registry.terraform.io/fractal-cloud/fc"
    }
  }
  required_version = ">= 1.1.0"
}

provider "fc" {
}
