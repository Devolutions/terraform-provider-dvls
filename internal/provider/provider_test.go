package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
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

// Shared literals used by every ephemeral acceptance test. Centralized so
// the HCL config and the TestCheckResourceAttr assertions reference one
// source of truth.
const (
	testAccEphFolder      = "tf_test_folder"
	testAccEphDescription = "test entry for ephemeral resource"
)

var testAccEphTags = []string{"acceptance", "tf-test"}

// ephemeralLookupBlock builds an ephemeral { ... } HCL block that looks up
// `resourceType.test` by either "name" or "id" (whichever lookupField is).
// All 6 credential ephemeral tests share this exact shape — only the
// resource type and lookup field differ.
func ephemeralLookupBlock(resourceType, lookupField string) string {
	return fmt.Sprintf(`
ephemeral %[1]q "test" {
  vault_id = dvls_vault.test.id
  %[2]s = %[1]s.test.%[2]s
}
`, resourceType, lookupField)
}

// testAccEntryCredentialEphemeralConfig builds the HCL for a credential
// ephemeral acceptance step: provider + vault + credential resource + (when
// lookupField is non-empty) ephemeral lookup block + echo provider/resource.
//
// `fields` is the subtype-specific HCL fragment (e.g. `  secret = "..."`).
// `lookupField` is "" (creation step, no ephemeral), "name", or "id".
func testAccEntryCredentialEphemeralConfig(resourceType, vaultName, entryName, fields, lookupField string) string {
	ephemeralBlock := ""
	echoConfig := ""
	if lookupField != "" {
		ephemeralBlock = ephemeralLookupBlock(resourceType, lookupField)
		echoConfig = testAccEphemeralEchoConfig(fmt.Sprintf("ephemeral.%s.test", resourceType))
	}

	return fmt.Sprintf(`
%[1]s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource %[3]q "test" {
  vault_id = dvls_vault.test.id
  name = %[4]q
  description = %[5]q
  folder = %[6]q
  tags = [%[7]q, %[8]q]
%[9]s
}

%[10]s

%[11]s
`, testAccProviderConfig(), vaultName, resourceType, entryName,
		testAccEphDescription, testAccEphFolder, testAccEphTags[0], testAccEphTags[1],
		fields, ephemeralBlock, echoConfig)
}

// getTestAccClient returns a freshly authenticated DVLS client. It does not
// cache the client because tokens are short-lived: a long test run can outlive
// the session, and a stale client surfaces as 401s in CheckDestroy callbacks.
func getTestAccClient() (*dvls.Client, error) {
	client, err := dvls.NewClient(
		os.Getenv("TEST_DVLS_APP_ID"),
		os.Getenv("TEST_DVLS_APP_SECRET"),
		os.Getenv("TEST_DVLS_BASE_URI"),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create test client: %s", err)
	}
	return &client, nil
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

// testAccCreateFolderInVault posts a Folder entry to a vault. DVLS rejects
// credential creation with status 400 when its `path` references a folder
// that does not exist; tests that exercise the folder attribute must
// pre-create the folder via this helper.
func testAccCreateFolderInVault(vaultId, folderName string) error {
	client, err := getTestAccClient()
	if err != nil {
		return err
	}

	body, err := json.Marshal(map[string]any{
		"name":    folderName,
		"type":    "Folder",
		"subType": "Folder",
		"data":    map[string]any{},
	})
	if err != nil {
		return fmt.Errorf("failed to marshal folder body: %w", err)
	}

	reqURL, err := url.JoinPath(os.Getenv("TEST_DVLS_BASE_URI"), "/api/v1/vault/", vaultId, "/entry")
	if err != nil {
		return fmt.Errorf("failed to build folder url: %w", err)
	}

	if _, err := client.Request(reqURL, http.MethodPost, bytes.NewBuffer(body)); err != nil {
		return fmt.Errorf("failed to create folder %q in vault %s: %w", folderName, vaultId, err)
	}
	return nil
}

// testAccVaultWithFoldersStep returns a TestStep that creates dvls_vault.test
// and pre-creates the given folders inside it via the DVLS API. Prepend it to
// a test's Steps when subsequent steps set `folder` on a credential entry.
func testAccVaultWithFoldersStep(vaultName string, folderNames ...string) resource.TestStep {
	return resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "dvls_vault" "test" {
  name = %q
}
`, testAccProviderConfig(), vaultName),
		Check: func(s *terraform.State) error {
			rs, ok := s.RootModule().Resources["dvls_vault.test"]
			if !ok {
				return fmt.Errorf("dvls_vault.test not found in state")
			}
			for _, name := range folderNames {
				if err := testAccCreateFolderInVault(rs.Primary.ID, name); err != nil {
					return err
				}
			}
			return nil
		},
	}
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
