package rating

import "fmt"

const (
	RATING_ELO RatingType = "ELO"
)

type RatingType string

type RatingProvider interface {
	GetRatingCalculator() RatingCalculator

	GetRatingEntityKindPrefix() string
}

func RatingProviderFactory(ratingType RatingType) RatingProvider {
	if ratingType == RATING_ELO {
		return EloRatingProvider{}
	}
	panic(fmt.Sprintf("Unknown rating type: %s", ratingType))
}
