angular.module('GameScoresApp').factory('PlayerService', function ($http, $q) {
    var playerMap = null;

    return {
        getAllPlayers: function () {
            var deferred = $q.defer();

            if (playerMap != null) {
                deferred.resolve(playerMap);
            } else {
                this._loadPlayers('/api/players').then(function(players) {
                    playerMap = {};

                    angular.forEach(players, function(player) {
                       playerMap[player.id] = player;
                    });

                    deferred.resolve(playerMap);
                });
            }

            return deferred.promise;
        },

        _loadPlayers: function(playersLink, playerArray) {
            var players = playerArray;
            if (angular.isUndefined(players)) {
                players = [];
            }

            var ps = this;

            return $http.get(playersLink).then(function (playersData) {
                var playerList = playersData.data;

                angular.forEach(playerList.players, function(player) {
                    players.push(player);
                });

                if (playerList._links.next && playerList._links.next.href) {
                    return ps._loadPlayers(playerList._links.next.href, players);
                } else {
                    return players;
                }
            });
        },

        getPlayersByIdLink: function (link) {
            var deferred = $q.defer();
            $http.get(link).then(function (response) {
                var playerMap = {};

                angular.forEach(response.data.players, function (player) {
                    playerMap[player.id] = player;
                });

                deferred.resolve(playerMap);
            }, deferred.reject);
            return deferred.promise;
        }
    };
});
