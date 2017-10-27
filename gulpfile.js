const gulp = require('gulp');
const through = require('through2')
const PluginError = require('gulp-util').PluginError;

const PLUGIN_NAME = 'proto2typescript-interfaces';


function buildProto(file, path, callback) {
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
        "Error.Error": "Error",
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
    file.contents = Buffer(parsed)
    file.path = file.path.replace('proto', 'd.ts')
    callback(null, file)
  }

}

gulp.task('build', () => {
  let stream = gulp.src('Gigamunch-Proto/**/*.proto')
    .pipe(
      through.obj(buildProto)
    )
    .pipe(
      gulp.dest('Gigamunch-Proto')
    )
  return new Promise((resolve, reject) => {
    stream.on('end', resolve);
    stream.on('error', reject);
  });
});