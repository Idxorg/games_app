package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"game-platform/internal/handler"
	"game-platform/internal/model"
	"game-platform/internal/service"
	"game-platform/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const testJWTSecret = "test-secret-key-for-unit-tests-32chars!"

// helper to create a valid JWT token for testing
func makeTestToken(secret, sid, email, name string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sid":   sid,
		"email": email,
		"name":  name,
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	})
	s, _ := token.SignedString([]byte(secret))
	return s
}

// helper to create an expired JWT token
func makeExpiredToken(secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sid":  "emp_12345",
		"exp":  time.Now().Add(-1 * time.Hour).Unix(),
	})
	s, _ := token.SignedString([]byte(secret))
	return s
}

func TestAuthHandler_VerifyToken_NoSecret(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userRepo := mocks.NewMockUserRepo()
	h := handler.NewAuthHandler(userRepo, nil)
	router := gin.New()
	router.POST("/api/v1/auth/verify", h.VerifyToken)

	// Without secret configured, any token is accepted (dev mode)
	req, _ := http.NewRequest("POST", "/api/v1/auth/verify", nil)
	req.Header.Set("Authorization", "Bearer any-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 in dev mode, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["valid"] != true {
		t.Errorf("Expected valid=true, got %v", resp["valid"])
	}
}

func TestAuthHandler_VerifyToken_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userRepo := mocks.NewMockUserRepo()
	portalAPI := service.NewPortalAPI("", "") // empty base URL = skip portal verify
	h := handler.NewAuthHandler(userRepo, portalAPI)
	h.SetJWTSecret(testJWTSecret)

	router := gin.New()
	router.POST("/api/v1/auth/verify", h.VerifyToken)

	token := makeTestToken(testJWTSecret, "emp_12345", "test@test.com", "Test User")
	req, _ := http.NewRequest("POST", "/api/v1/auth/verify", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["valid"] != true {
		t.Errorf("Expected valid=true, got %v", resp["valid"])
	}
	if resp["sid"] != "emp_12345" {
		t.Errorf("Expected sid=emp_12345, got %v", resp["sid"])
	}
	if resp["email"] != "test@test.com" {
		t.Errorf("Expected email=test@test.com, got %v", resp["email"])
	}
}

func TestAuthHandler_VerifyToken_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userRepo := mocks.NewMockUserRepo()
	h := handler.NewAuthHandler(userRepo, nil)
	h.SetJWTSecret(testJWTSecret)

	router := gin.New()
	router.POST("/api/v1/auth/verify", h.VerifyToken)

	req, _ := http.NewRequest("POST", "/api/v1/auth/verify", nil)
	req.Header.Set("Authorization", "Bearer invalid-garbage-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401, got %d", w.Code)
	}
}

func TestAuthHandler_VerifyToken_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userRepo := mocks.NewMockUserRepo()
	h := handler.NewAuthHandler(userRepo, nil)
	h.SetJWTSecret(testJWTSecret)

	router := gin.New()
	router.POST("/api/v1/auth/verify", h.VerifyToken)

	token := makeExpiredToken(testJWTSecret)
	req, _ := http.NewRequest("POST", "/api/v1/auth/verify", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for expired token, got %d", w.Code)
	}
}

func TestAuthHandler_VerifyToken_NoAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userRepo := mocks.NewMockUserRepo()
	h := handler.NewAuthHandler(userRepo, nil)
	h.SetJWTSecret(testJWTSecret)

	router := gin.New()
	router.POST("/api/v1/auth/verify", h.VerifyToken)

	req, _ := http.NewRequest("POST", "/api/v1/auth/verify", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for missing header, got %d", w.Code)
	}
}

