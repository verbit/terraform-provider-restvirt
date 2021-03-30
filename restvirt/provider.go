package restvirt

import (
	"context"
	"fmt"

	"github.com/verbit/restvirt-client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"ca": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"restvirt_domain":     resourceDomain(),
			"restvirt_dns_record": resourceDNSRecord(),
			"restvirt_forwarding": resourceForwarding(),
			"restvirt_volume":     resourceVolume(),
			"restvirt_attachment": resourceAttachment(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"restvirt_domain": dataSourceDomain(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// host := d.Get("host").(string)
	// username := d.Get("username").(string)
	// password := d.Get("password").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// c, err := restvirt.NewClient(host, username, password)
	c, err := restvirt.NewClientFromEnvironment()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create restvirt client",
			Detail:   fmt.Sprintf("Unable to create restvirt client: %v", err),
		})
		return nil, diags
	}

	return c, diags
}
