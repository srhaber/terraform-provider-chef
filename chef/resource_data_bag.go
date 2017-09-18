package chef

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform/helper/schema"

	chefc "github.com/go-chef/chef"
)

func resourceChefDataBag() *schema.Resource {
	return &schema.Resource{
		Create: CreateDataBag,
		Read:   ReadDataBag,
		Delete: DeleteDataBag,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

func CreateDataBag(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefc.Client)

	dataBag := &chefc.DataBag{
		Name: d.Get("name").(string),
	}

	result, err := client.DataBags.Create(dataBag)
	if err != nil {
		return err
	}

	d.SetId(dataBag.Name)
	d.Set("api_uri", result.URI)
	d.Set("items", &[]interface{}{})
	return nil
}

func ReadDataBag(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefc.Client)

	// The Chef API provides no API to read a data bag's metadata,
	// but we can try to read its items and use that as a proxy for
	// whether it still exists.

	name := d.Id()

	items, err := client.DataBags.ListItems(name)
	if err != nil {
		if errRes, ok := err.(*chefc.ErrorResponse); ok {
			if errRes.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}

	if _, ok := d.GetOk("api_uri"); !ok {
		d.Set("api_uri", fmt.Sprintf("%sdata/%s", client.BaseURL, name))
	}

	itemList := make([]string, 0, len(*items))
	for k := range *items {
		itemList = append(itemList, k)
	}

	sort.StringSlice(itemList).Sort()
	d.Set("items", itemList)

	return err
}

func DeleteDataBag(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefc.Client)

	name := d.Id()

	_, err := client.DataBags.Delete(name)
	if err == nil {
		d.SetId("")
	}
	return err
}
