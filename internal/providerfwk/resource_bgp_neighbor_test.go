package providerfwk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccJunosBgpNeighbor_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosBgpNeighborConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"routing_instance", "testacc_bgpneighbor"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"group", "testacc_bgpneighbor"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"advertise_inactive", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"advertise_peer_as", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"as_override", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"cluster", "192.0.2.3"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"damping", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"log_updown", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"mtu_discovery", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"remove_private", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"passive", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"hold_time", "30"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"local_as", "65001"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"local_as_private", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"local_as_loops", "1"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"local_preference", "100"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"metric_out", "100"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"out_delay", "30"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"peer_as", "65002"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"preference", "100"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"authentication_algorithm", "md5"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"local_address", "192.0.2.3"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"export.#", "1"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"export.0", "testacc_bgpneighbor"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"import.#", "1"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"import.0", "testacc_bgpneighbor"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"bfd_liveness_detection.detection_time_threshold", "60"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"bfd_liveness_detection.transmit_interval_threshold", "30"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"bfd_liveness_detection.transmit_interval_minimum_interval", "10"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"bfd_liveness_detection.holddown_interval", "10"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"bfd_liveness_detection.minimum_interval", "10"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"bfd_liveness_detection.minimum_receive_interval", "10"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"bfd_liveness_detection.multiplier", "2"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"bfd_liveness_detection.session_mode", "automatic"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.#", "2"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.0.nlri_type", "unicast"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.0.accepted_prefix_limit.maximum", "2"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.0.accepted_prefix_limit.teardown", "50"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.0.accepted_prefix_limit.teardown_idle_timeout", "30"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.0.prefix_limit.maximum", "2"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.0.prefix_limit.teardown", "50"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.0.prefix_limit.teardown_idle_timeout", "30"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.1.nlri_type", "multicast"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.1.accepted_prefix_limit.teardown_idle_timeout_forever", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet.1.prefix_limit.teardown_idle_timeout_forever", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"family_inet6.#", "2"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"graceful_restart.disable", "true"),
					),
				},
				{
					ResourceName:      "junos_bgp_neighbor.testacc_bgpneighbor",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosBgpNeighborConfigUpdate(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectSensitiveValue("junos_bgp_neighbor.testacc_bgpneighbor",
								tfjsonpath.New("authentication_key")),
						},
					},
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"routing_instance", "testacc_bgpneighbor"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"group", "testacc_bgpneighbor"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"advertise_external_conditional", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"bgp_multipath.multiple_as", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"no_advertise_peer_as", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"metric_out_igp_offset", "-10"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"metric_out_igp_delay_med_update", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"authentication_key", "password"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"graceful_restart.restart_time", "10"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor",
							"graceful_restart.stale_route_time", "10"),
					),
				},
				{
					Config: testAccJunosBgpNeighborConfigUpdate2(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor2",
							"advertise_external", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor2",
							"accept_remote_nexthop", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor2",
							"multihop", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor2",
							"local_as_no_prepend_global_as", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor2",
							"metric_out_minimum_igp_offset", "-10"),
					),
				},
				{
					Config: testAccJunosBgpNeighborConfigUpdate3(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor2",
							"local_as_alias", "true"),
						resource.TestCheckResourceAttr("junos_bgp_neighbor.testacc_bgpneighbor2",
							"metric_out_minimum_igp", "true"),
					),
				},
			},
		})
	}
}

