angular.module('defaultapp').controller('CreateContextCtrl', function ($scope, $q, user, ContextDefinitionService) {

    $scope.loading = true;
    $scope.available = false;
    $scope.validId = false;
    $scope.user = user;

    if (user._links.prepare && user._links.prepare.href) {
        ContextDefinitionService.prepareContext(user._links.prepare.href).then(function (context) {
            $scope.context = context;
            $scope.loading = false;
        });
    }

    $scope.create = function () {
        ContextDefinitionService.create($scope.context._links.create.href, $scope.context).then(function (data) {
            $scope.created = true;
            console.log(data);
            $scope.contextLink = data._links.self.href;
        });
    };

    $scope.$watch('context.id', function (newValue) {
        if ($scope.context && newValue) {
            ContextDefinitionService.checkId($scope.context._links.checkid.href + newValue).then(function (result) {
                $scope.available = result.available;
                $scope.validId = result.valid;
            });
        } else {
            $scope.available = false;
            $scope.validId = false;
        }
    });

    $scope.getButtonLabel = function () {
        if ($scope.canCreate()) {
            return 'Create';
        } else if(!$scope.validId && $scope.avaliable) {
            return 'Domain should be lower case and no spaces (a-z)';
        } else if ($scope.validId && !$scope.avaliable) {
            return 'Domain not available!';
        } else {
            return 'Enter a domain name'
        }
    };

    $scope.canCreate = function() {
        return $scope.available && $scope.validId;
    };

});
