angular.module('GameScoresApp').controller('AddGameButtonCtrl', function($scope, $stateParams, LeagueService, $state) {

    LeagueService.getLeague($stateParams.leagueId).then(function(league) {
       $scope.league = league;
    });

    $scope.openDialog = function(ev) {
        ev.preventDefault();
        $state.go('games.add', {
            ev: ev
        });
    };

});