func testAccJunosBgpNeighborConfigCreate() string {
	return `
resource "junos_routing_options" "testacc_bgpneighbor" {
  clean_on_destroy = true
  autonomous_system {
    number = "65001"
  }
  graceful_restart {}
}
resource "junos_routing_instance" "testacc_bgpneighbor" {
  name = "testacc_bgpneighbor"
  as   = "65000"
}
resource "junos_policyoptions_policy_statement" "testacc_bgpneighbor" {
  lifecycle {
    create_before_destroy = true
  }
  name = "testacc_bgpneighbor"
  then {
    action = "accept"
  }
}
resource "junos_bgp_group" "testacc_bgpneighbor" {
  name             = "testacc_bgpneighbor"
  routing_instance = junos_routing_instance.testacc_bgpneighbor.name
}
resource "junos_bgp_neighbor" "testacc_bgpneighbor" {
  depends_on = [
    junos_routing_options.testacc_bgpneighbor
  ]
  ip                 = "192.0.2.4"
  routing_instance   = junos_routing_instance.testacc_bgpneighbor.name
  group              = junos_bgp_group.testacc_bgpneighbor.name
  advertise_inactive = true
  advertise_peer_as  = true
  as_override        = true
  bgp_multipath {}
  cluster                  = "192.0.2.3"
  damping                  = true
  log_updown               = true
  mtu_discovery            = true
  remove_private           = true
  passive                  = true
  hold_time                = 30
  keep_all                 = true
  local_as                 = "65001"
  local_as_private         = true
  local_as_loops           = 1
  local_preference         = 100
  metric_out               = 100
  out_delay                = 30
  peer_as                  = "65002"
  preference               = 100
  authentication_algorithm = "md5"
  local_address            = "192.0.2.3"
  export                   = [junos_policyoptions_policy_statement.testacc_bgpneighbor.name]
  import                   = [junos_policyoptions_policy_statement.testacc_bgpneighbor.name]
  bfd_liveness_detection {
    detection_time_threshold           = 60
    transmit_interval_threshold        = 30
    transmit_interval_minimum_interval = 10
    holddown_interval                  = 10
    minimum_interval                   = 10
    minimum_receive_interval           = 10
    multiplier                         = 2
    session_mode                       = "automatic"
  }
  family_inet {
    nlri_type = "unicast"
    accepted_prefix_limit {
      maximum               = 2
      teardown              = 50
      teardown_idle_timeout = 30
    }
    prefix_limit {
      maximum               = 2
      teardown              = 50
      teardown_idle_timeout = 30
    }
  }
  family_inet {
    nlri_type = "multicast"
    accepted_prefix_limit {
      maximum                       = 2
      teardown_idle_timeout_forever = true
    }
    prefix_limit {
      maximum                       = 2
      teardown_idle_timeout_forever = true
    }
  }
  family_inet6 {
    nlri_type = "unicast"
    accepted_prefix_limit {
      maximum               = 2
      teardown              = 50
      teardown_idle_timeout = 30
    }
    prefix_limit {
      maximum               = 2
      teardown              = 50
      teardown_idle_timeout = 30
    }
  }
  family_inet6 {
    nlri_type = "multicast"
  }
  graceful_restart {
    disable = true
  }
}
`
}

func testAccJunosBgpNeighborConfigUpdate() string {
	return `
resource "junos_routing_options" "testacc_bgpneighbor" {
  clean_on_destroy = true
  autonomous_system {
    number = "65001"
  }
  graceful_restart {}
}
resource "junos_routing_instance" "testacc_bgpneighbor" {
  name = "testacc_bgpneighbor"
  as   = "65000"
}
resource "junos_bgp_group" "testacc_bgpneighbor" {
  name             = "testacc_bgpneighbor"
  routing_instance = junos_routing_instance.testacc_bgpneighbor.name
  type             = "internal"
}
resource "junos_bgp_neighbor" "testacc_bgpneighbor" {
  depends_on = [
    junos_routing_options.testacc_bgpneighbor
  ]
  ip                              = "192.0.2.4"
  routing_instance                = junos_routing_instance.testacc_bgpneighbor.name
  group                           = junos_bgp_group.testacc_bgpneighbor.name
  description                     = "peer 2.4"
  advertise_external_conditional  = true
  keep_none                       = true
  no_advertise_peer_as            = true
  metric_out_igp_offset           = -10
  metric_out_igp_delay_med_update = true
  authentication_key              = "password"
  bgp_multipath {
    multiple_as = true
  }
  graceful_restart {
    restart_time     = 10
    stale_route_time = 10
  }
  tcp_aggressive_transmission = true
  bgp_error_tolerance {}
}
`
}

