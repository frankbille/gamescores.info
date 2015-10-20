'use strict';

module.exports = function(grunt) {

  // Load grunt tasks automatically
  require('load-grunt-tasks')(grunt);

  // Define the configuration for all the tasks
  grunt.initConfig({
    cfg: {
      app: 'client',
      dist: 'dist'
    },

    clean: {
      dist: {
        files: [{
          dot: true,
          src: [
            '.tmp',
            '<%= cfg.dist %>/*',
            '!<%= cfg.dist %>/.git*'
          ]
        }]
      }
    },

    // Automatically inject Bower components into the app
    wiredep: {
      app: {
        src: '<%= cfg.app %>/index.html',
        ignorePath: '<%= cfg.app %>/'
      }
    },

    // Add vendor prefixed styles
    autoprefixer: {
      options: {
        browsers: ['last 1 version']
      },
      dist: {
        files: [{
          expand: true,
          cwd: '.tmp/styles/',
          src: '{,*/}*.css',
          dest: '.tmp/styles/'
        }]
      }
    },

    // Renames files for browser caching purposes
    rev: {
      files: {
        src: [
          '<%= cfg.dist %>/{,*/}*.js',
          '<%= cfg.dist %>/styles/{,*/}*.css',
          '<%= cfg.dist %>/styles/fonts/*'
        ]
      }
    },

    // Reads HTML for usemin blocks to enable smart builds that automatically
    // concat, minify and revision files. Creates configurations in memory so
    // additional tasks can operate on them
    useminPrepare: {
      html: '<%= cfg.app %>/index.html',
      options: {
        dest: '<%= cfg.dist %>'
      }
    },

    // Performs rewrites based on rev and the useminPrepare configuration
    usemin: {
      html: ['<%= cfg.dist %>/{,*/}*.html'],
      css: ['<%= cfg.dist %>/styles/{,*/}*.css'],
      options: {
        assetsDirs: ['<%= cfg.dist %>']
      }
    },

    // Copies remaining files to places other tasks can use
    copy: {
      dist: {
        files: [{
          expand: true,
          dot: true,
          cwd: '<%= cfg.app %>',
          dest: '<%= cfg.dist %>',
          src: [
            '*.{ico,png,txt}',
            '*.html',
            'bower_components/**/*',
            'fonts/*'
          ]
        }, {
          expand: true,
          cwd: '<%= cfg.app %>/images',
          dest: '<%= cfg.dist %>/images',
          src: ['*.svg']
        }, {
          expand: true,
          cwd: '.tmp/images',
          dest: '<%= cfg.dist %>/images',
          src: ['generated/*']
        }]
      },
      styles: {
        expand: true,
        cwd: '<%= cfg.app %>/styles',
        dest: '.tmp/styles/',
        src: '{,*/}*.css'
      }
    },

    htmlmin: {
      distviews: {
        options: {
          collapseWhitespace: true,
          collapseBooleanAttributes: true,
          removeCommentsFromCDATA: true,
          removeOptionalTags: true
        },
        files: [{
          expand: true,
          cwd: '<%= cfg.app %>',
          src: ['**/*.html'],
          dest: '.tmp'
        }]
      }
    },

    // Allow the use of non-minsafe AngularJS files. Automatically makes it
    // minsafe compatible so Uglify does not destroy the ng references
    ngAnnotate: {
      dist: {
        files: [{
          expand: true,
          cwd: '.tmp/concat/scripts',
          src: '*.js',
          dest: '.tmp/concat/scripts'
        }]
      }
    },

    ngtemplates: {
      dist: {
        options: {
          module: 'GameScoresApp',
          usemin: 'scripts/app.js',
          prefix: '/'
        },
        cwd: '.tmp',
        src: ['**/*.html', '!index.html'],
        dest: '.tmp/templates.js'
      }
    },

    manifest: {
      generate: {
        options: {
          basePath: '<%= cfg.dist %>',
          network: ['*'],
          preferOnline: true,
          timestamp: true
        },
        src: [
          'index.html',
          'scripts/*',
          'styles/*'
        ],
        dest: '<%= cfg.dist %>/manifest.appcache'
      }
    }
  });

  grunt.registerTask('build', [
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

  grunt.registerTask('default', [
    'clean',
    'build'
  ]);

};
