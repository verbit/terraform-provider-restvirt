package restvirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/verbit/restvirt-client"
	"github.com/verbit/restvirt-client/pb"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		DeleteContext: resourceNetworkDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	network := &pb.Network{
		Name: d.Get("name").(string),
		Cidr: d.Get("cidr").(string),
	}

	network, err := c.DomainServiceClient.CreateNetwork(ctx, &pb.CreateNetworkRequest{
		Network: network,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(network.Uuid)

	resourceNetworkRead(ctx, d, m)

	return diags
}

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	network, err := c.DomainServiceClient.GetNetwork(ctx, &pb.GetNetworkRequest{
		Uuid: d.Id(),
	})
	if _, notFound := err.(*restvirt.NotFoundError); notFound {
		d.SetId("")
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(network.Uuid)
	d.Set("name", network.Name)
	d.Set("cidr", network.Cidr)

	return diags
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	_, err := c.DomainServiceClient.DeleteNetwork(ctx, &pb.DeleteNetworkRequest{
		Uuid: d.Id(),
	})
	if _, notFound := err.(*restvirt.NotFoundError); notFound {
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
