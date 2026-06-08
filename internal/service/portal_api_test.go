package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPortalAPI_GetUser(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.URL.Path == "/api/users/s1" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"sid":"s1","name":"Test","email":"t@t.com","department":"IT","groups":["admin"]}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	p := NewPortalAPIWithTimeout(srv.URL, "test-key", 0)
	user, err := p.GetUser(context.Background(), "s1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.SID != "s1" || user.Name != "Test" {
		t.Errorf("unexpected user: %+v", user)
	}

	// Not found
	_, err = p.GetUser(context.Background(), "missing")
	if err == nil {
		t.Error("expected error for missing user")
	}
}

func TestPortalAPI_GetUserGroups(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"groups":["admin","dev"]}`))
	}))
	defer srv.Close()

	p := NewPortalAPI(srv.URL, "key")
	groups, err := p.GetUserGroups(context.Background(), "s1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 2 || groups[0] != "admin" {
		t.Errorf("unexpected groups: %v", groups)
	}
}

func TestPortalAPI_HasAccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"groups":["admin","dev"]}`))
	}))
	defer srv.Close()

	p := NewPortalAPI(srv.URL, "key")

	ok, err := p.HasAccess(context.Background(), "s1", "admin")
	if err != nil || !ok {
		t.Errorf("expected has access to admin, got %v %v", ok, err)
	}

	ok, err = p.HasAccess(context.Background(), "s1", "root")
	if err != nil || ok {
		t.Errorf("expected no access to root, got %v %v", ok, err)
	}
}

func TestPortalAPI_BaseURL(t *testing.T) {
	p := NewPortalAPI("http://example.com", "k")
	if p.BaseURL() != "http://example.com" {
		t.Errorf("expected http://example.com, got %s", p.BaseURL())
	}
}
