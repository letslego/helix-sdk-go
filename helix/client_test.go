package helix_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/letslego/helix-sdk-go/helix"
)

func TestChatPostsSessions(t *testing.T) {
	var gotPath string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"sessionId": "s1",
			"reply":     "hello",
			"usage":     map[string]any{},
			"toolCalls": []any{},
			"events":    []any{},
		})
	}))
	defer srv.Close()

	client := helix.New(srv.URL + "/")
	res, err := client.Chat("Plan a trip", "", true)
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/helix/v1/sessions" {
		t.Fatalf("path = %s", gotPath)
	}
	if gotBody["message"] != "Plan a trip" || gotBody["autoApprove"] != true {
		t.Fatalf("body = %#v", gotBody)
	}
	if res.SessionID != "s1" || res.Reply != "hello" {
		t.Fatalf("result = %#v", res)
	}
}

func TestBearerAndError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer secret" {
			t.Fatalf("missing bearer")
		}
		if r.URL.Path == "/helix/v1/stack" {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`{"error":"not_found"}`))
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"config": map[string]any{"model": "mock"}})
	}))
	defer srv.Close()

	client := helix.New(srv.URL)
	client.Token = "secret"
	agent, err := client.Agent()
	if err != nil {
		t.Fatal(err)
	}
	cfg := agent["config"].(map[string]any)
	if cfg["model"] != "mock" {
		t.Fatalf("agent = %#v", agent)
	}
	_, err = client.Stack()
	he, ok := err.(*helix.Error)
	if !ok || he.Status != 404 || he.Message != "not_found" {
		t.Fatalf("err = %#v", err)
	}
}

func TestApprovalsAndSchedule(t *testing.T) {
	var paths []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		if r.URL.Path == "/helix/v1/approvals" {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "s1", "status": "active", "messages": []any{}, "pendingApprovals": []any{}, "usage": map[string]any{},
			})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"sessionId": "s1", "reply": "done", "usage": map[string]any{}, "toolCalls": []any{}, "events": []any{},
		})
	}))
	defer srv.Close()

	client := helix.New(srv.URL)
	if _, err := client.ResolveApproval("s1", "a1", true); err != nil {
		t.Fatal(err)
	}
	if _, err := client.RunSchedule("weekend_watch"); err != nil {
		t.Fatal(err)
	}
	if len(paths) != 2 || paths[0] != "/helix/v1/approvals" || paths[1] != "/helix/v1/schedules/run" {
		t.Fatalf("paths = %#v", paths)
	}
}
