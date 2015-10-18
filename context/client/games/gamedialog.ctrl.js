angular
    .module('GameScoresApp')
    .factory('GameDialog', function ($mdDialog, $state, PlayerService) {
        return function (leagueId, gameId) {
            $mdDialog.show({
                controller: 'GameDialogCtrl',
                templateUrl: '/games/gamedialog.html',
                parent: angular.element(document.body),
                clickOutsideToClose: true,
                locals: {
                    leagueId: leagueId,
                    gameId: gameId
                }
            }).then(function () {
                $state.go('^', {}, {
                    reload: true
                });
            }, function () {
                $state.go('^');
            });
        };
    })
    .controller('GameDialogCtrl', function ($scope, $mdDialog, leagueId, gameId, GameService, PlayerService) {
        PlayerService.getAllPlayers().then(function (playerMap) {
            $scope.playerIdMap = playerMap;

            if (gameId != null) {
                $scope.title = 'Edit game';
                GameService.getGame(leagueId, gameId).then(function (game) {
                    game.gameDate = moment(game.gameDate).toDate();

                    $scope.game = game;
                });
            } else {
                $scope.title = 'Add game';
                $scope.game = {
                    gameDate: new Date(),
                    team1: {
                        players: [],
                        score: 0
                    },
                    team2: {
                        players: [],
                        score: 0
                    },
                    _links: {
                        update: {
                            href: '/api/leagues/' + leagueId + '/games'
                        }
                    }
                }
            }
        });

        $scope.playerName = function (playerId) {
            return $scope.playerIdMap[playerId].name;
        };

        var isNotInGamePlayers = function (playerId) {
            return $scope.game.team1.players.indexOf(playerId) == -1
                && $scope.game.team2.players.indexOf(playerId) == -1;
        };

        $scope.searchPlayers = function (searchText) {
            var foundPlayerIds = [];
            angular.forEach($scope.playerIdMap, function (player, playerId) {
                if (isNotInGamePlayers(playerId)) {
                    if (player.name.toLowerCase().indexOf(searchText.toLowerCase()) > -1) {
                        foundPlayerIds.push(playerId);
                    }
                }
            });
            return foundPlayerIds;
        };

        $scope.save = function () {
            $scope.saving = true;

            var convertPlayerIdsToNumber = function(playerIds) {
                var numberPlayerIds = [];
                angular.forEach(playerIds, function(playerId) {
                    numberPlayerIds.push(Number(playerId));
                });
                return numberPlayerIds;
            };

            $scope.game.team1.players = convertPlayerIdsToNumber($scope.game.team1.players);
            $scope.game.team2.players = convertPlayerIdsToNumber($scope.game.team2.players);

            GameService.saveGame($scope.game).then($mdDialog.hide);
        };

        $scope.cancel = function () {
            $mdDialog.cancel();
        };
    });