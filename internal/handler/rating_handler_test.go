package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"game-platform/internal/service"
	"game-platform/tests/mocks"

	"github.com/gin-gonic/gin"
)

func setupRatingTestRouter(rh *RatingHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/ratings/:game_type", rh.GetRatings)
	r.GET("/ratings/:game_type/leaderboard/department", rh.GetLeaderboardByDepartment)
	return r
}

func testReq(method, url string) *http.Request {
	req, _ := http.NewRequest(method, url, nil)
	return req
}

func TestNewRatingHandler(t *testing.T) {
	h := NewRatingHandler(nil, nil, nil)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestGetRatings(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	ratingRepo := mocks.NewMockRatingRepo()
	ratingSvc := service.NewRatingService(ratingRepo, matchRepo, nil)
	userRepo := mocks.NewMockUserRepo()
	rh := NewRatingHandler(ratingRepo, ratingSvc, userRepo)
	router := setupRatingTestRouter(rh)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, testReq("GET", "/ratings/chess"))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestGetRatings_InvalidGameType(t *testing.T) {
	rh := NewRatingHandler(nil, nil, nil)
	router := setupRatingTestRouter(rh)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, testReq("GET", "/ratings/monopoly"))

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetRatings_WithLimit(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	ratingRepo := mocks.NewMockRatingRepo()
	ratingSvc := service.NewRatingService(ratingRepo, matchRepo, nil)
	rh := NewRatingHandler(ratingRepo, ratingSvc, nil)
	router := setupRatingTestRouter(rh)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, testReq("GET", "/ratings/chess?limit=10"))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetRatings_LimitCapped(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	ratingRepo := mocks.NewMockRatingRepo()
	ratingSvc := service.NewRatingService(ratingRepo, matchRepo, nil)
	rh := NewRatingHandler(ratingRepo, ratingSvc, nil)
	router := setupRatingTestRouter(rh)

	// Request limit > 200 should be capped
	w := httptest.NewRecorder()
	router.ServeHTTP(w, testReq("GET", "/ratings/chess?limit=500"))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetLeaderboardByDepartment(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	ratingRepo := mocks.NewMockRatingRepo()
	ratingSvc := service.NewRatingService(ratingRepo, matchRepo, nil)
	rh := NewRatingHandler(ratingRepo, ratingSvc, nil)
	router := setupRatingTestRouter(rh)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, testReq("GET", "/ratings/chess/leaderboard/department?department=Engineering"))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestGetLeaderboardByDepartment_InvalidGameType(t *testing.T) {
	rh := NewRatingHandler(nil, nil, nil)
	router := setupRatingTestRouter(rh)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, testReq("GET", "/ratings/monopoly/leaderboard/department?department=Engineering"))

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetLeaderboardByDepartment_MissingDepartment(t *testing.T) {
	matchRepo := mocks.NewMockMatchRepo()
	ratingRepo := mocks.NewMockRatingRepo()
	ratingSvc := service.NewRatingService(ratingRepo, matchRepo, nil)
	rh := NewRatingHandler(ratingRepo, ratingSvc, nil)
	router := setupRatingTestRouter(rh)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, testReq("GET", "/ratings/chess/leaderboard/department"))

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
