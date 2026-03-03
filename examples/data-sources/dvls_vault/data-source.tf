# Lookup by ID
data "dvls_vault" "by_id" {
  id = "00000000-0000-0000-0000-000000000000"
}

# Lookup by name
data "dvls_vault" "by_name" {
  name = "foo"
}
