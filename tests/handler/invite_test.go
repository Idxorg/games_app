package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"game-platform/internal/handler"
	"game-platform/internal/model"
	"game-platform/tests/mocks"

	"github.com/gin-gonic/gin"
)

func setupInviteRouter(ih *handler.InviteHandler, sid string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("sid", sid)
		c.Next()
	})
	r.POST("/api/v1/games/invite", ih.CreateInvite)
	r.POST("/api/v1/games/invite/:id/accept", ih.AcceptInvite)
	r.POST("/api/v1/games/invite/:id/decline", ih.DeclineInvite)
	r.GET("/api/v1/games/invite/pending", ih.GetPendingInvites)
	return r
}

func TestCreateInvite_Success(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := handler.NewInviteHandler(inviteRepo, matchRepo)
	router := setupInviteRouter(ih, "emp_111")

	body := `{"game_type":"chess","recipient_sid":"emp_222"}`
	req, _ := http.NewRequest("POST", "/api/v1/games/invite", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d; body: %s", w.Code, w.Body.String())
	}

	var invite model.GameInvite
	json.Unmarshal(w.Body.Bytes(), &invite)
	if invite.GameType != "chess" {
		t.Errorf("Expected game_type chess, got %s", invite.GameType)
	}
	if invite.InviterSID != "emp_111" {
		t.Errorf("Expected inviter emp_111, got %s", invite.InviterSID)
	}
	if invite.RecipientSID != "emp_222" {
		t.Errorf("Expected recipient emp_222, got %s", invite.RecipientSID)
	}
	if invite.Status != "pending" {
		t.Errorf("Expected status pending, got %s", invite.Status)
	}
	if invite.ID == "" {
		t.Error("Expected invite ID to be set")
	}

	// Verify inviter_name field exists (will be empty without userRepo)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if _, ok := resp["inviter_name"]; !ok {
		t.Error("Expected inviter_name field in response")
	}
}

