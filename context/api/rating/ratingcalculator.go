package rating

import "api/domain"

type RatingCalculator interface {
	CalculateRating(previousGameRating domain.GameRating, previousPlayerRatings []domain.PlayerRating, currentGame domain.Game) domain.GameRating
}
