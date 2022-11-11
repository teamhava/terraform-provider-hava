---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "hava_source_aws_key_resource Resource - terraform-provider-hava"
subcategory: ""
description: |-
  A Source in Hava using AWS access key id and secret key to authenticate to the AWS account that will be imported.
---

# hava_source_aws_key_resource (Resource)

A Source in Hava using AWS access key id and secret key to authenticate to the AWS account that will be imported.

## Example Usage

```terraform
resource "hava_source_aws_key_resource" "example" {
  name       = "Example Source"
  access_key = "xxx"
  secret_key = "xxx"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `access_key` (String, Sensitive) The aws access key id of the account that will be used to access the source for import
- `name` (String) Display name of the source
- `secret_key` (String, Sensitive) The aws secret key of the account that will be used to access the source for import

### Read-Only

- `id` (String) The ID of this resource.
- `state` (String) State of the Source

