angular.module('GameScoresApp').controller('MenuCtrl', function($scope,
  UserService, LeagueService) {

  $scope.userLoading = true;
  UserService.getCurrentUser().then(function(user) {
    $scope.user = user;
    $scope.userLoading = false;
  });

  LeagueService.getLeagueList().then(function(leagueList) {
    $scope.leagueList = leagueList;
  })
});
