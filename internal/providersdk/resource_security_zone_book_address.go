package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type zoneBookAddressOptions struct {
	dnsIPv4Only bool
	dnsIPv6Only bool
	cidr        string
	description string
	dnsName     string
	name        string
	rangeFrom   string
	rangeTo     string
	wildcard    string
	zone        string
}

func resourceSecurityZoneBookAddress() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityZoneBookAddressCreate,
		ReadWithoutTimeout:   resourceSecurityZoneBookAddressRead,
		UpdateWithoutTimeout: resourceSecurityZoneBookAddressUpdate,
		DeleteWithoutTimeout: resourceSecurityZoneBookAddressDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityZoneBookAddressImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
			},
			"zone": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
			},
			"cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsCIDRNetwork(0, 128),
				ExactlyOneOf: []string{"cidr", "dns_name", "range_from", "wildcard"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_ipv4_only": {
				Type:          schema.TypeBool,
				Optional:      true,
				RequiredWith:  []string{"dns_name"},
				ConflictsWith: []string{"dns_ipv6_only"},
			},
			"dns_ipv6_only": {
				Type:          schema.TypeBool,
				Optional:      true,
				RequiredWith:  []string{"dns_name"},
				ConflictsWith: []string{"dns_ipv4_only"},
			},
			"dns_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"cidr", "dns_name", "range_from", "wildcard"},
			},
			"range_from": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
				RequiredWith: []string{"range_to"},
				ExactlyOneOf: []string{"cidr", "dns_name", "range_from", "wildcard"},
			},
			"range_to": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
				RequiredWith: []string{"range_from"},
			},
			"wildcard": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateWildcardFunc(),
				ExactlyOneOf:     []string{"cidr", "dns_name", "range_from", "wildcard"},
			},
		},
	}
}

func resourceSecurityZoneBookAddressCreate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSecurityZoneBookAddress(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("zone").(string) + junos.IDSeparator + d.Get("name").(string))

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if !junSess.CheckCompatibilitySecurity() {
		return diag.FromErr(fmt.Errorf("security zone address-book address not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	zonesExists, err := checkSecurityZonesExists(d.Get("zone").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if !zonesExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security zone %v doesn't exist", d.Get("zone").(string)))...)
	}
	securityZoneBookAddressExists, err := checkSecurityZoneBookAddresssExists(
		d.Get("zone").(string),
		d.Get("name").(string),
		junSess,
	)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityZoneBookAddressExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security zone address-book address %v already exists in zone %s",
			d.Get("name").(string), d.Get("zone").(string)))...)
	}

	if err := setSecurityZoneBookAddress(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_security_zone_book_address")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	securityZoneBookAddressExists, err = checkSecurityZoneBookAddresssExists(
		d.Get("zone").(string),
		d.Get("name").(string),
		junSess,
	)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if securityZoneBookAddressExists {
		d.SetId(d.Get("zone").(string) + junos.IDSeparator + d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf(
			"security zone address-book address %v not exists in zone %s after commit "+
				"=> check your config", d.Get("name").(string), d.Get("zone").(string)))...)
	}

	return append(diagWarns, resourceSecurityZoneBookAddressReadWJunSess(d, junSess)...)
}

func resourceSecurityZoneBookAddressRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSecurityZoneBookAddressReadWJunSess(d, junSess)
}

func resourceSecurityZoneBookAddressReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	zoneBookAddressOptions, err := readSecurityZoneBookAddress(
		d.Get("zone").(string),
		d.Get("name").(string),
		junSess,
	)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if zoneBookAddressOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityZoneBookAddressData(d, zoneBookAddressOptions)
	}

	return nil
}

