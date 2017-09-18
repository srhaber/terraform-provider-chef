package chef

import (
	"encoding/json"
	"fmt"

	chefc "github.com/go-chef/chef"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/srhaber/chefutil/datacrypt"
)

func dataSourceChefDataBagItem() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceChefDataBagItemRead,

		Schema: map[string]*schema.Schema{
			"data_bag": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"item_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"encryption_key": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"content_json": &schema.Schema{
				Type:      schema.TypeString,
				Computed:  true,
				StateFunc: jsonStateFunc,
			},
		},
	}
}

func dataSourceChefDataBagItemRead(d *schema.ResourceData, meta interface{}) error {
	dataBag := d.Get("data_bag").(string)
	itemId := d.Get("item_id").(string)

	d.SetId(fmt.Sprintf("%s/%s", dataBag, itemId))

	client := meta.(*chefc.Client)
	item, err := client.DataBags.GetItem(dataBag, itemId)
	if err != nil {
		if errRes, ok := err.(*chefc.ErrorResponse); ok {
			if errRes.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		} else {
			return err
		}
	}

	if key, ok := d.GetOk("encryption_key"); ok {
		// Build a DataDecryptor
		decryptor := &datacrypt.DataDecryptor{
			Item:   item.(map[string]interface{}),
			Secret: []byte(key.(string)),
		}

		// Decrypt the item
		value, err := decryptor.Decrypt()
		if err != nil {
			return err
		}

		jsonContent, err := json.Marshal(value)
		if err != nil {
			return err
		}

		d.Set("content_json", string(jsonContent))
	} else {
		jsonContent, err := json.Marshal(item)
		if err != nil {
			return err
		}

		d.Set("content_json", string(jsonContent))
	}

	return nil
}
