output "account_id" {
  value = objectstorage_account.terraform-account.id
}

output "account_id2" {
  value = objectstorage_account.terraform-account2.id
}

output "no_of_buckets" {
  value = objectstorage_account.terraform-account.no_of_buckets
}

output "no_of_objects" {
  value = objectstorage_account.terraform-account.no_of_objects
}