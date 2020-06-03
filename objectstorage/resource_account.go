package objectstorage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stoewer/go-strcase"
	"io/ioutil"
	"log"
	"os"
)

type Account struct {
	Source      string `json:"source,omitempty"`
	Group       string `json:"group,omitempty"`
	AskId       string `json:"askId,omitempty"`
	Description string `json:"description,omitempty"`
}

type AccountResponse struct {
	Id string `json:"id,omitempty"`
	Account
	NoOfBuckets int `json:"noOfBuckets,omitempty"`
	NoOfObjects int `json:"noOfObjects,omitempty"`
}

type UpdateElement struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "terraform",
			},
			"group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ask_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"no_of_buckets": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"no_of_objects": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func newAccountFromResource(d *schema.ResourceData) *Account {
	return &Account{
		Source:      d.Get("source").(string),
		Group:       d.Get("group").(string),
		AskId:       d.Get("ask_id").(string),
		Description: d.Get("description").(string),
	}
}

func newUpdateElementFromResource(d *schema.ResourceData, element string) *UpdateElement {
	return &UpdateElement{
		Name:  strcase.LowerCamelCase(element),
		Value: d.Get(element).(string),
	}
}

func resourceAccountCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] MS ID: %v, MS Password: %v", os.Getenv("MS_ID"), os.Getenv("MS_PASSWORD"))
	client := m.(*Client)
	account := newAccountFromResource(d)

	bytedata, err := json.Marshal(account)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating a new account resource")
	resp, err := client.Post("accounts", bytes.NewBuffer(bytedata))
	log.Printf("[DEBUG] Creating a new account resource: %s", err)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Account resource got created successfully")
	if resp.StatusCode == 201 {
		var accountResponse AccountResponse

		body, readerr := ioutil.ReadAll(resp.Body)
		log.Printf("[DEBUG] Create account Resp: %s, Err: %s", body, err)
		if readerr != nil {
			return readerr
		}

		log.Printf("[DEBUG] Unmarshalling the create response: %s", body)
		decodeerr := json.Unmarshal(body, &accountResponse)
		if decodeerr != nil {
			return decodeerr
		}

		log.Printf("[DEBUG] Account resource id: %s", accountResponse.Id)
		d.SetId(accountResponse.Id)
		d.Set("account_id", accountResponse.Id)
	}

	return resourceAccountRead(d, m)
}

func resourceAccountUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(true)

	elements := []string{"source", "group", "ask_id", "description"}
	for _, element := range elements {
		UpdateValue(d, m, element)
	}

	d.Partial(false)

	return nil
}

func UpdateValue(d *schema.ResourceData, m interface{}, element string) error {
	if d.HasChange(element) {

		if err := Update(d, m, element); err != nil {
			return err
		}

		d.SetPartial(element)
	}

	return nil
}

func Update(d *schema.ResourceData, m interface{}, element string) error {
	client := m.(*Client)
	updateElement := newUpdateElementFromResource(d, element)

	bytedata, err := json.Marshal(updateElement)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating %s field with value: %s", element, d.Get(element))
	resp, err := client.Put(fmt.Sprintf("accounts/%s", d.Id()), bytes.NewBuffer(bytedata))
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		log.Printf("[INFO] Account resource got updated successfully")
	}

	return nil
}

func resourceAccountRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("ID SHOULD NOT BE BLANK")
	}

	client := m.(*Client)
	resp, err := client.Get(fmt.Sprintf("accounts/%s", id))
	if err != nil {
		d.SetId("")
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

	log.Printf("[DEBUG] No Of Buckets: %v, No Of Objects: %v", accountResponse.NoOfBuckets, accountResponse.NoOfObjects)
	d.Set("account_id", accountResponse.Id)
	d.Set("source", accountResponse.Source)
	d.Set("group", accountResponse.Group)
	d.Set("ask_id", accountResponse.AskId)
	d.Set("description", accountResponse.Description)
	d.Set("no_of_buckets", accountResponse.NoOfBuckets)
	d.Set("no_of_objects", accountResponse.NoOfObjects)

	return nil
}

func resourceAccountDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("ID SHOULD NOT BE BLANK")
	}

	client := m.(*Client)
	resp, err := client.Delete(fmt.Sprintf("accounts/%s", id))
	if err != nil {
		return nil
	}

	if resp.StatusCode == 200 {
		log.Printf("[INFO] Account resource got deleted successfully")
	}

	d.SetId("")
	return nil
}
