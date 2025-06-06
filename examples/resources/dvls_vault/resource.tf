resource "dvls_vault" "example" {
  name            = "foo"
  description     = "bar"
  visibility      = "private"
  security_level  = "high"
  master_password = "foo!"
}
