package providersdk_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccJunosStaticRoute_basic(t *testing.T) {
	if os.Getenv("TESTACC_SWITCH") == "" {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccJunosStaticRouteConfigCreate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"routing_instance", "testacc_staticRoute"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"preference", "100"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"metric", "100"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"next_hop.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"next_hop.0", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"qualified_next_hop.#", "2"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"qualified_next_hop.0.next_hop", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"qualified_next_hop.0.preference", "101"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"qualified_next_hop.0.metric", "101"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"qualified_next_hop.1.next_hop", "192.0.2.250"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"qualified_next_hop.1.interface", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"community.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"community.0", "no-advertise"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"active", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"install", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"readvertise", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"no_resolve", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"retain", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"preference", "100"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"metric", "100"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"next_hop.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"next_hop.0", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"qualified_next_hop.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"qualified_next_hop.0.next_hop", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"qualified_next_hop.0.preference", "101"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"qualified_next_hop.0.metric", "101"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"community.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"community.0", "no-advertise"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"active", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"install", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"readvertise", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"no_resolve", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"retain", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"routing_instance", "testacc_staticRoute"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"preference", "100"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"metric", "100"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"next_hop.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"next_hop.0", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"qualified_next_hop.#", "2"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"qualified_next_hop.0.next_hop", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"qualified_next_hop.0.preference", "101"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"qualified_next_hop.0.metric", "101"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"qualified_next_hop.1.next_hop", "2001:db8:85a4::1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"qualified_next_hop.1.interface", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"community.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"community.0", "no-advertise"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"active", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"install", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"readvertise", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"no_resolve", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"retain", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"preference", "100"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"metric", "100"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"next_hop.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"next_hop.0", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"qualified_next_hop.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"qualified_next_hop.0.next_hop", "st0.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"qualified_next_hop.0.preference", "101"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"qualified_next_hop.0.metric", "101"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"community.#", "1"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"community.0", "no-advertise"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"active", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"install", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"readvertise", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"no_resolve", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"retain", "true"),
					),
				},
				{
					Config: testAccJunosStaticRouteConfigUpdate(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"qualified_next_hop.#", "2"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"qualified_next_hop.1.next_hop", "dsc.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"qualified_next_hop.1.preference", "102"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"qualified_next_hop.1.metric", "102"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"passive", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"no_install", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"no_readvertise", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"no_retain", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"qualified_next_hop.#", "2"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"qualified_next_hop.1.next_hop", "dsc.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"qualified_next_hop.1.preference", "102"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"qualified_next_hop.1.metric", "102"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"passive", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"no_install", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"no_readvertise", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"no_retain", "true"),
					),
				},
				{
					ResourceName:      "junos_static_route.testacc_staticRoute_instance",
					ImportState:       true,
					ImportStateVerify: true,
				},
				{
					Config: testAccJunosStaticRouteConfigCreate2(),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"receive", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance",
							"resolve", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"receive", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default",
							"resolve", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_instance2",
							"discard", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_default2",
							"discard", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default",
							"reject", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance",
							"reject", "true"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_default2",
							"next_table", "testacc_staticRoute2.inet.0"),
						resource.TestCheckResourceAttr("junos_static_route.testacc_staticRoute_ipv6_instance2",
							"next_table", "testacc_staticRoute2.inet6.0"),
					),
				},
			},
		})
	}
}

