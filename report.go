package gapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// ReportSchedule represents the schedule from a Grafana report.
type ReportSchedule struct {
	StartDate         *time.Time `json:"startDate,omitempty"`
	EndDate           *time.Time `json:"endDate,omitempty"`
	Frequency         string     `json:"frequency"`
	IntervalFrequency string     `json:"intervalFrequency"`
	IntervalAmount    int64      `json:"intervalAmount"`
	WorkdaysOnly      bool       `json:"workdaysOnly"`
	TimeZone          string     `json:"timeZone"`
	DayOfMonth        string     `json:"dayOfMonth,omitempty"`
}

// ReportTimeRange represents the time range from a Grafana report.
type ReportTimeRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// ReportOptions represents the options for a Grafana report.
type ReportOptions struct {
	Orientation string          `json:"orientation"`
	Layout      string          `json:"layout"`
	TimeRange   ReportTimeRange `json:"timeRange"`
}

// Report represents a Grafana report.
type Report struct {
	// ReadOnly
	ID     int64  `json:"id,omitempty"`
	UserID int64  `json:"userId,omitempty"`
	OrgID  int64  `json:"orgId,omitempty"`
	State  string `json:"state,omitempty"`

	DashboardID        int64          `json:"dashboardId"`
	DashboardUID       string         `json:"dashboardUid"`
	Name               string         `json:"name"`
	Recipients         string         `json:"recipients"`
	ReplyTo            string         `json:"replyTo"`
	Message            string         `json:"message"`
	Schedule           ReportSchedule `json:"schedule"`
	Options            ReportOptions  `json:"options"`
	EnableDashboardURL bool           `json:"enableDashboardUrl"`
	EnableCSV          bool           `json:"enableCsv"`
}

// Report fetches and returns a Grafana report.
func (c *Client) Report(id int64) (*Report, error) {
	return c.ReportContext(context.Background(), id)
}

// ReportContext does the same thing as Report(), but also takes in a context.
func (c *Client) ReportContext(ctx context.Context, id int64) (*Report, error) {
	path := fmt.Sprintf("/api/reports/%d", id)
	report := &Report{}
	err := c.request(ctx, "GET", path, nil, nil, report)
	if err != nil {
		return nil, err
	}

	return report, nil
}

// NewReport creates a new Grafana report.
func (c *Client) NewReport(report Report) (int64, error) {
	return c.NewReportContext(context.Background(), report)
}

// NewReportContext does the same thing as NewReport(), but also takes in a context.
func (c *Client) NewReportContext(ctx context.Context, report Report) (int64, error) {
	data, err := json.Marshal(report)
	if err != nil {
		return 0, err
	}

	result := struct {
		ID int64
	}{}

	err = c.request(ctx, "POST", "/api/reports", nil, bytes.NewBuffer(data), &result)
	if err != nil {
		return 0, err
	}

	return result.ID, nil
}

// UpdateReport updates a Grafana report.
func (c *Client) UpdateReport(report Report) error {
	return c.UpdateReportContext(context.Background(), report)
}

// UpdateReportContext does the same thing as UpdateReport(), but also takes in a context.
func (c *Client) UpdateReportContext(ctx context.Context, report Report) error {
	path := fmt.Sprintf("/api/reports/%d", report.ID)
	data, err := json.Marshal(report)
	if err != nil {
		return err
	}

	return c.request(ctx, "PUT", path, nil, bytes.NewBuffer(data), nil)
}

// DeleteReport deletes the Grafana report whose ID it's passed.
func (c *Client) DeleteReport(id int64) error {
	return c.DeleteReportContext(context.Background(), id)
}

// DeleteReportContext does the same thing as DeleteReport(), but also takes in a context.
func (c *Client) DeleteReportContext(ctx context.Context, id int64) error {
	path := fmt.Sprintf("/api/reports/%d", id)

	return c.request(ctx, "DELETE", path, nil, nil, nil)
}
