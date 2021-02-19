terraform {
  required_providers {
    restvirt = {
      version = "~> 0.0.1"
      source  = "github.com/verbit/restvirt"
    }
  }
}

provider "restvirt" {
}
