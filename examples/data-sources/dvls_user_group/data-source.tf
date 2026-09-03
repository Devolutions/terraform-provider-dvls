# Lookup by ID
data "dvls_user_group" "by_id" {
  id = "00000000-0000-0000-0000-000000000000"
}

# Lookup by name
data "dvls_user_group" "by_name" {
  name = "example"
}
