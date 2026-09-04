package provider

import (
	"fmt"
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

// Shared literals used by every entry acceptance test (resource, data
// source, and ephemeral). Centralized so the HCL config and the
// TestCheckResourceAttr assertions reference one source of truth.
const (
	testAccTestFolder      = "tf_test_folder"
	testAccTestDescription = "test entry for ephemeral resource"
)

var testAccTestTags = []string{"acceptance", "tf-test"}

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

// testAccEntryCredentialDataSourceCheck asserts that the data source mirrors
// the underlying credential resource on common attributes plus the given
// subtype-specific fields.
func testAccEntryCredentialDataSourceCheck(resourceType string, fields ...string) resource.TestCheckFunc {
	dataAddr := "data." + resourceType + ".test"
	resAddr := resourceType + ".test"
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(dataAddr, "id", resAddr, "id"),
		resource.TestCheckResourceAttrPair(dataAddr, "vault_id", resAddr, "vault_id"),
		resource.TestCheckResourceAttrPair(dataAddr, "name", resAddr, "name"),
		resource.TestCheckResourceAttrPair(dataAddr, "description", resAddr, "description"),
		resource.TestCheckResourceAttrPair(dataAddr, "folder", resAddr, "folder"),
		resource.TestCheckResourceAttrPair(dataAddr, "tags.#", resAddr, "tags.#"),
	}
	for _, f := range fields {
		checks = append(checks, resource.TestCheckResourceAttrPair(dataAddr, f, resAddr, f))
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

// testAccEntryCredentialResourceConfig builds the HCL for a credential
// resource acceptance step. Both folders (`tf_test_folder` and
// `tf_test_folder_updated`) are declared so the Update step can switch
// between them. `fields` is the subtype-specific HCL fragment.
func testAccEntryCredentialResourceConfig(resourceType, vaultName, entryName, description, folder, fields string) string {
	return fmt.Sprintf(`
%[1]s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_folder" "default" {
  vault_id = dvls_vault.test.id
  name     = "tf_test_folder"
}

resource "dvls_entry_folder" "updated" {
  vault_id = dvls_vault.test.id
  name     = "tf_test_folder_updated"
}

resource %[3]q "test" {
  vault_id = dvls_vault.test.id
  name = %[4]q
  description = %[5]q
  folder = %[6]q
  tags = [%[7]q, %[8]q]
%[9]s

  depends_on = [dvls_entry_folder.default, dvls_entry_folder.updated]
}
`, testAccProviderConfig(), vaultName, resourceType, entryName,
		description, folder, testAccTestTags[0], testAccTestTags[1], fields)
}

// testAccEntryCredentialDataSourceConfig builds the HCL for a credential
// data source acceptance step: provider + vault + folder + credential
// resource + data source. `lookupField` is "name" or "id".
func testAccEntryCredentialDataSourceConfig(resourceType, vaultName, entryName, fields, lookupField string) string {
	return fmt.Sprintf(`
%[1]s

resource "dvls_vault" "test" {
  name = %[2]q
}

resource "dvls_entry_folder" "default" {
  vault_id = dvls_vault.test.id
  name     = %[6]q
}

resource %[3]q "test" {
  vault_id = dvls_vault.test.id
  name = %[4]q
  description = %[5]q
  folder = %[6]q
  tags = [%[7]q, %[8]q]
%[9]s

  depends_on = [dvls_entry_folder.default]
}

data %[3]q "test" {
  vault_id = dvls_vault.test.id
  %[10]s = %[3]s.test.%[10]s
}
`, testAccProviderConfig(), vaultName, resourceType, entryName,
		testAccTestDescription, testAccTestFolder, testAccTestTags[0], testAccTestTags[1],
		fields, lookupField)
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

resource "dvls_entry_folder" "default" {
  vault_id = dvls_vault.test.id
  name     = %[6]q
}

resource %[3]q "test" {
  vault_id = dvls_vault.test.id
  name = %[4]q
  description = %[5]q
  folder = %[6]q
  tags = [%[7]q, %[8]q]
%[9]s

  depends_on = [dvls_entry_folder.default]
}

%[10]s

%[11]s
`, testAccProviderConfig(), vaultName, resourceType, entryName,
		testAccTestDescription, testAccTestFolder, testAccTestTags[0], testAccTestTags[1],
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

// testAccPreCheckRoleAssignment runs before resource.Test because the test
// config needs a live principal id. It mirrors the TF_ACC gate of
// resource.Test, then skips when the server does not expose the
// administrative role endpoints (DVLS 2026.3+).
func testAccPreCheckRoleAssignment(t *testing.T) {
	t.Helper()

	client := testAccLiveClient(t)

	_, err := client.AdministrativeRoles.List()
	if dvls.IsNotFound(err) {
		t.Skip("administrative roles require DVLS 2026.3 or later")
	}
	if err != nil {
		t.Fatalf("unable to list administrative roles: %s", err)
	}
}

// testAccLiveClient mirrors the TF_ACC gate of resource.Test for tests that
// must query the server before building their config, and returns a client.
func testAccLiveClient(t *testing.T) *dvls.Client {
	t.Helper()

	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skipf("Acceptance tests skipped unless env '%s' set", resource.EnvTfAcc)
	}

	testAccPreCheck(t)

	client, err := getTestAccClient()
	if err != nil {
		t.Fatal(err)
	}

	return client
}

// testAccFindAssigneeId returns a principal id usable in permission tests:
// TEST_DVLS_ASSIGNEE_ID when set, otherwise the first member of the built-in
// Administrator role (the test application itself is one of them).
func testAccFindAssigneeId(t *testing.T) string {
	t.Helper()

	if assigneeId := os.Getenv("TEST_DVLS_ASSIGNEE_ID"); assigneeId != "" {
		return assigneeId
	}

	client := testAccLiveClient(t)

	members, err := client.AdministrativeRoleAssignments.GetMembers(dvls.BuiltinRoleBuiltinAdministratorId, dvls.AdministrativeRoleScopeGlobal, "")
	if err != nil {
		t.Fatalf("unable to list administrator members: %s", err)
	}
	if len(members) == 0 {
		t.Fatal("no administrator member found, set TEST_DVLS_ASSIGNEE_ID")
	}

	return members[0].AssigneeId
}

func testAccProviderConfig() string {
	return fmt.Sprintf(`
provider "dvls" {
  base_uri = %q
}
`, os.Getenv("TEST_DVLS_BASE_URI"))
}

func testAccEntryImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
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

func testAccCheckAdministrativeRoleAssignmentDestroy(s *terraform.State) error {
	client, err := getTestAccClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "dvls_administrative_role_assignment" {
			continue
		}

		scopeType, err := lookupMapValue(administrativeRoleScopeTypes, rs.Primary.Attributes["scope_type"])
		if err != nil {
			return err
		}

		members, err := client.AdministrativeRoleAssignments.GetMembers(rs.Primary.Attributes["role_id"], scopeType, rs.Primary.Attributes["scope_id"])
		if err != nil {
			return fmt.Errorf("unexpected error checking role assignment %s: %s", rs.Primary.ID, err)
		}

		for _, member := range members {
			if member.Id == rs.Primary.ID {
				return fmt.Errorf("role assignment %s still exists", rs.Primary.ID)
			}
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

// testAccFirstOf lists a server-side collection and returns its first item,
// skipping the test when the collection is empty.
func testAccFirstOf[T any](t *testing.T, kind string, list func() ([]T, error)) T {
	t.Helper()

	items, err := list()
	if err != nil {
		t.Fatalf("unable to list %ss: %s", kind, err)
	}
	if len(items) == 0 {
		t.Skipf("no %s available on the server", kind)
	}

	return items[0]
}

func testAccDataSourceBlock(dataSourceType, name, lookup string) string {
	return fmt.Sprintf(`
data %q %q {
  %s
}
`, dataSourceType, name, lookup)
}
