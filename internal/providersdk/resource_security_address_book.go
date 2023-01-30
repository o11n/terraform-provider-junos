package providersdk

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
)

type addressBookOptions struct {
	name            string
	description     string
	attachZone      []string
	networkAddress  []map[string]interface{}
	wildcardAddress []map[string]interface{}
	dnsName         []map[string]interface{}
	rangeAddress    []map[string]interface{}
	addressSet      []map[string]interface{}
}

func resourceSecurityAddressBook() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSecurityAddressBookCreate,
		ReadWithoutTimeout:   resourceSecurityAddressBookRead,
		UpdateWithoutTimeout: resourceSecurityAddressBookUpdate,
		DeleteWithoutTimeout: resourceSecurityAddressBookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityAddressBookImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "global",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attach_zone": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatDefault),
				},
			},
			"network_address": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
						},
						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDRNetwork(0, 128),
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},
			"wildcard_address": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
						},
						"value": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateWildcardFunc(),
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},
			"dns_name": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"ipv4_only": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"ipv6_only": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"range_address": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
						},
						"from": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"to": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPAddress,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},
			"address_set": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
						},
						"address": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
							},
						},
						"address_set": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validateNameObjectJunos([]string{}, 64, formatAddressName),
							},
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},
		},
	}
}

func resourceSecurityAddressBookCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setSecurityAddressBook(d, junSess); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(d.Get("name").(string))

		return nil
	}
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()
	if !junSess.CheckCompatibilitySecurity() {
		return diag.FromErr(fmt.Errorf("security policy not compatible with Junos device %s",
			junSess.SystemInformation.HardwareModel))
	}
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	addressBookExists, err := checkSecurityAddressBookExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if addressBookExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("security address book %v already exists", d.Get("name").(string)))...)
	}
	if err := setSecurityAddressBook(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_security_address_book")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	addressBookExists, err = checkSecurityAddressBookExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if addressBookExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("security address book  %v does not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceSecurityAddressBookReadWJunSess(d, junSess)...)
}

func resourceSecurityAddressBookRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceSecurityAddressBookReadWJunSess(d, junSess)
}

func resourceSecurityAddressBookReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	mutex.Lock()
	addressOptions, err := readSecurityAddressBook(d.Get("name").(string), junSess)
	mutex.Unlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if addressOptions.name == "" {
		d.SetId("")
	} else {
		fillSecurityAddressBookData(d, addressOptions)
	}

	return nil
}

func resourceSecurityAddressBookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityAddressBook(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setSecurityAddressBook(d, junSess); err != nil {
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
	if err := delSecurityAddressBook(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setSecurityAddressBook(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_security_address_book")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceSecurityAddressBookReadWJunSess(d, junSess)...)
}

func resourceSecurityAddressBookDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delSecurityAddressBook(d.Get("name").(string), junSess); err != nil {
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
	if err := delSecurityAddressBook(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_security_address_book")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceSecurityAddressBookImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)
	securityAddressBookExists, err := checkSecurityAddressBookExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !securityAddressBookExists {
		return nil, fmt.Errorf("don't find address book with id '%v' (id must be <name>)", d.Id())
	}
	addressOptions, err := readSecurityAddressBook(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillSecurityAddressBookData(d, addressOptions)

	result[0] = d

	return result, nil
}

func checkSecurityAddressBookExists(addrBook string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security address-book " + addrBook + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setSecurityAddressBook(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)
	setPrefix := "set security address-book " + d.Get("name").(string)

	if d.Get("description").(string) != "" {
		configSet = append(configSet, setPrefix+" description \""+d.Get("description").(string)+"\"")
	}
	for _, v := range d.Get("attach_zone").([]interface{}) {
		if d.Get("name").(string) == "global" {
			return fmt.Errorf("cannot attach global address book to a zone")
		}
		attachZone := v.(string)
		configSet = append(configSet, setPrefix+" attach zone "+attachZone)
	}
	addressNameList := make([]string, 0)
	for _, v := range d.Get("network_address").(*schema.Set).List() {
		address := v.(map[string]interface{})
		if bchk.InSlice(address["name"].(string), addressNameList) {
			return fmt.Errorf("multiple addresses with the same name %s", address["name"].(string))
		}
		addressNameList = append(addressNameList, address["name"].(string))
		setPrefixAddr := setPrefix + " address " + address["name"].(string) + " "
		configSet = append(configSet, setPrefixAddr+address["value"].(string))
		if address["description"].(string) != "" {
			configSet = append(configSet, setPrefixAddr+"description \""+address["description"].(string)+"\"")
		}
	}
	for _, v := range d.Get("wildcard_address").(*schema.Set).List() {
		address := v.(map[string]interface{})
		if bchk.InSlice(address["name"].(string), addressNameList) {
			return fmt.Errorf("multiple addresses with the same name %s", address["name"].(string))
		}
		addressNameList = append(addressNameList, address["name"].(string))
		setPrefixAddr := setPrefix + " address " + address["name"].(string) + " "
		configSet = append(configSet, setPrefixAddr+"wildcard-address "+address["value"].(string))
		if address["description"].(string) != "" {
			configSet = append(configSet, setPrefixAddr+"description \""+address["description"].(string)+"\"")
		}
	}
	for _, v := range d.Get("dns_name").(*schema.Set).List() {
		address := v.(map[string]interface{})
		if bchk.InSlice(address["name"].(string), addressNameList) {
			return fmt.Errorf("multiple addresses with the same name %s", address["name"].(string))
		}
		addressNameList = append(addressNameList, address["name"].(string))
		setPrefixAddr := setPrefix + " address " + address["name"].(string) + " "
		configSet = append(configSet, setPrefixAddr+"dns-name "+address["value"].(string))
		if address["ipv4_only"].(bool) {
			configSet = append(configSet, setPrefixAddr+"dns-name "+address["value"].(string)+" ipv4-only")
		}
		if address["ipv6_only"].(bool) {
			configSet = append(configSet, setPrefixAddr+"dns-name "+address["value"].(string)+" ipv6-only")
		}
		if address["description"].(string) != "" {
			configSet = append(configSet, setPrefixAddr+"description \""+address["description"].(string)+"\"")
		}
	}
	for _, v := range d.Get("range_address").(*schema.Set).List() {
		address := v.(map[string]interface{})
		if bchk.InSlice(address["name"].(string), addressNameList) {
			return fmt.Errorf("multiple addresses with the same name %s", address["name"].(string))
		}
		addressNameList = append(addressNameList, address["name"].(string))
		setPrefixAddr := setPrefix + " address " + address["name"].(string) + " "
		configSet = append(configSet, setPrefixAddr+"range-address "+address["from"].(string)+" to "+address["to"].(string))
		if address["description"].(string) != "" {
			configSet = append(configSet, setPrefixAddr+"description \""+address["description"].(string)+"\"")
		}
	}
	for _, v := range d.Get("address_set").(*schema.Set).List() {
		addressSet := v.(map[string]interface{})
		if bchk.InSlice(addressSet["name"].(string), addressNameList) {
			return fmt.Errorf("multiple addresses or address-sets with the same name %s", addressSet["name"].(string))
		}
		addressNameList = append(addressNameList, addressSet["name"].(string))
		setPrefixAddrSet := setPrefix + " address-set " + addressSet["name"].(string) + " "
		if len(addressSet["address"].(*schema.Set).List()) == 0 &&
			len(addressSet["address_set"].(*schema.Set).List()) == 0 {
			return fmt.Errorf("at least one of address or address_set is required "+
				"in address_set %s", addressSet["name"].(string))
		}
		for _, addr := range sortSetOfString(addressSet["address"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixAddrSet+"address "+addr)
		}
		for _, addrSet := range sortSetOfString(addressSet["address_set"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixAddrSet+"address-set "+addrSet)
		}
		if addressSet["description"].(string) != "" {
			configSet = append(configSet, setPrefixAddrSet+"description \""+addressSet["description"].(string)+"\"")
		}
	}

	return junSess.ConfigSet(configSet)
}

func readSecurityAddressBook(addrBook string, junSess *junos.Session,
) (confRead addressBookOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"security address-book " + addrBook + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	descMap := make(map[string]string)
	if showConfig != junos.EmptyW {
		confRead.name = addrBook
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
			case balt.CutPrefixInString(&itemTrim, "address "):
				itemTrimFields := strings.Split(itemTrim, " ")
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "description "):
					descMap[itemTrimFields[0]] = strings.Trim(itemTrim, "\"")
				case balt.CutPrefixInString(&itemTrim, "wildcard-address "):
					confRead.wildcardAddress = append(confRead.wildcardAddress, map[string]interface{}{
						"name":        itemTrimFields[0],
						"value":       itemTrim,
						"description": descMap[itemTrimFields[0]],
					})
				case balt.CutPrefixInString(&itemTrim, "range-address "):
					rangeAddressFields := strings.Split(itemTrim, " ")
					if len(rangeAddressFields) < 3 { // <from> to <to>
						return confRead, fmt.Errorf(junos.CantReadValuesNotEnoughFields, "range-address", itemTrim)
					}
					confRead.rangeAddress = append(confRead.rangeAddress, map[string]interface{}{
						"name":        itemTrimFields[0],
						"from":        rangeAddressFields[0],
						"to":          rangeAddressFields[2],
						"description": descMap[itemTrimFields[0]],
					})
				case balt.CutPrefixInString(&itemTrim, "dns-name "):
					switch {
					case balt.CutSuffixInString(&itemTrim, " ipv4-only"):
						confRead.dnsName = append(confRead.dnsName, map[string]interface{}{
							"name":        itemTrimFields[0],
							"value":       itemTrim,
							"description": descMap[itemTrimFields[0]],
							"ipv4_only":   true,
							"ipv6_only":   false,
						})
					case balt.CutSuffixInString(&itemTrim, " ipv6-only"):
						confRead.dnsName = append(confRead.dnsName, map[string]interface{}{
							"name":        itemTrimFields[0],
							"value":       itemTrim,
							"description": descMap[itemTrimFields[0]],
							"ipv4_only":   false,
							"ipv6_only":   true,
						})
					default:
						confRead.dnsName = append(confRead.dnsName, map[string]interface{}{
							"name":        itemTrimFields[0],
							"value":       itemTrim,
							"description": descMap[itemTrimFields[0]],
							"ipv4_only":   false,
							"ipv6_only":   false,
						})
					}
				default:
					confRead.networkAddress = append(confRead.networkAddress, map[string]interface{}{
						"name":        itemTrimFields[0],
						"value":       itemTrim,
						"description": descMap[itemTrimFields[0]],
					})
				}
			case balt.CutPrefixInString(&itemTrim, "address-set "):
				itemTrimFields := strings.Split(itemTrim, " ")
				adSet := map[string]interface{}{
					"name":        itemTrimFields[0],
					"address":     make([]string, 0),
					"address_set": make([]string, 0),
					"description": "",
				}
				confRead.addressSet = copyAndRemoveItemMapList("name", adSet, confRead.addressSet)
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				switch {
				case balt.CutPrefixInString(&itemTrim, "description "):
					adSet["description"] = strings.Trim(itemTrim, "\"")
				case balt.CutPrefixInString(&itemTrim, "address "):
					adSet["address"] = append(adSet["address"].([]string), itemTrim)
				case balt.CutPrefixInString(&itemTrim, "address-set "):
					adSet["address_set"] = append(adSet["address_set"].([]string), itemTrim)
				}
				confRead.addressSet = append(confRead.addressSet, adSet)
			case balt.CutPrefixInString(&itemTrim, "attach zone "):
				confRead.attachZone = append(confRead.attachZone, itemTrim)
			}
		}
	}
	copySecurityAddressBookAddressDescriptions(descMap, confRead.networkAddress)
	copySecurityAddressBookAddressDescriptions(descMap, confRead.dnsName)
	copySecurityAddressBookAddressDescriptions(descMap, confRead.rangeAddress)
	copySecurityAddressBookAddressDescriptions(descMap, confRead.wildcardAddress)

	return confRead, nil
}

func delSecurityAddressBook(addrBook string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete security address-book "+addrBook)

	return junSess.ConfigSet(configSet)
}

func fillSecurityAddressBookData(d *schema.ResourceData, addressOptions addressBookOptions) {
	if tfErr := d.Set("name", addressOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("description", addressOptions.description); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("attach_zone", addressOptions.attachZone); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("network_address", addressOptions.networkAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("wildcard_address", addressOptions.wildcardAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("dns_name", addressOptions.dnsName); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("range_address", addressOptions.rangeAddress); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("address_set", addressOptions.addressSet); tfErr != nil {
		panic(tfErr)
	}
}

func copySecurityAddressBookAddressDescriptions(descMap map[string]string, addrList []map[string]interface{}) {
	for _, ele := range addrList {
		ele["description"] = descMap[ele["name"].(string)]
	}
}
