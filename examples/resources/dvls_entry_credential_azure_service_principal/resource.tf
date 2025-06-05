resource "dvls_entry_credential_azure_service_principal" "example" {
  vault_id    = "00000000-0000-0000-0000-000000000000"
  name        = "foo"
  folder      = "foo\\bar"
  description = "foo"
  tags        = ["foo"]

  client_id     = "foo"
  client_secret = "bar"
  tenant_id     = "foo"
}
