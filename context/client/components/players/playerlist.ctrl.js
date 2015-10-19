angular.module('GameScoresApp').controller('PlayerListCtrl', function ($scope, PlayerService, $state) {
    console.log($state);
    PlayerService.getAllPlayers(true).then(function (playerMap) {
        $scope.players = [];
        angular.forEach(playerMap, function(player) {
            $scope.players.push(player);
        });
    });

    $scope.editPlayer = function (playerId) {
        $state.go('players.edit', {
            playerId: playerId
        });
    };
});