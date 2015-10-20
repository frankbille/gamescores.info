angular.module('GameScoresApp').factory('PlayerService', function ($http, $q) {
    var playerMap = null;
    var createLink = null;
    var createLinkResolved = false;

    return {
        getAllPlayers: function (reload) {
            var deferred = $q.defer();

            if (reload) {
                playerMap = null;
            }

            if (playerMap != null) {
                deferred.resolve(playerMap);
            } else {
                this._loadPlayers('/api/players').then(function (players) {
                    playerMap = {};

                    angular.forEach(players, function (player) {
                        playerMap[player.id] = player;
                    });

                    deferred.resolve(playerMap);
                }, deferred.reject);
            }

            return deferred.promise;
        },

        _loadPlayers: function (playersLink, playerArray) {
            var players = playerArray;
            if (angular.isUndefined(players)) {
                players = [];
            }

            var ps = this;

            return $http.get(playersLink).then(function (playersData) {
                var playerList = playersData.data;

                ps._resolveCreateLink(playerList);

                angular.forEach(playerList.players, function (player) {
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
        },

        canCreate: function () {
            var deferred = $q.defer();
            if (createLinkResolved) {
                deferred.resolve(createLink != null);
            } else {
                var ps = this;
                $http.get('/api/players').then(function (response) {
                    var playerList = response.data;
                    ps._resolveCreateLink(playerList);
                    deferred.resolve(createLink != null);
                }, deferred.reject);
            }
            return deferred.promise;
        },

        _resolveCreateLink: function (playerList) {
            if (!createLinkResolved) {
                if (playerList._links.create && playerList._links.create.href) {
                    createLink = playerList._links.create.href;
                }
                createLinkResolved = true;
            }
        },

        getPlayer: function(playerId) {
            return $http.get('/api/players/'+playerId).then(function(response) {
                return response.data;
            });
        },

        savePlayer: function(player) {
            return $http.post(player._links.update.href, player);
        }
    };
});
