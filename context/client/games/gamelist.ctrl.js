angular.module('GameScoresApp').controller('GameListCtrl', function ($scope, GameService, LeagueService, $stateParams) {
    LeagueService.getLeague($stateParams.leagueId).then(function(league) {
       $scope.league = league;
    });

    $scope.gameDates = [];
    var gameDateMap = {};
    var nextLink = null;
    var hasNextLink = true;

    var processGameList = function(gameList) {
        // Next link
        if (gameList._links.next && gameList._links.next.href) {
            nextLink = gameList._links.next.href;
            hasNextLink = true;
        } else {
            nextLink = null;
            hasNextLink = false;
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
    };

    $scope.loadMore = function() {
        if (hasNextLink) {
            if (nextLink != null) {
                GameService.getGamesForLink(nextLink).then(processGameList);
            } else {
                GameService.getGamesForLeague($stateParams.leagueId).then(processGameList);
            }
        }
    };

    $scope.loadMore();


});
