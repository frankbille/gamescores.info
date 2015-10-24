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
                    templateUrl: '/components/games/gamelist.html',
                    controller: 'GameListCtrl'
                },
                footer: {
                    templateUrl: '/components/games/addgamebutton.html',
                    controller: 'AddGameButtonCtrl'
                }
            }
        })
        .state('games.edit', {
            url: '/{gameId:int}',
            // Override onEnter to show a modal dialog instead
            // See https://github.com/angular-ui/ui-router/wiki/Frequently-Asked-Questions#how-to-open-a-dialogmodal-at-a-certain-state
            onEnter: function ($stateParams, GameDialog) {
                GameDialog($stateParams.leagueId, $stateParams.gameId, event);
            }
        })
        .state('games.add', {
            url: '/add',
            // Override onEnter to show a modal dialog instead
            // See https://github.com/angular-ui/ui-router/wiki/Frequently-Asked-Questions#how-to-open-a-dialogmodal-at-a-certain-state
            onEnter: function ($stateParams, GameDialog) {
                GameDialog($stateParams.leagueId, null, event);
            }
        })
        .state('players', {
            url: '/players',
            views: {
                main: {
                    templateUrl: '/components/players/playerlist.html',
                    controller: 'PlayerListCtrl'
                },
                footer: {
                    templateUrl: '/components/players/addplayerbutton.html',
                    controller: 'AddPlayerButtonCtrl'
                }
            }
        })
        .state('players.add', {
            url: '/add',
            onEnter: function(PlayerDialog) {
                PlayerDialog(null, event);
            }
        })
        .state('players.edit', {
            url: '/{playerId:int}',
            onEnter: function($stateParams, PlayerDialog) {
                PlayerDialog($stateParams.playerId, event);
            }
        })
        .state('about', {
            url: '/about',
            onEnter: function(AboutDialog) {
                AboutDialog(event);
            }
        })
        .state('adminimport', {
            url: '/admin/import',
            views: {
                main: {
                    templateUrl: '/components/admin/import.html',
                    controller: 'ImportCtrl'
                }
            }
        });
});
