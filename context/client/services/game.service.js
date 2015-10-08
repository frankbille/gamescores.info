angular.module('GameScoresApp').factory('GameService', function($http,
  PlayerService, $q) {

  var addPlayerIdsToMap = function(playerIdList, playerIdMap) {
    angular.forEach(playerIdList, function(playerId) {
      playerIdMap[playerId] = true;
    });
  };

  var replacePlayerWithObject = function(playerIdList, playerMap) {
    var newPlayerList = [];
    angular.forEach(playerIdList, function(playerId) {
      newPlayerList.push(playerMap[playerId]);
    });
    return newPlayerList;
  };

  return {
    getGamesForLeague: function(leagueId) {
      return this.getGamesForLink('/api/leagues/' + leagueId + '/games');
    },

    getGamesForLink: function(hyperlink) {
      var deferred = $q.defer();
      $http.get(hyperlink).then(function(
        gameListData) {
        var gameList = gameListData.data;

        PlayerService.getPlayersByIdLink(gameList._links.playerlist.href)
          .then(function(
            playerMap) {

            angular.forEach(gameList.games, function(game) {
              game.team1.players = replacePlayerWithObject(game
                .team1.players, playerMap);
              game.team2.players = replacePlayerWithObject(game
                .team2.players, playerMap);

              game.gameDate = moment(game.gameDate);
            });

            deferred.resolve(gameList);
          }, deferred.reject);
      }, deferred.reject);
      return deferred.promise;
    }
  };
});
