angular
    .module('GameScoresApp')
    .factory('PlayerDialog', function ($mdDialog, $state) {
        return function (playerId, event) {
            $mdDialog.show({
                controller: 'PlayerDialogCtrl',
                templateUrl: '/components/players/playerdialog.html',
                parent: angular.element(document.body),
                clickOutsideToClose: true,
                targetEvent: event,
                locals: {
                    playerId: playerId
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
    .controller('PlayerDialogCtrl', function ($scope, $mdDialog, playerId, PlayerService) {
        $scope.loading = true;

        if (playerId != null) {
            PlayerService.getPlayer(playerId).then(function (player) {
                $scope.player = player;
                $scope.loading = false;
            });
        } else {
            $scope.player = {
                active: true,
                _links: {
                    update: {
                        href: '/api/players'
                    }
                }
            };
            $scope.loading = false;
        }


        $scope.save = function () {
            $scope.saving = true;
            PlayerService.savePlayer($scope.player).then($mdDialog.hide);
        };

        $scope.cancel = function () {
            $mdDialog.cancel();
        };
    });