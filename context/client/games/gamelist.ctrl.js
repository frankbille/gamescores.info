angular.module('GameScoresApp').controller('GameListCtrl', function($scope,
  GameService, $stateParams) {
  GameService.getGamesForLeague($stateParams.leagueId).then(
    function(gameList) {
      $scope.gameList = gameList;
    })
});
