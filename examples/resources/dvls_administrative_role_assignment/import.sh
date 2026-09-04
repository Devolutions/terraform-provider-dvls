# This resource can be imported using `<role_id>/<scope_type>/<scope_id>/<assignee_id>` format, e.g.
terraform import dvls_administrative_role_assignment.example 00000000-0000-0000-0000-00000000000f/vault/00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000

# Leave <scope_id> empty for a global assignment
terraform import dvls_administrative_role_assignment.example 00000000-0000-0000-0000-000000000004/global//00000000-0000-0000-0000-000000000000
