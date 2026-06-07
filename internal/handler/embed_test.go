package handler

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"game-platform/internal/model"
	"game-platform/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	testHandoffSecret = "test-embed-secret-12345"
	testJWTSecret     = "test-jwt-secret-abcde"
)

// setupEmbedTestRouter creates a gin router with the embed endpoint registered.
// Returns the handler so the caller can inspect configuration.
func setupEmbedTestRouter(handoffSecret string, jwtExpiryH int) (*EmbedHandler, *gin.Engine) {
	return setupEmbedTestRouterWithRepo(handoffSecret, jwtExpiryH, nil)
}

// setupEmbedTestRouterWithRepo creates a gin router with the embed endpoint
// and an optional user repo registered.
func setupEmbedTestRouterWithRepo(handoffSecret string, jwtExpiryH int, userRepo model.UserRepo) (*EmbedHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewEmbedHandler(handoffSecret, testJWTSecret, jwtExpiryH)
	if userRepo != nil {
		h.SetUserRepo(userRepo)
	}
	r.POST("/api/v1/auth/embed", h.EmbedAuth)
	return h, r
}

// doEmbedRequest is a helper that sends a POST to the embed endpoint.
func doEmbedRequest(r *gin.Engine, body []byte, secret string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/embed", bytes.NewReader(body))
	if secret != "" {
		req.Header.Set("X-Erlink-Embed-Secret", secret)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// validEmbedBody returns a JSON-encoded embed request with the given fields.
func validEmbedBody(sid, email, name, department string) []byte {
	b, _ := json.Marshal(embedRequest{
		SID:        sid,
		Email:      email,
		Name:       name,
		Department: department,
	})
	return b
}

// ---------------------------------------------------------------------------
// Success path
// ---------------------------------------------------------------------------

func TestEmbedAuth_Success(t *testing.T) {
	_, r := setupEmbedTestRouter(testHandoffSecret, 24)
	body := validEmbedBody("emp_001", "alice@example.com", "Alice Smith", "Engineering")

	w := doEmbedRequest(r, body, testHandoffSecret)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp embedResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected non-empty token")
	}
	if resp.SID != "emp_001" {
		t.Fatalf("expected sid emp_001, got %q", resp.SID)
	}
	if !resp.Valid {
		t.Fatal("expected valid=true")
	}
}

// ---------------------------------------------------------------------------
// Error paths — secret configuration
// ---------------------------------------------------------------------------

func TestEmbedAuth_NoSecret_DevMode(t *testing.T) {
	// When handoffSecret is empty, handler should work in dev mode (skip check).
	_, r := setupEmbedTestRouter("", 24)
	body := validEmbedBody("emp_001", "a@b.com", "A", "D")

	w := doEmbedRequest(r, body, testHandoffSecret)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 in dev mode, got %d; body: %s", w.Code, w.Body.String())
	}
	var resp embedResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected non-empty token in dev mode")
	}
}

func TestEmbedAuth_WrongSecret(t *testing.T) {
	_, r := setupEmbedTestRouter(testHandoffSecret, 24)
	body := validEmbedBody("emp_001", "a@b.com", "A", "D")

	w := doEmbedRequest(r, body, "totally-wrong-secret")

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestEmbedAuth_MissingHeader(t *testing.T) {
	// No X-Erlink-Embed-Secret header at all → constant-time compare fails → 401
	_, r := setupEmbedTestRouter(testHandoffSecret, 24)
	body := validEmbedBody("emp_001", "a@b.com", "A", "D")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/embed", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// Intentionally NOT setting X-Erlink-Embed-Secret
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing header, got %d; body: %s", w.Code, w.Body.String())
	}
}

// ---------------------------------------------------------------------------
// Error paths — request body
// ---------------------------------------------------------------------------

func TestEmbedAuth_MissingSID(t *testing.T) {
	_, r := setupEmbedTestRouter(testHandoffSecret, 24)
	body := validEmbedBody("", "a@b.com", "A", "D") // empty SID

	w := doEmbedRequest(r, body, testHandoffSecret)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d; body: %s", w.Code, w.Body.String())
	}
	var errResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &errResp)
	if errResp["error"] != "sid is required" {
		t.Fatalf("unexpected error: %v", errResp["error"])
	}
}

func TestEmbedAuth_InvalidBody(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"garbage bytes", "not-json-at-all"},
		{"invalid json", `{invalid}`},
		{"empty body", ""},
		{"partial json", `{\"sid\":\"emp_001\"`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, r := setupEmbedTestRouter(testHandoffSecret, 24)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/embed",
				strings.NewReader(tc.body))
			req.Header.Set("X-Erlink-Embed-Secret", testHandoffSecret)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected 400 for %s, got %d; body: %s", tc.name, w.Code, w.Body.String())
			}
		})
	}
}

// ---------------------------------------------------------------------------
// JWT claims verification
// ---------------------------------------------------------------------------

