angular.module('GameScoresApp').controller('AddPlayerButtonCtrl', function($scope, $stateParams, PlayerService, $state) {

    PlayerService.canCreate().then(function(canCreate) {
       $scope.canCreate = canCreate;
    });

    $scope.openDialog = function() {
        $state.go('players.add');
    };

});