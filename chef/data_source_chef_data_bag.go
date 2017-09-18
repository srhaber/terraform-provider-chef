package chef

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceChefDataBag() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceChefDataBagRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "(Required) The unique name to assign to the data bag. This is the name that other server clients will use to find and retrieve data from the data bag.",
				Required:    true,
			},
			"items": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"api_uri": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceChefDataBagRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	d.SetId(name)
	return ReadDataBag(d, meta)
}