func TestEmbedAuth_JWTClaims(t *testing.T) {
	_, r := setupEmbedTestRouter(testHandoffSecret, 24)
	sid, email, name, dept := "emp_042", "bob@example.com", "Bob Jones", "Marketing"
	body := validEmbedBody(sid, email, name, dept)

	w := doEmbedRequest(r, body, testHandoffSecret)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp embedResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	// Parse the returned JWT
	token, err := jwt.Parse(resp.Token, func(tok *jwt.Token) (interface{}, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			t.Fatalf("unexpected signing method: %v", tok.Header["alg"])
		}
		return []byte(testJWTSecret), nil
	})
	if err != nil {
		t.Fatalf("failed to parse JWT: %v", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("failed to cast claims to MapClaims")
	}

	// Verify each claim
	assertClaimString(t, claims, "sid", sid)
	assertClaimString(t, claims, "email", email)
	assertClaimString(t, claims, "name", name)
	assertClaimString(t, claims, "department", dept)

	// Verify groups is ["games"]
	groups, ok := claims["groups"].([]interface{})
	if !ok {
		t.Fatalf("expected groups to be array, got %T", claims["groups"])
	}
	if len(groups) != 1 || groups[0] != "games" {
		t.Fatalf("expected groups=[\"games\"], got %v", groups)
	}

	// iat and exp should be present and numeric
	assertClaimNumeric(t, claims, "iat")
	assertClaimNumeric(t, claims, "exp")
}

func TestEmbedAuth_JWTExpiry(t *testing.T) {
	h, r := setupEmbedTestRouter(testHandoffSecret, 24)
	body := validEmbedBody("emp_099", "x@y.com", "X", "Z")

	before := time.Now()
	w := doEmbedRequest(r, body, testHandoffSecret)
	after := time.Now()

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp embedResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	token, _ := jwt.Parse(resp.Token, func(tok *jwt.Token) (interface{}, error) {
		return []byte(testJWTSecret), nil
	})
	claims := token.Claims.(jwt.MapClaims)

	expFloat, _ := claims["exp"].(float64)
	iatFloat, _ := claims["iat"].(float64)
	exp := int64(expFloat)
	iat := int64(iatFloat)

	// exp - iat should be ~24 hours (86400 seconds), within a few seconds of skew
	diff := exp - iat
	if math.Abs(float64(diff-86400)) > 2.0 {
		t.Fatalf("expected exp-iat ≈ 86400s (24h), got %d", diff)
	}

	// Verify iat is between before and after
	if iat < before.Unix()-1 || iat > after.Unix()+1 {
		t.Fatalf("iat %d outside expected range [%d, %d]", iat, before.Unix(), after.Unix())
	}

	// Verify handler used configured expiry
	if h.jwtExpiryH != 24 {
		t.Fatalf("handler jwtExpiryH=%d, want 24", h.jwtExpiryH)
	}
}

// ---------------------------------------------------------------------------
// Dev mode / configuration
// ---------------------------------------------------------------------------

func TestEmbedAuth_DevMode(t *testing.T) {
	// Dev mode: empty handoff secret → skip check → works normally
	_, r := setupEmbedTestRouter("", 24)
	body := validEmbedBody("emp_001", "a@b.com", "A", "D")

	// Even with empty secret header, should succeed (dev mode)
	w := doEmbedRequest(r, body, "")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 in dev mode (no handoff secret), got %d; body: %s", w.Code, w.Body.String())
	}
	var resp embedResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected non-empty token in dev mode")
	}
}

// ---------------------------------------------------------------------------
// User upsert
// ---------------------------------------------------------------------------

func TestEmbedAuth_UpsertUser(t *testing.T) {
	userRepo := mocks.NewMockUserRepo()
	h, r := setupEmbedTestRouterWithRepo(testHandoffSecret, 24, userRepo)

	body := validEmbedBody("emp_upsert", "upsert@test.com", "Upsert User", "DevOps")
	w := doEmbedRequest(r, body, testHandoffSecret)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}

	// Verify user was created
	u, err := userRepo.GetBySID(nil, "emp_upsert")
	if err != nil {
		t.Fatalf("GetBySID error: %v", err)
	}
	if u == nil {
		t.Fatal("expected user to be created")
	}
	if u.Email != "upsert@test.com" {
		t.Fatalf("expected email upsert@test.com, got %q", u.Email)
	}
	if u.Name != "Upsert User" {
		t.Fatalf("expected name 'Upsert User', got %q", u.Name)
	}

	// Now update the same user via another embed request
	body2 := validEmbedBody("emp_upsert", "updated@test.com", "Updated Name", "Engineering")
	w2 := doEmbedRequest(r, body2, testHandoffSecret)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200 on update, got %d", w2.Code)
	}

	u2, _ := userRepo.GetBySID(nil, "emp_upsert")
	if u2.Email != "updated@test.com" {
		t.Fatalf("expected updated email, got %q", u2.Email)
	}
	if u2.Name != "Updated Name" {
		t.Fatalf("expected updated name, got %q", u2.Name)
	}
	if u2.Department != "Engineering" {
		t.Fatalf("expected updated department, got %q", u2.Department)
	}

	// Verify handler has userRepo set
	_ = h
}

