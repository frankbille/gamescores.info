angular.module('GameScoresApp').factory('ImportService', function ($q, $http) {
    return {
        getImportStatus: function () {
            return $http.get('/api/admin/import/scoreboardv1/status').then(function (response) {
                return response.data;
            });
        },

        prepareImport: function () {
            return $http.get('/api/admin/import/preparescoreboardv1').then(function (response) {
                return response.data;
            });
        },

        doImport: function (importDefinition) {
            return $http.post(importDefinition._links.import.href, importDefinition).then(function (response) {
                return response.data;
            });
        }
    };
});
