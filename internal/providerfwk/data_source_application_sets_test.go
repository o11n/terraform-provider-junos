package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceApplicationSets_basic(t *testing.T) {
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccDataSourceApplicationSetsPre(),
				},
				{
					Config: testAccDataSourceApplicationSetsConfig(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.junos_application_sets.testacc_ssh_without_telnet",
							"application_sets.#", "0"),
						resource.TestCheckResourceAttr("data.junos_application_sets.testacc_ssh_with_telnet",
							"application_sets.#", "1"),
						resource.TestCheckResourceAttr("data.junos_application_sets.testacc_name",
							"application_sets.#", "2"),
						resource.TestCheckResourceAttr("data.junos_application_sets.testacc_appsets",
							"application_sets.#", "1"),
						resource.TestCheckResourceAttr("data.junos_application_sets.testacc_appsets",
							"application_sets.0.description", "test-data-source-appSet"),
					),
				},
			},
		})
	}
}

func testAccDataSourceApplicationSetsPre() string {
	return `
resource "junos_application_set" "testacc_app_set" {
  name         = "testacc_app_set"
  applications = ["junos-ssh", "junos-telnet"]
}
resource "junos_application_set" "testacc_app_set2" {
  name            = "testacc_app_set2"
  application_set = [junos_application_set.testacc_app_set.name]
  description     = "test-data-source-appSet"
}
`
}

func testAccDataSourceApplicationSetsConfig() string {
	return `
resource "junos_application_set" "testacc_app_set" {
  name         = "testacc_app_set"
  applications = ["junos-ssh", "junos-telnet"]
}

data "junos_application_sets" "testacc_ssh_without_telnet" {
  match_applications = ["junos-ssh"]
}
data "junos_application_sets" "testacc_ssh_with_telnet" {
  match_applications = ["junos-telnet", "junos-ssh"]
}
data "junos_application_sets" "testacc_default_cifs" {
  match_applications = ["junos-netbios-session", "junos-smb-session"]
}
data "junos_application_sets" "testacc_name" {
  match_name = "testacc_.*"
}
data "junos_application_sets" "testacc_appsets" {
  match_application_sets = ["testacc_app_set"]
}
`
}
