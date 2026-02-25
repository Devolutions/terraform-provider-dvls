# Lookup by ID
data "dvls_entry_credential_ssh_key" "example" {
  vault_id = "00000000-0000-0000-0000-000000000000"
  id       = "00000000-0000-0000-0000-000000000000"
}

# Lookup by name
data "dvls_entry_credential_ssh_key" "example" {
  vault_id = "00000000-0000-0000-0000-000000000000"
  name     = "foo"
}

# Lookup by name in a specific folder
data "dvls_entry_credential_ssh_key" "example" {
  vault_id = "00000000-0000-0000-0000-000000000000"
  name     = "foo"
  folder   = "foo\\bar"
}
