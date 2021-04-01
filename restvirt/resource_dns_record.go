package restvirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/verbit/restvirt-client"
)

func resourceDNSRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSRecordCreate,
		UpdateContext: resourceDNSRecordCreate,
		ReadContext:   resourceDNSRecordRead,
		DeleteContext: resourceDNSRecordDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"records": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				Set:      schema.HashString,
			},
		},
	}
}

func resourceDNSRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	tfRecords := d.Get("records").(*schema.Set).List()
	records := make([]string, len(tfRecords))
	for i, rec := range tfRecords {
		records[i] = rec.(string)
	}
	dns := restvirt.DNSRecord{
		Name:    d.Get("name").(string),
		Type:    d.Get("type").(string),
		TTL:     d.Get("ttl").(int),
		Records: records,
	}

	err := c.SetDNSRecord(dns)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dns.ID())

	resourceDNSRecordRead(ctx, d, m)

	return diags
}

func resourceDNSRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	mapping, err := c.GetDNSRecordByID(d.Id())
	if _, notFound := err.(*restvirt.NotFoundError); notFound {
		d.SetId("")
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Id())
	d.Set("name", mapping.Name)
	d.Set("type", mapping.Type)
	d.Set("ttl", mapping.TTL)
	d.Set("records", mapping.Records)

	return diags
}

func resourceDNSRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	err := c.DeleteDNSRecordByID(d.Id())
	if _, notFound := err.(*restvirt.NotFoundError); notFound {
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
