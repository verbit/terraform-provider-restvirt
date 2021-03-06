# Terraform Provider restvirt

## Build provider

Run the following command to build the provider

```shell
$ go build -o terraform-provider-restvirt
```

## Test sample configuration

First, build and install the provider.

```shell
$ make install
$ make && (cd examples && rm -rf .terraform .terraform.lock.hcl && terraform init)
```

Then, navigate to the `examples` directory. 

```shell
$ cd examples
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ terraform init && terraform apply
```


## Installation

If you have `jq` installed, you can use this one-liner to install the latest version of this plugin to `~/.terraform.d/plugins/`:
```shell
curl -sfL https://raw.githubusercontent.com/verbit/terraform-provider-restvirt/main/hack/install_latest.sh | sh -
```
Take a look at [hack/install_latest.sh](hack/install_latest.sh) to see what the script does.
