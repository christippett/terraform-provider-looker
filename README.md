# Looker Terraform Provider (⚠️ WIP)

This provider enables various parts of a Looker installation to be configured via Terraform.

The provider leverages the Go Looker SDK: [**GoLook**](https://github.com/looker-open-source/sdk-codegen/tree/main/go).

## Build provider

Run the following command to build the provider

```shell
$ go build -o terraform-provider-looker
```

## Test sample configuration

First, build and install the provider.

```shell
$ make install
```

Then, navigate to the `examples` directory.

```shell
$ cd examples
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ terraform init && terraform apply
```
