package restvirt

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/verbit/restvirt-client"
)

func resourceAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAttachmentCreate,
		ReadContext:   resourceAttachmentRead,
		DeleteContext: resourceAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"disk_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAttachmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	domainID := d.Get("domain_id").(string)
	volumeID := d.Get("volume_id").(string)

	_, err := c.CreateAttachment(domainID, volumeID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%s", domainID, volumeID))

	resourceAttachmentRead(ctx, d, m)

	return diags
}

func resourceAttachmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	attachmentID := strings.Split(d.Id(), ":")
	domainID := attachmentID[0]
	volumeID := attachmentID[1]

	attachment, err := c.GetAttachment(domainID, volumeID)
	if _, notFound := err.(*restvirt.NotFoundError); notFound {
		d.SetId("")
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Id())
	d.Set("domain_id", domainID)
	d.Set("volume_id", volumeID)
	d.Set("disk_address", attachment.DiskAddress)

	return diags
}

func resourceAttachmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	attachmentID := strings.Split(d.Id(), ":")
	domainID := attachmentID[0]
	volumeID := attachmentID[1]

	err := c.DeleteAttachment(domainID, volumeID)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
