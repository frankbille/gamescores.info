angular.module('defaultapp').config(function ($stateProvider, $urlRouterProvider) {
    //$urlRouterProvider.rule(function($injector, $location) {
    //    console.log($injector, $location);
    //    return $location.path();
    //});

    //$urlRouterProvider.when('/_ah', function($match, $stateParams) {
    //    console.log($match, $stateParams);
    //});
    $urlRouterProvider.otherwise('/');

    $stateProvider
        .state('createcontext', {
            url: '/',
            controller: 'CreateContextCtrl',
            templateUrl: '/components/context/createcontext.html',
            resolve: {
                user: function (UserService) {
                    return UserService.getCurrentUser();
                }
            },
            onEnter: function(user, $state) {
                if (!user.loggedIn) {
                    $state.go('login', {}, {
                        reload: true
                    })
                }
            }
        })
        .state('login', {
            url: '/login',
            controller: 'LoginCtrl',
            templateUrl: '/components/login/login.html',
            resolve: {
                user: function (UserService) {
                    return UserService.getCurrentUser();
                }
            },
            onEnter: function(user, $state) {
                if (user.loggedIn) {
                    $state.go('createcontext', {}, {
                        reload: true
                    })
                }
            }
        });
});
