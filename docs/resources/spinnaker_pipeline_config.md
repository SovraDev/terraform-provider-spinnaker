---
page_title: "spinnaker_pipeline_config"
---

# spinnaker_pipeline_config Resource

Manage spinnaker pipeline configuration

> [! Note]
> As of this version, only the `webhook` trigger type is supported.


## Example Usage

```
provider "spinnaker" {
    server = "http://spinnaker-gate.myorg.io"
}

resource "spinnaker_application" "terraform_example" {
    application = "terraformexample"
    email       = "user@example.com"
}

resource "spinnaker_pipeline" "terraform_example" {
    application = spinnaker_application.terraform_example.application
    name        = "Example Pipeline"
    pipeline    = file("pipelines/example.json")
}

resource "spinnaker_pipeline_config" "example" {
    application = spinnaker_application.terraform_example.name
    pipeline    = spinnaker_pipeline.terraform_example.name

    trigger {
        type                = "webhook"
        enabled             = true
        source              = "deploy"
        payload_constraints = {}
    }
}
```

## Argument Reference

- `application` - (Required) Spinnaker application name.
- `pipeline` - (Required) Pipeline name.
- [trigger](#trigger-configuration-block) - (Optional) Configuration block used to specify trigger information about the related pipeline. Detailed below.

## trigger Configuration Block
Supported nested attributes for the `trigger` block:

- `type` - (Required) The type of trigger. e.g. `webhook`, `cron`, `pipeline` or `docker_registry`.
- `enabled` - (Optional) Whether the trigger is enabled. (Default: `true`)
- `source` - (Optional) Determines the target URL required to trigger this pipeline. Only used when `type` is `webhook`.
- `payload_constraints` - (Optional) When provided, only a webhook with a payload containing at least the specified key/value pairs will be allowed to trigger this pipeline. Only used when `type` is `webhook`. (Default: `{}`)


## Attribute Reference

In addition to the above, the following attributes are exported:

- `triggers` - List of pipeline trigger configuration
