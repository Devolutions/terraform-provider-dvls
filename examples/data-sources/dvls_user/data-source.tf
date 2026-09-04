# Lookup by ID
data "dvls_user" "by_id" {
  id = "00000000-0000-0000-0000-000000000000"
}

# Lookup by login name
data "dvls_user" "by_name" {
  name = "jdoe"
}

# Lookup by display name
data "dvls_user" "by_full_name" {
  full_name = "John Doe"
}
