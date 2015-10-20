module.exports = function (grunt) {
    grunt.registerTask('defaultbuild', [
        'wiredep',
        'useminPrepare',
        'copy:styles',
        'autoprefixer',
        'htmlmin',
        'ngtemplates',
        'concat',
        'ngAnnotate',
        'copy:dist',
        'cssmin',
        'uglify',
        'rev',
        'usemin',
        'manifest'
    ]);
};