angular.module('GameScoresApp').controller('LeagueListCtrl', function ($scope, LeagueService, $state) {
    LeagueService.getLeagueList().then(function (leagues) {
        $scope.leagues = leagues;
    });

    $scope.editLeague = function (leagueId) {
        $state.go('leagues.edit', {
            leagueId: leagueId
        });
    };
});