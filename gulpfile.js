'use strict';

const gulp = require('gulp');
const through = require('through2')
const PluginError = require('plugin-error');
const replace = require('gulp-replace');
const rename = require('gulp-rename');

const PLUGIN_NAME = 'proto2typescript-interfaces';

function buildProtoTS(file, path, callback) {
  if (file.isNull()) {
    // nothing to do
    return callback(null, file);
  }

  if (file.isStream()) {
    // file.contents is a Stream - https://nodejs.org/api/stream.html
    this.emit('error', new PluginError(PLUGIN_NAME, 'Streams not supported!'));

    // or, if you can handle Streams:
    //file.contents = file.contents.pipe(...
    //return callback(null, file);
    return callback(null, file);
  } else if (file.isBuffer()) {

    function protoToTsType(p) {
      const types = {
        "double": "number",
        "float": "number",
        "int32": "number",
        "int64": "number",
        "uint32": "number",
        "uint64": "number",
        "sint32": "number",
        "sint64": "number",
        "fixed32": "number",
        "fixed64": "number",
        "sfixed32": "number",
        "sfixed64": "number",
        "bool": "boolean",
        "string": "string",
        "bytes": "string",
        "Code": "number",
        "interface Error {": "",
      };
      if (types[p]) {
        return types[p];
      }
      return p;
    }

    function parseProtobufLine(line) {
      if (!line) {
        return "";
      }
      const indent = line.length - line.trimLeft().length
      const indentChar = line[0];
      const tokens = line.trim().split(" ").filter(Boolean);
      let isRepeated = false;
      // debugger;
      switch (tokens[0]) {
        case "//":
          return line;
        case "}":
          return "}";
        case "message":
          return "interface " + tokens[1] + " {";
        case "repeated":
          isRepeated = true;
          tokens.shift();
      }
      return `${indentChar.repeat(indent)}${tokens[1]}: ${protoToTsType(tokens[0])}${ isRepeated ? "[]" : ""}`;
    }
    let parsed = "";
    if (file.path.indexOf('Common.proto') > 0) {
      parsed += "declare namespace Common {\n";
    }
    if (file.path.indexOf('SubAPI.proto') > 0) {
      parsed += "declare namespace SubAPI {\n";
    }
    if (file.path.indexOf('AdminAPI.proto') > 0) {
      parsed += "declare namespace AdminAPI {\n";
    }
    let pb = file.contents.toString();

    let inMessage = false;
    for (const line of pb.split("\n")) {
      if (line.indexOf('message') !== -1) {
        inMessage = true;
      }
      if (inMessage) {
        parsed += parseProtobufLine(line);
        parsed += "\n";
      }
      if (line.indexOf('}') !== -1) {
        inMessage = false;
      }
    }
    if (file.path.indexOf('Common.proto') > 0 || file.path.indexOf('SubAPI.proto') > 0 || file.path.indexOf('AdminAPI.proto') > 0) {
      parsed += "}";
    }
    file.contents = Buffer(parsed)
    file.path = file.path.replace('proto', 'd.ts')
    callback(null, file)
  }
}

gulp.task('build-proto-ts', () => {
  let stream = gulp.src('Gigamunch-Proto/**/*.proto')
    .pipe(
      through.obj(buildProtoTS)
    )
    .pipe(
      gulp.dest('Gigamunch-Proto')
    )
    .pipe(rename({
      dirname: ''
    }))
    .pipe(
      gulp.dest('admin/app/ts/prototypes')
    )
  return new Promise((resolve, reject) => {
    stream.on('end', resolve);
    stream.on('error', reject);
  });
});

gulp.task('build-proto-go', () => {
  let stream = gulp.src('Gigamunch-Proto/*/*.pb.go')
    .pipe(replace('../pbcommon', 'github.com/atishpatel/Gigamunch-Backend/Gigamunch-Proto/pbcommon'))
    .pipe(replace(',omitempty', ''))
    .pipe(replace('Url', 'URL'))
    .pipe(replace('Id', 'ID'))
    .pipe(replace('Sms', 'SMS'))
    .pipe(replace('Option_1', 'Option1'))
    .pipe(replace('Option_2', 'Option2'))
    .pipe(replace('Instructions_1', 'Instructions1'))
    .pipe(replace('Instructions_2', 'Instructions2'))
    .pipe(replace('Time_1', 'Time1'))
    .pipe(replace('Time_2', 'Time2'))
    .pipe(gulp.dest('Gigamunch-Proto/'));
  return new Promise((resolve, reject) => {
    stream.on('end', resolve);
    stream.on('error', reject);
  });
});

gulp.task('build-proto-swagger', () => {
  let stream = gulp.src(['*/*/*.swagger.json', '*/*.swagger.json'])
    .pipe(replace('"http",', ''))
    .pipe(replace('"2.0",', '"2.0",\n"securityDefinitions": {"auth-token": {"type": "apiKey","in": "header","name": "auth-token"}},"security": [{"auth-token": []}],'))
    .pipe(gulp.dest('.'));
  return new Promise((resolve, reject) => {
    stream.on('end', resolve);
    stream.on('error', reject);
  });
});

gulp.task('build', gulp.parallel('build-proto-ts', 'build-proto-go', 'build-proto-swagger'));
