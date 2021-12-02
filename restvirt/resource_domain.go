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
			"private_ip": {
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
			"nested_virtualization": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
		},
	}
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	resp, err := c.DomainServiceClient.CreateDomain(context.Background(), &pb.CreateDomainRequest{
		Domain: &pb.Domain{
			Name:                 d.Get("name").(string),
			Vcpu:                 uint32(d.Get("vcpu").(int)),
			Memory:               uint64(d.Get("memory").(int)),
			UserData:             d.Get("user_data").(string),
			PrivateIp:            d.Get("private_ip").(string),
			NestedVirtualization: d.Get("nested_virtualization").(bool),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.Uuid)

	resourceDomainRead(ctx, d, m)

	return diags
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	domain, err := c.GetDomain(d.Id())
	if _, notFound := err.(*restvirt.NotFoundError); notFound {
		d.SetId("")
		return diags
	}
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domain.UUID)
	d.Set("name", domain.Name)
	d.Set("vcpu", domain.VCPU)
	d.Set("memory", domain.MemoryMiB)
	d.Set("user_data", userDataHashSum(domain.UserData))
	d.Set("private_ip", domain.PrivateIP)

	return diags
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*restvirt.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	err := c.DeleteDomain(d.Id())
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
