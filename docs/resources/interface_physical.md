---
page_title: "Junos: junos_interface_physical"
---

# junos_interface_physical

Provides a physical interface resource.

## Example Usage

```hcl
# Configure interface of switch
resource "junos_interface_physical" "interface_switch_demo" {
  name         = "ge-0/0/0"
  description  = "interfaceSwitchDemo"
  trunk        = true
  vlan_members = ["100"]
}
# Prepare physical interface for L3 logical interfaces on Junos Router or firewall
resource "junos_interface_physical" "interface_fw_demo" {
  name         = "ge-0/0/1"
  description  = "interfaceFwDemo"
  vlan_tagging = true
}
```

## Argument Reference

~> **NOTE:** This resource computes the maximum number of aggregate interfaces required with the
current configuration (searches lines `ether-options 802.3ad` and `ae` interfaces set) then
add/remove `chassis aggregated-devices ethernet device-count` line with this maximum.

The following arguments are supported:

- **name** (Required, String, Forces new resource)  
  Name of physical interface (without dot).
- **no_disable_on_destroy** (Optional, Boolean)  
  When destroy this resource, delete all configurations => do not add
  `disable` + `description NC` or `apply-groups` with `group_interface_delete` provider argument on interface.
- **description** (Optional, String)  
  Description for interface.
- **disable** (Optional, Boolean)  
  Disable this interface.
- **encapsulation** (Optional, String)  
  Physical link-layer encapsulation.