func TestUserHandler_GetProfile_Found(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userRepo := mocks.NewMockUserRepo()
	userRepo.Create(context.Background(), "emp_12345", "test@test.com", "Test User", "male", "IT", "Dev", "")

	h := handler.NewUserHandler(userRepo, nil)
	router := gin.New()
	router.GET("/api/v1/users/:sid/profile", h.GetProfile)

	req, _ := http.NewRequest("GET", "/api/v1/users/emp_12345/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var user model.User
	json.Unmarshal(w.Body.Bytes(), &user)
	if user.SID != "emp_12345" {
		t.Errorf("Expected SID emp_12345, got %s", user.SID)
	}
	if user.Name != "Test User" {
		t.Errorf("Expected Name Test User, got %s", user.Name)
	}
}

func TestUserHandler_GetProfile_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userRepo := mocks.NewMockUserRepo()

	h := handler.NewUserHandler(userRepo, nil)
	router := gin.New()
	router.GET("/api/v1/users/:sid/profile", h.GetProfile)

	req, _ := http.NewRequest("GET", "/api/v1/users/nonexistent/profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}

func TestUserHandler_GetStats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userRepo := mocks.NewMockUserRepo()

	h := handler.NewUserHandler(userRepo, nil)
	router := gin.New()
	router.GET("/api/v1/users/:sid/stats", h.GetStats)

	req, _ := http.NewRequest("GET", "/api/v1/users/emp_12345/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var stats map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &stats)
	if stats["sid"] != "emp_12345" {
		t.Errorf("Expected sid emp_12345, got %v", stats["sid"])
	}
}

// ---------- TournamentHandler tests ----------

func setupTournamentRouter(th *handler.TournamentHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("sid", "emp_12345")
		c.Next()
	})
	r.GET("/api/v1/tournaments", th.ListTournaments)
	r.POST("/api/v1/tournaments", th.CreateTournament)
	r.GET("/api/v1/tournaments/:id", th.GetTournament)
	r.PUT("/api/v1/tournaments/:id", th.UpdateTournament)
	r.DELETE("/api/v1/tournaments/:id", th.DeleteTournament)
	r.POST("/api/v1/tournaments/:id/join", th.JoinTournament)
	r.POST("/api/v1/tournaments/:id/leave", th.LeaveTournament)
	r.GET("/api/v1/tournaments/:id/players", th.GetTournamentPlayers)
	return r
}

func TestCreateTournament_Valid(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()
	th := handler.NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentRouter(th)

	body := `{"name":"Chess Cup","game_type":"chess","max_players":32}`
	req, _ := http.NewRequest("POST", "/api/v1/tournaments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d; body: %s", w.Code, w.Body.String())
	}

	var created model.Tournament
	json.Unmarshal(w.Body.Bytes(), &created)
	if created.Name != "Chess Cup" {
		t.Errorf("Expected name 'Chess Cup', got %s", created.Name)
	}
	if created.GameType != "chess" {
		t.Errorf("Expected game_type chess, got %s", created.GameType)
	}
	if created.Status != "upcoming" {
		t.Errorf("Expected status upcoming, got %s", created.Status)
	}
	if created.CreatedBy != "emp_12345" {
		t.Errorf("Expected created_by emp_12345, got %s", created.CreatedBy)
	}
}

func TestCreateTournament_InvalidGameType(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()
	th := handler.NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentRouter(th)

	body := `{"name":"Bad Game","game_type":"monopoly","max_players":10}`
	req, _ := http.NewRequest("POST", "/api/v1/tournaments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestCreateTournament_MissingFields(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()
	th := handler.NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentRouter(th)

	body := `{"name":"No Game Type"}`
	req, _ := http.NewRequest("POST", "/api/v1/tournaments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for missing fields, got %d", w.Code)
	}
}

func TestGetTournament_Existing(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()

	// Pre-populate
	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "t_001", Name: "Spring Chess", GameType: "chess", Status: "upcoming", MaxPlayers: 64,
	})

	th := handler.NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentRouter(th)

	req, _ := http.NewRequest("GET", "/api/v1/tournaments/t_001", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var tourney model.Tournament
	json.Unmarshal(w.Body.Bytes(), &tourney)
	if tourney.Name != "Spring Chess" {
		t.Errorf("Expected 'Spring Chess', got %s", tourney.Name)
	}
}

