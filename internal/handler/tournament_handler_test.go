package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"game-platform/internal/model"
	"game-platform/tests/mocks"

	"github.com/gin-gonic/gin"
)

func setupTournamentTestRouter(th *TournamentHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("sid", "emp_12345")
		c.Next()
	})
	r.GET("/tournaments", th.ListTournaments)
	r.POST("/tournaments", th.CreateTournament)
	r.GET("/tournaments/:id", th.GetTournament)
	r.PUT("/tournaments/:id", th.UpdateTournament)
	r.DELETE("/tournaments/:id", th.DeleteTournament)
	r.POST("/tournaments/:id/join", th.JoinTournament)
	r.POST("/tournaments/:id/leave", th.LeaveTournament)
	r.GET("/tournaments/:id/players", th.GetTournamentPlayers)
	return r
}

func TestNewTournamentHandler(t *testing.T) {
	h := NewTournamentHandler(nil, nil)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestListTournaments(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()

	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "t1", Name: "Chess Tour", GameType: "chess", Status: "active",
	})
	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "t2", Name: "Poker Tour", GameType: "poker", Status: "upcoming",
	})

	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	req, _ := http.NewRequest("GET", "/tournaments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["count"].(float64)) != 2 {
		t.Errorf("expected 2 tournaments, got %v", resp["count"])
	}
}

func TestListTournaments_Empty(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	th := NewTournamentHandler(tournamentRepo, nil)
	router := setupTournamentTestRouter(th)

	req, _ := http.NewRequest("GET", "/tournaments", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestListTournaments_WithFilters(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()

	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "t_a", Name: "Chess", GameType: "chess", Status: "active",
	})
	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "t_b", Name: "Poker", GameType: "poker", Status: "upcoming",
	})

	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	req, _ := http.NewRequest("GET", "/tournaments?game_type=chess", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["count"].(float64)) != 1 {
		t.Errorf("expected 1, got %v", resp["count"])
	}
}

func TestListTournaments_InvalidGameType(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	th := NewTournamentHandler(tournamentRepo, nil)
	router := setupTournamentTestRouter(th)

	req, _ := http.NewRequest("GET", "/tournaments?game_type=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateTournament_Valid(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()
	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	body := `{"name":"Chess Cup","game_type":"chess","max_players":32}`
	req, _ := http.NewRequest("POST", "/tournaments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d; body: %s", w.Code, w.Body.String())
	}
	var created model.Tournament
	json.Unmarshal(w.Body.Bytes(), &created)
	if created.Name != "Chess Cup" {
		t.Errorf("expected Chess Cup, got %s", created.Name)
	}
	if created.Status != "upcoming" {
		t.Errorf("expected upcoming, got %s", created.Status)
	}
	if created.CreatedBy != "emp_12345" {
		t.Errorf("expected emp_12345, got %s", created.CreatedBy)
	}
}

