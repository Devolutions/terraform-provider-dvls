resource "dvls_entry_folder" "parent" {
  vault_id    = "00000000-0000-0000-0000-000000000000"
  name        = "foo"
  description = "parent folder at the vault root"
}

resource "dvls_entry_folder" "child" {
  vault_id      = "00000000-0000-0000-0000-000000000000"
  name          = "bar"
  parent_folder = dvls_entry_folder.parent.name
}
