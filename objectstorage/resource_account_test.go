package objectstorage

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccStorageAccount_basic(t *testing.T) {
	var account Account

	testStorageAccountConfig := `
	resource "objectstorage_account" "test_account" {
	  source = "terraform"
	  group = "Test_Functionality_Non_Privileged"
	  ask_id = "NEW_ASK_ID_AN"
	  description = "storage account for testing purposes"
	}
	`
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testStorageAccountConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testStorageAccountExists("objectstorage_account.test_account", &account),
				),
			},
		},
	})
}

func testStorageAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	rs, ok := s.RootModule().Resources["objectstorage_account.test_account"]
	if !ok {
		return fmt.Errorf("not found %s", "objectstorage_account.test_account")
	}

	response, _ := client.Get(fmt.Sprintf("accounts/%s", rs.Primary.Attributes["account_id"]))
	if response.StatusCode != 404 {
		return fmt.Errorf("Account still exists")
	}

	return nil
}

func testStorageAccountExists(n string, account *Account) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Account ID is set")
		}
		return nil
	}
}

/**
func TestCamelCase(t *testing.T) {
	fmt.Printf("Camel Case: %s", strcase.LowerCamelCase("ask_id"))
}
*/
