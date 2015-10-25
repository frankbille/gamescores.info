angular
    .module('GameScoresApp')
    .factory('LeagueDialog', function ($mdDialog, $state) {
        return function (leagueId, event) {
            $mdDialog.show({
                controller: 'LeagueDialogCtrl',
                templateUrl: '/components/leagues/leaguedialog.html',
                parent: angular.element(document.body),
                clickOutsideToClose: true,
                targetEvent: event,
                locals: {
                    leagueId: leagueId
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
    .controller('LeagueDialogCtrl', function ($scope, $mdDialog, leagueId, LeagueService) {
        $scope.loading = true;

        if (leagueId != null) {
            LeagueService.getLeague(leagueId).then(function (league) {
                $scope.league = league;
                $scope.loading = false;
            });
        } else {
            $scope.league = {
                active: true,
                _links: {
                    update: {
                        href: '/api/leagues'
                    }
                }
            };
            $scope.loading = false;
        }


        $scope.save = function () {
            $scope.saving = true;
            console.log($scope.league);
            LeagueService.saveLeague($scope.league).then($mdDialog.hide);
        };

        $scope.cancel = function () {
            $mdDialog.cancel();
        };
    });