func testAccJunosStaticRouteConfigCreate() string {
	return `
resource "junos_routing_instance" "testacc_staticRoute" {
  name = "testacc_staticRoute"
}
resource "junos_static_route" "testacc_staticRoute_instance" {
  destination      = "192.0.2.0/24"
  routing_instance = junos_routing_instance.testacc_staticRoute.name
  preference       = 100
  metric           = 100
  next_hop         = ["st0.0"]
  active           = true
  install          = true
  readvertise      = true
  no_resolve       = true
  retain           = true
  qualified_next_hop {
    next_hop   = "st0.0"
    preference = 101
    metric     = 101
  }
  qualified_next_hop {
    next_hop  = "192.0.2.250"
    interface = "st0.0"
  }
  community = ["no-advertise"]
}
resource "junos_static_route" "testacc_staticRoute_default" {
  destination = "192.0.2.0/24"
  preference  = 100
  metric      = 100
  next_hop    = ["st0.0"]
  active      = true
  install     = true
  readvertise = true
  no_resolve  = true
  retain      = true
  qualified_next_hop {
    next_hop   = "st0.0"
    preference = 101
    metric     = 101
  }
  community                    = ["no-advertise"]
  as_path_aggregator_as_number = "65000"
  as_path_aggregator_address   = "192.0.2.1"
  as_path_atomic_aggregate     = true
  as_path_origin               = "igp"
  as_path_path                 = "65000 65000"
}
resource "junos_static_route" "testacc_staticRoute_ipv6_default" {
  destination = "2001:db8:85a3::/48"
  preference  = 100
  metric      = 100
  next_hop    = ["st0.0"]
  active      = true
  install     = true
  readvertise = true
  no_resolve  = true
  retain      = true
  qualified_next_hop {
    next_hop   = "st0.0"
    preference = 101
    metric     = 101
  }
  community                    = ["no-advertise"]
  as_path_aggregator_as_number = "65000"
  as_path_aggregator_address   = "192.0.2.1"
  as_path_atomic_aggregate     = true
  as_path_origin               = "igp"
  as_path_path                 = "65000 65000"
}
resource "junos_static_route" "testacc_staticRoute_ipv6_instance" {
  destination      = "2001:db8:85a3::/48"
  routing_instance = junos_routing_instance.testacc_staticRoute.name
  preference       = 100
  metric           = 100
  next_hop         = ["st0.0"]
  active           = true
  install          = true
  readvertise      = true
  no_resolve       = true
  retain           = true
  qualified_next_hop {
    next_hop   = "st0.0"
    preference = 101
    metric     = 101
  }
  qualified_next_hop {
    next_hop  = "2001:db8:85a4::1"
    interface = "st0.0"
  }
  community = ["no-advertise"]
}
`
}

func testAccJunosStaticRouteConfigUpdate() string {
	return `
resource "junos_routing_instance" "testacc_staticRoute" {
  name = "testacc_staticRoute"
}
resource "junos_static_route" "testacc_staticRoute_instance" {
  destination      = "192.0.2.0/24"
  routing_instance = junos_routing_instance.testacc_staticRoute.name
  preference       = 100
  metric           = 100
  passive          = true
  no_install       = true
  no_readvertise   = true
  no_retain        = true
  qualified_next_hop {
    next_hop   = "st0.0"
    preference = 101
    metric     = 101
  }
  qualified_next_hop {
    next_hop   = "dsc.0"
    preference = 102
    metric     = 102
  }
}
resource "junos_static_route" "testacc_staticRoute_ipv6_default" {
  destination    = "2001:db8:85a3::/48"
  preference     = 100
  metric         = 100
  passive        = true
  no_install     = true
  no_readvertise = true
  no_retain      = true
  qualified_next_hop {
    next_hop   = "st0.0"
    preference = 101
    metric     = 101
  }
  qualified_next_hop {
    next_hop   = "dsc.0"
    preference = 102
    metric     = 102
  }
  community = ["no-advertise"]
}
`
}

func testAccJunosStaticRouteConfigCreate2() string {
	return `
resource "junos_routing_instance" "testacc_staticRoute" {
  name = "testacc_staticRoute"
}
resource "junos_routing_instance" "testacc_staticRoute2" {
  name = "testacc_staticRoute2"
}
resource "junos_static_route" "testacc_staticRoute_instance" {
  destination      = "192.0.2.0/25"
  routing_instance = junos_routing_instance.testacc_staticRoute.name
  receive          = true
  resolve          = true
}
resource "junos_static_route" "testacc_staticRoute_ipv6_default" {
  destination = "2001:db8:85a3::/50"
  receive     = true
  resolve     = true
}
resource "junos_static_route" "testacc_staticRoute_instance2" {
  destination      = "192.0.2.0/26"
  routing_instance = junos_routing_instance.testacc_staticRoute.name
  discard          = true
}
resource "junos_static_route" "testacc_staticRoute_ipv6_default2" {
  destination = "2001:db8:85a3::/52"
  discard     = true
}
resource "junos_static_route" "testacc_staticRoute_default" {
  destination = "192.0.2.0/27"
  reject      = true
}
resource "junos_static_route" "testacc_staticRoute_ipv6_instance" {
  destination      = "2001:db8:85a3::/54"
  routing_instance = junos_routing_instance.testacc_staticRoute.name
  reject           = true
}
resource "junos_static_route" "testacc_staticRoute_default2" {
  destination = "192.0.2.0/28"
  next_table  = "${junos_routing_instance.testacc_staticRoute2.name}.inet.0"
}
resource "junos_static_route" "testacc_staticRoute_ipv6_instance2" {
  destination      = "2001:db8:85a3::/56"
  routing_instance = junos_routing_instance.testacc_staticRoute.name
  next_table       = "${junos_routing_instance.testacc_staticRoute2.name}.inet6.0"
}
`
}
