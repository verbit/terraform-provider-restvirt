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
