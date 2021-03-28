package restvirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/verbit/restvirt-client"
)

func resourceForwarding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceForwardingCreate,
		ReadContext:   resourceForwardingRead,
		DeleteContext: resourceForwardingDelete,
		Schema: map[string]*schema.Schema{
			"source_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"target_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"target_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceForwardingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	forwarding := restvirt.PortForwarding{
		SourcePort: uint16(d.Get("source_port").(int)),
		TargetIP:   d.Get("target_ip").(string),
		TargetPort: uint16(d.Get("target_port").(int)),
		Protocol:   d.Get("protocol").(string),
	}

	id, err := c.CreatePortForwarding(forwarding)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	resourceForwardingRead(ctx, d, m)

	return diags
}

func resourceForwardingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	forwarding, err := c.GetPortForwarding(d.Id())
	if _, notFound := err.(*restvirt.NotFoundError); notFound {
		d.SetId("")
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Id())
	d.Set("source_port", forwarding.SourcePort)
	d.Set("target_ip", forwarding.TargetIP)
	d.Set("target_port", forwarding.TargetPort)
	d.Set("protocol", forwarding.Protocol)

	return diags
}

func resourceForwardingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	err := c.DeletePortForwarding(d.Id())
	if _, notFound := err.(*restvirt.NotFoundError); notFound {
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
