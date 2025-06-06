resource "dvls_entry_credential_ssh_key" "example" {
  vault_id    = "00000000-0000-0000-0000-000000000000"
  name        = "foo"
  folder      = "foo\\bar"
  description = "bar"
  tags        = ["foo"]

  password         = "foo"
  passphrase       = "bar"
  private_key_data = "foo"
  public_key       = "bar"
}
