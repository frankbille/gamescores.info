angular.module('GameScoresApp').controller('ImportCtrl', function ($scope, ImportService, $timeout) {

    $scope.loading = true;

    var handleStatus = function (status) {
        $scope.doneImporting = !status.importing;
        $scope.importStatus = status;

        var total = status.totalPlayerCount + status.totalLeagueCount + status.totalGameCount;
        var step = status.importedPlayerCount + status.importedLeagueCount + status.importedGameCount;

        if (total > 0) {
            $scope.progress = step / total * 100;
        } else {
            $scope.progress = 0;
        }

        if (status.importing) {
            $timeout(function () {
                ImportService.getImportStatus().then(handleStatus);
            }, 1000);
        }
    };

    ImportService.getImportStatus().then(function (status) {

        if (status.importing) {
            $scope.importing = true;
            $scope.loading = false;
            handleStatus(status);
        } else {
            $scope.importing = false;
            ImportService.prepareImport().then(function (importDefinition) {
                $scope.importDefinition = importDefinition;
                $scope.loading = false;
            });

            $scope.import = function () {
                $scope.importing = true;
                ImportService.doImport($scope.importDefinition).then(handleStatus);
            };
        }
    });

});