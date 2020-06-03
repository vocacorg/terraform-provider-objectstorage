resource "objectstorage_account" "terraform-account" {
  source = "terraform"
  group = "Test_Functionality_Non_Privileged"
  ask_id = "NEW_ASK_ID_AN"
  description = "storage account for testing purposes - terraform - learning"
}

resource "objectstorage_account" "terraform-account2" {
  source = "terraform"
  group = "Test_Functionality_Non_Privileged"
  ask_id = "ASK_ID"
  description = "storage account for testing purposes - terraform"
}