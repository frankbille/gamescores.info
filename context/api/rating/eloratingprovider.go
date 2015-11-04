package rating

type EloRatingProvider struct {
}

func (erp EloRatingProvider) GetRatingCalculator() RatingCalculator {
	return CreateEloRatingCalculator()
}

func (erp EloRatingProvider) GetRatingEntityKindPrefix() string {
	return "Elo"
}