func TestGetTournament_NotFound(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()
	th := handler.NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentRouter(th)

	req, _ := http.NewRequest("GET", "/api/v1/tournaments/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}

func TestJoinTournament_Success(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()

	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "t_join1", Name: "Open Cup", GameType: "chess",
		Status: "upcoming", MaxPlayers: 10, CurrentPlayers: 0,
	})

	th := handler.NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentRouter(th)

	req, _ := http.NewRequest("POST", "/api/v1/tournaments/t_join1/join", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestJoinTournament_Full(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()

	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "t_full", Name: "Full Cup", GameType: "chess",
		Status: "upcoming", MaxPlayers: 2, CurrentPlayers: 2,
	})

	th := handler.NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentRouter(th)

	req, _ := http.NewRequest("POST", "/api/v1/tournaments/t_full/join", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for full tournament, got %d", w.Code)
	}
}

func TestJoinTournament_GroupRestriction(t *testing.T) {
	t.Skip("Requires mock with configurable groups — test against real repo with integration test")
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()


	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "t_grp", Name: "Group Cup", GameType: "chess",
		Status: "upcoming", MaxPlayers: 10, CurrentPlayers: 0, RequiresGroup: "tournaments",
	})

	th := handler.NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentRouter(th)

	req, _ := http.NewRequest("POST", "/api/v1/tournaments/t_grp/join", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected 403 for missing group, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestListTournaments_WithFilters(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()

	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "t_a", Name: "Chess Tour", GameType: "chess", Status: "active", CreatedBy: "emp_12345",
	})
	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "t_b", Name: "Poker Tour", GameType: "poker", Status: "upcoming", CreatedBy: "emp_12345",
	})

	th := handler.NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentRouter(th)

	// Filter by game_type=chess
	req, _ := http.NewRequest("GET", "/api/v1/tournaments?game_type=chess", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	count := int(resp["count"].(float64))
	if count != 1 {
		t.Errorf("Expected 1 chess tournament, got %d", count)
	}
}

func TestListTournaments_InvalidGameType(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()
	th := handler.NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentRouter(th)

	req, _ := http.NewRequest("GET", "/api/v1/tournaments?game_type=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

// ---------- GameHandler tests ----------

func setupGameRouter(gh *handler.GameHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("sid", "emp_12345")
		c.Next()
	})
	r.GET("/api/v1/games/available", gh.AvailableGames)
	r.POST("/api/v1/games/match/start", gh.StartMatch)
	r.POST("/api/v1/games/match/complete", gh.CompleteMatch)
	r.GET("/api/v1/games/match/history", gh.GetMatchHistory)
	return r
}

