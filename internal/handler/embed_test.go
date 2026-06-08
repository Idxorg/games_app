package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"game-platform/tests/mocks"
)

func TestNewEmbedHandler_ZeroExpiry(t *testing.T) {
	h := NewEmbedHandler("secret", "jwt-secret", 0)
	if h.jwtExpiryH != 24 {
		t.Errorf("expected default 24 hours, got %d", h.jwtExpiryH)
	}
}

func TestNewEmbedHandler_NegativeExpiry(t *testing.T) {
	h := NewEmbedHandler("secret", "jwt-secret", -5)
	if h.jwtExpiryH != 24 {
		t.Errorf("expected default 24 hours, got %d", h.jwtExpiryH)
	}
}

func TestEmbedAuth_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "test-secret"
	jwtSecret := "test-jwt-secret-32charslong!!!!!!"
	h := NewEmbedHandler(secret, jwtSecret, 24)

	r := gin.New()
	r.POST("/embed", h.EmbedAuth)

	body, _ := json.Marshal(embedRequest{
		SID:        "emp_001",
		Email:      "alice@example.com",
		Name:       "Alice",
		Department: "Engineering",
	})
	req := httptest.NewRequest(http.MethodPost, "/embed", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Erlink-Embed-Secret", secret)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp embedResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Token == "" {
		t.Fatal("expected non-empty token")
	}
	if resp.SID != "emp_001" {
		t.Errorf("expected sid emp_001, got %s", resp.SID)
	}
	if !resp.Valid {
		t.Error("expected valid=true")
	}

	// Verify JWT is parseable
	parsed, err := jwt.Parse(resp.Token, func(tok *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil || !parsed.Valid {
		t.Fatal("JWT should be valid and parseable")
	}
}

func TestEmbedAuth_InvalidSecret(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewEmbedHandler("correct-secret", "jwt-secret", 24)
	r := gin.New()
	r.POST("/embed", h.EmbedAuth)

	body, _ := json.Marshal(embedRequest{SID: "emp_001"})
	req := httptest.NewRequest(http.MethodPost, "/embed", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Erlink-Embed-Secret", "wrong-secret")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestEmbedAuth_DevMode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewEmbedHandler("", "jwt-secret", 24)
	r := gin.New()
	r.POST("/embed", h.EmbedAuth)

	body, _ := json.Marshal(embedRequest{SID: "emp_001"})
	req := httptest.NewRequest(http.MethodPost, "/embed", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// No secret header needed in dev mode
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 in dev mode, got %d", w.Code)
	}
}

func TestEmbedAuth_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewEmbedHandler("secret", "jwt-secret", 24)
	r := gin.New()
	r.POST("/embed", h.EmbedAuth)

	req := httptest.NewRequest(http.MethodPost, "/embed", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Erlink-Embed-Secret", "secret")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestEmbedAuth_MissingSID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewEmbedHandler("secret", "jwt-secret", 24)
	r := gin.New()
	r.POST("/embed", h.EmbedAuth)

	body, _ := json.Marshal(embedRequest{Email: "a@b.com"})
	req := httptest.NewRequest(http.MethodPost, "/embed", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Erlink-Embed-Secret", "secret")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing SID, got %d", w.Code)
	}
}

func TestDecodeEmbedRequest(t *testing.T) {
	data := []byte(`{"sid":"s1","email":"e@e.com","name":"N","department":"D"}`)
	req, err := decodeEmbedRequest(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.SID != "s1" || req.Email != "e@e.com" || req.Name != "N" || req.Department != "D" {
		t.Errorf("unexpected parsed values: %+v", req)
	}

	// Invalid JSON
	_, err = decodeEmbedRequest([]byte("bad"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestEmbedHandler_SetUserRepo(t *testing.T) {
	h := NewEmbedHandler("secret", "key", 3600)
	h.SetUserRepo(mocks.NewMockUserRepo())
	if h.userRepo == nil {
		t.Error("userRepo should not be nil after SetUserRepo")
	}
}
