const del = require('del');
const gulp = require('gulp');
// const gulpif = require('gulp-if');
// const mergeStream = require('merge-stream');
// const polymerBuild = require('polymer-build');
const importsInliner = require('gulp-js-text-imports');
const ts = require('gulp-typescript');
const watch = require('gulp-watch');
const rollup = require('gulp-better-rollup');
let rename = require("gulp-rename");



// Additional plugins can be used to optimize your source files after splitting.
// Before using each plugin, install with `npm i --save-dev <package-name>`
const uglify = require('gulp-uglify-es').default;
// const cssSlam = require('css-slam').gulp;
const htmlMinifier = require('html-minifier').minify;

// const swPrecacheConfig = require('./sw-precache-config.js');
const buildDirectory = 'js';
const protoDirectory = '../../Gigamunch-Proto';

function buildApp() {
  return new Promise((resolve, reject) => { // eslint-disable-line no-unused-vars
    // Okay, so first thing we do is clear the build directory
    console.log(`Deleting ${buildDirectory} directory...`);
    del([buildDirectory])
      .then(() => {
        console.log(`Compiling typescript...`);
        let stream = gulp.src('./ts/*.ts')
          .pipe(importsInliner({
            parserOptions: {
              allowImportExportEverywhere: true,
            },
            handlers: {
              html: (content, path, callback) => {
                let result;
                try {
                  result = htmlMinifier(content, {
                    collapseWhitespace: true,
                  });
                } catch (err) {
                  return callback(err, null);
                }
                return callback(null, result);
              },
            },
          }))
          .pipe(ts({
            noImplicitAny: true,
            allowJs: true,
            target: 'ES2015',
            module: 'es2015',
            removeComments: true,
            typeRoots: [`${protoDirectory}/shared`, `${protoDirectory}/admin`, `${protoDirectory}/driver`, `${protoDirectory}/sub`],
          }))
          .pipe(
            gulp.dest(buildDirectory)
          );
        return new Promise((resolve, reject) => {
          stream.on('end', resolve);
          stream.on('error', reject);
        });
      })
      .then(() => {
        console.log(`Rolling up...`);
        let stream = gulp.src(`${buildDirectory}/app.js`)
          .pipe(
            rollup({
              treeshake: false,
              format: 'es',
            })
          )
          .pipe(
            gulp.dest(buildDirectory)
          );
        return new Promise((resolve, reject) => {
          stream.on('end', resolve);
          stream.on('error', reject);
        });
      })
      .then(() => {
        console.log(`Minifying...`);
        let stream = gulp.src(`${buildDirectory}/app.js`)
          .pipe(
            uglify()
          )
          .pipe(
            rename("app.min.js")
          )
          .pipe(
            gulp.dest(buildDirectory)
          );
        return new Promise((resolve, reject) => {
          stream.on('end', resolve);
          stream.on('error', reject);
        });
      })
      .then(() => {
        // You did it!
        console.log('Build complete!');
        resolve();
      });
  });
}

gulp.task('build', buildApp);

gulp.task('watch', () => {
  return watch(['src/**/*.ts', 'src/**/*.html'], gulp.parallel('build'));
});