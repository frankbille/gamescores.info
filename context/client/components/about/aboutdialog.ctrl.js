angular
    .module('GameScoresApp')
    .factory('AboutDialog', function ($mdDialog, $state) {
        return function (event) {
            var returnState = $state.current.name;
            var returnParams = $state.params;
            console.log('State:', returnState, 'Params:', returnParams);

            $mdDialog.show({
                controller: 'AboutDialogCtrl',
                templateUrl: '/components/about/aboutdialog.html',
                parent: angular.element(document.body),
                clickOutsideToClose: true,
                targetEvent: event
            }).then(function () {
                if (returnState !== '') {
                    $state.go(returnState, returnParams);
                } else {
                    $state.go('leagues');
                }
            }, function () {
                if (returnState !== '') {
                    $state.go(returnState, returnParams);
                } else {
                    $state.go('leagues');
                }
            });
        };
    })
    .controller('AboutDialogCtrl', function ($scope, $mdDialog) {
        $scope.cancel = $mdDialog.cancel;
    });