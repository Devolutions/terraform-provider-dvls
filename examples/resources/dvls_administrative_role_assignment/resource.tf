resource "dvls_vault" "example" {
  name = "example"
}

# Grant the built-in Vault User role on a vault
resource "dvls_administrative_role_assignment" "vault_user" {
  role_name   = "Vault User"
  assignee_id = "00000000-0000-0000-0000-000000000000"
  scope_type  = "vault"
  scope_id    = dvls_vault.example.id
}

# Global assignment by role id
resource "dvls_administrative_role_assignment" "vaults_administrator" {
  role_id     = "00000000-0000-0000-0000-000000000004"
  assignee_id = "00000000-0000-0000-0000-000000000000"
  scope_type  = "global"
}
