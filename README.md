# Terraform Provider for Devolutions Server (Terraform Plugin Framework)
:warning: **This provider is a work in progress, expect breaking changes between releases** :warning:

Terraform Provider for managing your Devolutions Server instance.

## Compatibility

| Provider version | DVLS version                                   |
|------------------|------------------------------------------------|
| 0.8.0+           | the first DVLS release with the public API v2 |
| 0.6.x - 0.7.x    | 2026.x                                         |
| 0.5.x            | 2025.x                                         |

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0 (>= 1.10 required to use ephemeral resources)
- [Go](https://golang.org/doc/install) >= 1.26
- [Devolutions Server](https://devolutions.net/server) with the public API v2 (see [Compatibility](#compatibility))

## Upgrading to 0.8.0

The provider uses [go-dvls](https://github.com/Devolutions/go-dvls) v1.0.0, which talks to the DVLS public API v2
only. Every entry is reached through its vault, so lookups that took an entry ID alone now take the vault ID too.

Configuration changes:

- `dvls_entry_permissions`: `vault_id` is required. `view_override` and `view_roles` are removed: set the `view`
  right in `permissions` like any other right (`{ right = "view", roles = [...] }`), except with the `default`
  override, which DVLS does not store: leave the right out instead. The `sealed` right is removed.
  The import ID is `<vault_id>/<entry_id>`. Existing state is upgraded: the view settings move into `permissions` (a `default` view is dropped),
  and the first apply after adding `vault_id` updates the resource in place. Until that apply the resource is not
  refreshed, and destroying it fails.
- `dvls_entry_certificate`, `dvls_entry_host` and `dvls_entry_website` data sources and the
  `dvls_entry_certificate` ephemeral resource: `vault_id` is required.
- `dvls_entry_certificate` resource: the import ID is `<vault_id>/<entry_id>`.
- `dvls_administrative_role_assignment`: the `gateway` and `pam_provider` scope types are removed; remove such
  assignments from the state with `terraform state rm` before upgrading. Creating an assignment that already exists
  adopts it, unless the principal holds the role at the scope through several assignments (a temporary grant next to
  the permanent one): the create then fails, and the one to manage can be imported by its assignment ID. A refresh
  reads the assignment by its ID.
- `dvls_entry_folder`: changing `parent_folder` moves the folder, with its content, in place (same ID) instead of
  replacing it. The new parent must already exist: declare it as its own `dvls_entry_folder` and reference it; the
  provider does not create it, and a move into a missing folder fails without changing anything.

Data source attribute changes:

| Data source                | Before                                                     | After                                    |
|----------------------------|------------------------------------------------------------|------------------------------------------|
| `dvls_administrative_role` | `supported_scopes`, `is_assignable`, `is_pam`, `is_privileged` | `assignable_scopes` (empty when the role cannot be assigned) |
| `dvls_user`                | `authentication_type`, `is_administrator`, `user_groups`      | `identity_provider`, `is_workspace_owner` |
| `dvls_user_group`          | `is_administrator`, `type` value `office_365`                 | `is_workspace_owner`, `is_enabled`, `type` values `entra_id` and `keycloak` |
| `dvls_application`         | `is_administrator`, `user_groups`                             | `is_workspace_owner` |
| `dvls_entry_website`       | `web_browser_application`                                   | removed |

Behaviour changes: vault and entry name lookups match ignoring case, and updating a nested `dvls_entry_folder` (a
rename or a description change) no longer moves it to the vault root. DVLS answers 404 for an entry in a vault the
provider's identity cannot open, as for a deleted entry: a refresh or destroy removes an entry from the state only when
the content of its vault is readable, and fails otherwise (if the vault was deleted, remove its entries with
`terraform state rm`). Removing `password` from a `dvls_entry_certificate` clears the stored password. Destroying a
`dvls_entry_folder` that still holds entries fails (DVLS answers 409): delete or move them first.

## Building The Provider

Until go-dvls v1.0.0 is tagged, `go.mod` pins go-dvls to a pseudo-version of a commit on its `feat/api-v2` branch.
Bump it to v1.0.0 once tagged: `go get github.com/Devolutions/go-dvls@v1.0.0 && go mod tidy`, then commit `go.mod` and
`go.sum`.

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Visit the Terraform Registry at https://registry.terraform.io/providers/Devolutions/dvls/latest for usage information.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## Testing the Provider while developing

In order to perform a test of the provider, create/update the `~/.terraformrc` (`%APPDATA%\terraform.rc` on Windows) file with the following content:

```hcl
provider_installation {
  dev_overrides {
    "devolutions/dvls" = "/Users/USERNAME/go/bin"
  }
  direct {
  }
}
```

Replace `/Users/USERNAME/go/bin` with the path to the compiled provider binary according to your operating system and environment.

Then run the following command, assuming you are on macOS or Linux:

```shell
go build -o ~/go/bin/terraform-provider-dvls_vX.Y.Z
chmod +x ~/go/bin/terraform-provider-dvls_vX.Y.Z
```
Notes: the version number should match the version of the provider you are working on.

You can then create a test.tf or example.tf file with the required content; here is a sample:

```hcl
provider "dvls" {
  base_uri   = "https://your-dvls-instance.com/"
  app_id     = "your-app-id"
  app_secret = "your-app-secret"
  # or, instead of app_id / app_secret:
  # api_key  = "your-api-key"
}

data "dvls_entry_website" "example" {
  vault_id = "id-of-vault"
  id       = "id-of-website-entry"
}

output "website_name" {
  value = data.dvls_entry_website.example.name
}

terraform {
  required_providers {
    dvls = {
      source = "devolutions/dvls"
    }
  }
}
```

Then run the following command:

```shell
terraform plan
```

This will be the output:

```shell
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - devolutions/dvls in /Users/USERNAME/go/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible
│ with published releases.
╵
data.dvls_entry_website.example: Reading...
data.dvls_entry_website.example: Read complete after 1s [id=123e4567-e89b-12d3-a456-426614174000]

Changes to Outputs:
  + website_name = "TestWebsite"
```

Please note that the `.gitignore` already ignores the `dev.tfrc`, `.terraform.lock.hcl`, `test.tf`, `example.tf`, and `terraform.tfstate` files and the folder `.terraform/`.
