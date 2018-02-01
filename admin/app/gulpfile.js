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
// const htmlMinifier = require('html-minifier').minify;

const buildTSDirectory = 'js';
const buildOldTSDirectory = 'js-old';

function buildOldTS() {
  return new Promise((resolve, reject) => { // eslint-disable-line no-unused-vars
    // Okay, so first thing we do is clear the build directory
    console.log(`Deleting ${buildOldTSDirectory} directory...`);
    del([buildOldTSDirectory])
      .then(() => {
        console.log(`Compiling old typescript...`);
        let stream = gulp.src('./ts-old/*.ts')
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
          }))
          .pipe(
            gulp.dest(buildOldTSDirectory)
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

function buildTS() {
  return new Promise((resolve, reject) => { // eslint-disable-line no-unused-vars
    // Okay, so first thing we do is clear the build directory
    console.log(`Deleting ${buildTSDirectory} directory...`);
    del([buildTSDirectory])
      .then(() => {
        console.log(`Compiling typescript...`);
        let stream = gulp.src('./ts/**/*.ts')
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
          }))
          .pipe(
            gulp.dest(buildTSDirectory)
          );
        return new Promise((resolve, reject) => {
          stream.on('end', resolve);
          stream.on('error', reject);
        });
      })
      .then(() => {
        console.log(`Rolling up...`);
        let stream = gulp.src(`${buildTSDirectory}/app.js`)
          .pipe(
            rollup('es')
          )
          .pipe(
            gulp.dest(buildTSDirectory)
          );
        return new Promise((resolve, reject) => {
          stream.on('end', resolve);
          stream.on('error', reject);
        });
      })
    //   .then(() => {
    //     console.log(`Minifying...`);
    //     let stream = gulp.src(`${buildDirectory}/app-shell.js`)
    //       .pipe(
    //         uglify()
    //       )
    //       .pipe(
    //         rename("app-shell.min.js")
    //       )
    //       .pipe(
    //         gulp.dest(buildDirectory)
    //       );
    //     return new Promise((resolve, reject) => {
    //       stream.on('end', resolve);
    //       stream.on('error', reject);
    //     });
    //   })
      // .then(() => {
      //   // Okay, now let's generate the Service Worker
      //   console.log('Generating the Service Worker...');
      //   return polymerBuild.addServiceWorker({
      //     project: polymerProject,
      //     buildRoot: buildDirectory,
      //     bundled: true,
      //     swPrecacheConfig: swPrecacheConfig
      //   });
      // })
      .then(() => {
        // You did it!
        console.log('Build complete!');
        resolve();
      });
  });
}

gulp.task('build', buildTS);
gulp.task('build-old', buildOldTS);

gulp.task('watch', () => {
  return watch(['ts/**/*.ts'], gulp.parallel('build'));
});