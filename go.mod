module github.com/verbit/terraform-provider-restvirt

go 1.16

require (
	github.com/hashicorp/terraform-plugin-docs v0.4.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.4.4
	github.com/verbit/restvirt-client v0.5.0
)

replace github.com/verbit/restvirt-client => ../restvirt-client
