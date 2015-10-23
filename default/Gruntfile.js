'use strict';

module.exports = function (grunt) {

    require('load-grunt-tasks')(grunt);

    // Load default config
    var config = require('../shared/gruntdefaultconfig');

    config.ngtemplates.dist.options.module = 'defaultapp';

    grunt.initConfig(config);

    // Load shared tasks
    grunt.loadTasks('../shared/grunttasks');

    // Define default
    grunt.registerTask('default', [
        'clean',
        'defaultbuild'
    ]);

};
