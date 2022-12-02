package gapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// AlertRule represents a Grafana Alert Rule.
type AlertRule struct {
	Annotations  map[string]string `json:"annotations,omitempty"`
	Condition    string            `json:"condition"`
	Data         []*AlertQuery     `json:"data"`
	ExecErrState ExecErrState      `json:"execErrState"`
	FolderUID    string            `json:"folderUid"`
	ID           int64             `json:"id,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
	NoDataState  NoDataState       `json:"noDataState"`
	OrgID        int64             `json:"orgId"`
	RuleGroup    string            `json:"ruleGroup"`
	Title        string            `json:"title"`
	UID          string            `json:"uid,omitempty"`
	Updated      time.Time         `json:"updated"`
	For          string            `json:"for"`
	ForDuration  time.Duration     `json:"-"`
	Provenance   string            `json:"provenance"`
}

// RuleGroup represents a group of rules in Grafana Alerting.
type RuleGroup struct {
	Title     string      `json:"title"`
	FolderUID string      `json:"folderUid"`
	Interval  int64       `json:"interval"`
	Rules     []AlertRule `json:"rules"`
}

// AlertQuery represents a single query stage associated with an alert definition.
type AlertQuery struct {
	DatasourceUID     string            `json:"datasourceUid,omitempty"`
	Model             interface{}       `json:"model"`
	QueryType         string            `json:"queryType,omitempty"`
	RefID             string            `json:"refId,omitempty"`
	RelativeTimeRange RelativeTimeRange `json:"relativeTimeRange"`
}

type ExecErrState string
type NoDataState string

const (
	ErrOK          ExecErrState = "OK"
	ErrError       ExecErrState = "Error"
	ErrAlerting    ExecErrState = "Alerting"
	NoDataOk       NoDataState  = "OK"
	NoData         NoDataState  = "NoData"
	NoDataAlerting NoDataState  = "Alerting"
)

// RelativeTimeRange represents the time range for an alert query.
type RelativeTimeRange struct {
	From time.Duration `json:"from"`
	To   time.Duration `json:"to"`
}

// AlertRule fetches a single alert rule, identified by its UID.
func (c *Client) AlertRule(uid string) (AlertRule, error) {
	return c.AlertRuleContext(context.Background(), uid)
}

// AlertRuleContext does the same thing as AlertRule(), but also takes in a context.
func (c *Client) AlertRuleContext(ctx context.Context, uid string) (AlertRule, error) {
	path := fmt.Sprintf("/api/v1/provisioning/alert-rules/%s", uid)
	result := AlertRule{}
	err := c.request(ctx, "GET", path, nil, nil, &result)
	if err != nil {
		return AlertRule{}, err
	}
	return result, err
}

// AlertRuleGroup fetches a group of alert rules, identified by its name and the UID of its folder.
func (c *Client) AlertRuleGroup(folderUID string, name string) (RuleGroup, error) {
	return c.AlertRuleGroupContext(context.Background(), folderUID, name)
}

// AlertRuleGroupContext does the same thing as AlertRuleGroup(), but also takes in a context.
func (c *Client) AlertRuleGroupContext(ctx context.Context, folderUID string, name string) (RuleGroup, error) {
	path := fmt.Sprintf("/api/v1/provisioning/folder/%s/rule-groups/%s", folderUID, name)
	result := RuleGroup{}
	err := c.request(ctx, "GET", path, nil, nil, &result)
	return result, err
}

// SetAlertRuleGroup overwrites an existing rule group on the server.
func (c *Client) SetAlertRuleGroup(group RuleGroup) error {
	return c.SetAlertRuleGroupContext(context.Background(), group)
}

// SetAlertRuleGroupContext does the same thing as SetAlertRuleGroup(), but also takes in a context.
func (c *Client) SetAlertRuleGroupContext(ctx context.Context, group RuleGroup) error {
	syncCalculatedRuleGroupFields(&group)
	folderUID := group.FolderUID
	name := group.Title
	req, err := json.Marshal(group)
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("/api/v1/provisioning/folder/%s/rule-groups/%s", folderUID, name)
	return c.request(ctx, "PUT", uri, nil, bytes.NewBuffer(req), nil)
}

// NewAlertRule creates a new alert rule and returns its UID.
func (c *Client) NewAlertRule(ar *AlertRule) (string, error) {
	return c.NewAlertRuleContext(context.Background(), ar)
}

// NewAlertRuleContext does the same thing as NewAlertRule(), but also takes in a context.
func (c *Client) NewAlertRuleContext(ctx context.Context, ar *AlertRule) (string, error) {
	syncCalculatedRuleFields(ar)
	req, err := json.Marshal(ar)
	if err != nil {
		return "", err
	}
	result := AlertRule{}
	err = c.request(ctx, "POST", "/api/v1/provisioning/alert-rules", nil, bytes.NewBuffer(req), &result)
	if err != nil {
		return "", err
	}
	return result.UID, nil
}

// UpdateAlertRule replaces an alert rule, identified by the alert rule's UID.
func (c *Client) UpdateAlertRule(ar *AlertRule) error {
	return c.UpdateAlertRuleContext(context.Background(), ar)
}

// UpdateAlertRuleContext does the same thing as UpdateAlertRule(), but also takes in a context.
func (c *Client) UpdateAlertRuleContext(ctx context.Context, ar *AlertRule) error {
	syncCalculatedRuleFields(ar)
	uri := fmt.Sprintf("/api/v1/provisioning/alert-rules/%s", ar.UID)
	req, err := json.Marshal(ar)
	if err != nil {
		return err
	}

	return c.request(ctx, "PUT", uri, nil, bytes.NewBuffer(req), nil)
}

// DeleteAlertRule deletes a alert rule, identified by the alert rule's UID.
func (c *Client) DeleteAlertRule(uid string) error {
	return c.DeleteAlertRuleContext(context.Background(), uid)
}

// DeleteAlertRuleContext does the same thing as DeleteAlertRule(), but also takes in a context.
func (c *Client) DeleteAlertRuleContext(ctx context.Context, uid string) error {
	uri := fmt.Sprintf("/api/v1/provisioning/alert-rules/%s", uid)
	return c.request(ctx, "DELETE", uri, nil, nil, nil)
}

func syncCalculatedRuleGroupFields(group *RuleGroup) {
	for i := range group.Rules {
		syncCalculatedRuleFields(&group.Rules[i])
	}
}

func syncCalculatedRuleFields(rule *AlertRule) {
	// rule.For is the newer field. Older systems may not provide the value.
	// If the user provided a For, prefer that over whatever we might calculate.
	// Otherwise, translate the time.Duration-based field to the format that the rule API expects.
	if rule.For == "" {
		rule.For = timeDurationToRuleDuration(rule.ForDuration)
	}
}

// timeDurationToRuleDuration converts a typical time.Duration to the string-based format that alert rules expect.
// This code is adapted from Prometheus: https://github.com/prometheus/common/blob/dfbc25bd00225c70aca0d94c3c4bb7744f28ace0/model/time.go#L236
func timeDurationToRuleDuration(d time.Duration) string {
	ms := int64(d / time.Millisecond)
	if ms == 0 {
		return "0s"
	}

	r := ""
	f := func(unit string, mult int64, exact bool) {
		if exact && ms%mult != 0 {
			return
		}
		if v := ms / mult; v > 0 {
			r += fmt.Sprintf("%d%s", v, unit)
			ms -= v * mult
		}
	}

	// Only format years and weeks if the remainder is zero, as it is often
	// easier to read 90d than 12w6d.
	f("y", 1000*60*60*24*365, true)
	f("w", 1000*60*60*24*7, true)

	f("d", 1000*60*60*24, false)
	f("h", 1000*60*60, false)
	f("m", 1000*60, false)
	f("s", 1000, false)
	f("ms", 1, false)

	return r
}
