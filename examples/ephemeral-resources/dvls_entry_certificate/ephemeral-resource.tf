# Certificate file content and password are fetched only during plan/apply
# and never stored in Terraform state. Lookup is by ID only.
ephemeral "dvls_entry_certificate" "by_id" {
  id = "00000000-0000-0000-0000-000000000000"
}
