package rating

import "api/domain"

type RatingCalculator interface {
	CalculateRating(latestPlayerRatings map[int64]domain.PlayerRating, currentGame domain.Game) domain.GameRating
}
