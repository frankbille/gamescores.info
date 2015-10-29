angular.module('GameScoresApp').controller('GameListCtrl', function ($scope, GameService, RatingService, LeagueService, $stateParams, $state) {
    LeagueService.getLeague($stateParams.leagueId).then(function(league) {
       $scope.league = league;
    });

    RatingService.getLeagueResult($stateParams.leagueId).then(function(leagueResult) {
        $scope.playerRankings = leagueResult.players;
    });

    $scope.gameDates = [];
    $scope.hasNextLink = true;
    $scope.loading = false;
    var gameDateMap = {};
    var nextLink = null;

    var processGameList = function(gameList) {
        // Next link
        if (gameList._links.next && gameList._links.next.href) {
            nextLink = gameList._links.next.href;
            $scope.hasNextLink = true;
        } else {
            nextLink = null;
            $scope.hasNextLink = false;
        }

        // Process games
        angular.forEach(gameList.games, function (game) {
            var gameDateKey = game.gameDate.format('YYYY-MM-DD');

            var gameDateObject = gameDateMap[gameDateKey];
            if (angular.isUndefined(gameDateObject)) {
                gameDateObject = {
                    date: game.gameDate,
                    games: []
                };
                gameDateMap[gameDateKey] = gameDateObject;
                $scope.gameDates.push(gameDateObject);
            }

            gameDateObject.games.push(game);
        });

        $scope.loading = false;
    };

    $scope.loadMore = function() {
        if ($scope.hasNextLink) {
            $scope.loading = true;

            if (nextLink != null) {
                GameService.getGamesForLink(nextLink).then(processGameList);
            } else {
                GameService.getGamesForLeague($stateParams.leagueId).then(processGameList);
            }
        }
    };

    $scope.loadMore();

    $scope.editGame = function(gameId) {
        $state.go('games.edit', {
           gameId: gameId
        });
    };
});
