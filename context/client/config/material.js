angular.module('GameScoresApp').config(function ($mdThemingProvider, $mdIconProvider, $mdDateLocaleProvider) {
    $mdThemingProvider.theme('default')
        .primaryPalette('blue-grey')
        .accentPalette('green');

    $mdIconProvider
        .iconSet('action', '/images/action-icons.svg', 24)
        .iconSet('content', '/images/content-icons.svg', 24)
        .iconSet('editor', '/images/editor-icons.svg', 24)
        .iconSet('navigation', '/images/navigation-icons.svg', 24);

    $mdDateLocaleProvider.firstDayOfWeek = 1;
    $mdDateLocaleProvider.parseDate = function(dateString) {
      var m = moment(dateString, 'YYYY-MM-DD', true);
        return m.isValid() ? m.toDate() : new Date(NaN);
    };
    $mdDateLocaleProvider.formatDate = function(date) {
      return moment(date).format('YYYY-MM-DD')
    };
});
