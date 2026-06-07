package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"game-platform/internal/handler"
	"game-platform/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	embedTestHandoffSecret = "test-embed-secret-g0"
	embedTestJWTSecret     = "test-jwt-secret-g0-32chars!!"
)

// ---------------------------------------------------------------------------
// Helper: create a full test router with embed + health endpoints
// ---------------------------------------------------------------------------

func setupEmbedRouter(handoffSecret string, userRepo *mocks.MockUserRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Embed handler
	eh := handler.NewEmbedHandler(handoffSecret, embedTestJWTSecret, 24)
	if userRepo != nil {
		eh.SetUserRepo(userRepo)
	}
	r.POST("/api/v1/auth/embed", eh.EmbedAuth)

	// Health endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	return r
}

// ---------------------------------------------------------------------------
// TestAuthEmbed_Success: valid secret + body → 200 + JWT
// ---------------------------------------------------------------------------

func TestAuthEmbed_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userRepo := mocks.NewMockUserRepo()
	r := setupEmbedRouter(embedTestHandoffSecret, userRepo)

	body, _ := json.Marshal(map[string]string{
		"sid":        "emp_001",
		"email":      "alice@example.com",
		"name":       "Alice Smith",
		"department": "Engineering",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/embed", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Erlink-Embed-Secret", embedTestHandoffSecret)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	token, ok := resp["token"].(string)
	if !ok || token == "" {
		t.Fatal("expected non-empty token in response")
	}
	if resp["sid"] != "emp_001" {
		t.Fatalf("expected sid=emp_001, got %v", resp["sid"])
	}
	if resp["valid"] != true {
		t.Fatal("expected valid=true")
	}

	// Verify JWT can be parsed with the correct secret
	parsed, err := jwt.Parse(token, func(tok *jwt.Token) (interface{}, error) {
		return []byte(embedTestJWTSecret), nil
	})
	if err != nil {
		t.Fatalf("failed to parse returned JWT: %v", err)
	}
	if !parsed.Valid {
		t.Fatal("parsed JWT should be valid")
	}

	// Verify user was upserted
	u, err := userRepo.GetBySID(nil, "emp_001")
	if err != nil {
		t.Fatalf("GetBySID error: %v", err)
	}
	if u == nil {
		t.Fatal("expected user to be created in DB")
	}
	if u.Email != "alice@example.com" {
		t.Fatalf("expected email alice@example.com, got %q", u.Email)
	}
}

// ---------------------------------------------------------------------------
// TestAuthEmbed_MissingSecret: no header when secret configured → 401
// ---------------------------------------------------------------------------

func TestAuthEmbed_MissingSecret(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupEmbedRouter(embedTestHandoffSecret, nil)

	body, _ := json.Marshal(map[string]string{
		"sid":        "emp_002",
		"email":      "bob@example.com",
		"name":       "Bob",
		"department": "Sales",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/embed", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// Intentionally NOT setting X-Erlink-Embed-Secret
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d; body: %s", w.Code, w.Body.String())
	}

	var errResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &errResp)
	if errResp["error"] == nil {
		t.Fatal("expected error field in response")
	}
}

// ---------------------------------------------------------------------------
// TestAuthEmbed_InvalidBody: bad JSON → 400
// ---------------------------------------------------------------------------

func TestAuthEmbed_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupEmbedRouter(embedTestHandoffSecret, nil)

	tests := []struct {
		name string
		body string
	}{
		{"garbage", "not-json"},
		{"invalid json", `{bad}`},
		{"empty", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/embed",
				bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Erlink-Embed-Secret", embedTestHandoffSecret)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected 400 for %s, got %d; body: %s", tc.name, w.Code, w.Body.String())
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestAuthEmbed_DevMode: no secret configured → still works
// ---------------------------------------------------------------------------

func TestAuthEmbed_DevMode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupEmbedRouter("", nil) // empty handoffSecret = dev mode

	body, _ := json.Marshal(map[string]string{
		"sid":   "emp_dev",
		"email": "dev@example.com",
		"name":  "Dev User",
	})

	// No secret header at all — should still succeed in dev mode
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/embed", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 in dev mode, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] == nil {
		t.Fatal("expected token in dev mode response")
	}
	if resp["valid"] != true {
		t.Fatal("expected valid=true in dev mode")
	}

	// Also try with wrong secret — should still succeed
	body2, _ := json.Marshal(map[string]string{
		"sid":   "emp_dev2",
		"email": "dev2@example.com",
		"name":  "Dev User 2",
	})
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/embed", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-Erlink-Embed-Secret", "wrong-secret")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200 in dev mode with wrong secret, got %d; body: %s", w2.Code, w2.Body.String())
	}
}

// ---------------------------------------------------------------------------
// TestHealth: GET /health → 200 + status ok
// ---------------------------------------------------------------------------

func TestHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := setupEmbedRouter(embedTestHandoffSecret, nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse health response: %v", err)
	}

	if resp["status"] != "ok" {
		t.Fatalf("expected status=ok, got %v", resp["status"])
	}

	timeStr, ok := resp["time"].(string)
	if !ok || timeStr == "" {
		t.Fatal("expected non-empty time field")
	}

	// Verify time is parseable RFC3339
	if _, err := time.Parse(time.RFC3339, timeStr); err != nil {
		t.Fatalf("time field is not valid RFC3339: %v", err)
	}
}
