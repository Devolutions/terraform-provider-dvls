resource "dvls_entry_credential_secret" "example" {
  vault_id = "00000000-0000-0000-0000-000000000000"
  name     = "foo"
  folder   = "foo\\bar"

  secret = "bar"

  description = "bar"
  tags        = ["foo"]
}
