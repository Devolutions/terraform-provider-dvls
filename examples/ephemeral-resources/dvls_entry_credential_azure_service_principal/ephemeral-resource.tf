# Lookup by ID — the client secret is fetched only during plan/apply
# and never stored in Terraform state.
ephemeral "dvls_entry_credential_azure_service_principal" "by_id" {
  vault_id = "00000000-0000-0000-0000-000000000000"
  id       = "00000000-0000-0000-0000-000000000000"
}

# Lookup by name
ephemeral "dvls_entry_credential_azure_service_principal" "by_name" {
  vault_id = "00000000-0000-0000-0000-000000000000"
  name     = "foo"
}

# Lookup by name in a specific folder
ephemeral "dvls_entry_credential_azure_service_principal" "by_name_in_folder" {
  vault_id = "00000000-0000-0000-0000-000000000000"
  name     = "foo"
  folder   = "foo\\bar"
}
