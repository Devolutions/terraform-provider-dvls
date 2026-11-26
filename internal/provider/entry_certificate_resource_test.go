package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEntryCertificateResource_url(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntryCertificateResourceConfig("tf_test_certificate", "2030-12-31T23:59:59-05:00", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("dvls_entry_certificate.test", "id"),
					resource.TestCheckResourceAttr("dvls_entry_certificate.test", "expiration", "2030-12-31T23:59:59-05:00"),
					resource.TestCheckResourceAttr("dvls_entry_certificate.test", "url.url", "https://example.com/cert.pfx"),
					resource.TestCheckResourceAttr("dvls_entry_certificate.test", "url.use_default_credentials", "true"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_certificate.test", "name", "dvls_entry_certificate.test", "name"),
					resource.TestCheckResourceAttrPair("data.dvls_entry_certificate.test", "url.url", "dvls_entry_certificate.test", "url.url"),
					resource.TestCheckResourceAttr("data.dvls_entry_certificate.test", "url.use_default_credentials", "true"),
					resource.TestCheckResourceAttr("data.dvls_entry_certificate.test", "expiration", "2031-01-01T04:59:59Z"),
				),
			},
			{
				Config:   testAccEntryCertificateResourceConfig("tf_test_certificate", "2030-12-31T23:59:59-05:00", true),
				PlanOnly: true,
			},
			{
				Config: testAccEntryCertificateResourceConfig("tf_test_certificate", "2031-06-30T12:00:00Z", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dvls_entry_certificate.test", "expiration", "2031-06-30T12:00:00Z"),
					resource.TestCheckResourceAttr("data.dvls_entry_certificate.test", "expiration", "2031-06-30T12:00:00Z"),
					resource.TestCheckResourceAttr("data.dvls_entry_certificate.test", "url.use_default_credentials", "false"),
				),
			},
			{
				ResourceName:      "dvls_entry_certificate.test",
				ImportState:       true,
				ImportStateIdFunc: testAccEntryImportStateIdFunc("dvls_entry_certificate.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEntryCertificateResource_file(t *testing.T) {
	config := fmt.Sprintf(`
%[1]s

resource "dvls_vault" "test" {
  name = "tf_test_certificate_file"
}

resource "dvls_entry_certificate" "test" {
  vault_id   = dvls_vault.test.id
  name       = "tf_test_certificate_file"
  expiration = "2030-12-31T23:59:59Z"
  password   = "certificate-password"

  file = {
    name        = "test.p12"
    content_b64 = %[2]q
  }
}

data "dvls_entry_certificate" "test" {
  vault_id = dvls_vault.test.id
  id       = dvls_entry_certificate.test.id
}
`, testAccProviderConfig(), testAccCertificateContentB64)

	configWithoutPassword := strings.Replace(config, "  password   = \"certificate-password\"\n", "", 1)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVaultDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dvls_entry_certificate.test", "file.name", "test.p12"),
					resource.TestCheckResourceAttr("dvls_entry_certificate.test", "file.content_b64", testAccCertificateContentB64),
					resource.TestCheckResourceAttr("data.dvls_entry_certificate.test", "file.content_b64", testAccCertificateContentB64),
					resource.TestCheckResourceAttr("data.dvls_entry_certificate.test", "password", "certificate-password"),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
			{
				ResourceName:      "dvls_entry_certificate.test",
				ImportState:       true,
				ImportStateIdFunc: testAccEntryImportStateIdFunc("dvls_entry_certificate.test"),
				ImportStateVerify: true,
			},
			{
				Config: configWithoutPassword,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("dvls_entry_certificate.test", "password"),
					resource.TestCheckNoResourceAttr("data.dvls_entry_certificate.test", "password"),
					resource.TestCheckResourceAttr("data.dvls_entry_certificate.test", "file.content_b64", testAccCertificateContentB64),
				),
			},
			{
				Config:   configWithoutPassword,
				PlanOnly: true,
			},
		},
	})
}

const testAccCertificateContentB64 = "dGVycmFmb3JtLXByb3ZpZGVyLWR2bHMgdGVzdCBjZXJ0aWZpY2F0ZQ=="

func testAccEntryCertificateResourceConfig(vaultName, expiration string, useDefaultCredentials bool) string {
	return fmt.Sprintf(`
%[1]s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_certificate" "test" {
  vault_id   = dvls_vault.test.id
  name       = "tf_test_certificate"
  expiration = %[3]q
  password   = "certificate-password"

  url = {
    url                     = "https://example.com/cert.pfx"
    use_default_credentials = %[4]t
  }
}

data "dvls_entry_certificate" "test" {
  vault_id = dvls_vault.test.id
  id       = dvls_entry_certificate.test.id

  depends_on = [dvls_entry_certificate.test]
}
`, testAccProviderConfig(), vaultName, expiration, useDefaultCredentials)
}
