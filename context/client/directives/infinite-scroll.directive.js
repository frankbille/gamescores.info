angular.module('GameScoresApp').directive('infiniteScroll', function ($timeout) {
    var throttler = function (delay, func) {
        var throttleTimeout;
        var throttled = false;
        return function () {
            if (!throttled) {
                throttled = true;
                func.apply();
                throttleTimeout = $timeout(function() {
                   throttled = false;
                }, delay);
            }
        };
    };

    return {
        restrict: 'A',
        link: function ($scope, element, attrs) {
            var raw = element[0];

            element.bind('scroll', throttler(200, function () {
                if (raw.scrollTop + raw.offsetHeight >= raw.scrollHeight - 200) {
                    $scope.$apply(attrs.infiniteScroll);
                }
            }));

        }
    };
});