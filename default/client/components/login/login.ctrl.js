angular.module('defaultapp').controller('LoginCtrl', function($scope, user) {
    $scope.loginUrl = user._links.login.href;
});