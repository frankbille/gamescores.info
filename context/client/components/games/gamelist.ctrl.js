angular.module('GameScoresApp').controller('GameListCtrl', function ($scope, GameService, RatingService, LeagueService, $stateParams, $state) {
    LeagueService.getLeague($stateParams.leagueId).then(function (league) {
        $scope.league = league;
    });

    RatingService.getLeagueResult($stateParams.leagueId).then(function (leagueResult) {
        $scope.playerRankings = leagueResult.players;
    });

    $scope.gameDates = [];
    $scope.hasNextLink = true;
    $scope.loading = false;
    var gameDateMap = {};
    var gameIdMap = {};
    var nextLink = null;
    //var gameRatings = {};

    var processGameList = function (gameList) {
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

            gameIdMap[game.id] = game;
        });

        // Ratings
        if (gameList._links.gameratinglist && gameList._links.gameratinglist.href) {
            RatingService.getGameRatings(gameList._links.gameratinglist.href).then(function (gameRatingList) {

                angular.forEach(gameRatingList, function (gameRating) {
                    var game = gameIdMap[gameRating.gameId];

                    var findPlayerInTeam = function (team, playerId) {
                        var foundPlayer = null;
                        angular.forEach(team.players, function (player) {
                            if (player.id == playerId) {
                                foundPlayer = player;
                            }
                        });
                        return foundPlayer;
                    };

                    var findPlayer = function (playerId) {
                        var foundPlayer = findPlayerInTeam(game.team1, playerId);
                        if (foundPlayer == null) {
                            foundPlayer = findPlayerInTeam(game.team2, playerId);
                        }
                        return foundPlayer;
                    };

                    var forEachPlayerRating = function(playerRatingCallback) {
                        angular.forEach(gameRating.loosingTeamRating.playerRatings, playerRatingCallback);
                        angular.forEach(gameRating.winningTeamRating.playerRatings, playerRatingCallback);
                    };

                    forEachPlayerRating(function(playerRating) {
                        playerRating.diff = playerRating.newrating-playerRating.oldrating;

                        var player = findPlayer(playerRating.playerId);
                        player.playerRating = playerRating;
                    });
                });

                console.log(gameIdMap);

                $scope.loading = false;
            });
        } else {
            $scope.loading = false;
        }
    };

    $scope.loadMore = function () {
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

    $scope.editGame = function (gameId) {
        $state.go('games.edit', {
            gameId: gameId
        });
    };

    //$scope.getPlayerRating = function(gameId, playerId) {
    //    var gameRating = gameRatings[gameId];
    //    if (gameRating) {
    //        var playerRating = gameRating.playerRatings[playerId];
    //        return playerRating.newrating;
    //    }
    //
    //    return null;
    //};
    //
    //$scope.getPlayerRatingDiff = function(gameId, playerId) {
    //    var gameRating = gameRatings[gameId];
    //    if (gameRating) {
    //        var playerRating = gameRating.playerRatings[playerId];
    //        return playerRating.newrating - playerRating.oldrating;
    //    }
    //
    //    return null;
    //};
});
