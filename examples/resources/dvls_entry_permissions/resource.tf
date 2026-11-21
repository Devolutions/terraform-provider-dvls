resource "dvls_entry_folder" "secrets" {
  vault_id = "00000000-0000-0000-0000-000000000000"
  name     = "secrets"
}

# Make the folder visible to everyone
resource "dvls_entry_permissions" "everyone" {
  vault_id      = dvls_entry_folder.secrets.vault_id
  entry_id      = dvls_entry_folder.secrets.id
  role_override = "everyone"
}

# Custom permissions: only the listed principals can view, one of them can edit
resource "dvls_entry_permissions" "custom" {
  vault_id      = "00000000-0000-0000-0000-000000000000"
  entry_id      = "00000000-0000-0000-0000-000000000000"
  role_override = "custom"

  permissions = [
    {
      right = "view"
      roles = ["00000000-0000-0000-0000-000000000000"]
    },
    {
      right = "edit"
      roles = ["00000000-0000-0000-0000-000000000000"]
    },
    {
      right    = "view_password"
      override = "everyone"
    },
  ]
}
