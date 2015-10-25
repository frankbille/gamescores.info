angular.module('GameScoresApp').controller('AddLeagueButtonCtrl', function($scope, $stateParams, LeagueService, $state) {

    LeagueService.canCreate().then(function(canCreate) {
       $scope.canCreate = canCreate;
    });

    $scope.openDialog = function() {
        $state.go('leagues.add');
    };

});