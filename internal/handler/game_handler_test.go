package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"game-platform/internal/model"
	"game-platform/internal/service"
	"game-platform/tests/mocks"

	"github.com/gin-gonic/gin"
)

func setupGameTestRouter(gh *GameHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("sid", "emp_12345")
		c.Next()
	})
	r.GET("/games/available", gh.AvailableGames)
	r.POST("/games/match/start", gh.StartMatch)
	r.POST("/games/match/complete", gh.CompleteMatch)
	r.GET("/games/match/history", gh.GetMatchHistory)
	r.GET("/ratings/:game_type/leaderboard", gh.GetLeaderboard)
	r.GET("/ratings/:game_type/me", gh.GetUserRating)
	return r
}

func TestNewGameHandler(t *testing.T) {
	h := NewGameHandler(nil, nil)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestAvailableGames(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	gh := NewGameHandler(matchRepo, nil)
	router := setupGameTestRouter(gh)

	req, _ := http.NewRequest("GET", "/games/available", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["count"].(float64)) != 7 {
		t.Errorf("expected 7 games, got %v", resp["count"])
	}
}

func TestStartMatch_Valid(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	gh := NewGameHandler(matchRepo, nil)
	router := setupGameTestRouter(gh)

	body := `{"game_type":"chess","player2_sid":"emp_67890"}`
	req, _ := http.NewRequest("POST", "/games/match/start", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d; body: %s", w.Code, w.Body.String())
	}
	var match model.Match
	json.Unmarshal(w.Body.Bytes(), &match)
	if match.GameType != "chess" {
		t.Errorf("expected chess, got %s", match.GameType)
	}
	if match.Player1SID != "emp_12345" {
		t.Errorf("expected emp_12345, got %s", match.Player1SID)
	}
	if match.Status != "in_progress" {
		t.Errorf("expected in_progress, got %s", match.Status)
	}
}

func TestStartMatch_SelfPlay(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	gh := NewGameHandler(matchRepo, nil)
	router := setupGameTestRouter(gh)

	body := `{"game_type":"chess","player2_sid":"emp_12345"}`
	req, _ := http.NewRequest("POST", "/games/match/start", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestStartMatch_InvalidGameType(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	gh := NewGameHandler(matchRepo, nil)
	router := setupGameTestRouter(gh)

	body := `{"game_type":"monopoly","player2_sid":"emp_67890"}`
	req, _ := http.NewRequest("POST", "/games/match/start", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestStartMatch_MissingFields(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	gh := NewGameHandler(matchRepo, nil)
	router := setupGameTestRouter(gh)

	body := `{"game_type":"chess"}`
	req, _ := http.NewRequest("POST", "/games/match/start", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCompleteMatch_Valid(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	ratingRepo := mocks.NewMockRatingRepo()
	ratingSvc := service.NewRatingService(ratingRepo, matchRepo, nil)
	gh := NewGameHandler(matchRepo, ratingSvc)
	router := setupGameTestRouter(gh)

	matchRepo.Create(context.Background(), &model.Match{
		ID: "m1", GameType: "chess",
		Player1SID: "emp_12345", Player2SID: "emp_67890",
		Status: "in_progress",
	})

	body := `{"match_id":"m1","winner_sid":"emp_12345","score":"1-0"}`
	req, _ := http.NewRequest("POST", "/games/match/complete", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestCompleteMatch_NotFound(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	gh := NewGameHandler(matchRepo, nil)
	router := setupGameTestRouter(gh)

	body := `{"match_id":"nonexistent","winner_sid":"emp_12345","score":"1-0"}`
	req, _ := http.NewRequest("POST", "/games/match/complete", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestCompleteMatch_AlreadyCompleted(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	gh := NewGameHandler(matchRepo, nil)
	router := setupGameTestRouter(gh)

	matchRepo.Create(context.Background(), &model.Match{
		ID: "m2", GameType: "chess",
		Player1SID: "emp_12345", Player2SID: "emp_67890",
		Status: "completed", WinnerSID: "emp_12345",
	})

	body := `{"match_id":"m2","winner_sid":"emp_12345","score":"1-0"}`
	req, _ := http.NewRequest("POST", "/games/match/complete", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCompleteMatch_NotParticipant(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	gh := NewGameHandler(matchRepo, nil)
	router := setupGameTestRouter(gh)

	matchRepo.Create(context.Background(), &model.Match{
		ID: "m3", GameType: "chess",
		Player1SID: "emp_11111", Player2SID: "emp_22222",
		Status: "in_progress",
	})

	body := `{"match_id":"m3","winner_sid":"emp_11111","score":"1-0"}`
	req, _ := http.NewRequest("POST", "/games/match/complete", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestCompleteMatch_InvalidWinner(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	gh := NewGameHandler(matchRepo, nil)
	router := setupGameTestRouter(gh)

	matchRepo.Create(context.Background(), &model.Match{
		ID: "m4", GameType: "chess",
		Player1SID: "emp_12345", Player2SID: "emp_67890",
		Status: "in_progress",
	})

	body := `{"match_id":"m4","winner_sid":"emp_99999","score":"1-0"}`
	req, _ := http.NewRequest("POST", "/games/match/complete", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetMatchHistory_Empty(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	gh := NewGameHandler(matchRepo, nil)
	router := setupGameTestRouter(gh)

	req, _ := http.NewRequest("GET", "/games/match/history", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetMatchHistory_Filtered(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	gh := NewGameHandler(matchRepo, nil)
	router := setupGameTestRouter(gh)

	matchRepo.Create(context.Background(), &model.Match{
		ID: "mh1", GameType: "chess",
		Player1SID: "emp_12345", Player2SID: "emp_67890",
		Status: "completed",
	})
	matchRepo.Create(context.Background(), &model.Match{
		ID: "mh2", GameType: "poker",
		Player1SID: "emp_12345", Player2SID: "emp_67890",
		Status: "completed",
	})

	req, _ := http.NewRequest("GET", "/games/match/history?game_type=chess", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["count"].(float64)) != 1 {
		t.Errorf("expected 1 chess match, got %v", resp["count"])
	}
}

func TestGetLeaderboard(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	ratingRepo := mocks.NewMockRatingRepo()
	ratingSvc := service.NewRatingService(ratingRepo, matchRepo, nil)
	gh := NewGameHandler(matchRepo, ratingSvc)
	router := setupGameTestRouter(gh)

	req, _ := http.NewRequest("GET", "/ratings/chess/leaderboard", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["game_type"] != "chess" {
		t.Errorf("expected chess, got %v", resp["game_type"])
	}
}

func TestGetLeaderboard_InvalidGameType(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	gh := NewGameHandler(matchRepo, nil)
	router := setupGameTestRouter(gh)

	req, _ := http.NewRequest("GET", "/ratings/monopoly/leaderboard", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetLeaderboard_WithLimit(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	ratingRepo := mocks.NewMockRatingRepo()
	ratingSvc := service.NewRatingService(ratingRepo, matchRepo, nil)
	gh := NewGameHandler(matchRepo, ratingSvc)
	router := setupGameTestRouter(gh)

	req, _ := http.NewRequest("GET", "/ratings/chess/leaderboard?limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetUserRating(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	ratingRepo := mocks.NewMockRatingRepo()
	ratingSvc := service.NewRatingService(ratingRepo, matchRepo, nil)
	gh := NewGameHandler(matchRepo, ratingSvc)
	router := setupGameTestRouter(gh)

	req, _ := http.NewRequest("GET", "/ratings/chess/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Rating not found returns 200 with null
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for no rating, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestGetUserRating_InvalidGameType(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	ratingRepo := mocks.NewMockRatingRepo()
	ratingSvc := service.NewRatingService(ratingRepo, matchRepo, nil)
	gh := NewGameHandler(matchRepo, ratingSvc)
	router := setupGameTestRouter(gh)

	req, _ := http.NewRequest("GET", "/ratings/monopoly/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
