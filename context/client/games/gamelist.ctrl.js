angular.module('GameScoresApp').controller('GameListCtrl', function($scope,
  GameService, $stateParams, $window) {
  $scope.getListHeight = function() {
    return {
      height: '' + ($window.innerHeight - 72) + 'px'
    };
  };

  function onResize() {
    $scope.$digest();
  }
  $window.addEventListener('resize', onResize);
  $scope.$on('$destroy', function() {
    $window.removeEventListener('resize', onResize);
  });

  GameService.getGamesForLeague($stateParams.leagueId).then(
    function(gameList) {
      $scope.gameList = {
        gameList: gameList.games,
        total: gameList.total,
        numLoaded: gameList.games.length,
        nextLink: gameList._links.next.href,
        loading: false,

        getLength: function() {
          var length = this.numLoaded + 5;
          if (length > this.total) {
            length = this.total;
          }
          return length;
        },

        getItemAtIndex: function(index) {
          if (index >= this.numLoaded) {
            if (this.nextLink != null && !this.loading) {
              var ga = this;
              ga.loading = true;
              GameService.getGamesForLink(this.nextLink).then(
                function(
                  gameList) {
                  ga.total = gameList.total;
                  ga.numLoaded += gameList.games.length;
                  if (angular.isDefined(gameList._links.next)) {
                    ga.nextLink = gameList._links.next.href;
                  } else {
                    ga.nextLink = null;
                  }
                  angular.forEach(gameList.games, function(game) {
                    ga.gameList.push(game);
                  });
                  ga.loading = false;
                });
            }
            return null;
          }
          return this.gameList[index];
        }
      };
    });
});
