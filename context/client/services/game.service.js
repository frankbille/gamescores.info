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
      var deferred = $q.defer();
      $http.get('/api/leagues/' + leagueId + '/games').then(function(
        gameListData) {
        var gameList = gameListData.data;
        var playerIds = {};
        angular.forEach(gameList.games, function(game) {
          addPlayerIdsToMap(game.team1.players, playerIds);
          addPlayerIdsToMap(game.team2.players, playerIds);
        });

        var playerIdList = [];
        angular.forEach(playerIds, function(value, playerId) {
          playerIdList.push(playerId);
        })

        PlayerService.getPlayersByIdList(playerIdList).then(function(
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
