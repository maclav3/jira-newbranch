package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type JiraIssue struct {
	Key    string `json:"key"`
	Fields struct {
		Summary string `json:"summary"`
		Updated string `json:"updated"`
	} `json:"fields"`
}

type JiraSearchResponse struct {
	Issues []JiraIssue `json:"issues"`
}

type Client struct {
	BaseURL string
	User    string
	Token   string
	HTTP    *http.Client
}

func NewClient(baseURL, user, token string) *Client {
	return &Client{
		BaseURL: strings.TrimSuffix(baseURL, "/"),
		User:    user,
		Token:   token,
		HTTP:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) GetActiveIssues() ([]JiraIssue, error) {
	// JQL: assignee = 'user' AND statusCategory != 'Done' ORDER BY updated DESC
	jql := fmt.Sprintf("assignee = '%s' AND statusCategory != 'Done' ORDER BY updated DESC", c.User)
	searchURL := fmt.Sprintf("%s/rest/api/3/search/jql", c.BaseURL)

	payload := map[string]any{
		"jql":    jql,
		"fields": []string{"summary", "updated"},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search payload: %w", err)
	}

	req, err := http.NewRequest("POST", searchURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(c.User, c.Token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch issues: status %s", resp.Status)
	}

	var searchResp JiraSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return searchResp.Issues, nil
}

func ParseJiraTime(timeStr string) (time.Time, error) {
	// Jira sometimes returns time as "2026-02-16T13:51:20.182+0000"
	// time.RFC3339Nano expects "Z07:00"
	const jiraTimeLayout = "2006-01-02T15:04:05.000-0700"
	updated, err := time.Parse(jiraTimeLayout, timeStr)
	if err != nil {
		// Fallback to RFC3339 if custom layout fails
		updated, err = time.Parse(time.RFC3339, timeStr)
	}
	return updated, err
}
