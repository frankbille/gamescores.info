angular.module('defaultapp').factory('ContextDefinitionService', function ($http) {
    return {
        prepareContext: function (prepareLink) {
            return $http.get(prepareLink).then(function (response) {
                return response.data;
            });
        },

        checkId: function (checkIdLink) {
            return $http.get(checkIdLink).then(function (response) {
                return response.data;
            });
        },

        create: function (createLink, context) {
            return $http.post(createLink, context).then(function (response) {
                return response.data;
            });
        }
    };
});
