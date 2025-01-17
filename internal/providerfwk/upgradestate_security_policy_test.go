package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccJunosSecurityPolicyUpgradeStateV0toV1_basic(t *testing.T) {
	if os.Getenv("TESTACC_UPGRADE_STATE") == "" {
		return
	}
	if os.Getenv("TESTACC_SRX") != "" {
		resource.Test(t, resource.TestCase{
			Steps: []resource.TestStep{
				{
					ExternalProviders: map[string]resource.ExternalProvider{
						"junos": {
							VersionConstraint: "1.33.0",
							Source:            "jeremmfr/junos",
						},
					},
					Config: testAccJunosSecurityPolicyConfigV0(),
				},
				{
					ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
					Config:                   testAccJunosSecurityPolicyConfigV0(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectEmptyPlan(),
						},
					},
				},
			},
		})
	}
}

func testAccJunosSecurityPolicyConfigV0() string {
	return `
resource "junos_services_advanced_anti_malware_policy" "testacc_v0to1_secuPolicy" {
  name                     = "testacc_v0to1_secuPolicy"
  verdict_threshold        = "recommended"
  default_notification_log = true
}
resource "junos_security_idp_policy" "testacc_v0to1_secuPolicy" {
  name = "testacc_v0to1_secuPolicy"
}
resource "junos_security_policy" "testacc_v0to1_secuPolicy" {
  from_zone = junos_security_zone.testacc_v0to1_secuPolicy.name
  to_zone   = junos_security_zone.testacc_v0to1_secuPolicy.name
  policy {
    name                          = "testacc_v0to1_secuPolicy_1"
    match_source_address          = ["testacc_address1"]
    match_destination_address     = ["any"]
    match_application             = ["junos-ssh"]
    match_source_address_excluded = true
    log_init                      = true
    log_close                     = true
    count                         = true
    permit_application_services {
      advanced_anti_malware_policy = junos_services_advanced_anti_malware_policy.testacc_v0to1_secuPolicy.name
      idp_policy                   = junos_security_idp_policy.testacc_v0to1_secuPolicy.name
      redirect_wx                  = true
      ssl_proxy {}
      uac_policy {}
    }
  }
  policy {
    name                               = "testacc_Policy_2"
    match_source_address               = ["testacc_address1"]
    match_destination_address          = ["testacc_address1"]
    match_destination_address_excluded = true
    match_application                  = ["any"]
    then                               = "reject"
  }
}

resource "junos_security_zone" "testacc_v0to1_secuPolicy" {
  name = "testacc_v0to1_secuPolicy"
  address_book {
    name    = "testacc_address1"
    network = "192.0.2.0/25"
  }
}
`
}
