package provider

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"dvls": providerserver.NewProtocol6WithError(New("test")()),
}

// echoprovider surfaces ephemeral values into a managed resource so they can be asserted.
var testAccProtoV6ProviderFactoriesWithEcho = map[string]func() (tfprotov6.ProviderServer, error){
	"dvls": providerserver.NewProtocol6WithError(New("test")()),
	"echo": echoprovider.NewProviderServer(),
}

var testAccEphemeralTerraformVersionCheck = []tfversion.TerraformVersionCheck{
	tfversion.SkipBelow(tfversion.Version1_10_0),
}

// testAccEphemeralEchoConfig wires the echo provider/resource around a
// reference expression (e.g. "ephemeral.dvls_entry_credential_secret.test")
// so its attributes can be asserted via "echo.test.data.<field>".
func testAccEphemeralEchoConfig(refExpr string) string {
	return fmt.Sprintf(`
provider "echo" {
  data = %s
}

resource "echo" "test" {}
`, refExpr)
}

var (
	testAccClient     *dvls.Client
	testAccClientOnce sync.Once
	testAccClientErr  error
)

func getTestAccClient() (*dvls.Client, error) {
	testAccClientOnce.Do(func() {
		client, err := dvls.NewClient(
			os.Getenv("TEST_DVLS_APP_ID"),
			os.Getenv("TEST_DVLS_APP_SECRET"),
			os.Getenv("TEST_DVLS_BASE_URI"),
		)
		if err != nil {
			testAccClientErr = fmt.Errorf("unable to create test client: %s", err)
			return
		}

		testAccClient = &client
	})

	return testAccClient, testAccClientErr
}

func testAccPreCheck(t *testing.T) {
	t.Helper()

	envVars := []string{"TEST_DVLS_BASE_URI", "TEST_DVLS_APP_ID", "TEST_DVLS_APP_SECRET"}

	for _, env := range envVars {
		if os.Getenv(env) == "" {
			t.Fatalf("%s must be set for acceptance tests", env)
		}
	}

	t.Setenv("DVLS_APP_ID", os.Getenv("TEST_DVLS_APP_ID"))
	t.Setenv("DVLS_APP_SECRET", os.Getenv("TEST_DVLS_APP_SECRET"))
}

func testAccProviderConfig() string {
	return fmt.Sprintf(`
provider "dvls" {
  base_uri = %q
}
`, os.Getenv("TEST_DVLS_BASE_URI"))
}

func testAccEntryCredentialImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["vault_id"], rs.Primary.ID), nil
	}
}

func testAccCheckVaultDestroy(s *terraform.State) error {
	client, err := getTestAccClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "dvls_vault" {
			continue
		}

		_, err := client.Vaults.Get(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("vault %s still exists", rs.Primary.ID)
		}

		if !dvls.IsNotFound(err) {
			return fmt.Errorf("unexpected error checking vault %s: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

var credentialResourceTypes = map[string]bool{
	"dvls_entry_credential_username_password":       true,
	"dvls_entry_credential_api_key":                 true,
	"dvls_entry_credential_secret":                  true,
	"dvls_entry_credential_ssh_key":                 true,
	"dvls_entry_credential_azure_service_principal": true,
	"dvls_entry_credential_connection_string":       true,
}

func testAccCheckEntryCredentialDestroy(s *terraform.State) error {
	client, err := getTestAccClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if !credentialResourceTypes[rs.Type] {
			continue
		}

		vaultId := rs.Primary.Attributes["vault_id"]
		entryId := rs.Primary.ID

		_, err := client.Entries.Credential.GetById(vaultId, entryId)
		if err == nil {
			return fmt.Errorf("entry %s/%s still exists", vaultId, entryId)
		}

		if !dvls.IsNotFound(err) {
			return fmt.Errorf("unexpected error checking entry %s/%s: %s", vaultId, entryId, err)
		}
	}

	return nil
}
