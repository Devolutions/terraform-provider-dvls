resource "dvls_entry_credential_api_key" "example" {
  vault_id = "00000000-0000-0000-0000-000000000000"
  name     = "foo"
  folder   = "foo\\bar"

  api_id    = "foo"
  api_key   = "bar"
  tenant_id = "foo"

  description = "bar"
  tags        = ["foo"]
}
