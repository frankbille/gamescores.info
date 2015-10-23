angular.module('defaultapp').factory('UserService', function($q, $http) {
  var user = null;

  return {
    getCurrentUser: function() {
      var deferred = $q.defer();
      if (user === null) {
        $http.get('/api/me').then(function(response) {
          user = response.data;
          deferred.resolve(user);
        });
      } else {
        deferred.resolve(user);
      }
      return deferred.promise;
    }
  };
});
