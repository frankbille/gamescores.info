angular.module('GameScoresApp').factory('PlayerService', function($http, $q) {
  return {
    getPlayersByIdLink: function(link) {
      var deferred = $q.defer();
      $http.get(link).then(function(response) {
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
