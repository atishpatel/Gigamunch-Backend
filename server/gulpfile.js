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
const htmlMinifier = require('html-minifier').minify;

const buildDirectory = 'js';

function buildTS() {
  return new Promise((resolve, reject) => {
    console.log(`Deleting ${buildDirectory} directory...`);
    del([buildDirectory])
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
            rollup('es')
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

gulp.task('build', buildTS);

gulp.task('watch', () => {
  return watch(['ts/**/*.ts'], gulp.parallel('build'));
});
