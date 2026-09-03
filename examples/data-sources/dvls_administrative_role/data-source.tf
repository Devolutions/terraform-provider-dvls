# Lookup by ID
data "dvls_administrative_role" "by_id" {
  id = "00000000-0000-0000-0000-00000000000f"
}

# Lookup by name
data "dvls_administrative_role" "by_name" {
  name = "Vault User"
}
