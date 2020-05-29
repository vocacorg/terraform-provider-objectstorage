package objectstorage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
)

type Account struct {
	Source string `json:"source,omitempty"`
	Group string `json:"group,omitempty"`
	AskId string `json:"askId,omitempty"`
	Description string `json:"description,omitempty"`
}

type AccountResponse struct {
	Id string `json:"id,omitempty"`
	Account
	NoOfBuckets int `json:"noOfBuckets,omitempty"`
	NoOfObjects int `json:"noOfObjects,omitempty"`
}

func resourceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountCreate,
		Update: resourceAccountUpdate,
		Read:   resourceAccountRead,
		Delete: resourceAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type: schema.TypeString,
				Optional: true,
			},
			"source": {
				Type: schema.TypeString,
				Optional: true,
				Default: "terraform",
			},
			"group": {
				Type: schema.TypeString,
				Optional: false,
			},
			"askId": {
				Type: schema.TypeString,
				Optional: false,
				Default: "terraform",
			},
			"description": {
				Type: schema.TypeString,
				Optional: true,
				Default: "description",
			},
			"noOfBuckets": {
				Type: schema.TypeInt,
				Optional: true,
				Default: 0,
			},
			"noOfObjects": {
				Type: schema.TypeInt,
				Optional: true,
				Default: 0,
			},
		},
	}
}

func newAccountFromResource(d *schema.ResourceData) *Account {
	return &Account{
		Source: d.Get("source").(string),
		Group: d.Get("group").(string),
		AskId: d.Get("askId").(string),
		Description: d.Get("description").(string),
	}
}

func resourceAccountCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	account := newAccountFromResource(d)

	bytedata, err := json.Marshal(account)
	if err != nil {
		return err
	}

	_, err = client.Post("accounts", bytes.NewBuffer(bytedata))
	if err != nil {
		return err
	}

	return resourceAccountRead(d, m)
}

func resourceAccountUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceAccountRead(d *schema.ResourceData, m interface{}) error {
	accountId := d.Get("account_id")
	if accountId == "" {
		return fmt.Errorf("ACCOUNT ID SHOULD NOT BE BLANK")
	}

	client := m.(*Client)
	resp, err := client.Get(fmt.Sprintf("accounts/%s", accountId))
	if err != nil {
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var accountResponse AccountResponse
	if err := json.Unmarshal(body, &accountResponse); err != nil {
		d.SetId("")
		return err
	}

	d.SetId(accountResponse.Id)
	d.Set("account_id", accountResponse.Id)
	d.Set("source", accountResponse.Source)
	d.Set("group", accountResponse.Group)
	d.Set("askId", accountResponse.AskId)
	d.Set("description", accountResponse.Description)
	d.Set("noOfBuckets", accountResponse.NoOfBuckets)
	d.Set("noOfObjects", accountResponse.NoOfObjects)

	return nil
}

func resourceAccountDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
