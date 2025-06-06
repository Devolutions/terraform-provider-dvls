resource "dvls_entry_credential_username_password" "example" {
  vault_id    = "00000000-0000-0000-0000-000000000000"
  name        = "foo"
  folder      = "foo\\bar"
  description = "bar"
  tags        = ["foo"]

  username = "foo"
  domain   = "foo.bar"
  password = "bar"
}