- **esi** (Optional, Block)  
  Define ESI Config parameters.  
  See [below for nested schema](#esi-arguments).
- **ether_opts** (Optional, Block)  
  Declare `ether-options` configuration.  
  Conflict with `gigether_opts`.
  - **ae_8023ad** (Optional, String)  
    Name of an aggregated Ethernet interface to join.
  - **auto_negotiation** (Optional, Boolean)  
    Enable auto-negotiation.
  - **no_auto_negotiation** (Optional, Boolean)  
    Don't enable auto-negotiation.
  - **flow_control** (Optional, Boolean)  
    Enable flow control.
  - **no_flow_control** (Optional, Boolean)  
    Don't enable flow control.
  - **loopback** (Optional, Boolean)  
    Enable loopback.
  - **no_loopback** (Optional, Boolean)  
    Don't enable loopback.
  - **redundant_parent** (Optional, String)  
    Name of a redundant ethernet interface to join.
- **flexible_vlan_tagging** (Optional, Boolean)  
  Support for no tagging, or single and double 802.1q VLAN tagging.
- **gigether_opts** (Optional, Block)  
  Declare `gigether-options` configuration.  
  Conflict with `ether_opts`.
  - **ae_8023ad** (Optional, String)  
    Name of an aggregated Ethernet interface to join.
  - **auto_negotiation** (Optional, Boolean)  
    Enable auto-negotiation.
  - **no_auto_negotiation** (Optional, Boolean)  
    Don't enable auto-negotiation.
  - **flow_control** (Optional, Boolean)  
    Enable flow control.
  - **no_flow_control** (Optional, Boolean)  
    Don't enable flow control.
  - **loopback** (Optional, Boolean)  
    Enable loopback.
  - **no_loopback** (Optional, Boolean)  
    Don't enable loopback.
  - **redundant_parent** (Optional, String)  
    Name of a redundant ethernet interface to join.
- **gratuitous_arp_reply** (Optional, Boolean)  
  Enable gratuitous ARP reply.
- **hold_time_down** (Optional, Number)  
  Link down hold time (milliseconds).  
  `hold_time_up` must also be specified.
- **hold_time_up** (Optional, Number)  
  Link up hold time (milliseconds).  
  `hold_time_down` must also be specified.
- **link_mode** (Optional, String)  
  Link operational mode.  
  Need to be `automatic`, `full-duplex` or `half-duplex`.
- **mtu** (Optional, Number)  
  Maximum transmission unit.
- **no_gratuitous_arp_reply** (Optional, Boolean)  
  Don't enable gratuitous ARP reply.
- **no_gratuitous_arp_request** (Optional, Boolean)  
  Ignore gratuitous ARP request.
- **parent_ether_opts** (Optional, Block)  
  Declare `aggregated-ether-options` or `redundant-ether-options` configuration
  (it depends on the interface `name`).  
  See [below for nested schema](#parent_ether_opts-arguments).
- **speed** (Optional, String)  
  Link speed.  
  Must be a valid speed (10m | 100m | 1g ...)
- **trunk** (Optional, Boolean)  
  Interface mode is trunk.
- **vlan_members** (Optional, List of String)  
  List of vlan for membership for this interface.
- **vlan_native** (Optional, Number)  
  Vlan for untagged frames.
- **vlan_tagging** (Optional, Boolean)  
  Add 802.1q VLAN tagging support.

---

### esi arguments

- **mode** (Required, String)  
  ESI Mode.
- **auto_derive_lacp** (Optional, Boolean)  
  Auto-derive ESI value for the interface.
- **df_election_type** (Optional, String)  
  DF Election Type.
- **identifier** (Optional, String)  
  The ESI value for the interface.
- **source_bmac** (Optional, String)  
  Unicast Source B-MAC address per ESI for PBB-EVPN.

---

### parent_ether_opts arguments

- **bfd_liveness_detection** (Optional, Block)  
  Declare `bfd-liveness-detection` in `aggregated-ether-options` configuration.  
  See [below for nested schema](#bfd_liveness_detection-arguments-in-parent_ether_opts).
- **flow_control** (Optional, Boolean)  
  Enable flow control.
- **no_flow_control** (Optional, Boolean)  
  Don't enable flow control.
- **lacp** (Optional, Block)  
  Declare `lacp` configuration.
  - **mode** (Required, String)  
    Active or passive.
  - **admin_key** (Optional, Number)  
    Node's administrative key.
  - **periodic** (Optional, String)  
    Timer interval for periodic transmission of LACP packets.  
    Need to be `fast` or `slow`.
  - **sync_reset** (Optional, String)  
    On minimum-link failure notify out of sync to peer.  
    Need to be `disable` or `enable`.
  - **system_id** (Optional, String)  
    Node's System ID, encoded as a MAC address
  - **system_priority** (Optional, Number)  
    Priority of the system (0 ... 65535).
- **loopback** (Optional, Boolean)  
  Enable loopback.
- **no_loopback** (Optional, Boolean)  
  Don't enable loopback.
- **link_speed** (Optional, String)  
  Link speed of individual interface that joins the AE.
- **minimum_bandwidth** (Optional, String)  
  Minimum bandwidth configured for aggregated bundle.  
  Need to be `N (k|g|m)?bps` format.
- **minimum_links** (Optional, Number)  
  Minimum number of aggregated/active links (1..64).
- **redundancy_group** (Optional, Number)  
  Redundancy group of this interface (1..128) for reth interface.
- **source_address_filter** (Optional, List of String)  
  Source address filters.
- **source_filtering** (Optional, Boolean)  
  Enable source address filtering.

---

### bfd_liveness_detection arguments in parent_ether_opts

- **local_address** (Required, String)  
  BFD local address.
- **authentication_algorithm** (Optional, String)  
  Authentication algorithm name.
- **authentication_key_chain** (Optional, String)  
  Authentication Key chain name.
- **authentication_loose_check** (Optional, Boolean)  
  Verify authentication only if authentication is negotiated.
- **detection_time_threshold** (Optional, Number)  
  High detection-time triggering a trap (milliseconds).
- **holddown_interval** (Optional, Number)  
  Time to hold the session-UP notification to the client (0..255000 milliseconds).
- **minimum_interval** (Optional, Number)  
  Minimum transmit and receive interval (1..255000 milliseconds).
- **minimum_receive_interval** (Optional, Number)  
  Minimum receive interval (1..255000 milliseconds).
- **multiplier** (Optional, Number)  
  Detection time multiplier (1..255).
- **neighbor** (Optional, String)  
  BFD neighbor address.
- **no_adaptation** (Optional, Boolean)  
  Disable adaptation.
- **transmit_interval_minimum_interval** (Optional, Number)  
  Minimum transmit interval (1..255000 milliseconds).
- **transmit_interval_threshold** (Optional, Number)  
  High transmit interval triggering a trap (milliseconds).
- **version** (Optional, String)  
  BFD protocol version number.

## Attributes Reference

The following attributes are exported:

- **id** (String)  
  An identifier for the resource with format `<name>`.

## Import

Junos interface can be imported using an id made up of `<name>`, e.g.

```shell
$ terraform import junos_interface_physical.interface_switch_demo ge-0/0/0
$ terraform import junos_interface_physical.interface_fw_demo_100 ge-0/0/1
```