func TestCreateInvite_InvalidGameType(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := handler.NewInviteHandler(inviteRepo, matchRepo)
	router := setupInviteRouter(ih, "emp_111")

	body := `{"game_type":"monopoly","recipient_sid":"emp_222"}`
	req, _ := http.NewRequest("POST", "/api/v1/games/invite", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestAcceptInvite_Success(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := handler.NewInviteHandler(inviteRepo, matchRepo)

	// Create an invite first — use repo directly
	invite := &model.GameInvite{
		GameType:     "chess",
		InviterSID:   "emp_111",
		RecipientSID: "emp_222",
		Status:       "pending",
		ExpiresAt:    time.Now().Add(5 * time.Minute),
	}
	inviteRepo.Create(nil, invite)

	// Router as the recipient (emp_222)
	router := setupInviteRouter(ih, "emp_222")

	req, _ := http.NewRequest("POST", "/api/v1/games/invite/"+invite.ID+"/accept", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	// Should have invite and match in response
	if resp["message"] != "invite accepted, match created" {
		t.Errorf("Expected success message, got %v", resp["message"])
	}

	// Verify invite was updated
	updated, _ := inviteRepo.GetByID(nil, invite.ID)
	if updated.Status != "accepted" {
		t.Errorf("Expected invite status accepted, got %s", updated.Status)
	}
	if updated.MatchID == nil {
		t.Error("Expected match_id to be set on invite")
	}
}

func TestDeclineInvite_Success(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := handler.NewInviteHandler(inviteRepo, matchRepo)

	invite := &model.GameInvite{
		GameType:     "checkers",
		InviterSID:   "emp_111",
		RecipientSID: "emp_222",
		Status:       "pending",
		ExpiresAt:    time.Now().Add(5 * time.Minute),
	}
	inviteRepo.Create(nil, invite)

	router := setupInviteRouter(ih, "emp_222")

	req, _ := http.NewRequest("POST", "/api/v1/games/invite/"+invite.ID+"/decline", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != "invite declined" {
		t.Errorf("Expected decline message, got %v", resp["message"])
	}

	updated, _ := inviteRepo.GetByID(nil, invite.ID)
	if updated.Status != "declined" {
		t.Errorf("Expected status declined, got %s", updated.Status)
	}
}

func TestAcceptInvite_AlreadyAccepted(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := handler.NewInviteHandler(inviteRepo, matchRepo)

	invite := &model.GameInvite{
		GameType:     "chess",
		InviterSID:   "emp_111",
		RecipientSID: "emp_222",
		Status:       "accepted",
		MatchID:      ptrString("match_001"),
		ExpiresAt:    time.Now().Add(5 * time.Minute),
	}
	inviteRepo.Create(nil, invite)

	router := setupInviteRouter(ih, "emp_222")

	req, _ := http.NewRequest("POST", "/api/v1/games/invite/"+invite.ID+"/accept", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected 409 for already accepted, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestGetPendingInvites(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := handler.NewInviteHandler(inviteRepo, matchRepo)

	// Create two pending invites for emp_222
	inviteRepo.Create(nil, &model.GameInvite{
		GameType:     "chess", InviterSID: "emp_111", RecipientSID: "emp_222",
		Status: "pending", ExpiresAt: time.Now().Add(5 * time.Minute),
	})
	inviteRepo.Create(nil, &model.GameInvite{
		GameType:     "checkers", InviterSID: "emp_333", RecipientSID: "emp_222",
		Status: "pending", ExpiresAt: time.Now().Add(5 * time.Minute),
	})
	// One invite for a different user
	inviteRepo.Create(nil, &model.GameInvite{
		GameType:     "backgammon", InviterSID: "emp_111", RecipientSID: "emp_444",
		Status: "pending", ExpiresAt: time.Now().Add(5 * time.Minute),
	})

	router := setupInviteRouter(ih, "emp_222")

	req, _ := http.NewRequest("GET", "/api/v1/games/invite/pending", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	count := int(resp["count"].(float64))
	if count != 2 {
		t.Errorf("Expected 2 pending invites, got %d", count)
	}
	// Verify invites have inviter_name field
	invites := resp["invites"].([]interface{})
	if len(invites) == 0 {
		t.Fatal("expected non-empty invites")
	}
	first := invites[0].(map[string]interface{})
	if _, ok := first["inviter_name"]; !ok {
		t.Error("Expected inviter_name field in pending invite items")
	}
}

func TestGetPendingInvites_Empty(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := handler.NewInviteHandler(inviteRepo, matchRepo)
	router := setupInviteRouter(ih, "emp_111")

	req, _ := http.NewRequest("GET", "/api/v1/games/invite/pending", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	count := int(resp["count"].(float64))
	if count != 0 {
		t.Errorf("Expected 0 pending invites, got %d", count)
	}

	invites := resp["invites"].([]interface{})
	if len(invites) != 0 {
		t.Errorf("Expected empty invites array, got %d items", len(invites))
	}
}

func TestCreateInvite_WithInviterName(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	userRepo := mocks.NewMockUserRepo()
	userRepo.Create(nil, "emp_111", "emp_111@test.com", "Alice Smith", "female", "eng", "dev", "")

	ih := handler.NewInviteHandler(inviteRepo, matchRepo)
	ih.SetUserRepo(userRepo)

	pub := &mocks.MockNotificationPublisher{}
	ih.SetNotificationsClient(pub)

	router := setupInviteRouter(ih, "emp_111")
	body := `{"game_type":"chess","recipient_sid":"emp_222"}`
	req, _ := http.NewRequest("POST", "/api/v1/games/invite", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected 201, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["inviter_name"] != "Alice Smith" {
		t.Errorf("Expected inviter_name 'Alice Smith', got %v", resp["inviter_name"])
	}

	// Verify notification published
	if pub.EventCount() != 1 {
		t.Fatalf("Expected 1 notification event, got %d", pub.EventCount())
	}
	evt := pub.Events[0]
	if evt["type"] != "game.invite" {
		t.Errorf("Expected notification type game.invite, got %v", evt["type"])
	}
}

func TestAcceptInvite_PublishesNotification(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	pub := &mocks.MockNotificationPublisher{}

	ih := handler.NewInviteHandler(inviteRepo, matchRepo)
	ih.SetNotificationsClient(pub)

	invite := &model.GameInvite{
		GameType:     "chess", InviterSID: "emp_111", RecipientSID: "emp_222",
		Status: "pending", ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	inviteRepo.Create(nil, invite)

	router := setupInviteRouter(ih, "emp_222")
	req, _ := http.NewRequest("POST", "/api/v1/games/invite/"+invite.ID+"/accept", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
	if pub.EventCount() != 1 {
		t.Fatalf("Expected 1 notification, got %d", pub.EventCount())
	}
	if pub.Events[0]["type"] != "game.invite.accepted" {
		t.Errorf("Expected game.invite.accepted, got %v", pub.Events[0]["type"])
	}
}

func TestDeclineInvite_PublishesNotification(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	pub := &mocks.MockNotificationPublisher{}

	ih := handler.NewInviteHandler(inviteRepo, matchRepo)
	ih.SetNotificationsClient(pub)

	invite := &model.GameInvite{
		GameType:     "checkers", InviterSID: "emp_111", RecipientSID: "emp_222",
		Status: "pending", ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	inviteRepo.Create(nil, invite)

	router := setupInviteRouter(ih, "emp_222")
	req, _ := http.NewRequest("POST", "/api/v1/games/invite/"+invite.ID+"/decline", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
	if pub.EventCount() != 1 {
		t.Fatalf("Expected 1 notification, got %d", pub.EventCount())
	}
	if pub.Events[0]["type"] != "game.invite.declined" {
		t.Errorf("Expected game.invite.declined, got %v", pub.Events[0]["type"])
	}
}

// Helper
func ptrString(s string) *string {
	return &s
}
