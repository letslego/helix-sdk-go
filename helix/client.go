package helix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Error is returned for non-2xx Helix responses.
type Error struct {
	Status  int
	Message string
	Body    any
}

func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("helix HTTP %d", e.Status)
}

// Client talks to a running Helix console (`helix console`).
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// New creates a Client for baseURL (e.g. http://127.0.0.1:8787).
func New(baseURL string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Agent GET /helix/v1/agent
func (c *Client) Agent() (map[string]any, error) {
	var out map[string]any
	err := c.do("GET", "/helix/v1/agent", nil, &out)
	return out, err
}

// Stack GET /helix/v1/stack
func (c *Client) Stack() (map[string]any, error) {
	var out map[string]any
	err := c.do("GET", "/helix/v1/stack", nil, &out)
	return out, err
}

// Sessions GET /helix/v1/sessions
func (c *Client) Sessions() ([]SessionRecord, error) {
	var out []SessionRecord
	err := c.do("GET", "/helix/v1/sessions", nil, &out)
	return out, err
}

// Workflows GET /helix/v1/workflows
func (c *Client) Workflows(sessionID string) ([]WorkflowRun, error) {
	path := "/helix/v1/workflows"
	if sessionID != "" {
		path += "?sessionId=" + url.QueryEscape(sessionID)
	}
	var out []WorkflowRun
	err := c.do("GET", path, nil, &out)
	return out, err
}

// Events GET /helix/v1/events
func (c *Client) Events(sessionID string) ([]RuntimeEvent, error) {
	path := "/helix/v1/events"
	if sessionID != "" {
		path += "?sessionId=" + url.QueryEscape(sessionID)
	}
	var out []RuntimeEvent
	err := c.do("GET", path, nil, &out)
	return out, err
}

// Run POST /helix/v1/sessions
func (c *Client) Run(opts RunOptions) (*RunResult, error) {
	channel := opts.Channel
	if channel == "" {
		channel = "http"
	}
	body := map[string]any{
		"message":     opts.Message,
		"autoApprove": opts.AutoApprove,
		"channel":     channel,
	}
	if opts.SessionID != "" {
		body["sessionId"] = opts.SessionID
	}
	var out RunResult
	if err := c.do("POST", "/helix/v1/sessions", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Chat is a convenience wrapper around Run.
func (c *Client) Chat(message, sessionID string, autoApprove bool) (*RunResult, error) {
	return c.Run(RunOptions{
		Message:     message,
		SessionID:   sessionID,
		AutoApprove: autoApprove,
	})
}

// ResolveApproval POST /helix/v1/approvals
func (c *Client) ResolveApproval(sessionID, approvalID string, approve bool) (*SessionRecord, error) {
	body := map[string]any{
		"sessionId":  sessionID,
		"approvalId": approvalID,
		"approve":    approve,
	}
	var out SessionRecord
	if err := c.do("POST", "/helix/v1/approvals", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RunSchedule POST /helix/v1/schedules/run
func (c *Client) RunSchedule(name string) (*RunResult, error) {
	var out RunResult
	if err := c.do("POST", "/helix/v1/schedules/run", map[string]any{"name": name}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) do(method, path string, body any, out any) error {
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(method, c.BaseURL+path, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		var parsed any
		_ = json.Unmarshal(raw, &parsed)
		msg := fmt.Sprintf("helix HTTP %d", res.StatusCode)
		if m, ok := parsed.(map[string]any); ok {
			if e, ok := m["error"].(string); ok {
				msg = e
			}
		}
		return &Error{Status: res.StatusCode, Message: msg, Body: parsed}
	}
	if out == nil || len(raw) == 0 {
		return nil
	}
	return json.Unmarshal(raw, out)
}
