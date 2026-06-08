package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"

	"game-platform/internal/model"
)

// RatingService handles Elo rating calculations and leaderboard queries.
type RatingService struct {
	ratingRepo model.RatingRepo
	matchRepo  model.MatchRepo
	userRepo   model.UserRepo
}

// NewRatingService creates a new RatingService.
func NewRatingService(
	ratingRepo model.RatingRepo,
	matchRepo model.MatchRepo,
	userRepo model.UserRepo,
) *RatingService {
	return &RatingService{
		ratingRepo: ratingRepo,
		matchRepo:  matchRepo,
		userRepo:   userRepo,
	}
}

// validGameTypes enumerates the supported game types.
var validGameTypes = map[string]bool{
	"chess":      true,
	"checkers":   true,
	"backgammon": true,
	"snake":      true,
	"mines":      true,
	"arena":      true,
	"poker":      true,
}

// ValidGameType checks whether a game type string is recognised.
func ValidGameType(gt string) bool {
	return validGameTypes[gt]
}

// AllGameTypes returns the full list of supported game types.
func AllGameTypes() []string {
	types := make([]string, 0, len(validGameTypes))
	for gt := range validGameTypes {
		types = append(types, gt)
	}
	return types
}

// getKFactor returns the K-factor for a player based on their game count and rating.
//   - 32 for new players (< 30 games)
//   - 24 for established players
//   - 16 for high-rated players (> 2000)
func getKFactor(gamesPlayed, elo int) float64 {
	if elo > 2000 {
		return 16.0
	}
	if gamesPlayed < 30 {
		return 32.0
	}
	return 24.0
}

// CalculateElo computes new Elo ratings for two players given a match outcome.
//
// Parameters:
//   - player1Elo, player2Elo: current Elo ratings (as float64 for precision)
//   - winner: "player1", "player2", or "draw"
//   - kFactor: K-factor override; pass 0 to auto-select per player
//
// Returns the new Elo values for player1 and player2.
//
// Standard Elo formula:
//
//	expected = 1 / (1 + 10^((opponentElo - playerElo)/400))
//	newElo  = oldElo + K * (actual - expected)
func CalculateElo(player1Elo, player2Elo float64, winner string, kFactor ...float64) (newElo1, newElo2 float64) {
	// Determine actual scores
	var actual1, actual2 float64
	switch winner {
	case "player1":
		actual1, actual2 = 1.0, 0.0
	case "player2":
		actual1, actual2 = 0.0, 1.0
	case "draw":
		actual1, actual2 = 0.5, 0.5
	default:
		// Treat unknown as draw
		actual1, actual2 = 0.5, 0.5
	}

	// Expected scores
	expected1 := 1.0 / (1.0+pow10((player2Elo-player1Elo)/400.0))
	expected2 := 1.0 - expected1 // symmetric property

	// Per-player K-factors (allow caller override for both via first variadic)
	k1, k2 := getKFactor(int(player1Elo), int(player1Elo)), getKFactor(int(player2Elo), int(player2Elo))
	if len(kFactor) > 0 && kFactor[0] > 0 {
		k1 = kFactor[0]
		k2 = kFactor[0]
	}

	newElo1 = player1Elo + k1*(actual1-expected1)
	newElo2 = player2Elo + k2*(actual2-expected2)

	return newElo1, newElo2
}

// pow10 computes 10^x.
func pow10(x float64) float64 {
	// Use math.Pow; import below.
	return _pow10(x)
}

func _pow10(x float64) float64 {
	// Implemented without math import to keep imports minimal.
	// 10^x = e^(x*ln(10))
	ln10 := 2.302585092994046
	return exp(x * ln10)
}

func exp(x float64) float64 {
	// Taylor-series expansion of e^x, sufficient for Elo precision.
	sum := 1.0
	term := 1.0
	for i := 1; i < 30; i++ {
		term *= x / float64(i)
		sum += term
	}
	return sum
}

// UpdateMatchRatings processes a completed match and updates both players' Elo
// ratings, win/draw/loss counters, and games-played counts atomically.
func (s *RatingService) UpdateMatchRatings(ctx context.Context, match *model.Match) error {
	if match.WinnerSID == "" && match.Status != "completed" {
		return nil
	}

	p1Sid := match.Player1SID
	p2Sid := match.Player2SID

	// Fetch or initialise ratings for both players.
	rating1, err := s.ratingRepo.Get(ctx, p1Sid, match.GameType)
	if err != nil || rating1 == nil {
		// No rating yet — create fresh entry.
		rating1 = &model.PlayerRating{
			SID:      p1Sid,
			GameType: match.GameType,
			ELO:      1000, // default starting Elo
		}
	}
	rating2, err := s.ratingRepo.Get(ctx, p2Sid, match.GameType)
	if err != nil || rating2 == nil {
		rating2 = &model.PlayerRating{
			SID:      p2Sid,
			GameType: match.GameType,
			ELO:      1000,
		}
	}

	// Determine winner string for Elo calculation.
	winner := "draw"
	if match.WinnerSID == p1Sid {
		winner = "player1"
	} else if match.WinnerSID == p2Sid {
		winner = "player2"
	}

	// Calculate new Elo.
	newElo1, newElo2 := CalculateElo(
		float64(rating1.ELO), float64(rating2.ELO), winner,
	)

	// Update counters.
	rating1.GamesPlayed++
	rating2.GamesPlayed++

	switch winner {
	case "player1":
		rating1.Wins++
		rating2.Losses++
	case "player2":
		rating2.Wins++
		rating1.Losses++
	case "draw":
		rating1.Draws++
		rating2.Draws++
	}

	rating1.ELO = int(roundFloat(newElo1))
	rating2.ELO = int(roundFloat(newElo2))

	// Persist both ratings.
	if err := s.ratingRepo.Upsert(ctx, rating1); err != nil {
		return err
	}
	if err := s.ratingRepo.Upsert(ctx, rating2); err != nil {
		return err
	}

	log.Printf("Ratings updated: %s=%d %s=%d (game=%s)", p1Sid, rating1.ELO, p2Sid, rating2.ELO, match.GameType)
	return nil
}

// GetLeaderboard returns the top-rated players for a game type.
func (s *RatingService) GetLeaderboard(ctx context.Context, gameType string, limit int) ([]model.PlayerRating, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	return s.ratingRepo.GetLeaderboard(ctx, gameType, limit)
}

// GetRating retrieves a single player's rating.
func (s *RatingService) GetRating(ctx context.Context, sid, gameType string) (*model.PlayerRating, error) {
	return s.ratingRepo.Get(ctx, sid, gameType)
}

// GetUserRatings retrieves all ratings for a player across game types.
func (s *RatingService) GetUserRatings(ctx context.Context, sid string) ([]model.PlayerRating, error) {
	ratings := make([]model.PlayerRating, 0)
	for gt := range validGameTypes {
		r, err := s.ratingRepo.Get(ctx, sid, gt)
		if err != nil || r == nil {
			continue // player may not have played this game type
		}
		ratings = append(ratings, *r)
	}
	return ratings, nil
}

// roundFloat rounds a float64 to the nearest integer.
func roundFloat(f float64) float64 {
	if f >= 0 {
		return float64(int(f + 0.5))
	}
	return float64(int(f - 0.5))
}

// GenerateID creates a new UUID v4 string using crypto/rand.
func GenerateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 2
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
