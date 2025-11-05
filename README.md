# Terraform Provider

The [Terraform provider](https://registry.terraform.io/providers/SovraDev/spinnaker/latest/docs) allowing [Terraform](https://developer.hashicorp.com/terraform) to manage [Spinnaker](https://spinnaker.io/docs/) resources.

## Spinnaker

https://spinnaker.io/concepts/

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.x
- [Go](https://golang.org/doc/install) 1.24 (to build the provider plugin)

## Building and Developing The Provider

```sh
$ git clone git@github.com:sovradev/terraform-provider-spinnaker.git
$ cd terraform-provider-spinnaker/
$ go build
$ go test ./...
```

### Testing the provider

If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory, run `terraform init` to initialize it.

## Status
This provider is in early active development. The goal of this project is to provide a fully and up-to-date featured Spinnaker provider for Terraform.

Please open issues for bugs or feature requests.