func TestAvailableGames(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	gh := handler.NewGameHandler(matchRepo, nil)
	router := setupGameRouter(gh)

	req, _ := http.NewRequest("GET", "/api/v1/games/available", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	count := int(resp["count"].(float64))
	if count != 7 {
		t.Errorf("Expected 7 games, got %d", count)
	}

	games := resp["games"].([]interface{})
	if len(games) != 7 {
		t.Errorf("Expected 7 games in list, got %d", len(games))
	}
}

func TestStartMatch_Valid(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	gh := handler.NewGameHandler(matchRepo, nil)
	router := setupGameRouter(gh)

	body := `{"game_type":"chess","player2_sid":"emp_67890"}`
	req, _ := http.NewRequest("POST", "/api/v1/games/match/start", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d; body: %s", w.Code, w.Body.String())
	}

	var match model.Match
	json.Unmarshal(w.Body.Bytes(), &match)
	if match.GameType != "chess" {
		t.Errorf("Expected game_type chess, got %s", match.GameType)
	}
	if match.Player1SID != "emp_12345" {
		t.Errorf("Expected player1 emp_12345, got %s", match.Player1SID)
	}
	if match.Player2SID != "emp_67890" {
		t.Errorf("Expected player2 emp_67890, got %s", match.Player2SID)
	}
	if match.Status != "in_progress" {
		t.Errorf("Expected status in_progress, got %s", match.Status)
	}
}

func TestStartMatch_SelfPlay(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	gh := handler.NewGameHandler(matchRepo, nil)
	router := setupGameRouter(gh)

	body := `{"game_type":"chess","player2_sid":"emp_12345"}`
	req, _ := http.NewRequest("POST", "/api/v1/games/match/start", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for self-play, got %d", w.Code)
	}
}

func TestStartMatch_InvalidGameType(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	gh := handler.NewGameHandler(matchRepo, nil)
	router := setupGameRouter(gh)

	body := `{"game_type":"monopoly","player2_sid":"emp_67890"}`
	req, _ := http.NewRequest("POST", "/api/v1/games/match/start", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid game_type, got %d", w.Code)
	}
}

func TestCompleteMatch_Valid(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	ratingRepo := mocks.NewMockRatingRepo()
	ratingSvc := service.NewRatingService(ratingRepo, matchRepo, nil)
	gh := handler.NewGameHandler(matchRepo, ratingSvc)
	router := setupGameRouter(gh)

	// Create a match first
	matchRepo.Create(context.Background(), &model.Match{
		ID: "m_comp1", GameType: "chess",
		Player1SID: "emp_12345", Player2SID: "emp_67890",
		Status: "in_progress",
	})

	body := `{"match_id":"m_comp1","winner_sid":"emp_12345","score":"1-0"}`
	req, _ := http.NewRequest("POST", "/api/v1/games/match/complete", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d; body: %s", w.Code, w.Body.String())
	}

	// Verify match was completed
	m, _ := matchRepo.GetByID(context.Background(), "m_comp1")
	if m.Status != "completed" {
		t.Errorf("Expected status completed, got %s", m.Status)
	}
	if m.WinnerSID != "emp_12345" {
		t.Errorf("Expected winner emp_12345, got %s", m.WinnerSID)
	}
}

func TestCompleteMatch_NotParticipant(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	gh := handler.NewGameHandler(matchRepo, nil)
	router := setupGameRouter(gh)

	matchRepo.Create(context.Background(), &model.Match{
		ID: "m_comp2", GameType: "chess",
		Player1SID: "emp_11111", Player2SID: "emp_22222",
		Status: "in_progress",
	})

	body := `{"match_id":"m_comp2","winner_sid":"emp_11111","score":"1-0"}`
	req, _ := http.NewRequest("POST", "/api/v1/games/match/complete", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected 403 for non-participant, got %d", w.Code)
	}
}

func TestGetMatchHistory_Empty(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	gh := handler.NewGameHandler(matchRepo, nil)
	router := setupGameRouter(gh)

	req, _ := http.NewRequest("GET", "/api/v1/games/match/history", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	count := int(resp["count"].(float64))
	if count != 0 {
		t.Errorf("Expected 0 matches, got %d", count)
	}
}

func TestGetMatchHistory_WithData(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	gh := handler.NewGameHandler(matchRepo, nil)
	router := setupGameRouter(gh)

	matchRepo.Create(context.Background(), &model.Match{
		ID: "m_h1", GameType: "chess",
		Player1SID: "emp_12345", Player2SID: "emp_67890",
		Status: "completed", WinnerSID: "emp_12345",
	})
	matchRepo.Create(context.Background(), &model.Match{
		ID: "m_h2", GameType: "poker",
		Player1SID: "emp_99999", Player2SID: "emp_12345",
		Status: "completed",
	})

	req, _ := http.NewRequest("GET", "/api/v1/games/match/history", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	count := int(resp["count"].(float64))
	if count != 2 {
		t.Errorf("Expected 2 matches, got %d", count)
	}
}

func TestGetMatchHistory_FilterByGameType(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	gh := handler.NewGameHandler(matchRepo, nil)
	router := setupGameRouter(gh)

	matchRepo.Create(context.Background(), &model.Match{
		ID: "m_f1", GameType: "chess",
		Player1SID: "emp_12345", Player2SID: "emp_67890",
		Status: "completed",
	})
	matchRepo.Create(context.Background(), &model.Match{
		ID: "m_f2", GameType: "poker",
		Player1SID: "emp_99999", Player2SID: "emp_12345",
		Status: "completed",
	})

	req, _ := http.NewRequest("GET", "/api/v1/games/match/history?game_type=chess", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	count := int(resp["count"].(float64))
	if count != 1 {
		t.Errorf("Expected 1 chess match, got %d", count)
	}
}
