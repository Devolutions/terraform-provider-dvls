# Certificate file content and password are fetched only during plan/apply
# and never stored in Terraform state. Lookup is by vault and entry ID.
ephemeral "dvls_entry_certificate" "by_id" {
  vault_id = "00000000-0000-0000-0000-000000000000"
  id       = "00000000-0000-0000-0000-000000000000"
}
