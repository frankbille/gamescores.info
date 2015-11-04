angular.module('GameScoresApp').factory('GameService', function ($http, PlayerService, $q) {

    var replacePlayerWithObject = function (playerIdList, playerMap) {
        var newPlayerList = [];
        angular.forEach(playerIdList, function (playerId) {
            var player = playerMap[playerId];
            // Clone player to have unique references for each use. This means we can
            // modify the player object per game
            player = JSON.parse(JSON.stringify(player));
            newPlayerList.push(player);
        });
        return newPlayerList;
    };

    return {
        getGamesForLeague: function (leagueId) {
            return this.getGamesForLink('/api/leagues/' + leagueId + '/games');
        },

        getGamesForLink: function (hyperlink) {
            var deferred = $q.defer();
            $http.get(hyperlink).then(function (gameListData) {
                var gameList = gameListData.data;

                PlayerService.getPlayersByIdLink(gameList._links.playerlist.href).then(function (playerMap) {
                    angular.forEach(gameList.games, function (game) {
                        game.team1.players = replacePlayerWithObject(game.team1.players, playerMap);
                        game.team2.players = replacePlayerWithObject(game.team2.players, playerMap);

                        game.gameDate = moment(game.gameDate);
                    });

                    deferred.resolve(gameList);
                }, deferred.reject);
            }, deferred.reject);
            return deferred.promise;
        },

        getGame: function(leagueId, gameId) {
            return $http.get('/api/leagues/'+leagueId+'/games/'+gameId).then(function(gameData) {
              return gameData.data;
            });
        },

        saveGame: function(game) {
            return $http.post(game._links.update.href, game);
        }
    };
});
