# S3-compatible Terraform Provider

The S3 provider is used to interact with S3 resources supported by S3-compatible
cloud storage providers, including AWS, Ceph, Google storage, etc.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html)
- [Go](https://golang.org/doc/install) (to build the provider plugin)

## Developing the Provider

You need [Go](http://www.golang.org) fully installed and configured on your
machine before proceeding. The instructions that follow assume a directory in
your home directory outside of the standard GOPATH (e.g.
`$HOME/development/terraform-providers/`).

Clone repository to: `$HOME/development/terraform-providers/`

```sh
$ mkdir -p $HOME/development/terraform-providers/
$ cd $HOME/development/terraform-providers/
$ git clone git@github.com:b9company/terraform-provider-s3
...
```

To build the provider, run `make`.

```sh
$ make
...
$ ./terraform-provider-s3
...
```

## Using the Provider

To use a built provider in your Terraform environment (e.g. the provider binary
from the build instructions above), follow the instructions to [install it as a
plugin](https://www.terraform.io/docs/plugins/basics.html#installing-plugins).
After placing the custom-built provider into your plugins directory, run
`terraform init` to initialize it.
