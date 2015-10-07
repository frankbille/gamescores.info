angular.module('GameScoresApp').factory('PlayerService', function($http, $q) {
  return {
    getPlayersByIdList: function(playerIdList) {
      var deferred = $q.defer();
      $http.get('/api/players', {
        params: {
          ids: playerIdList
        }
      }).then(function(response) {
        var playerMap = {};

        angular.forEach(response.data.players, function(player) {
          playerMap[player.id] = player;
        });

        deferred.resolve(playerMap);
      }, deferred.reject);
      return deferred.promise;
    }
  };
});
