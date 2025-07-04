---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "dvls_entry_credential_secret Data Source - terraform-provider-dvls"
subcategory: ""
description: |-
  A DVLS Secret Credential Entry
---

# dvls_entry_credential_secret (Data Source)

A DVLS Secret Credential Entry

## Example Usage

```terraform
data "dvls_entry_credential_secret" "example" {
  id       = "00000000-0000-0000-0000-000000000000"
  vault_id = "00000000-0000-0000-0000-000000000000"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The ID of the entry.
- `vault_id` (String) The ID of the vault.

### Read-Only

- `description` (String) The description of the entry.
- `folder` (String) The folder path of the entry.
- `name` (String) The name of the entry.
- `secret` (String, Sensitive) The entry credential secret.
- `tags` (List of String) A list of tags added to the entry.
