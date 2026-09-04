# Lookup by ID
data "dvls_application" "by_id" {
  id = "00000000-0000-0000-0000-000000000000"
}

# Lookup by application key
data "dvls_application" "by_key" {
  name = "00000000-0000-0000-0000-000000000000"
}

# Lookup by display name
data "dvls_application" "by_full_name" {
  full_name = "Example Application"
}
