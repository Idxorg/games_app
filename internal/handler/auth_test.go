package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"game-platform/internal/model"
	"game-platform/internal/service"
	"game-platform/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ---------- AuthHandler tests ----------

func TestNewAuthHandler(t *testing.T) {
	h := NewAuthHandler(nil, nil)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestSetJWTSecret(t *testing.T) {
	h := NewAuthHandler(nil, nil)
	h.SetJWTSecret("my-secret")
	if h.jwtSecret != "my-secret" {
		t.Errorf("expected my-secret, got %s", h.jwtSecret)
	}
}

func TestVerifyToken_NoSecret(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewAuthHandler(nil, nil)
	r := gin.New()
	r.POST("/verify", h.VerifyToken)

	req := httptest.NewRequest(http.MethodPost, "/verify", nil)
	req.Header.Set("Authorization", "Bearer anything")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 in dev mode, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["valid"] != true {
		t.Error("expected valid=true")
	}
}

func TestVerifyToken_NoAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewAuthHandler(nil, nil)
	h.SetJWTSecret("secret")
	r := gin.New()
	r.POST("/verify", h.VerifyToken)

	req := httptest.NewRequest(http.MethodPost, "/verify", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestVerifyToken_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "test-secret-32charslong!!!!!!!!!!"
	portalAPI := service.NewPortalAPI("", "")
	h := NewAuthHandler(nil, portalAPI)
	h.SetJWTSecret(secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sid":   "emp_001",
		"email": "test@test.com",
		"name":  "Test",
		"groups": []string{"games"},
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
	})
	tokenStr, _ := token.SignedString([]byte(secret))

	r := gin.New()
	r.POST("/verify", h.VerifyToken)

	req := httptest.NewRequest(http.MethodPost, "/verify", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["valid"] != true {
		t.Error("expected valid=true")
	}
	if resp["sid"] != "emp_001" {
		t.Errorf("expected sid emp_001, got %v", resp["sid"])
	}
	if resp["email"] != "test@test.com" {
		t.Errorf("expected email test@test.com, got %v", resp["email"])
	}
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewAuthHandler(nil, nil)
	h.SetJWTSecret("secret")
	r := gin.New()
	r.POST("/verify", h.VerifyToken)

	req := httptest.NewRequest(http.MethodPost, "/verify", nil)
	req.Header.Set("Authorization", "Bearer garbage")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestVerifyToken_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "test-secret-32charslong!!!!!!!!!!"
	h := NewAuthHandler(nil, nil)
	h.SetJWTSecret(secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sid": "emp_001",
		"exp": time.Now().Add(-1 * time.Hour).Unix(),
	})
	tokenStr, _ := token.SignedString([]byte(secret))

	r := gin.New()
	r.POST("/verify", h.VerifyToken)

	req := httptest.NewRequest(http.MethodPost, "/verify", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for expired token, got %d", w.Code)
	}
}

func TestVerifyToken_BearerPrefixOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewAuthHandler(nil, nil)
	h.SetJWTSecret("secret")
	r := gin.New()
	r.POST("/verify", h.VerifyToken)

	req := httptest.NewRequest(http.MethodPost, "/verify", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for Basic auth, got %d", w.Code)
	}
}

// ---------- UserHandler tests ----------

func TestNewUserHandler(t *testing.T) {
	h := NewUserHandler(nil, nil, nil, nil)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestGetProfile_Found(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userRepo := mocks.NewMockUserRepo()
	userRepo.Create(context.Background(), "emp_001", "test@test.com", "Test User", "male", "IT", "Dev", "")

	h := NewUserHandler(userRepo, nil, nil, nil)
	r := gin.New()
	r.GET("/users/:sid/profile", h.GetProfile)

	req := httptest.NewRequest(http.MethodGet, "/users/emp_001/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
	var user model.User
	json.Unmarshal(w.Body.Bytes(), &user)
	if user.SID != "emp_001" {
		t.Errorf("expected SID emp_001, got %s", user.SID)
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userRepo := mocks.NewMockUserRepo()
	h := NewUserHandler(userRepo, nil, nil, nil)
	r := gin.New()
	r.GET("/users/:sid/profile", h.GetProfile)

	req := httptest.NewRequest(http.MethodGet, "/users/nonexistent/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetStats_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userRepo := mocks.NewMockUserRepo()
	matchRepo := mocks.NewMockMatchRepo()
	tournamentRepo := mocks.NewMockTournamentRepo()
	h := NewUserHandler(userRepo, nil, matchRepo, tournamentRepo)
	r := gin.New()
	r.GET("/users/:sid/stats", h.GetStats)

	req := httptest.NewRequest(http.MethodGet, "/users/emp_001/stats", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestGetStats_WithData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	matchRepo := mocks.NewMockMatchRepo()
	tournamentRepo := mocks.NewMockTournamentRepo()

	matchRepo.Create(context.Background(), &model.Match{
		ID: "m1", GameType: "chess",
		Player1SID: "emp_001", Player2SID: "emp_002",
		Status: "completed", WinnerSID: "emp_001",
	})
	matchRepo.Create(context.Background(), &model.Match{
		ID: "m2", GameType: "chess",
		Player1SID: "emp_001", Player2SID: "emp_002",
		Status: "completed", WinnerSID: "emp_002",
	})
	matchRepo.Create(context.Background(), &model.Match{
		ID: "m3", GameType: "chess",
		Player1SID: "emp_001", Player2SID: "emp_002",
		Status: "completed", WinnerSID: "",
	})
	tournamentRepo.AddPlayer(context.Background(), "t1", "emp_001")

	h := NewUserHandler(nil, nil, matchRepo, tournamentRepo)
	r := gin.New()
	r.GET("/users/:sid/stats", h.GetStats)

	req := httptest.NewRequest(http.MethodGet, "/users/emp_001/stats", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
	var stats map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &stats)
	if int(stats["games_played"].(float64)) != 3 {
		t.Errorf("expected 3 games, got %v", stats["games_played"])
	}
	if int(stats["wins"].(float64)) != 1 {
		t.Errorf("expected 1 win, got %v", stats["wins"])
	}
	if int(stats["draws"].(float64)) != 1 {
		t.Errorf("expected 1 draw, got %v", stats["draws"])
	}
	if int(stats["losses"].(float64)) != 1 {
		t.Errorf("expected 1 loss, got %v", stats["losses"])
	}
	if int(stats["tournaments_joined"].(float64)) != 1 {
		t.Errorf("expected 1 tournament, got %v", stats["tournaments_joined"])
	}
}
