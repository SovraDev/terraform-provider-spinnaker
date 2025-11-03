# Terraform Provider

Terraform provider for managing Spinnaker applications and pipelines.

- Website: https://www.terraform.io

## Spinnaker

https://spinnaker.io/concepts/

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.x
- [Go](https://golang.org/doc/install) 1.14 (to build the provider plugin)

## Building and Developing The Provider

```sh
$ git clone git@github.com:sovradev/terraform-provider-spinnaker.git
$ cd terraform-provider-spinnaker/
$ go build
$ go test ./...
```

## Using the provider

If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory, run `terraform init` to initialize it.
