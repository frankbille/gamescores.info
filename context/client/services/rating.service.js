angular.module('GameScoresApp').factory('RatingService', function ($q, $http) {
    return {
        getLeagueResult: function (leagueId) {
            return $http.get('/api/ratings/'+leagueId+'/result').then(function(response) {
                return response.data;
            });
        },

        getGameRatings: function (gameRatingLink) {
            return $http.get(gameRatingLink).then(function(response) {
                return response.data;
            });
        }
    };
});
