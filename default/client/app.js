angular.module('defaultapp', [
    'ui.router',
    'ngMaterial',
    'ui.gravatar'
]).config(function ($locationProvider) {
    $locationProvider.html5Mode(true);
});