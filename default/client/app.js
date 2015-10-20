angular.module('defaultapp', [
    'ui.router',
    'ngMaterial'
]).config(function ($locationProvider) {
    $locationProvider.html5Mode(true);
});