// ---------------------------------------------------------------------------
// decodeEmbedRequest helper
// ---------------------------------------------------------------------------

func TestDecodeEmbedRequest(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantSID string
		wantErr bool
	}{
		{
			name:    "valid full request",
			data:    `{"sid":"emp_123","email":"a@b.com","name":"Alice","department":"IT"}`,
			wantSID: "emp_123",
			wantErr: false,
		},
		{
			name:    "valid minimal (only sid)",
			data:    `{"sid":"emp_456"}`,
			wantSID: "emp_456",
			wantErr: false,
		},
		{
			name:    "empty json object",
			data:    `{}`,
			wantSID: "",
			wantErr: false,
		},
		{
			name:    "invalid json",
			data:    `{not valid}`,
			wantSID: "",
			wantErr: true,
		},
		{
			name:    "empty string",
			data:    "",
			wantSID: "",
			wantErr: true,
		},
		{
			name:    "wrong types",
			data:    `{"sid":12345,"email":true}`,
			wantSID: "",
			wantErr: true, // json.Unmarshal rejects number for string field
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := decodeEmbedRequest([]byte(tc.data))
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tc.wantErr && req.SID != tc.wantSID {
				t.Fatalf("expected sid=%q, got %q", tc.wantSID, req.SID)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Edge cases
// ---------------------------------------------------------------------------

func TestEmbedAuth_DefaultExpiry(t *testing.T) {
	// Passing 0 or negative expiry should default to 24h
	h := NewEmbedHandler(testHandoffSecret, testJWTSecret, 0)
	if h.jwtExpiryH != 24 {
		t.Fatalf("expected default expiry 24, got %d", h.jwtExpiryH)
	}

	h2 := NewEmbedHandler(testHandoffSecret, testJWTSecret, -5)
	if h2.jwtExpiryH != 24 {
		t.Fatalf("expected default expiry 24 for negative input, got %d", h2.jwtExpiryH)
	}
}

func TestEmbedAuth_EmptyOptionalFields(t *testing.T) {
	// email, name, department can be empty — handler should still work
	_, r := setupEmbedTestRouter(testHandoffSecret, 24)
	body, _ := json.Marshal(embedRequest{SID: "emp_777"})

	w := doEmbedRequest(r, body, testHandoffSecret)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 with empty optional fields, got %d; body: %s", w.Code, w.Body.String())
	}
	var resp embedResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.SID != "emp_777" {
		t.Fatalf("expected sid emp_777, got %q", resp.SID)
	}

	// Verify JWT contains empty strings for optional fields
	token, _ := jwt.Parse(resp.Token, func(tok *jwt.Token) (interface{}, error) {
		return []byte(testJWTSecret), nil
	})
	claims := token.Claims.(jwt.MapClaims)
	assertClaimString(t, claims, "email", "")
	assertClaimString(t, claims, "name", "")
	assertClaimString(t, claims, "department", "")
}

func TestEmbedAuth_CustomExpiry(t *testing.T) {
	// Use a 48h expiry and verify
	_, r := setupEmbedTestRouter(testHandoffSecret, 48)
	body := validEmbedBody("emp_888", "c@d.com", "C", "D")

	before := time.Now()
	w := doEmbedRequest(r, body, testHandoffSecret)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp embedResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	token, _ := jwt.Parse(resp.Token, func(tok *jwt.Token) (interface{}, error) {
		return []byte(testJWTSecret), nil
	})
	claims := token.Claims.(jwt.MapClaims)

	expFloat, _ := claims["exp"].(float64)
	iatFloat, _ := claims["iat"].(float64)
	diff := int64(expFloat) - int64(iatFloat)

	// 48h = 172800 seconds
	if diff != 172800 {
		t.Fatalf("expected exp-iat=172800 (48h), got %d", diff)
	}

	// iat should be approximately now
	iat := int64(iatFloat)
	if iat < before.Unix()-1 || iat > time.Now().Unix()+1 {
		t.Fatalf("iat %d outside expected range", iat)
	}
}

func TestEmbedAuth_ResponseContentType(t *testing.T) {
	_, r := setupEmbedTestRouter(testHandoffSecret, 24)
	body := validEmbedBody("emp_001", "a@b.com", "A", "D")

	w := doEmbedRequest(r, body, testHandoffSecret)

	ct := w.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}
}

// ---------------------------------------------------------------------------
// Claim assertion helpers
// ---------------------------------------------------------------------------

func assertClaimString(t *testing.T, claims jwt.MapClaims, key, want string) {
	t.Helper()
	val, ok := claims[key].(string)
	if !ok {
		t.Fatalf("claim %q missing or not string (got %T: %v)", key, claims[key], claims[key])
	}
	if val != want {
		t.Fatalf("claim %q: want %q, got %q", key, want, val)
	}
}

func assertClaimNumeric(t *testing.T, claims jwt.MapClaims, key string) {
	t.Helper()
	_, ok := claims[key].(float64)
	if !ok {
		t.Fatalf("claim %q missing or not numeric (got %T: %v)", key, claims[key], claims[key])
	}
}
