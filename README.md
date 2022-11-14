# Hava Terraform Provider

The Hava Terraform provider makes it simple to integrate Hava into your GitOps workflows, allowing you to add your cloud environments to Hava for automatic documentation together with the infrastructure of code you use to manage those environments.

This provider is built and tested with Terraform v1.3.4 or later

- [Terraform Website](https://www.terraform.io/)
- [Hava Provider Documentation](https://registry.terraform.io/providers/teamhava/hava/latest/docs)

## Usage Example
The below example show how the Terraform provider can be used to configure an AWS account source using a cross account role.

```hcl
terraform {
  required_providers {
    hava = {
      source = "teamhava/hava"
      version = "~> 0.1"
    }

    aws = {
      source = "hashicorp/aws"
      version = "~> 4.39"
    }
  }
}

// Get the ARN for the AWS Read Only Managed Policy
data "aws_iam_policy" "example" {
  name = "ReadOnlyAccess"
}

// Create the role that will be used for cross account role accesss
resource "aws_iam_role" "hava_ro" {
  name                = "hava-read-only-role"
  assume_role_policy  = jsonencode({
      "Version": "2012-10-17",
      "Statement": [
          {
              "Effect": "Allow",
              "Principal": {
                  // Hava CAR account
                  "AWS": "arn:aws:iam::281013829959:root"
              },
              "Action": "sts:AssumeRole",
              "Condition": {
                  "StringEquals": {
                      // unique id for your Hava account, 
                      "sts:ExternalId": var.external_id
                  }
              }
          }
      ]
    })
  
  managed_policy_arns = [data.aws_iam_policy.example.arn]
}

// 
resource "hava_source_aws_car_resource" "example" {
  name        = "Example Source"
  role_arn    = aws_iam_role.hava_ro.arn 
  external_id = var.external_id
}
```

## Developer Requirements

This repository uses developer containers to make sure everyone has the same development environment. It's highly recommended to use the included containers.

### Dependencies

- Golang 1.19
- Terraform v1.3.4

### Environment setup

To effectively develop a terraform provider, you need to be able to test it locally. Hashicorp has an config file `~/.terraformrc` that lives in your home directory, where you can override certain things, like where to download the provider from.

The bellow config file overrides the path to look for the `teamhava/hava provider` to look in the golang bin folder `/go/bin`. This allows us to run `go install` to compile our provider, and it will be compiled into the bin directory or us to use it when executing terraform.

```hcl
provider_installation {

  dev_overrides {
      "registry.terraform.io/teamhava/hava" = "/go/bin/"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}

```