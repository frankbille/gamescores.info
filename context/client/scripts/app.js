angular.module('GameScoresApp', [
  'ngMaterial',
  'ui.gravatar'
]).config(function($locationProvider) {
  $locationProvider.html5Mode(true);
});
