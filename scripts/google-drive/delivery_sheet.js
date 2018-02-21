function CopyFromMainToOptimized() {
  var sourceSheet = SpreadsheetApp.getActiveSpreadsheet().getSheetByName("Main");
  var destSheet = SpreadsheetApp.getActiveSpreadsheet().getSheetByName("Optimized");
  // clear 
  destSheet.clear();
  // get address column
  var addressColumn = 1;
  var driverColumn = 1;
  var header = sourceSheet.getRange("A1:J1");
  var headerValues = header.getValues();
  for (var i = 0; i < headerValues[0].length; i++) {
    if (headerValues[0][i] === 'Address') {
      addressColumn = i;
    }
    if (headerValues[0][i] === 'Driver') {
      driverColumn = i;
    }
  }
  // sort
  sourceSheet.sort(driverColumn + 1);
  // get driver ranges
  var body = sourceSheet.getRange("A2:J500")
  var bodyValues = body.getValues();

  var driverRanges = [];
  var startRange = 1;
  var driverName = bodyValues[1][driverColumn];
  for (var i = 0; i < bodyValues.length; i++) {
    if (driverName !== bodyValues[i][driverColumn]) {
      driverRanges.push({
        start: startRange,
        end: i,
        name: driverName,
      });
      startRange = i + 1;
      driverName = bodyValues[i][driverColumn];
    }
  }
  // find best route
  Logger.log('driverRanges: %s', JSON.stringify(driverRanges))
  for (var d; d < driverRanges.length; d++)
    Browser.msgBox(String(JSON.stringify(driverRanges)));
  //  var directions = Maps.newDirectionFinder()
  //     .setOrigin('166 Chesapeake Harbor Blvd, Hendersonville, TN')
  //     .setDestination('2410 Eugenia Ave, Nashville, TN')
  //     .setMode(Maps.DirectionFinder.Mode.DRIVING)
  //     .getDirections();
  //  
  //  var route = directions.routes[0];
  //  Logger.log('route: ', route);
  //  
  //  Maps.newDirectionFinder().setOptimizeWaypoints(true)
  //  Logger.log('testing logging', {obj: 1, obj2: 's'});

  // copy rows

}

function ColorTabs() {
  var sourceSheet = SpreadsheetApp.getActiveSpreadsheet().getSheetByName("Main");
  // get status column
  var statusColumn = 1;
  var vegColumn = 1;
  var nonVegColumn = 1;
  var header = sourceSheet.getRange("A1:J1");
  var headerValues = header.getValues();
  for (var i = 0; i < headerValues[0].length; i++) {
    if (headerValues[0][i] === 'Status') {
      statusColumn = i;
    }
    if (headerValues[0][i] === 'veg') {
      vegColumn = i;
    }
    if (headerValues[0][i] === 'non-veg') {
      nonVegColumn = i;
    }
  }
  // color body
  var body = sourceSheet.getRange("A1:J500")
  var bodyValues = body.getValues();

  var orange = '#fce5cd';
  var blue = '#cfe2f3';
  var purple = '#d9d2e9';
  var green = '#d9ead3';
  for (var i = 0; i < bodyValues.length; i++) {
    var row = bodyValues[i];
    var iPlus = i + 1;
    if (row[vegColumn] === 2) {
      var notation = 'A' + iPlus + ':J' + iPlus;
      sourceSheet.getRange(notation).setBackground(green);
    }
    if (row[vegColumn] === 4) {
      var notation = 'A' + iPlus + ':J' + iPlus;
      sourceSheet.getRange(notation).setBackground(purple);
    }
    if (row[nonVegColumn] === 4) {
      var notation = 'A' + iPlus + ':J' + iPlus;
      sourceSheet.getRange(notation).setBackground(blue);
    }
    if (row[statusColumn] === 'Free') {
      var notation = 'A' + iPlus + ':J' + iPlus;
      sourceSheet.getRange(notation).setBackground(orange);
    }
    if (row[vegColumn] === 2) {
      var notation = 'B' + iPlus + ':B' + iPlus;
      sourceSheet.getRange(notation).setBackground(green);
    }
    if (row[vegColumn] === 4) {
      var notation = 'B' + iPlus + ':B' + iPlus;
      sourceSheet.getRange(notation).setBackground(purple);
    }
    if (row[nonVegColumn] === 4) {
      var notation = 'B' + iPlus + ':B' + iPlus;
      sourceSheet.getRange(notation).setBackground(blue);
    }
  }

}
