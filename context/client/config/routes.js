angular.module('GameScoresApp').config(function ($stateProvider,
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
            url: '/{leagueId:int}'
        })
        .state('games', {
            url: '/leagues/{leagueId:int}/games',
            views: {
                main: {
                    templateUrl: '/games/games.html',
                    controller: 'GameListCtrl'
                },
                footer: {
                    templateUrl: '/games/addgamebutton.html',
                    controller: 'AddGameButtonCtrl'
                }
            }
        })
        .state('games.edit', {
            url: '/{gameId:int}',
            // Override onEnter to show a modal dialog instead
            // See https://github.com/angular-ui/ui-router/wiki/Frequently-Asked-Questions#how-to-open-a-dialogmodal-at-a-certain-state
            onEnter: function ($stateParams, GameDialog) {
                GameDialog($stateParams.leagueId, $stateParams.gameId);
            }
        })
        .state('games.add', {
            url: '/add',
            // Override onEnter to show a modal dialog instead
            // See https://github.com/angular-ui/ui-router/wiki/Frequently-Asked-Questions#how-to-open-a-dialogmodal-at-a-certain-state
            onEnter: function ($stateParams, GameDialog) {
                GameDialog($stateParams.leagueId, null);
            }
        });
});
