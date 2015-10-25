angular.module('GameScoresApp').factory('LeagueService', function ($q, $http) {
    var leagueList = null;
    var leagueMap = {};
    var createLink = null;
    var createLinkResolved = false;

    return {
        getLeagueList: function () {
            var deferred = $q.defer();
            if (leagueList === null) {
                var ls = this;
                $http.get('/api/leagues').then(function (response) {
                    leagueList = response.data;

                    ls._resolveCreateLink(leagueList);

                    for (var i = 0; i < leagueList.leagues.length; i++) {
                        var league = leagueList.leagues[i];
                        leagueMap[league.id] = league;
                    }
                    deferred.resolve(leagueList);
                });
            } else {
                deferred.resolve(leagueList);
            }
            return deferred.promise;
        },

        getLeague: function (leagueId) {
            var deferred = $q.defer();
            if (angular.isUndefined(leagueMap[leagueId])) {
                $http.get('/api/leagues/' + leagueId).then(function (response) {
                    var league = response.data;
                    leagueMap[league.id] = league;
                    deferred.resolve(league);
                });
            } else {
                deferred.resolve(leagueMap[leagueId]);
            }
            return deferred.promise;
        },

        canCreate: function () {
            var deferred = $q.defer();
            if (createLinkResolved) {
                deferred.resolve(createLink != null);
            } else {
                var ls = this;
                $http.get('/api/leagues').then(function (response) {
                    var leagueList = response.data;
                    ls._resolveCreateLink(leagueList);
                    deferred.resolve(createLink != null);
                }, deferred.reject);
            }
            return deferred.promise;
        },

        _resolveCreateLink: function (leagueList) {
            if (!createLinkResolved) {
                if (leagueList._links.create && leagueList._links.create.href) {
                    createLink = leagueList._links.create.href;
                }
                createLinkResolved = true;
            }
        },

        saveLeague: function(league) {
            return $http.post(league._links.update.href, league);
        }
    };
});
