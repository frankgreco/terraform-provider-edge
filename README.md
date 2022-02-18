[![acctest](https://github.com/frankgreco/terraform-provider-edge/actions/workflows/acctest.yml/badge.svg)](https://github.com/frankgreco/terraform-provider-edge/actions/workflows/acctest.yml)

# EdgeOS Terraform Provider (terraform-provider-edge)
> NOTE: This provider is under rapid active development.

Terraform wrapper for [edge-sdk-go](https://github.com/frankgreco/edge-sdk-go).

## Documentation
You can browse documentation on the [Terraform provider registry](https://registry.terraform.io/providers/frankgreco/edge/latest/docs).

## Supported EdgeOS Versions
The only version i've tested this against is the version that I use, `v2.0.9`. I plan on making a full compatability matrix as I get further into development.

## Using the Provider
I believe anything `v1.0` or newer will work.

```
terraform {
  required_providers {
    edge = {
      source  = "frankgreco/edge"
      version = "0.1.6"
    }
  }
}

provider "edge" {
    ...
}
```
