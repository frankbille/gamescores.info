angular.module('GameScoresApp').controller('MenuCtrl', function($scope,
  UserService) {

  $scope.userLoading = true;
  UserService.getCurrentUser().then(function(user) {
    $scope.user = user;
    $scope.userLoading = false;
  });
});
