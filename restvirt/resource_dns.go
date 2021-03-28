package restvirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/verbit/restvirt-client"
)

func resourceDNSMapping() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSMappingCreate,
		ReadContext:   resourceDNSMappingRead,
		DeleteContext: resourceDNSMappingDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDNSMappingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	dns := restvirt.DNSMapping{
		Name: d.Get("name").(string),
		IP:   d.Get("ip").(string),
	}

	err := c.SetDNSMapping(dns)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dns.Name)

	resourceDNSMappingRead(ctx, d, m)

	return diags
}

func resourceDNSMappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	mapping, err := c.GetDNSMapping(d.Id())
	if _, notFound := err.(*restvirt.NotFoundError); notFound {
		d.SetId("")
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Id())
	d.Set("name", mapping.Name)
	d.Set("ip", mapping.IP)

	return diags
}

func resourceDNSMappingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	err := c.DeleteDNSMapping(d.Id())
	if _, notFound := err.(*restvirt.NotFoundError); notFound {
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
