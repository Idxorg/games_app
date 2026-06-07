package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"game-platform/internal/model"

	"github.com/gin-gonic/gin"
)

func TestHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	})

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %v", response["status"])
	}
}

func TestTournamentAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("sid", "emp_12345")
		c.Next()
	})
	
	router.GET("/api/v1/tournaments", func(c *gin.Context) {
		tournaments := []model.Tournament{
			{
				ID:         "t_001",
				Name:       "Весенний чемпионат 2026",
				GameType:   "chess",
				Status:     "active",
				MaxPlayers: 128,
			},
		}
		c.JSON(http.StatusOK, tournaments)
	})

	req, _ := http.NewRequest("GET", "/api/v1/tournaments", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var tournaments []model.Tournament
	if err := json.Unmarshal(w.Body.Bytes(), &tournaments); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if len(tournaments) == 0 {
		t.Error("Expected at least one tournament")
	}
}

func TestRatingAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("sid", "emp_12345")
		c.Next()
	})
	
	router.GET("/api/v1/ratings/chess", func(c *gin.Context) {
		ratings := []model.PlayerRating{
			{
				SID:         "emp_11111",
				GameType:    "chess",
				ELO:         1650,
				GamesPlayed: 100,
				Wins:        70,
				Losses:      30,
			},
		}
		c.JSON(http.StatusOK, ratings)
	})

	req, _ := http.NewRequest("GET", "/api/v1/ratings/chess", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var ratings []model.PlayerRating
	if err := json.Unmarshal(w.Body.Bytes(), &ratings); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if len(ratings) == 0 {
		t.Error("Expected at least one rating")
	}
}

func TestMatchAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("sid", "emp_12345")
		c.Next()
	})
	
	router.POST("/api/v1/games/match/start", func(c *gin.Context) {
		var request struct {
			GameType  string `json:"game_type"`
			OpponentID string `json:"opponent_id"`
		}
		
		if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		response := map[string]interface{}{
			"match_id":   "m_001",
			"game_id":    "g_chess_123456",
			"status":     "waiting",
			"opponent":   request.OpponentID,
		}
		c.JSON(http.StatusOK, response)
	})

	body := bytes.NewBufferString(`{"game_type":"chess","opponent_id":"emp_67890"}`)
	req, _ := http.NewRequest("POST", "/api/v1/games/match/start", body)
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
		t.Logf("Response: %s", w.Body.String())
	}
}

func TestUserProfileAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("sid", "emp_12345")
		c.Next()
	})
	
	router.GET("/api/v1/users/:sid/profile", func(c *gin.Context) {
		user := model.User{
			SID:        "emp_12345",
			Email:      "ivanov@yakbson.digital",
			Name:       "Иванов Иван",
			Gender:     "male",
			Department: "IT",
			Position:   "Разработчик",
			PhotoURL:   "https://s3.yakbson.digital/avatars/emp_12345.jpg",
		}
		c.JSON(http.StatusOK, user)
	})

	req, _ := http.NewRequest("GET", "/api/v1/users/emp_12345/profile", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var user model.User
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if user.SID != "emp_12345" {
		t.Errorf("Expected SID emp_12345, got %s", user.SID)
	}

	if user.Department != "IT" {
		t.Errorf("Expected department IT, got %s", user.Department)
	}
}
