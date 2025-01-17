package providersdk

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	balt "github.com/jeremmfr/go-utils/basicalter"
	bchk "github.com/jeremmfr/go-utils/basiccheck"
)

type eventoptionsPolicyOptions struct {
	name            string
	events          []string
	attributesMatch []map[string]interface{}
	then            []map[string]interface{}
	within          []map[string]interface{}
}

func resourceEventoptionsPolicy() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceEventoptionsPolicyCreate,
		ReadWithoutTimeout:   resourceEventoptionsPolicyRead,
		UpdateWithoutTimeout: resourceEventoptionsPolicyUpdate,
		DeleteWithoutTimeout: resourceEventoptionsPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceEventoptionsPolicyImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"events": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"then": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"change_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							AtLeastOneOf: []string{
								"then.0.change_configuration",
								"then.0.event_script",
								"then.0.execute_commands",
								"then.0.ignore",
								"then.0.priority_override_facility",
								"then.0.priority_override_severity",
								"then.0.raise_trap",
								"then.0.upload",
							},
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"commands": {
										Type:     schema.TypeList,
										Required: true,
										MinItems: 1,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"commit_options_check": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"commit_options_check_synchronize": {
										Type:         schema.TypeBool,
										Optional:     true,
										RequiredWith: []string{"then.0.change_configuration.0.commit_options_check"},
									},
									"commit_options_force": {
										Type:          schema.TypeBool,
										Optional:      true,
										ConflictsWith: []string{"then.0.change_configuration.0.commit_options_check"},
									},
									"commit_options_log": {
										Type:          schema.TypeString,
										Optional:      true,
										ConflictsWith: []string{"then.0.change_configuration.0.commit_options_check"},
									},
									"commit_options_synchronize": {
										Type:          schema.TypeBool,
										Optional:      true,
										ConflictsWith: []string{"then.0.change_configuration.0.commit_options_check"},
									},
									"retry_count": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										RequiredWith: []string{"then.0.change_configuration.0.retry_interval"},
										ValidateFunc: validation.IntBetween(0, 10),
									},
									"retry_interval": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										RequiredWith: []string{"then.0.change_configuration.0.retry_count"},
										ValidateFunc: validation.IntBetween(0, 4294967295),
									},
									"user_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"event_script": {
							Type:     schema.TypeList,
							Optional: true,
							AtLeastOneOf: []string{
								"then.0.change_configuration",
								"then.0.event_script",
								"then.0.execute_commands",
								"then.0.ignore",
								"then.0.priority_override_facility",
								"then.0.priority_override_severity",
								"then.0.raise_trap",
								"then.0.upload",
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filename": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringDoesNotContainAny(" "),
									},
									"arguments": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringDoesNotContainAny(" "),
												},
												"value": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringDoesNotContainAny(" "),
												},
											},
										},
									},
									"destination": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringDoesNotContainAny(" "),
												},
												"retry_count": {
													Type:         schema.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 10),
												},
												"retry_interval": {
													Type:         schema.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 4294967295),
												},
												"transfer_delay": {
													Type:         schema.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 4294967295),
												},
											},
										},
									},
									"output_filename": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"output_format": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"text", "xml"}, false),
									},
									"user_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"execute_commands": {
							Type:     schema.TypeList,
							Optional: true,
							AtLeastOneOf: []string{
								"then.0.change_configuration",
								"then.0.event_script",
								"then.0.execute_commands",
								"then.0.ignore",
								"then.0.priority_override_facility",
								"then.0.priority_override_severity",
								"then.0.raise_trap",
								"then.0.upload",
							},
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"commands": {
										Type:     schema.TypeList,
										Required: true,
										MinItems: 1,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"destination": {
										Type:         schema.TypeList,
										Optional:     true,
										RequiredWith: []string{"then.0.execute_commands.0.output_filename"},
										MaxItems:     1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringDoesNotContainAny(" "),
												},
												"retry_count": {
													Type:         schema.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 10),
												},
												"retry_interval": {
													Type:         schema.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 4294967295),
												},
												"transfer_delay": {
													Type:         schema.TypeInt,
													Optional:     true,
													Default:      -1,
													ValidateFunc: validation.IntBetween(0, 4294967295),
												},
											},
										},
									},
									"output_filename": {
										Type:         schema.TypeString,
										Optional:     true,
										RequiredWith: []string{"then.0.execute_commands.0.destination"},
									},
									"output_format": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"text", "xml"}, false),
									},
									"user_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"ignore": {
							Type:     schema.TypeBool,
							Optional: true,
							ConflictsWith: []string{
								"then.0.change_configuration",
								"then.0.event_script",
								"then.0.execute_commands",
								"then.0.priority_override_facility",
								"then.0.priority_override_severity",
								"then.0.raise_trap",
								"then.0.upload",
							},
						},
						"priority_override_facility": {
							Type:     schema.TypeString,
							Optional: true,
							AtLeastOneOf: []string{
								"then.0.change_configuration",
								"then.0.event_script",
								"then.0.execute_commands",
								"then.0.ignore",
								"then.0.priority_override_facility",
								"then.0.priority_override_severity",
								"then.0.raise_trap",
								"then.0.upload",
							},
							ValidateFunc: validation.StringInSlice([]string{
								"authorization",
								"change-log",
								"conflict-log",
								"daemon",
								"dfc",
								"external",
								"firewall",
								"ftp",
								"interactive-commands",
								"kernel",
								"ntp",
								"pfe",
								"security",
								"user",
							}, false),
						},
						"priority_override_severity": {
							Type:     schema.TypeString,
							Optional: true,
							AtLeastOneOf: []string{
								"then.0.change_configuration",
								"then.0.event_script",
								"then.0.execute_commands",
								"then.0.ignore",
								"then.0.priority_override_facility",
								"then.0.priority_override_severity",
								"then.0.raise_trap",
								"then.0.upload",
							},
							ValidateFunc: validation.StringInSlice([]string{
								"alert",
								"critical",
								"emergency",
								"error",
								"info",
								"notice",
								"warning",
							}, false),
						},
						"raise_trap": {
							Type:     schema.TypeBool,
							Optional: true,
							AtLeastOneOf: []string{
								"then.0.change_configuration",
								"then.0.event_script",
								"then.0.execute_commands",
								"then.0.ignore",
								"then.0.priority_override_facility",
								"then.0.priority_override_severity",
								"then.0.raise_trap",
								"then.0.upload",
							},
						},
						"upload": {
							Type:     schema.TypeList,
							Optional: true,
							AtLeastOneOf: []string{
								"then.0.change_configuration",
								"then.0.event_script",
								"then.0.execute_commands",
								"then.0.ignore",
								"then.0.priority_override_facility",
								"then.0.priority_override_severity",
								"then.0.raise_trap",
								"then.0.upload",
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filename": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringDoesNotContainAny(" "),
									},
									"destination": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringDoesNotContainAny(" "),
									},
									"retry_count": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 10),
									},
									"retry_interval": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 4294967295),
									},
									"transfer_delay": {
										Type:         schema.TypeInt,
										Optional:     true,
										Default:      -1,
										ValidateFunc: validation.IntBetween(0, 4294967295),
									},
									"user_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"attributes_match": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringDoesNotContainAny(" "),
						},
						"compare": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"equals", "matches", "starts-with"}, false),
						},
						"to": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringDoesNotContainAny(" "),
						},
					},
				},
			},
			"within": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time_interval": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 604800),
						},
						"events": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"not_events": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"trigger_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      -1,
							ValidateFunc: validation.IntBetween(0, 4294967295),
						},
						"trigger_when": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"after", "on", "until"}, false),
						},
					},
				},
			},
		},
	}
}

func resourceEventoptionsPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeCreateSetFile() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := setEventoptionsPolicy(d, junSess); err != nil {
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
	if err := junSess.ConfigLock(ctx); err != nil {
		return diag.FromErr(err)
	}
	var diagWarns diag.Diagnostics
	eventoptionsPolicyExists, err := checkEventoptionsPolicyExists(d.Get("name").(string), junSess)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if eventoptionsPolicyExists {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns,
			diag.FromErr(fmt.Errorf("event-options policy %v already exists", d.Get("name").(string)))...)
	}

	if err := setEventoptionsPolicy(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("create resource junos_eventoptions_policy")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	eventoptionsPolicyExists, err = checkEventoptionsPolicyExists(d.Get("name").(string), junSess)
	if err != nil {
		return append(diagWarns, diag.FromErr(err)...)
	}
	if eventoptionsPolicyExists {
		d.SetId(d.Get("name").(string))
	} else {
		return append(diagWarns, diag.FromErr(fmt.Errorf("event-options policy %v not exists after commit "+
			"=> check your config", d.Get("name").(string)))...)
	}

	return append(diagWarns, resourceEventoptionsPolicyReadWJunSess(d, junSess)...)
}

func resourceEventoptionsPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer junSess.Close()

	return resourceEventoptionsPolicyReadWJunSess(d, junSess)
}