func TestCreateTournament_InvalidGameType(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()
	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	body := `{"name":"Bad","game_type":"monopoly","max_players":10}`
	req, _ := http.NewRequest("POST", "/tournaments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateTournament_MissingFields(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()
	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	body := `{"name":"No Game Type"}`
	req, _ := http.NewRequest("POST", "/tournaments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateTournament_NoSID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tournamentRepo := mocks.NewMockTournamentRepo()
	th := NewTournamentHandler(tournamentRepo, nil)

	r := gin.New()
	r.POST("/tournaments", th.CreateTournament)

	body := `{"name":"Test","game_type":"chess","max_players":10}`
	req, _ := http.NewRequest("POST", "/tournaments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGetTournament_Found(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()

	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "t001", Name: "Spring Chess", GameType: "chess", Status: "upcoming", MaxPlayers: 64,
	})

	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	req, _ := http.NewRequest("GET", "/tournaments/t001", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetTournament_NotFound(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()
	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	req, _ := http.NewRequest("GET", "/tournaments/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestUpdateTournament_Valid(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()

	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "tu1", Name: "Old Name", GameType: "chess", Status: "upcoming", MaxPlayers: 16,
	})

	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	body := `{"name":"New Name"}`
	req, _ := http.NewRequest("PUT", "/tournaments/tu1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
	var updated model.Tournament
	json.Unmarshal(w.Body.Bytes(), &updated)
	if updated.Name != "New Name" {
		t.Errorf("expected New Name, got %s", updated.Name)
	}
}

func TestUpdateTournament_NotFound(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()
	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	body := `{"name":"X"}`
	req, _ := http.NewRequest("PUT", "/tournaments/nonexistent", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestUpdateTournament_InvalidGameType(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()

	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "tu2", Name: "Test", GameType: "chess", Status: "upcoming", MaxPlayers: 16,
	})

	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	body := `{"game_type":"monopoly"}`
	req, _ := http.NewRequest("PUT", "/tournaments/tu2", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdateTournament_InvalidStatus(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()

	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "tu3", Name: "Test", GameType: "chess", Status: "upcoming", MaxPlayers: 16,
	})

	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	body := `{"status":"invalid_status"}`
	req, _ := http.NewRequest("PUT", "/tournaments/tu3", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestDeleteTournament(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()

	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "td1", Name: "Delete Me", GameType: "chess", Status: "upcoming", MaxPlayers: 16,
	})

	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	req, _ := http.NewRequest("DELETE", "/tournaments/td1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestDeleteTournament_NotFound(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()
	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	req, _ := http.NewRequest("DELETE", "/tournaments/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestJoinTournament_Success(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()

	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "tj1", Name: "Open Cup", GameType: "chess",
		Status: "upcoming", MaxPlayers: 10, CurrentPlayers: 0,
	})

	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	req, _ := http.NewRequest("POST", "/tournaments/tj1/join", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestJoinTournament_NotFound(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()
	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	req, _ := http.NewRequest("POST", "/tournaments/nonexistent/join", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestJoinTournament_Full(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()

	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "tf1", Name: "Full Cup", GameType: "chess",
		Status: "upcoming", MaxPlayers: 2, CurrentPlayers: 2,
	})

	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	req, _ := http.NewRequest("POST", "/tournaments/tf1/join", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestJoinTournament_WrongStatus(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()

	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "tw1", Name: "Done Cup", GameType: "chess",
		Status: "completed", MaxPlayers: 10, CurrentPlayers: 0,
	})

	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	req, _ := http.NewRequest("POST", "/tournaments/tw1/join", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestJoinTournament_GroupRestriction(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()

	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "tg1", Name: "Group Cup", GameType: "chess",
		Status: "upcoming", MaxPlayers: 10, CurrentPlayers: 0,
		RequiresGroup: "some_group",
	})

	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	req, _ := http.NewRequest("POST", "/tournaments/tg1/join", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// MockUserRepo.GetUserGroups returns ["games", "tournaments"], not "some_group"
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestJoinTournament_NoSID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tournamentRepo := mocks.NewMockTournamentRepo()
	th := NewTournamentHandler(tournamentRepo, nil)

	r := gin.New()
	r.POST("/tournaments/:id/join", th.JoinTournament)

	req, _ := http.NewRequest("POST", "/tournaments/tg1/join", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLeaveTournament_Success(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()

	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "tl1", Name: "Leave Cup", GameType: "chess",
		Status: "active", MaxPlayers: 10, CurrentPlayers: 5,
	})

	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	req, _ := http.NewRequest("POST", "/tournaments/tl1/leave", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestLeaveTournament_NotFound(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()
	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	req, _ := http.NewRequest("POST", "/tournaments/nonexistent/leave", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestLeaveTournament_Completed(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()

	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "tlc1", Name: "Done", GameType: "chess",
		Status: "completed", MaxPlayers: 10, CurrentPlayers: 5,
	})

	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	req, _ := http.NewRequest("POST", "/tournaments/tlc1/leave", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestLeaveTournament_Cancelled(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()

	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "tlx1", Name: "Cancelled", GameType: "chess",
		Status: "cancelled", MaxPlayers: 10, CurrentPlayers: 5,
	})

	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	req, _ := http.NewRequest("POST", "/tournaments/tlx1/leave", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetTournamentPlayers(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()

	tournamentRepo.Create(context.Background(), &model.Tournament{
		ID: "tp1", Name: "Players Cup", GameType: "chess",
		Status: "upcoming", MaxPlayers: 10, CurrentPlayers: 0,
	})

	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	req, _ := http.NewRequest("GET", "/tournaments/tp1/players", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestGetTournamentPlayers_NotFound(t *testing.T) {
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()
	th := NewTournamentHandler(tournamentRepo, userRepo)
	router := setupTournamentTestRouter(th)

	req, _ := http.NewRequest("GET", "/tournaments/nonexistent/players", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestContainsGroup(t *testing.T) {
	groups := []string{"games", "tournaments", "admin"}
	if !containsGroup(groups, "games") {
		t.Error("expected to find games")
	}
	if !containsGroup(groups, "admin") {
		t.Error("expected to find admin")
	}
	if containsGroup(groups, "nonexistent") {
		t.Error("expected not to find nonexistent")
	}
	if containsGroup(nil, "games") {
		t.Error("expected nil groups to return false")
	}
	if containsGroup([]string{}, "games") {
		t.Error("expected empty groups to return false")
	}
}

func TestCreateTournament_WithDates(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tournamentRepo := mocks.NewMockTournamentRepo()
	userRepo := mocks.NewMockUserRepo()
	th := NewTournamentHandler(tournamentRepo, userRepo)

	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("sid", "emp_12345"); c.Next() })
	r.POST("/tournaments", th.CreateTournament)

	now := time.Now().UTC()
	startDate := now.Add(24 * time.Hour).Format(time.RFC3339)
	endDate := now.Add(72 * time.Hour).Format(time.RFC3339)

	body := `{"name":"Dated Tournament","game_type":"chess","max_players":16,"start_date":"` + startDate + `","end_date":"` + endDate + `"}`
	req, _ := http.NewRequest("POST", "/tournaments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d; body: %s", w.Code, w.Body.String())
	}
}
