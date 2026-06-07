package service_test

import (
	"testing"

	"game-platform/internal/service"
)

func TestCalculateElo_Player1Wins(t *testing.T) {
	elo1, elo2 := service.CalculateElo(1200, 1200, "player1", 32)
	if elo1 <= 1200 {
		t.Errorf("Expected elo1 > 1200 for winner, got %.1f", elo1)
	}
	if elo2 >= 1200 {
		t.Errorf("Expected elo2 < 1200 for loser, got %.1f", elo2)
	}
	sum := elo1 + elo2
	if sum != 2400 {
		t.Errorf("Elo sum should be 2400, got %.1f", sum)
	}
}

func TestCalculateElo_Player2Wins(t *testing.T) {
	elo1, elo2 := service.CalculateElo(1200, 1200, "player2", 32)
	if elo1 >= 1200 {
		t.Errorf("Expected elo1 < 1200 for loser, got %.1f", elo1)
	}
	if elo2 <= 1200 {
		t.Errorf("Expected elo2 > 1200 for winner, got %.1f", elo2)
	}
}

func TestCalculateElo_Draw(t *testing.T) {
	elo1, elo2 := service.CalculateElo(1500, 1100, "draw", 32)
	if elo1 >= 1500 {
		t.Errorf("Expected elo1 <= 1500 (higher rated draws lower), got %.1f", elo1)
	}
	if elo2 <= 1100 {
		t.Errorf("Expected elo2 >= 1100 (lower rated draws higher), got %.1f", elo2)
	}
}

func TestCalculateElo_HigherRatedWins(t *testing.T) {
	elo1, elo2 := service.CalculateElo(2000, 1000, "player1", 24)
	if elo1 <= 2000 {
		t.Error("Winner should gain points")
	}
	if elo2 >= 1000 {
		t.Error("Loser should lose points")
	}
	gain := elo1 - 2000
	loss := 1000 - elo2
	if gain-loss > 0.01 || loss-gain > 0.01 {
		t.Errorf("Gain (%.1f) should equal loss (%.1f)", gain, loss)
	}
	if gain > float64(24) {
		t.Errorf("Gain (%.1f) should not exceed K-factor (24)", gain)
	}
}

func TestCalculateElo_Upset(t *testing.T) {
	elo1, elo2 := service.CalculateElo(1000, 2000, "player1", 24)
	gain := elo1 - 1000
	if gain < 15 {
		t.Errorf("Expected significant gain for upset, got %.1f", gain)
	}
	loss := 2000 - elo2
	if gain != loss {
		t.Errorf("Gain (%.1f) should equal loss (%.1f)", gain, loss)
	}
}

func TestValidGameType(t *testing.T) {
	valid := []string{"chess", "checkers", "backgammon", "snake", "mines", "arena", "poker"}
	for _, gt := range valid {
		if !service.ValidGameType(gt) {
			t.Errorf("Expected %s to be valid", gt)
		}
	}
}

func TestValidGameType_Invalid(t *testing.T) {
	invalid := []string{"", "soccer", "basketball", "tetris", "CHESS"}
	for _, gt := range invalid {
		if service.ValidGameType(gt) {
			t.Errorf("Expected %s to be invalid", gt)
		}
	}
}
