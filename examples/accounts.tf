resource "objectstorage_account" "terraform-account" {
  name = "terraform-account"
  source = "terraform"
  group = "Test_Functionality_Non_Privileged"
  askId = "ASK_ID"
  description = "terraform storage"
}