# S3 Provider

The S3 provider is used to interact with S3 resources supported by S3-compatible
cloud storage providers, including AWS, Ceph, Google storage, etc. The provider
needs to be configured with the proper credentials before it can be used.

## Example Usage

```hcl
# Configure the S3 provider.
provider "s3" {
  endpoint = "http://localhost:9000"
  region   = "eu-west-3"
}
```

## Authentication

The S3 provider offers a flexible means of providing credentials for
authentication. The following methods are supported, in this order, and
explained below:

- Static credentials
- Environment variables
- Shared credentials file

### Static Credentials

!> **WARNING**: Hard-coding credentials into any Terraform configuration is not
recommended, and risks secret leakage should this file ever be committed to a
public version control system.

Static credentials can be provided by adding both `access_key` and `secret_key`
in-line in the S3 provider block:

Usage:

```hcl
provider "s3" {
  region     = "eu-west-3"
  access_key = "my-access-key"
  secret_key = "my-secret-key"
}
```

### Environment variables

You can provide your credentials via the `AWS_ACCESS_KEY_ID` and
`AWS_SECRET_ACCESS_KEY`, environment variables, representing your AWS Access Key
and AWS Secret Key, respectively. Note that setting your AWS credentials using
either these (or legacy) environment variables will override the use of
`AWS_SHARED_CREDENTIALS_FILE` and `AWS_PROFILE`. The `AWS_DEFAULT_REGION` and
`AWS_SESSION_TOKEN` environment variables are also used, if applicable:

```hcl
provider "s3" {}
```

Usage:

```sh
$ export AWS_ACCESS_KEY_ID="my-access-key"
$ export AWS_SECRET_ACCESS_KEY="my-secret-key"
$ export AWS_DEFAULT_REGION="eu-west-3"
$ terraform plan
```

### Shared Credentials file

You can use an AWS credentials file to specify your credentials. The default
location is `$HOME/.aws/credentials` on Linux and OS X, or
`"%USERPROFILE%\.aws\credentials"` for Windows users. If we fail to detect
credentials inline, or in the environment, Terraform will check this location.
You can optionally specify a different location in the configuration by
providing the `shared_credentials_file` attribute, or in the environment with
the `AWS_SHARED_CREDENTIALS_FILE` variable. This method also supports a profile
configuration and matching `AWS_PROFILE` environment variable:

Usage:

```hcl
provider "aws" {
  region                  = "eu-west-3"
  shared_credentials_file = "/home/johndoe/.aws/creds"
  profile                 = "myprofile"
}
```

## Argument Reference

* `access_key` - (Optional) This is the AWS access key. It must be provided, but
  it can also be sourced from the `AWS_ACCESS_KEY_ID` environment variable, or
  via a shared credentials file if profile is specified.

* `secret_key` - (Optional) This is the AWS secret key. It must be provided, but
  it can also be sourced from the `AWS_SECRET_ACCESS_KEY` environment variable,
  or via a shared credentials file if profile is specified.

* `shared_credentials_file` - (Optional) This is the path to the shared
  credentials file. If this is not set and a profile is specified,
  `~/.aws/credentials` will be used.

* `profile` - (Optional) This is the AWS profile name as set in the shared
  credentials file.

* `token` - (Optional) Session token for validating temporary credentials.
  Typically provided after successful identity federation or Multi-Factor
  Authentication (MFA) login. With MFA login, this is the session token provided
  afterwards, not the 6 digit MFA code used to get temporary credentials. It can
  also be sourced from the `AWS_SESSION_TOKEN` environment variable.

* `region`: (Optional) This is the AWS region. It must be provided, but it can
  also be sourced from the `AWS_DEFAULT_REGION` environment variables, or via a
  shared credentials file if profile is specified.

* `endpoint` - (Required) The S3 endpoint.

* `force_path_style` - (Optional) Set this to true to force the request to use
  path-style addressing, i.e., `http://s3.amazonaws.com/BUCKET/KEY`. By default,
  the S3 client will use virtual hosted bucket addressing,
  `http://BUCKET.s3.amazonaws.com/KEY`, when possible.

* `tls_enabled` - (Optional) Set this to false to explicitly allow the provider
  to perform "insecure" SSL requests. If omitted, default value is true.

* `tls_skip_verify` - (Optional) Set this to true to accept any certificate
  presented by the endpoint and any host name in that certificate.

* `tls_cert_path` - (Optional) Set this to the path of the PEM file holding the
  certificates to be added to the system cert pool.

* `aws_signature_version` - (Optional) AWS API request signature version to use
  to sign API requests. Can be "v2" or "v4". Defaults to "v4".
