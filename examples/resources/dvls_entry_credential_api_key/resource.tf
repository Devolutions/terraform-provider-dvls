resource "dvls_entry_credential_api_key" "example" {
  vault_id    = "00000000-0000-0000-0000-000000000000"
  name        = "foo"
  folder      = "foo\\bar"
  description = "foo"
  tags        = ["foo"]

  api_id    = "foo"
  api_key   = "bar"
  tenant_id = "foo"
}
