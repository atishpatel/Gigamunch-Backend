const del = require('del');
const gulp = require('gulp');

const ts = require('gulp-typescript');
const tsProject = ts.createProject('../tsconfig.json')
const watch = require('gulp-watch');
// const babel = require('gulp-babel');
const rollup = require('gulp-better-rollup');
let rename = require("gulp-rename");

// Additional plugins can be used to optimize your source files after splitting.
// Before using each plugin, install with `npm i --save-dev <package-name>`
const uglify = require('gulp-uglify-es').default;
// const htmlMinifier = require('html-minifier').minify;

const buildDirectory = 'js';

function buildTS() {
  return new Promise((resolve, reject) => {
    // delete build dir
    console.log(`Deleting ${buildDirectory} directory...`);
    del([buildDirectory])
      // compile ts
      .then(() => {
        console.log(`Compiling typescript...`);
        let stream = gulp.src('./ts/**/*.ts')
          // compile ts
          .pipe(tsProject())
          // output files to build dir
          .pipe(
            gulp.dest(buildDirectory)
          );
        return new Promise((resolve, reject) => {
          stream.on('end', resolve);
          stream.on('error', reject);
        });
      })
      // roll up
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
      // minify
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
