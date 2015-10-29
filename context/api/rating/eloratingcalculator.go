package rating

import (
	"api/domain"
	"math"
)

type EloRatingCalculator struct {
	DefaultRating float64
	KFactor       float64
	RatingFactor  float64
	ScorePercent  float64
}

func CreateEloRatingCalculator() *EloRatingCalculator {
	return &EloRatingCalculator{
		DefaultRating: 1000,
		KFactor:       50,  // Max +25 rating points for a win (and -25 for losing - giving the sum of 50 points)
		RatingFactor:  400, // Rating +400 means 10 times as good
		ScorePercent:  50,  // Smallest winning margin will give at least 50% of the K_FACTOR
	}
}

func (erc *EloRatingCalculator) CalculateRating(latestPlayerRatings map[int64]domain.PlayerRating, currentGame domain.Game) domain.GameRating {
	// Find winning and losing team
	var winner, looser domain.GameTeam
	if currentGame.Team1.Score > currentGame.Team2.Score {
		winner = currentGame.Team1
		looser = currentGame.Team2
	} else {
		winner = currentGame.Team2
		looser = currentGame.Team1
	}

	previousWinnerRating := erc.calculateWinnerRating(winner, latestPlayerRatings)
	previousLooserRating := erc.calculateWinnerRating(looser, latestPlayerRatings)

	teamRating := erc.calculateTeamRating(previousWinnerRating, previousLooserRating, winner.Score, looser.Score)

	winnerTeamRating := domain.TeamRating{
		Rating:        teamRating,
		PlayerRatings: erc.createPlayerRatings(winner, latestPlayerRatings, teamRating),
	}

	looserTeamRating := domain.TeamRating{
		Rating:        -teamRating,
		PlayerRatings: erc.createPlayerRatings(looser, latestPlayerRatings, -teamRating),
	}

	return domain.GameRating{
		GameID:            currentGame.ID,
		WinningTeamRating: winnerTeamRating,
		LoosingTeamRating: looserTeamRating,
	}
}

func (erc *EloRatingCalculator) calculateWinnerRating(team domain.GameTeam, latestPlayerRatings map[int64]domain.PlayerRating) float64 {
	totalRating := float64(0)

	for _, playerId := range team.Players {
		previousPlayerRating, found := latestPlayerRatings[playerId]
		if !found {
			previousPlayerRating = domain.PlayerRating{
				PlayerID: playerId,
				Rating:   erc.DefaultRating,
			}
			latestPlayerRatings[playerId] = previousPlayerRating
		}

		totalRating += previousPlayerRating.Rating
	}

	return totalRating / float64(len(team.Players))
}

func (erc *EloRatingCalculator) createPlayerRatings(gameTeam domain.GameTeam, latestPlayerRatings map[int64]domain.PlayerRating, teamRating float64) []domain.PlayerRating {
	playerRatings := make([]domain.PlayerRating, len(gameTeam.Players))

	ratingPerPlayer := teamRating / float64(len(gameTeam.Players))

	for idx, playerId := range gameTeam.Players {
		previousPlayerRating := latestPlayerRatings[playerId]

		playerRating := domain.PlayerRating{
			PlayerID: playerId,
			Rating:   previousPlayerRating.Rating + ratingPerPlayer,
		}

		playerRatings[idx] = playerRating
		latestPlayerRatings[playerId] = playerRating
	}

	return playerRatings
}

// See the formula at http://en.wikipedia.org/wiki/Elo_rating_system#Mathematical_details
func (erc *EloRatingCalculator) calculateTeamRating(winnerRating, loserRating float64, winnerScore, loserScore int32) float64 {
	// Expected win ration (0.50 = 50% chance of winning)
	expected := math.Pow(10, winnerRating/erc.RatingFactor)
	expected = expected / (expected + math.Pow(10, loserRating/erc.RatingFactor))

	var maxRatingPoints float64

	if winnerScore > loserScore {
		// Max rating point that can be earned (Highest win margin will give 100% of K_FACTOR)
		winMargin := float64((float64(winnerScore) - float64(loserScore)) / float64(winnerScore))
		maxRatingPoints = winMargin*erc.KFactor*(erc.ScorePercent/100) + erc.KFactor*(100-erc.ScorePercent)/100
	} else {
		//The game was drawn, give 25% of the max score to the lowest rated team
		switch {
		case winnerRating > loserRating:
			//"winner" has the lowest rating
			maxRatingPoints = erc.KFactor * (100 - erc.ScorePercent) / 100 / 2
		case winnerRating < loserRating:
			//"loser" has the lowest rating
			maxRatingPoints = -erc.KFactor * (100 - erc.ScorePercent) / 100 / 2
		default:
			maxRatingPoints = 0
		}
	}

	return maxRatingPoints * (1 - expected)
}