func resourceEventoptionsPolicyReadWJunSess(d *schema.ResourceData, junSess *junos.Session,
) diag.Diagnostics {
	junos.MutexLock()
	eventoptionsPolicyOptions, err := readEventoptionsPolicy(d.Get("name").(string), junSess)
	junos.MutexUnlock()
	if err != nil {
		return diag.FromErr(err)
	}
	if eventoptionsPolicyOptions.name == "" {
		d.SetId("")
	} else {
		fillEventoptionsPolicyData(d, eventoptionsPolicyOptions)
	}

	return nil
}

func resourceEventoptionsPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.Partial(true)
	clt := m.(*junos.Client)
	if clt.FakeUpdateAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delEventoptionsPolicy(d.Get("name").(string), junSess); err != nil {
			return diag.FromErr(err)
		}
		if err := setEventoptionsPolicy(d, junSess); err != nil {
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
	if err := delEventoptionsPolicy(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	if err := setEventoptionsPolicy(d, junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("update resource junos_eventoptions_policy")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	d.Partial(false)

	return append(diagWarns, resourceEventoptionsPolicyReadWJunSess(d, junSess)...)
}

func resourceEventoptionsPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clt := m.(*junos.Client)
	if clt.FakeDeleteAlso() {
		junSess := clt.NewSessionWithoutNetconf(ctx)
		if err := delEventoptionsPolicy(d.Get("name").(string), junSess); err != nil {
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
	if err := delEventoptionsPolicy(d.Get("name").(string), junSess); err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}
	warns, err := junSess.CommitConf("delete resource junos_eventoptions_policy")
	appendDiagWarns(&diagWarns, warns)
	if err != nil {
		appendDiagWarns(&diagWarns, junSess.ConfigClear())

		return append(diagWarns, diag.FromErr(err)...)
	}

	return diagWarns
}

func resourceEventoptionsPolicyImport(ctx context.Context, d *schema.ResourceData, m interface{},
) ([]*schema.ResourceData, error) {
	clt := m.(*junos.Client)
	junSess, err := clt.StartNewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer junSess.Close()
	result := make([]*schema.ResourceData, 1)

	eventoptionsPolicyExists, err := checkEventoptionsPolicyExists(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	if !eventoptionsPolicyExists {
		return nil, fmt.Errorf("don't find event-options policy with id '%v' (id must be <name>)", d.Id())
	}
	eventoptionsPolicyOptions, err := readEventoptionsPolicy(d.Id(), junSess)
	if err != nil {
		return nil, err
	}
	fillEventoptionsPolicyData(d, eventoptionsPolicyOptions)

	result[0] = d

	return result, nil
}

func checkEventoptionsPolicyExists(name string, junSess *junos.Session) (bool, error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"event-options policy \"" + name + "\"" + junos.PipeDisplaySet)
	if err != nil {
		return false, err
	}
	if showConfig == junos.EmptyW {
		return false, nil
	}

	return true, nil
}

func setEventoptionsPolicy(d *schema.ResourceData, junSess *junos.Session) error {
	configSet := make([]string, 0)
	setPrefix := "set event-options policy \"" + d.Get("name").(string) + "\" "

	for _, v := range sortSetOfString(d.Get("events").(*schema.Set).List()) {
		configSet = append(configSet, setPrefix+"events \""+v+"\"")
	}
	for _, v := range d.Get("then").([]interface{}) {
		then := v.(map[string]interface{})
		for _, v2 := range then["change_configuration"].([]interface{}) {
			changeConfig := v2.(map[string]interface{})
			for _, command := range changeConfig["commands"].([]interface{}) {
				configSet = append(configSet, setPrefix+"then change-configuration commands \""+command.(string)+"\"")
			}
			if changeConfig["commit_options_check"].(bool) {
				configSet = append(configSet, setPrefix+"then change-configuration commit-options check")
				if changeConfig["commit_options_check_synchronize"].(bool) {
					configSet = append(configSet, setPrefix+"then change-configuration commit-options check synchronize")
				}
			} else if changeConfig["commit_options_check_synchronize"].(bool) {
				return fmt.Errorf("commit_options_check must be set to true if commit_options_check_synchronize is set to true")
			}
			if changeConfig["commit_options_force"].(bool) {
				configSet = append(configSet, setPrefix+"then change-configuration commit-options force")
			}
			if v3 := changeConfig["commit_options_log"].(string); v3 != "" {
				configSet = append(configSet, setPrefix+"then change-configuration commit-options log \""+v3+"\"")
			}
			if changeConfig["commit_options_synchronize"].(bool) {
				configSet = append(configSet, setPrefix+"then change-configuration commit-options synchronize")
			}
			if v3 := changeConfig["retry_count"].(int); v3 != -1 {
				configSet = append(configSet, setPrefix+"then change-configuration retry count "+
					strconv.Itoa(v3)+" interval "+strconv.Itoa(changeConfig["retry_interval"].(int)))
			}
			if v3 := changeConfig["user_name"].(string); v3 != "" {
				configSet = append(configSet, setPrefix+"then change-configuration user-name "+v3)
			}
		}
		eventScriptFilenameList := make([]string, 0)
		for _, v2 := range then["event_script"].([]interface{}) {
			eventScript := v2.(map[string]interface{})
			if bchk.InSlice(eventScript["filename"].(string), eventScriptFilenameList) {
				return fmt.Errorf("multiple blocks event_script with the same filename %s", eventScript["filename"].(string))
			}
			eventScriptFilenameList = append(eventScriptFilenameList, eventScript["filename"].(string))
			setPrefixThenEventScript := setPrefix + "then event-script \"" + eventScript["filename"].(string) + "\" "
			configSet = append(configSet, setPrefixThenEventScript)
			argumentsNameList := make([]string, 0)
			for _, v3 := range eventScript["arguments"].([]interface{}) {
				arguments := v3.(map[string]interface{})
				if bchk.InSlice(arguments["name"].(string), argumentsNameList) {
					return fmt.Errorf("multiple blocks arguments with the same name %s", arguments["name"].(string))
				}
				argumentsNameList = append(argumentsNameList, arguments["name"].(string))
				configSet = append(configSet, setPrefixThenEventScript+
					"arguments \""+arguments["name"].(string)+"\" \""+arguments["value"].(string)+"\"")
			}
			for _, v3 := range eventScript["destination"].([]interface{}) {
				destination := v3.(map[string]interface{})
				setPrefixDestination := setPrefixThenEventScript + "destination \"" + destination["name"].(string) + "\" "
				configSet = append(configSet, setPrefixDestination)
				if retryCount := destination["retry_count"].(int); retryCount != -1 {
					if retryInterval := destination["retry_interval"].(int); retryInterval != -1 {
						configSet = append(configSet, setPrefixDestination+
							"retry-count "+strconv.Itoa(retryCount)+" retry-interval "+strconv.Itoa(retryInterval))
					} else {
						return fmt.Errorf("retry_interval must be set with retry_count")
					}
				} else if destination["retry_interval"].(int) != -1 {
					return fmt.Errorf("retry_count must be set with retry_interval")
				}
				if transferDelay := destination["transfer_delay"].(int); transferDelay != -1 {
					configSet = append(configSet, setPrefixDestination+"transfer-delay "+strconv.Itoa(transferDelay))
				}
			}
			if v3 := eventScript["output_filename"].(string); v3 != "" {
				configSet = append(configSet, setPrefixThenEventScript+"output-filename \""+v3+"\"")
			}
			if v3 := eventScript["output_format"].(string); v3 != "" {
				configSet = append(configSet, setPrefixThenEventScript+"output-format "+v3)
			}
			if v3 := eventScript["user_name"].(string); v3 != "" {
				configSet = append(configSet, setPrefixThenEventScript+"user-name "+v3)
			}
		}
		for _, v2 := range then["execute_commands"].([]interface{}) {
			executeCommands := v2.(map[string]interface{})
			for _, command := range executeCommands["commands"].([]interface{}) {
				configSet = append(configSet, setPrefix+"then execute-commands commands \""+command.(string)+"\"")
			}
			for _, v3 := range executeCommands["destination"].([]interface{}) {
				destination := v3.(map[string]interface{})
				setPrefixDestination := setPrefix + "then execute-commands destination \"" + destination["name"].(string) + "\" "
				configSet = append(configSet, setPrefixDestination)
				if retryCount := destination["retry_count"].(int); retryCount != -1 {
					if retryInterval := destination["retry_interval"].(int); retryInterval != -1 {
						configSet = append(configSet, setPrefixDestination+
							"retry-count "+strconv.Itoa(retryCount)+" retry-interval "+strconv.Itoa(retryInterval))
					} else {
						return fmt.Errorf("retry_interval must be set with retry_count")
					}
				} else if destination["retry_interval"].(int) != -1 {
					return fmt.Errorf("retry_count must be set with retry_interval")
				}
				if transferDelay := destination["transfer_delay"].(int); transferDelay != -1 {
					configSet = append(configSet, setPrefixDestination+"transfer-delay "+strconv.Itoa(transferDelay))
				}
			}
			if v3 := executeCommands["output_filename"].(string); v3 != "" {
				configSet = append(configSet, setPrefix+"then execute-commands output-filename \""+v3+"\"")
			}
			if v3 := executeCommands["output_format"].(string); v3 != "" {
				configSet = append(configSet, setPrefix+"then execute-commands output-format "+v3)
			}
			if v3 := executeCommands["user_name"].(string); v3 != "" {
				configSet = append(configSet, setPrefix+"then execute-commands user-name "+v3)
			}
		}
		if then["ignore"].(bool) {
			configSet = append(configSet, setPrefix+"then ignore")
		}
		if v2 := then["priority_override_facility"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+"then priority-override facility "+v2)
		}
		if v2 := then["priority_override_severity"].(string); v2 != "" {
			configSet = append(configSet, setPrefix+"then priority-override severity "+v2)
		}
		if then["raise_trap"].(bool) {
			configSet = append(configSet, setPrefix+"then raise-trap")
		}
		uploadFileDestList := make([]string, 0)
		for _, v2 := range then["upload"].([]interface{}) {
			upload := v2.(map[string]interface{})
			setPrefixThenUpload := setPrefix + "then upload filename \"" + upload["filename"].(string) + "\" " +
				"destination \"" + upload["destination"].(string) + "\" "
			if bchk.InSlice(setPrefixThenUpload, uploadFileDestList) {
				return fmt.Errorf("multiple blocks upload with the same filename %s and destination %s",
					upload["filename"].(string), upload["destination"].(string))
			}
			uploadFileDestList = append(uploadFileDestList, setPrefixThenUpload)
			configSet = append(configSet, setPrefixThenUpload)
			if retryCount := upload["retry_count"].(int); retryCount != -1 {
				if retryInterval := upload["retry_interval"].(int); retryInterval != -1 {
					configSet = append(configSet, setPrefixThenUpload+
						"retry-count "+strconv.Itoa(retryCount)+" retry-interval "+strconv.Itoa(retryInterval))
				} else {
					return fmt.Errorf("retry_interval must be set with retry_count")
				}
			} else if upload["retry_interval"].(int) != -1 {
				return fmt.Errorf("retry_count must be set with retry_interval")
			}
			if transferDelay := upload["transfer_delay"].(int); transferDelay != -1 {
				configSet = append(configSet, setPrefixThenUpload+"transfer-delay "+strconv.Itoa(transferDelay))
			}
			if v3 := upload["user_name"].(string); v3 != "" {
				configSet = append(configSet, setPrefixThenUpload+"user-name "+v3)
			}
		}
	}
	attriMatchList := make([]string, 0)
	for _, v := range d.Get("attributes_match").([]interface{}) {
		attriMatch := v.(map[string]interface{})
		setAttriMatch := setPrefix + "attributes-match \"" + attriMatch["from"].(string) + "\" " +
			attriMatch["compare"].(string) + " \"" + attriMatch["to"].(string) + "\""
		if bchk.InSlice(setAttriMatch, attriMatchList) {
			return fmt.Errorf("multiple blocks attributes_match with the same from %s, compare %s and to %s",
				attriMatch["from"].(string), attriMatch["compare"].(string), attriMatch["to"].(string))
		}
		attriMatchList = append(attriMatchList, setAttriMatch)
		configSet = append(configSet, setAttriMatch)
	}
	withinTimeInterval := make([]int, 0)
	for _, v := range d.Get("within").([]interface{}) {
		within := v.(map[string]interface{})
		if bchk.InSlice(within["time_interval"].(int), withinTimeInterval) {
			return fmt.Errorf("multiple blocks within with the same time_interval %d", within["time_interval"].(int))
		}
		withinTimeInterval = append(withinTimeInterval, within["time_interval"].(int))
		setPrefixWithin := setPrefix + "within " + strconv.Itoa(within["time_interval"].(int)) + " "
		for _, v2 := range sortSetOfString(within["events"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixWithin+"events \""+v2+"\"")
		}
		for _, v2 := range sortSetOfString(within["not_events"].(*schema.Set).List()) {
			configSet = append(configSet, setPrefixWithin+"not events \""+v2+"\"")
		}
		if v2 := within["trigger_when"].(string); v2 != "" {
			if c := within["trigger_count"].(int); c != -1 {
				configSet = append(configSet, setPrefixWithin+"trigger "+v2+" "+strconv.Itoa(c))
			} else {
				return fmt.Errorf("trigger_count must be set with trigger_when")
			}
		} else if within["trigger_count"].(int) != -1 {
			return fmt.Errorf("trigger_when must be set with trigger_count")
		}
		if len(configSet) == 0 || !strings.HasPrefix(configSet[len(configSet)-1], setPrefixWithin) {
			return fmt.Errorf("missing argument for within (time_interval=%d)", within["time_interval"].(int))
		}
	}

	return junSess.ConfigSet(configSet)
}

func readEventoptionsPolicy(name string, junSess *junos.Session,
) (confRead eventoptionsPolicyOptions, err error) {
	showConfig, err := junSess.Command(junos.CmdShowConfig +
		"event-options policy \"" + name + "\"" + junos.PipeDisplaySetRelative)
	if err != nil {
		return confRead, err
	}
	if showConfig != junos.EmptyW {
		confRead.name = name
		for _, item := range strings.Split(showConfig, "\n") {
			if strings.Contains(item, junos.XMLStartTagConfigOut) {
				continue
			}
			if strings.Contains(item, junos.XMLEndTagConfigOut) {
				break
			}
			itemTrim := strings.TrimPrefix(item, junos.SetLS)
			switch {
			case balt.CutPrefixInString(&itemTrim, "events "):
				confRead.events = append(confRead.events, strings.Trim(itemTrim, "\""))
			case balt.CutPrefixInString(&itemTrim, "then "):
				if len(confRead.then) == 0 {
					confRead.then = append(confRead.then, map[string]interface{}{
						"change_configuration":       make([]map[string]interface{}, 0),
						"event_script":               make([]map[string]interface{}, 0),
						"execute_commands":           make([]map[string]interface{}, 0),
						"ignore":                     false,
						"priority_override_facility": "",
						"priority_override_severity": "",
						"raise_trap":                 false,
						"upload":                     make([]map[string]interface{}, 0),
					})
				}
				if err := readEventoptionsPolicyThen(itemTrim, confRead.then[0]); err != nil {
					return confRead, err
				}
			case balt.CutPrefixInString(&itemTrim, "attributes-match "):
				itemTrimFields := strings.Split(itemTrim, " ")
				if len(itemTrimFields) < 3 { // <from> <compare> <to>
					return confRead, fmt.Errorf(junos.CantReadValuesNotEnoughFields, "attributes-match", itemTrim)
				}
				confRead.attributesMatch = append(confRead.attributesMatch, map[string]interface{}{
					"from":    strings.Trim(itemTrimFields[0], "\""),
					"compare": itemTrimFields[1],
					"to":      strings.Trim(itemTrimFields[2], "\""),
				})
			case balt.CutPrefixInString(&itemTrim, "within "):
				itemTrimFields := strings.Split(itemTrim, " ")
				withinSeconds, err := strconv.Atoi(itemTrimFields[0])
				if err != nil {
					return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
				within := map[string]interface{}{
					"time_interval": withinSeconds,
					"events":        make([]string, 0),
					"not_events":    make([]string, 0),
					"trigger_count": -1,
					"trigger_when":  "",
				}
				balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
				confRead.within = copyAndRemoveItemMapList("time_interval", within, confRead.within)
				switch {
				case balt.CutPrefixInString(&itemTrim, "events "):
					within["events"] = append(within["events"].([]string), strings.Trim(itemTrim, "\""))
				case balt.CutPrefixInString(&itemTrim, "not events "):
					within["not_events"] = append(within["not_events"].([]string), strings.Trim(itemTrim, "\""))
				case balt.CutPrefixInString(&itemTrim, "trigger "):
					switch itemTrim {
					case "after", "on", "until":
						within["trigger_when"] = itemTrim
					default:
						within["trigger_count"], err = strconv.Atoi(itemTrim)
						if err != nil {
							return confRead, fmt.Errorf(failedConvAtoiError, itemTrim, err)
						}
					}
				}
				confRead.within = append(confRead.within, within)
			}
		}
	}

	return confRead, nil
}

func readEventoptionsPolicyThen(itemTrim string, then map[string]interface{}) (err error) {
	switch {
	case balt.CutPrefixInString(&itemTrim, "change-configuration "):
		if len(then["change_configuration"].([]map[string]interface{})) == 0 {
			then["change_configuration"] = append(
				then["change_configuration"].([]map[string]interface{}), map[string]interface{}{
					"commands":                         make([]string, 0),
					"commit_options_check":             false,
					"commit_options_check_synchronize": false,
					"commit_options_force":             false,
					"commit_options_log":               "",
					"commit_options_synchronize":       false,
					"retry_count":                      -1,
					"retry_interval":                   -1,
					"user_name":                        "",
				})
		}
		changeConfiguration := then["change_configuration"].([]map[string]interface{})[0]
		switch {
		case balt.CutPrefixInString(&itemTrim, "commands "):
			changeConfiguration["commands"] = append(changeConfiguration["commands"].([]string), strings.Trim(itemTrim, "\""))
		case itemTrim == "commit-options check":
			changeConfiguration["commit_options_check"] = true
		case itemTrim == "commit-options check synchronize":
			changeConfiguration["commit_options_check"] = true
			changeConfiguration["commit_options_check_synchronize"] = true
		case itemTrim == "commit-options force":
			changeConfiguration["commit_options_force"] = true
		case balt.CutPrefixInString(&itemTrim, "commit-options log "):
			changeConfiguration["commit_options_log"] = strings.Trim(itemTrim, "\"")
		case itemTrim == "commit-options synchronize":
			changeConfiguration["commit_options_synchronize"] = true
		case balt.CutPrefixInString(&itemTrim, "retry count "):
			changeConfiguration["retry_count"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "retry interval "):
			changeConfiguration["retry_interval"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "user-name "):
			changeConfiguration["user_name"] = itemTrim
		}
	case balt.CutPrefixInString(&itemTrim, "event-script "):
		itemTrimFields := strings.Split(itemTrim, " ")
		eventScript := map[string]interface{}{
			"filename":        strings.Trim(itemTrimFields[0], "\""),
			"arguments":       make([]map[string]interface{}, 0),
			"destination":     make([]map[string]interface{}, 0),
			"output_filename": "",
			"output_format":   "",
			"user_name":       "",
		}
		then["event_script"] = copyAndRemoveItemMapList(
			"filename", eventScript, then["event_script"].([]map[string]interface{}))
		balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
		switch {
		case balt.CutPrefixInString(&itemTrim, "arguments "):
			itemTrimArgsFields := strings.Split(itemTrim, " ")
			if len(itemTrimArgsFields) < 2 { // <name> <value>
				return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "arguments", itemTrim)
			}
			eventScript["arguments"] = append(eventScript["arguments"].([]map[string]interface{}), map[string]interface{}{
				"name":  strings.Trim(itemTrimArgsFields[0], "\""),
				"value": strings.Trim(itemTrimArgsFields[1], "\""),
			})
		case balt.CutPrefixInString(&itemTrim, "destination "):
			itemTrimDestFields := strings.Split(itemTrim, " ")
			if len(eventScript["destination"].([]map[string]interface{})) == 0 {
				eventScript["destination"] = append(eventScript["destination"].([]map[string]interface{}), map[string]interface{}{
					"name":           strings.Trim(itemTrimDestFields[0], "\""),
					"retry_count":    -1,
					"retry_interval": -1,
					"transfer_delay": -1,
				})
			}
			destination := eventScript["destination"].([]map[string]interface{})[0]
			balt.CutPrefixInString(&itemTrim, itemTrimDestFields[0]+" ")
			switch {
			case balt.CutPrefixInString(&itemTrim, "retry-count retry-interval "):
				destination["retry_interval"], err = strconv.Atoi(itemTrim)
				if err != nil {
					return fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "retry-count "):
				destination["retry_count"], err = strconv.Atoi(itemTrim)
				if err != nil {
					return fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "transfer-delay "):
				destination["transfer_delay"], err = strconv.Atoi(itemTrim)
				if err != nil {
					return fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			}
		case balt.CutPrefixInString(&itemTrim, "output-filename "):
			eventScript["output_filename"] = strings.Trim(itemTrim, "\"")
		case balt.CutPrefixInString(&itemTrim, "output-format "):
			eventScript["output_format"] = itemTrim
		case balt.CutPrefixInString(&itemTrim, "user-name "):
			eventScript["user_name"] = itemTrim
		}
		then["event_script"] = append(then["event_script"].([]map[string]interface{}), eventScript)
	case balt.CutPrefixInString(&itemTrim, "execute-commands "):
		if len(then["execute_commands"].([]map[string]interface{})) == 0 {
			then["execute_commands"] = append(
				then["execute_commands"].([]map[string]interface{}), map[string]interface{}{
					"commands":        make([]string, 0),
					"destination":     make([]map[string]interface{}, 0),
					"output_filename": "",
					"output_format":   "",
					"user_name":       "",
				})
		}
		executeCommands := then["execute_commands"].([]map[string]interface{})[0]
		switch {
		case balt.CutPrefixInString(&itemTrim, "commands "):
			executeCommands["commands"] = append(executeCommands["commands"].([]string), strings.Trim(itemTrim, "\""))
		case balt.CutPrefixInString(&itemTrim, "destination "):
			itemTrimFields := strings.Split(itemTrim, " ")
			if len(executeCommands["destination"].([]map[string]interface{})) == 0 {
				executeCommands["destination"] = append(
					executeCommands["destination"].([]map[string]interface{}), map[string]interface{}{
						"name":           strings.Trim(itemTrimFields[0], "\""),
						"retry_count":    -1,
						"retry_interval": -1,
						"transfer_delay": -1,
					})
			}
			destination := executeCommands["destination"].([]map[string]interface{})[0]
			balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" ")
			switch {
			case balt.CutPrefixInString(&itemTrim, "retry-count retry-interval "):
				destination["retry_interval"], err = strconv.Atoi(itemTrim)
				if err != nil {
					return fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "retry-count "):
				destination["retry_count"], err = strconv.Atoi(itemTrim)
				if err != nil {
					return fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			case balt.CutPrefixInString(&itemTrim, "transfer-delay "):
				destination["transfer_delay"], err = strconv.Atoi(itemTrim)
				if err != nil {
					return fmt.Errorf(failedConvAtoiError, itemTrim, err)
				}
			}
		case balt.CutPrefixInString(&itemTrim, "output-filename "):
			executeCommands["output_filename"] = strings.Trim(itemTrim, "\"")
		case balt.CutPrefixInString(&itemTrim, "output-format "):
			executeCommands["output_format"] = itemTrim
		case balt.CutPrefixInString(&itemTrim, "user-name "):
			executeCommands["user_name"] = itemTrim
		}
	case itemTrim == "ignore":
		then["ignore"] = true
	case balt.CutPrefixInString(&itemTrim, "priority-override facility "):
		then["priority_override_facility"] = itemTrim
	case balt.CutPrefixInString(&itemTrim, "priority-override severity "):
		then["priority_override_severity"] = itemTrim
	case itemTrim == "raise-trap":
		then["raise_trap"] = true
	case balt.CutPrefixInString(&itemTrim, "upload filename "):
		itemTrimFields := strings.Split(itemTrim, " ")
		if len(itemTrimFields) < 3 { // <filename> destination <destination>
			return fmt.Errorf(junos.CantReadValuesNotEnoughFields, "upload filename", itemTrim)
		}
		upload := map[string]interface{}{
			"filename":       strings.Trim(itemTrimFields[0], "\""),
			"destination":    strings.Trim(itemTrimFields[2], "\""),
			"retry_count":    -1,
			"retry_interval": -1,
			"transfer_delay": -1,
			"user_name":      "",
		}

		then["upload"] = copyAndRemoveItemMapList2(
			"filename", "destination", upload, then["upload"].([]map[string]interface{}))
		balt.CutPrefixInString(&itemTrim, itemTrimFields[0]+" destination "+itemTrimFields[2]+" ")
		switch {
		case balt.CutPrefixInString(&itemTrim, "retry-count retry-interval "):
			upload["retry_interval"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "retry-count "):
			upload["retry_count"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "transfer-delay "):
			upload["transfer_delay"], err = strconv.Atoi(itemTrim)
			if err != nil {
				return fmt.Errorf(failedConvAtoiError, itemTrim, err)
			}
		case balt.CutPrefixInString(&itemTrim, "user-name "):
			upload["user_name"] = itemTrim
		}
		then["upload"] = append(then["upload"].([]map[string]interface{}), upload)
	}

	return nil
}

func delEventoptionsPolicy(policy string, junSess *junos.Session) error {
	configSet := make([]string, 0, 1)
	configSet = append(configSet, "delete event-options policy \""+policy+"\"")

	return junSess.ConfigSet(configSet)
}

func fillEventoptionsPolicyData(d *schema.ResourceData, eventoptionsPolicyOptions eventoptionsPolicyOptions) {
	if tfErr := d.Set("name", eventoptionsPolicyOptions.name); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("events", eventoptionsPolicyOptions.events); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("then", eventoptionsPolicyOptions.then); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("attributes_match", eventoptionsPolicyOptions.attributesMatch); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("within", eventoptionsPolicyOptions.within); tfErr != nil {
		panic(tfErr)
	}
}
