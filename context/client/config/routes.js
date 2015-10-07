angular.module('GameScoresApp').config(function($stateProvider,
  $urlRouterProvider) {
  //
  // For any unmatched url, redirect to /state1
  $urlRouterProvider.otherwise('/');
  //
  // Now set up the states
  $stateProvider
    .state('leagues', {
      url: '/leagues'
    })
    .state('leagues.detail', {
      url: '/{leagueId}'
    })
    .state('games', {
      url: '/leagues/{leagueId:int}/games',
      templateUrl: '/games/games.html',
      controller: 'GameListCtrl'
    });
});
