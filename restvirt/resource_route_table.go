package restvirt

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/verbit/restvirt-client"
	"github.com/verbit/restvirt-client/pb"
)

func resourceRouteTable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRouteTableCreate,
		ReadContext:   resourceRouteTableRead,
		DeleteContext: resourceRouteTableDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceRouteTableCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	var diags diag.Diagnostics

	table := &pb.RouteTable{
		Name: d.Get("name").(string),
	}
	table, err := c.RouteServiceClient.CreateRouteTable(ctx, &pb.CreateRouteTableRequest{RouteTable: table})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(table.Id)))

	resourceRouteTableRead(ctx, d, m)

	return diags
}

func resourceRouteTableRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id, err := strconv.ParseUint(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	table, err := c.RouteServiceClient.GetRouteTable(ctx, &pb.RouteTableIdentifier{Id: uint32(id)})
	if _, notFound := err.(*restvirt.NotFoundError); notFound {
		d.SetId("")
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", table.Name)

	return diags
}

func resourceRouteTableDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id, err := strconv.ParseUint(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = c.RouteServiceClient.DeleteRouteTable(ctx, &pb.RouteTableIdentifier{Id: uint32(id)})
	if _, notFound := err.(*restvirt.NotFoundError); notFound {
		d.SetId("")
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
