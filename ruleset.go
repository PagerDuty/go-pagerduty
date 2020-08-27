package pagerduty

import (
	"fmt"
	"net/http"
)

// Ruleset represents a ruleset.
type Ruleset struct {
	ID          string         `json:"id,omitempty"`
	Name        string         `json:"name,omitempty"`
	Type        string         `json:"type,omitempty"`
	RoutingKeys []string       `json:"routing_keys,omitempty"`
	Team        *RulesetObject `json:"team,omitempty"`
	Updater     *RulesetObject `json:"updater,omitempty"`
	Creator     *RulesetObject `json:"creator,omitempty"`
}

// RulesetObject represents a generic object that is common within a ruleset object
type RulesetObject struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
}

// RulesetPayload represents payload with a ruleset object
type RulesetPayload struct {
	Ruleset *Ruleset `json:"ruleset,omitempty"`
}

// ListRulesetsResponse represents a list response of rulesets.
type ListRulesetsResponse struct {
	Total    uint       `json:"total,omitempty"`
	Rulesets []*Ruleset `json:"rulesets,omitempty"`
	Offset   uint       `json:"offset,omitempty"`
	More     bool       `json:"more,omitempty"`
	Limit    uint       `json:"limit,omitempty"`
}

// RulesetRule represents a Ruleset rule
type RulesetRule struct {
	ID         string          `json:"id,omitempty"`
	Position   *int            `json:"position,omitempty"`
	Disabled   bool            `json:"disabled,omitempty"`
	Conditions *RuleConditions `json:"conditions,omitempty"`
	Actions    *RuleActions    `json:"actions,omitempty"`
	Ruleset    *APIObject      `json:"ruleset,omitempty"`
	Self       string          `json:"self,omitempty"`
	CatchAll   bool            `json:"catch_all,omitempty"`
	TimeFrame  *RuleTimeFrame  `json:"time_frame,omitempty"`
}

// RulesetRulePayload represents a payload for ruleset rules
type RulesetRulePayload struct {
	Rule *RulesetRule `json:"rule,omitempty"`
}

// RuleConditions represents the conditions field for a Ruleset
type RuleConditions struct {
	Operator          string              `json:"operator,omitempty"`
	RuleSubconditions []*RuleSubcondition `json:"subconditions,omitempty"`
}

// RuleSubcondition represents a subcondition of a ruleset condition
type RuleSubcondition struct {
	Operator   string              `json:"operator,omitempty"`
	Parameters *ConditionParameter `json:"parameters,omitempty"`
}

// ConditionParameter represents  parameters in a rule condition
type ConditionParameter struct {
	Path  string `json:"path,omitempty"`
	Value string `json:"value,omitempty"`
}

// RuleTimeFrame represents a time_frame object on the rule object
type RuleTimeFrame struct {
	ScheduledWeekly *ScheduledWeekly `json:"scheduled_weekly,omitempty"`
	ActiveBetween   *ActiveBetween   `json:"active_between,omitempty"`
}

// ScheduledWeekly represents a time_frame object for scheduling rules weekly
type ScheduledWeekly struct {
	Weekdays  []int  `json:"weekdays,omitempty"`
	Timezone  string `json:"timezone,omitempty"`
	StartTime int    `json:"start_time,omitempty"`
	Duration  int    `json:"duration,omitempty"`
}

// ActiveBetween represents an active_between object for setting a timeline for rules
type ActiveBetween struct {
	StartTime int `json:"start_time,omitempty"`
	EndTime   int `json:"end_time,omitempty"`
}

// ListRulesetRulesResponse represents a list of rules in a ruleset
type ListRulesetRulesResponse struct {
	Total  uint           `json:"total,omitempty"`
	Rules  []*RulesetRule `json:"rules,omitempty"`
	Offset uint           `json:"offset,omitempty"`
	More   bool           `json:"more,omitempty"`
	Limit  uint           `json:"limit,omitempty"`
}

// RuleActions represents a rule action
type RuleActions struct {
	Suppress    *RuleActionSuppress     `json:"suppress,omitempty"`
	Annotate    *RuleActionParameter    `json:"annotate,omitempty"`
	Severity    *RuleActionParameter    `json:"severity,omitempty"`
	Priority    *RuleActionParameter    `json:"priority,omitempty"`
	Route       *RuleActionParameter    `json:"route,omitempty"`
	EventAction *RuleActionParameter    `json:"event_action,omitempty"`
	Extractions []*RuleActionExtraction `json:"extractions,omitempty"`
}

// RuleActionParameter represents a generic parameter object on a rule action
type RuleActionParameter struct {
	Value string `json:"value,omitempty"`
}

// RuleActionSuppress represents a rule suppress action object
type RuleActionSuppress struct {
	Value               bool   `json:"value,omitempty"`
	ThresholdValue      int    `json:"threshold_value,omitempty"`
	ThresholdTimeUnit   string `json:"threshold_time_unit,omitempty"`
	ThresholdTimeAmount int    `json:"threshold_time_amount,omitempty"`
}

// RuleActionExtraction represents a rule extraction action object
type RuleActionExtraction struct {
	Target string `json:"target,omitempty"`
	Source string `json:"source,omitempty"`
	Regex  string `json:"regex,omitempty"`
}