func resourceSecurityZoneBookAddressUpdate(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityZoneBookAddress(d.Get("zone").(string), d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityZoneBookAddress(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.Partial(false)

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSecurityZoneBookAddress(d.Get("zone").(string), d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityZoneBookAddress(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_security_zone_book_address")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityZoneBookAddressReadWJunSess(d, junSess)...)
}

func resourceSecurityZoneBookAddressDelete(ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityZoneBookAddress(d.Get("zone").(string), d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	if err := delSecurityZoneBookAddress(d.Get("zone").(string), d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_security_zone_book_address")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityZoneBookAddressImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	idList := strings.Split(d.Id(), junos.IDSeparator)
	if len(idList) < 2 {
		return nil, fmt.Errorf("missing element(s) in id with separator %v", junos.IDSeparator)
	}
	securityZoneBookAddressExists, err := checkSecurityZoneBookAddresssExists(idList[0], idList[1], junSess)
	if err != nil {
		return nil, err
	}
	if !securityZoneBookAddressExists {
		return nil, fmt.Errorf(
			"don't find zone address-book address with id '%v' (id must be <zone>"+junos.IDSeparator+"<name>)", d.Id())
	}
	zoneBookAddressOptions, err := readSecurityZoneBookAddress(idList[0], idList[1], junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityZoneBookAddressData(d, zoneBookAddressOptions)

	result[0] = d

	return result, nil
}

func checkSecurityZoneBookAddresssExists(zone, address string, junSess *junos.Session,
) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security zones security-zone " + zone + " address-book address " + address + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityZoneBookAddress(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)

	setPrefix := "set security zones security-zone " +
		d.Get("zone").(string) + " address-book address " + d.Get("name").(string) + " "

	if v := d.Get("cidr").(string); v != "" {
		configSet = append(configSet, setPrefix+v)
	}
	if v := d.Get("description").(string); v != "" {
		configSet = append(configSet, setPrefix+"description \""+v+"\"")
	}
	if v := d.Get("dns_name").(string); v != "" {
		configSet = append(configSet, setPrefix+"dns-name "+v)
		if d.Get("dns_ipv4_only").(bool) {
			configSet = append(configSet, setPrefix+"dns-name "+v+" ipv4-only")
		}
		if d.Get("dns_ipv6_only").(bool) {
			configSet = append(configSet, setPrefix+"dns-name "+v+" ipv6-only")
		}
	}
	if v := d.Get("range_from").(string); v != "" {
		configSet = append(configSet, setPrefix+"range-address "+v+" to "+d.Get("range_to").(string))
	}
	if v := d.Get("wildcard").(string); v != "" {
		configSet = append(configSet, setPrefix+"wildcard-address "+v)
	}

	return junSess.ConfigSet(configSet)
}

func readSecurityZoneBookAddress(zone, address string, junSess *junos.Session,
) (confRead zoneBookAddressOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security zones security-zone " + zone + " address-book address " + address + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = address
		confRead.zone = zone
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "description "):
				confRead.description = strings.Trim(itemTrim, "\"")
			case balt.CutPrefixInString(&itemTrim, "dns-name "):
				switch {
				case balt.CutSuffixInString(&itemTrim, " ipv4-only"):
					confRead.dnsIPv4Only = true
					confRead.dnsName = itemTrim
				case balt.CutSuffixInString(&itemTrim, " ipv6-only"):
					confRead.dnsIPv6Only = true
					confRead.dnsName = itemTrim
				default:
					confRead.dnsName = itemTrim
				}
			case balt.CutPrefixInString(&itemTrim, "range-address "):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(itemTrimFields) < 3 { // <from> to <to>
					return confRead, fmt.Errorf(junos.CantReadValuesNotEnoughFields, "range-address", itemTrim)
				}
				confRead.rangeFrom = itemTrimFields[0]
				confRead.rangeTo = itemTrimFields[2]
			case balt.CutPrefixInString(&itemTrim, "wildcard-address "):
				confRead.wildcard = itemTrim
			case strings.Contains(itemTrim, "/"):
				confRead.cidr = itemTrim
			}
		}
	}

	return confRead, nil
}

func delSecurityZoneBookAddress(zone, address string, junSess *junos.Session) error {
	configSet := []string{"delete security zones security-zone " + zone + " address-book address " + address}

	return junSess.ConfigSet(configSet)
}

func fillSecurityZoneBookAddressData(d *schema.ResourceData, zoneBookAddressOptions zoneBookAddressOptions) {
	if tfErr := d.Set("name", zoneBookAddressOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("zone", zoneBookAddressOptions.zone); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("cidr", zoneBookAddressOptions.cidr); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", zoneBookAddressOptions.description); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dns_ipv4_only", zoneBookAddressOptions.dnsIPv4Only); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dns_ipv6_only", zoneBookAddressOptions.dnsIPv6Only); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dns_name", zoneBookAddressOptions.dnsName); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("range_from", zoneBookAddressOptions.rangeFrom); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("range_to", zoneBookAddressOptions.rangeTo); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("wildcard", zoneBookAddressOptions.wildcard); tfErr != nil {
		panic(tfErr)
	}
}
