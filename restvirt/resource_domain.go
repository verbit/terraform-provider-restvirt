package restvirt

import (
	"context"
	"crypto/sha1"
	"encoding/hex"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/verbit/restvirt-client"
	"github.com/verbit/restvirt-client/pb"
)

func userDataHashSum(d interface{}) string {
	hash := sha1.Sum([]byte(d.(string)))
	return hex.EncodeToString(hash[:])
}

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		UpdateContext: resourceDomainUpdate,
		DeleteContext: resourceDomainDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vcpu": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ipv6_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user_data": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				StateFunc: userDataHashSum,
			},
		},
	}
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	domain := &pb.Domain{
		Name:        d.Get("name").(string),
		Vcpu:        uint32(d.Get("vcpu").(int)),
		Memory:      uint64(d.Get("memory").(int)),
		UserData:    d.Get("user_data").(string),
		Network:     d.Get("network_id").(string),
		PrivateIp:   d.Get("private_ip").(string),
		Ipv6Address: d.Get("ipv6_address").(string),
	}

	domain, err := c.DomainServiceClient.CreateDomain(ctx, &pb.CreateDomainRequest{
		Host:   d.Get("host").(string),
		Domain: domain,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domain.Uuid)

	resourceDomainRead(ctx, d, m)

	return diags
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	domain, err := c.DomainServiceClient.GetDomain(ctx, &pb.GetDomainRequest{
		Host: d.Get("host").(string),
		Uuid: d.Id(),
	})
	if _, notFound := err.(*restvirt.NotFoundError); notFound {
		d.SetId("")
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domain.Uuid)
	d.Set("host", d.Get("host"))
	d.Set("name", domain.Name)
	d.Set("vcpu", domain.Vcpu)
	d.Set("memory", domain.Memory)
	d.Set("user_data", userDataHashSum(domain.UserData))
	d.Set("network_id", domain.Network)
	d.Set("private_ip", domain.PrivateIp)
	d.Set("ipv6_address", domain.Ipv6Address)

	return diags
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	_, err := c.DomainServiceClient.DeleteDomain(ctx, &pb.DeleteDomainRequest{
		Host: d.Get("host").(string),
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

func resourceDomainUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	return nil
}
