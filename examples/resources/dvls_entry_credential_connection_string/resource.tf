resource "dvls_entry_credential_connection_string" "example" {
  vault_id    = "00000000-0000-0000-0000-000000000000"
  name        = "foo"
  folder      = "foo\\bar"
  description = "bar"
  tags        = ["foo"]

  connection_string = "bar"
}