func testAccJunosBgpNeighborConfigUpdate2() string {
	return `
resource "junos_routing_options" "testacc_bgpneighbor" {
  clean_on_destroy = true
  autonomous_system {
    number = "65001"
  }
  graceful_restart {}
}
resource "junos_routing_instance" "testacc_bgpneighbor2" {
  name = "testacc_bgpneighbor2"
  as   = "65000"
}
resource "junos_bgp_group" "testacc_bgpneighbor2" {
  name             = "testacc_bgpneighbor2"
  routing_instance = junos_routing_instance.testacc_bgpneighbor2.name
  type             = "internal"
}
resource "junos_bgp_neighbor" "testacc_bgpneighbor2" {
  depends_on = [
    junos_routing_options.testacc_bgpneighbor
  ]
  ip                            = "192.0.2.4"
  routing_instance              = junos_routing_instance.testacc_bgpneighbor2.name
  group                         = junos_bgp_group.testacc_bgpneighbor2.name
  advertise_external            = true
  accept_remote_nexthop         = true
  multihop                      = true
  local_as                      = "65000"
  local_as_no_prepend_global_as = true
  metric_out_minimum_igp_offset = -10
}
resource "junos_bgp_group" "testacc_bgpneighbor2b" {
  depends_on = [
    junos_routing_options.testacc_bgpneighbor
  ]
  name = "testacc_bgpneighbor2b"
  type = "internal"
}
resource "junos_bgp_neighbor" "testacc_bgpneighbor2b" {
  ip    = "192.0.2.5"
  group = junos_bgp_group.testacc_bgpneighbor2b.name
  family_evpn {
    accepted_prefix_limit {
      maximum               = 2
      teardown              = 50
      teardown_idle_timeout = 30
    }
    prefix_limit {
      maximum               = 2
      teardown              = 50
      teardown_idle_timeout = 30
    }
  }
  bgp_error_tolerance {
    malformed_route_limit         = 234
    malformed_update_log_interval = 567
  }
}
`
}

func testAccJunosBgpNeighborConfigUpdate3() string {
	return `
resource "junos_routing_options" "testacc_bgpneighbor" {
  clean_on_destroy = true
  autonomous_system {
    number = "65001"
  }
  graceful_restart {}
}
resource "junos_routing_instance" "testacc_bgpneighbor2" {
  name = "testacc_bgpneighbor2"
  as   = "65000"
}
resource "junos_bgp_group" "testacc_bgpneighbor2" {
  name             = "testacc_bgpneighbor2"
  routing_instance = junos_routing_instance.testacc_bgpneighbor2.name
  type             = "internal"
}
resource "junos_bgp_neighbor" "testacc_bgpneighbor2" {
  depends_on = [
    junos_routing_options.testacc_bgpneighbor
  ]
  ip                     = "192.0.2.4"
  routing_instance       = junos_routing_instance.testacc_bgpneighbor2.name
  group                  = junos_bgp_group.testacc_bgpneighbor2.name
  local_as               = "65000"
  local_as_alias         = true
  metric_out_minimum_igp = true
}
resource "junos_bgp_group" "testacc_bgpneighbor2b" {
  depends_on = [
    junos_routing_options.testacc_bgpneighbor
  ]
  name = "testacc_bgpneighbor2b"
  type = "internal"
}
resource "junos_bgp_neighbor" "testacc_bgpneighbor2b" {
  ip    = "192.0.2.5"
  group = junos_bgp_group.testacc_bgpneighbor2b.name
  family_evpn {}
  bgp_error_tolerance {
    no_malformed_route_limit = true
  }
}
`
}