// ListRulesets gets all rulesets.
func (c *Client) ListRulesets() (*ListRulesetsResponse, error) {
	rulesetResponse := new(ListRulesetsResponse)
	rulesets := make([]*Ruleset, 0)

	// Create a handler closure capable of parsing data from the rulesets endpoint
	// and appending resultant rulesets to the return slice.
	responseHandler := func(response *http.Response) (APIListObject, error) {
		var result ListRulesetsResponse
		if err := c.decodeJSON(response, &result); err != nil {
			return APIListObject{}, err
		}

		rulesets = append(rulesets, result.Rulesets...)

		// Return stats on the current page. Caller can use this information to
		// adjust for requesting additional pages.
		return APIListObject{
			More:   result.More,
			Offset: result.Offset,
			Limit:  result.Limit,
		}, nil
	}

	// Make call to get all pages associated with the base endpoint.
	if err := c.pagedGet("/rulesets/", responseHandler); err != nil {
		return nil, err
	}
	rulesetResponse.Rulesets = rulesets

	return rulesetResponse, nil
}

// CreateRuleset creates a new user.
func (c *Client) CreateRuleset(r *Ruleset) (*Ruleset, *http.Response, error) {
	data := make(map[string]*Ruleset)
	data["ruleset"] = r
	resp, err := c.post("/rulesets", data, nil)
	return getRulesetFromResponse(c, resp, err)
}

// DeleteRuleset deletes a ruleset.
func (c *Client) DeleteRuleset(id string) error {
	_, err := c.delete("/rulesets/" + id)
	return err
}

// GetRuleset gets details about a ruleset.
func (c *Client) GetRuleset(id string) (*Ruleset, *http.Response, error) {
	resp, err := c.get("/rulesets/" + id)
	return getRulesetFromResponse(c, resp, err)
}

// UpdateRuleset updates a ruleset.
func (c *Client) UpdateRuleset(r *Ruleset) (*Ruleset, *http.Response, error) {
	v := make(map[string]*Ruleset)
	v["ruleset"] = r
	resp, err := c.put("/rulesets/"+r.ID, v, nil)
	return getRulesetFromResponse(c, resp, err)
}

func getRulesetFromResponse(c *Client, resp *http.Response, err error) (*Ruleset, *http.Response, error) {
	if err != nil {
		return nil, nil, err
	}
	var target map[string]Ruleset
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	t, nodeOK := target["ruleset"]
	if !nodeOK {
		return nil, nil, fmt.Errorf("JSON response does not have ruleset field")
	}
	return &t, resp, nil
}

// ListRulesetRules gets all rules for a ruleset.
func (c *Client) ListRulesetRules(rulesetID string) (*ListRulesetRulesResponse, error) {
	rulesResponse := new(ListRulesetRulesResponse)
	rules := make([]*RulesetRule, 0)

	// Create a handler closure capable of parsing data from the ruleset rules endpoint
	// and appending resultant ruleset rules to the return slice.
	responseHandler := func(response *http.Response) (APIListObject, error) {
		var result ListRulesetRulesResponse

		if err := c.decodeJSON(response, &result); err != nil {
			return APIListObject{}, err
		}

		rules = append(rules, result.Rules...)

		// Return stats on the current page. Caller can use this information to
		// adjust for requesting additional pages.
		return APIListObject{
			More:   result.More,
			Offset: result.Offset,
			Limit:  result.Limit,
		}, nil
	}

	// Make call to get all pages associated with the base endpoint.
	if err := c.pagedGet("/rulesets/"+rulesetID+"/rules", responseHandler); err != nil {
		return nil, err
	}
	rulesResponse.Rules = rules

	return rulesResponse, nil
}

// GetRulesetRule gets an event rule
func (c *Client) GetRulesetRule(rulesetID, ruleID string) (*RulesetRule, *http.Response, error) {
	resp, err := c.get("/rulesets/" + rulesetID + "/rules/" + ruleID)
	return getRuleFromResponse(c, resp, err)
}

// DeleteRulesetRule deletes a rule.
func (c *Client) DeleteRulesetRule(rulesetID, ruleID string) error {
	_, err := c.delete("/rulesets/" + rulesetID + "/rules/" + ruleID)
	return err
}

// CreateRulesetRule creates a new rule for a ruleset.
func (c *Client) CreateRulesetRule(rulesetID string, rule *RulesetRule) (*RulesetRule, *http.Response, error) {
	data := make(map[string]*RulesetRule)
	data["rule"] = rule
	resp, err := c.post("/rulesets/"+rulesetID+"/rules/", data, nil)
	return getRuleFromResponse(c, resp, err)
}

// UpdateRulesetRule updates a rule.
func (c *Client) UpdateRulesetRule(rulesetID, ruleID string, r *RulesetRule) (*RulesetRule, *http.Response, error) {
	v := make(map[string]*RulesetRule)
	v["rule"] = r
	resp, err := c.put("/rulesets/"+rulesetID+"/rules/"+ruleID, v, nil)
	return getRuleFromResponse(c, resp, err)
}

func getRuleFromResponse(c *Client, resp *http.Response, err error) (*RulesetRule, *http.Response, error) {
	if err != nil {
		return nil, nil, err
	}
	var target map[string]RulesetRule
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	rootNode := "rule"
	t, nodeOK := target[rootNode]
	if !nodeOK {
		return nil, nil, fmt.Errorf("JSON response does not have %s field", rootNode)
	}
	return &t, resp, nil
}
