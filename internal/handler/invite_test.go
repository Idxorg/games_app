package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"game-platform/internal/model"
	"game-platform/tests/mocks"

	"github.com/gin-gonic/gin"
)

// ---------- InviteHandler tests ----------

func TestNewInviteHandler(t *testing.T) {
	h := NewInviteHandler(nil, nil)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestSetUserRepo(t *testing.T) {
	h := NewInviteHandler(nil, nil)
	repo := mocks.NewMockUserRepo()
	h.SetUserRepo(repo)
	if h.userRepo != repo {
		t.Error("expected userRepo to be set")
	}
}

func TestSetNotificationsClient(t *testing.T) {
	h := NewInviteHandler(nil, nil)
	pub := &mocks.MockNotificationPublisher{}
	h.SetNotificationsClient(pub)
	if h.notifier != pub {
		t.Error("expected notifier to be set")
	}
}

func setupInviteTestRouter(ih *InviteHandler, sid string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("sid", sid)
		c.Next()
	})
	r.POST("/invite", ih.CreateInvite)
	r.POST("/invite/:id/accept", ih.AcceptInvite)
	r.POST("/invite/:id/decline", ih.DeclineInvite)
	r.GET("/invite/pending", ih.GetPendingInvites)
	return r
}

func TestCreateInvite_Success(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := NewInviteHandler(inviteRepo, matchRepo)
	router := setupInviteTestRouter(ih, "emp_111")

	body := `{"game_type":"chess","recipient_sid":"emp_222"}`
	req, _ := http.NewRequest("POST", "/invite", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestCreateInvite_InvalidGameType(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := NewInviteHandler(inviteRepo, matchRepo)
	router := setupInviteTestRouter(ih, "emp_111")

	body := `{"game_type":"monopoly","recipient_sid":"emp_222"}`
	req, _ := http.NewRequest("POST", "/invite", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateInvite_MissingFields(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := NewInviteHandler(inviteRepo, matchRepo)
	router := setupInviteTestRouter(ih, "emp_111")

	body := `{"game_type":"chess"}`
	req, _ := http.NewRequest("POST", "/invite", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateInvite_SelfInvite(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := NewInviteHandler(inviteRepo, matchRepo)
	router := setupInviteTestRouter(ih, "emp_111")

	body := `{"game_type":"chess","recipient_sid":"emp_111"}`
	req, _ := http.NewRequest("POST", "/invite", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for self-invite, got %d", w.Code)
	}
}

func TestCreateInvite_NoSID(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := NewInviteHandler(inviteRepo, matchRepo)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/invite", ih.CreateInvite)

	body := `{"game_type":"chess","recipient_sid":"emp_222"}`
	req, _ := http.NewRequest("POST", "/invite", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for missing SID, got %d", w.Code)
	}
}

func TestCreateInvite_WithInviterName(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	userRepo := mocks.NewMockUserRepo()
	userRepo.Create(nil, "emp_111", "e@e.com", "Alice", "f", "eng", "dev", "")

	ih := NewInviteHandler(inviteRepo, matchRepo)
	ih.SetUserRepo(userRepo)

	pub := &mocks.MockNotificationPublisher{}
	ih.SetNotificationsClient(pub)

	router := setupInviteTestRouter(ih, "emp_111")
	body := `{"game_type":"chess","recipient_sid":"emp_222"}`
	req, _ := http.NewRequest("POST", "/invite", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["inviter_name"] != "Alice" {
		t.Errorf("expected inviter_name Alice, got %v", resp["inviter_name"])
	}
	if pub.EventCount() != 1 {
		t.Errorf("expected 1 notification, got %d", pub.EventCount())
	}
}

func TestAcceptInvite_Success(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := NewInviteHandler(inviteRepo, matchRepo)

	invite := &model.GameInvite{
		GameType: "chess", InviterSID: "emp_111", RecipientSID: "emp_222",
		Status: "pending", ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	inviteRepo.Create(nil, invite)

	router := setupInviteTestRouter(ih, "emp_222")
	req, _ := http.NewRequest("POST", "/invite/"+invite.ID+"/accept", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != "invite accepted, match created" {
		t.Errorf("expected success message, got %v", resp["message"])
	}
}

func TestAcceptInvite_NotFound(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := NewInviteHandler(inviteRepo, matchRepo)

	router := setupInviteTestRouter(ih, "emp_222")
	req, _ := http.NewRequest("POST", "/invite/nonexistent/accept", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestAcceptInvite_NotRecipient(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := NewInviteHandler(inviteRepo, matchRepo)

	invite := &model.GameInvite{
		GameType: "chess", InviterSID: "emp_111", RecipientSID: "emp_222",
		Status: "pending", ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	inviteRepo.Create(nil, invite)

	router := setupInviteTestRouter(ih, "emp_333")
	req, _ := http.NewRequest("POST", "/invite/"+invite.ID+"/accept", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestAcceptInvite_AlreadyAccepted(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := NewInviteHandler(inviteRepo, matchRepo)

	invite := &model.GameInvite{
		GameType: "chess", InviterSID: "emp_111", RecipientSID: "emp_222",
		Status: "accepted", ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	inviteRepo.Create(nil, invite)

	router := setupInviteTestRouter(ih, "emp_222")
	req, _ := http.NewRequest("POST", "/invite/"+invite.ID+"/accept", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestAcceptInvite_Expired(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := NewInviteHandler(inviteRepo, matchRepo)

	invite := &model.GameInvite{
		GameType: "chess", InviterSID: "emp_111", RecipientSID: "emp_222",
		Status: "pending", ExpiresAt: time.Now().Add(-1 * time.Minute),
	}
	inviteRepo.Create(nil, invite)

	router := setupInviteTestRouter(ih, "emp_222")
	req, _ := http.NewRequest("POST", "/invite/"+invite.ID+"/accept", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusGone {
		t.Errorf("expected 410 for expired, got %d", w.Code)
	}
}

func TestDeclineInvite_Success(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := NewInviteHandler(inviteRepo, matchRepo)

	invite := &model.GameInvite{
		GameType: "checkers", InviterSID: "emp_111", RecipientSID: "emp_222",
		Status: "pending", ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	inviteRepo.Create(nil, invite)

	router := setupInviteTestRouter(ih, "emp_222")
	req, _ := http.NewRequest("POST", "/invite/"+invite.ID+"/decline", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestDeclineInvite_NotFound(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := NewInviteHandler(inviteRepo, matchRepo)

	router := setupInviteTestRouter(ih, "emp_222")
	req, _ := http.NewRequest("POST", "/invite/nonexistent/decline", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestDeclineInvite_NotRecipient(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := NewInviteHandler(inviteRepo, matchRepo)

	invite := &model.GameInvite{
		GameType: "checkers", InviterSID: "emp_111", RecipientSID: "emp_222",
		Status: "pending", ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	inviteRepo.Create(nil, invite)

	router := setupInviteTestRouter(ih, "emp_333")
	req, _ := http.NewRequest("POST", "/invite/"+invite.ID+"/decline", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestDeclineInvite_NotPending(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := NewInviteHandler(inviteRepo, matchRepo)

	invite := &model.GameInvite{
		GameType: "checkers", InviterSID: "emp_111", RecipientSID: "emp_222",
		Status: "declined", ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	inviteRepo.Create(nil, invite)

	router := setupInviteTestRouter(ih, "emp_222")
	req, _ := http.NewRequest("POST", "/invite/"+invite.ID+"/decline", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestGetPendingInvites(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := NewInviteHandler(inviteRepo, matchRepo)

	inviteRepo.Create(nil, &model.GameInvite{
		GameType: "chess", InviterSID: "emp_111", RecipientSID: "emp_222",
		Status: "pending", ExpiresAt: time.Now().Add(5 * time.Minute),
	})
	inviteRepo.Create(nil, &model.GameInvite{
		GameType: "checkers", InviterSID: "emp_333", RecipientSID: "emp_222",
		Status: "pending", ExpiresAt: time.Now().Add(5 * time.Minute),
	})

	router := setupInviteTestRouter(ih, "emp_222")
	req, _ := http.NewRequest("GET", "/invite/pending", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["count"].(float64)) != 2 {
		t.Errorf("expected 2 invites, got %v", resp["count"])
	}
}

func TestGetPendingInvites_Empty(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	ih := NewInviteHandler(inviteRepo, matchRepo)

	router := setupInviteTestRouter(ih, "emp_111")
	req, _ := http.NewRequest("GET", "/invite/pending", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["count"].(float64)) != 0 {
		t.Errorf("expected 0, got %v", resp["count"])
	}
}

func TestGetPendingInvites_WithInviterNames(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	userRepo := mocks.NewMockUserRepo()
	userRepo.Create(nil, "emp_111", "e@e.com", "Alice", "f", "eng", "dev", "")

	ih := NewInviteHandler(inviteRepo, matchRepo)
	ih.SetUserRepo(userRepo)

	inviteRepo.Create(nil, &model.GameInvite{
		GameType: "chess", InviterSID: "emp_111", RecipientSID: "emp_222",
		Status: "pending", ExpiresAt: time.Now().Add(5 * time.Minute),
	})

	router := setupInviteTestRouter(ih, "emp_222")
	req, _ := http.NewRequest("GET", "/invite/pending", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	invites := resp["invites"].([]interface{})
	if len(invites) == 0 {
		t.Fatal("expected at least one invite")
	}
	first := invites[0].(map[string]interface{})
	if first["inviter_name"] != "Alice" {
		t.Errorf("expected inviter_name Alice, got %v", first["inviter_name"])
	}
}

func TestAcceptInvite_WithNotification(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	pub := &mocks.MockNotificationPublisher{}

	ih := NewInviteHandler(inviteRepo, matchRepo)
	ih.SetNotificationsClient(pub)

	invite := &model.GameInvite{
		GameType: "chess", InviterSID: "emp_111", RecipientSID: "emp_222",
		Status: "pending", ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	inviteRepo.Create(nil, invite)

	router := setupInviteTestRouter(ih, "emp_222")
	req, _ := http.NewRequest("POST", "/invite/"+invite.ID+"/accept", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if pub.EventCount() != 1 {
		t.Errorf("expected 1 notification, got %d", pub.EventCount())
	}
	if pub.Events[0]["type"] != "game.invite.accepted" {
		t.Errorf("expected game.invite.accepted, got %v", pub.Events[0]["type"])
	}
}

func TestDeclineInvite_WithNotification(t *testing.T) {
	inviteRepo := mocks.NewMockInviteRepo()
	matchRepo := mocks.NewMockMatchRepo()
	pub := &mocks.MockNotificationPublisher{}

	ih := NewInviteHandler(inviteRepo, matchRepo)
	ih.SetNotificationsClient(pub)

	invite := &model.GameInvite{
		GameType: "checkers", InviterSID: "emp_111", RecipientSID: "emp_222",
		Status: "pending", ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	inviteRepo.Create(nil, invite)

	router := setupInviteTestRouter(ih, "emp_222")
	req, _ := http.NewRequest("POST", "/invite/"+invite.ID+"/decline", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if pub.EventCount() != 1 {
		t.Errorf("expected 1 notification, got %d", pub.EventCount())
	}
	if pub.Events[0]["type"] != "game.invite.declined" {
		t.Errorf("expected game.invite.declined, got %v", pub.Events[0]["type"])
	}
}
