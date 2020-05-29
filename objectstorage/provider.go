package objectstorage

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"net/http"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Required: true,
				Type: schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_USERNAME", nil),
			},
			"password": {
				Required: true,
				Sensitive: true,
				Type: schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_PASSWORD", nil),
			},
		},
		ConfigureFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"objectstorage_account": resourceAccount(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := &Client{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		HTTPClient: &http.Client{},
	}

	return client, nil
}
