function ColorTabs() {
  var sourceSheet = SpreadsheetApp.getActiveSpreadsheet().getSheetByName("Main");
  // get status column
  var statusColumn = 1;
  var vegColumn = 1;
  var nonVegColumn = 1;
  var driverNameColumn = 1;
  var header = sourceSheet.getRange("A1:J1");
  var headerValues = header.getValues();
  for (var i = 0; i < headerValues[0].length; i++) {
    if (headerValues[0][i].toUpperCase() === 'STATUS') {
      statusColumn = i;
    }

    if (headerValues[0][i].toUpperCase() === 'VEG') {
      vegColumn = i;
    }
    if (headerValues[0][i].toUpperCase() === 'NON-VEG') {
      nonVegColumn = i;
    }
    if (headerValues[0][i].toUpperCase() === 'DRIVER') {
      driverNameColumn = i;
    }
  }
  // color body

  sourceSheet.getRange("A2:K1000").setBackground('#fff');
  var body = sourceSheet.getRange("A1:K1000")
  var bodyValues = body.getValues();

  var orange = '#fce5cd';
  var blue = '#cfe2f3';
  var purple = '#d9d2e9';
  var green = '#d9ead3';
  var red = '#ea9999';
  for (var i = 0; i < bodyValues.length; i++) {
    var row = bodyValues[i];
    var iPlus = i + 1;
    if (row[vegColumn] === 2) {
      var notation = 'A' + iPlus + ':K' + iPlus;
      sourceSheet.getRange(notation).setBackground(green);
    }
    if (row[vegColumn] === 4) {
      var notation = 'A' + iPlus + ':K' + iPlus;
      sourceSheet.getRange(notation).setBackground(purple);
    }
    if (row[nonVegColumn] === 4) {
      var notation = 'A' + iPlus + ':K' + iPlus;
      sourceSheet.getRange(notation).setBackground(blue);
    }
    if (row[statusColumn] === 'Free') {
      var notation = 'A' + iPlus + ':K' + iPlus;
      sourceSheet.getRange(notation).setBackground(orange);
    }
    if (row[vegColumn] === 2) {
      var notation = 'B' + iPlus + ':B' + iPlus;
      sourceSheet.getRange(notation).setBackground(green);
      notation = 'D' + iPlus + ':D' + iPlus;
      sourceSheet.getRange(notation).setBackground(green);
      notation = 'F' + iPlus + ':F' + iPlus;
      sourceSheet.getRange(notation).setBackground(green);
    }
    if (row[vegColumn] === 4) {
      var notation = 'B' + iPlus + ':B' + iPlus;
      sourceSheet.getRange(notation).setBackground(purple);
      notation = 'D' + iPlus + ':D' + iPlus;
      sourceSheet.getRange(notation).setBackground(purple);
      notation = 'F' + iPlus + ':F' + iPlus;
      sourceSheet.getRange(notation).setBackground(purple);
    }
    if (row[nonVegColumn] === 4) {
      var notation = 'B' + iPlus + ':B' + iPlus;
      sourceSheet.getRange(notation).setBackground(blue);
      notation = 'D' + iPlus + ':D' + iPlus;
      sourceSheet.getRange(notation).setBackground(blue);
      notation = 'F' + iPlus + ':F' + iPlus;
      sourceSheet.getRange(notation).setBackground(blue);
    }
    var total = row[nonVegColumn] + row[vegColumn];
    if (total > 4 || (row[nonVegColumn] > 0 && row[vegColumn] > 0)) {
      var notation = 'A' + iPlus + ':K' + iPlus;
      sourceSheet.getRange(notation).setBackground(red);
      notation = 'D' + iPlus + ':D' + iPlus;
      sourceSheet.getRange(notation).setBackground(red);
      notation = 'F' + iPlus + ':F' + iPlus;
      sourceSheet.getRange(notation).setBackground(red);
    }
  }

}

function InsertReminders() {
  var name = "WithReminders";
  // COPY Main to WithReminders
  var ss = SpreadsheetApp.getActiveSpreadsheet();
  var sheet = ss.getSheetByName('Main').copyTo(ss);

  /* Before cloning the sheet, delete any previous copy */
  var old = ss.getSheetByName(name);
  if (old) ss.deleteSheet(old);

  SpreadsheetApp.flush();
  sheet.setName(name);
  ss.setActiveSheet(sheet);
  // END COPY

  var sourceSheet = SpreadsheetApp.getActiveSpreadsheet().getSheetByName(name);
  // get status column
  var statusColumn = 1;
  var vegColumn = 1;
  var nonVegColumn = 1;
  var driverNameColumn = 1;
  var header = sourceSheet.getRange("A1:J1");
  var headerValues = header.getValues();
  for (var i = 0; i < headerValues[0].length; i++) {
    if (headerValues[0][i].toUpperCase() === 'STATUS') {
      statusColumn = i;
    }
    if (headerValues[0][i].toUpperCase() === 'VEG') {
      vegColumn = i;
    }
    if (headerValues[0][i].toUpperCase() === 'NON-VEG') {
      nonVegColumn = i;
    }
    if (headerValues[0][i].toUpperCase() === 'DRIVER') {
      driverNameColumn = i;
    }
  }
  // color body
  sourceSheet.setFrozenColumns(0);
  sourceSheet.getRange("A2:K1000").setFontWeight('normal');
  sourceSheet.getRange("D2:D1000").setFontWeight('bold');
  var body = sourceSheet.getRange("A1:K1000");
  var bodyValues = body.getValues();

  var gray = '#f3f3f3';

  // Add Text Chris messages
  var previousDriverEnd = 0;
  var previousDriver = '';
  for (var i = bodyValues.length - 1; i >= 0; i--) {
    var row = bodyValues[i];
    var driverName = row[driverNameColumn];
    var iPlus = i + 2;
    if (previousDriver !== driverName) {
      if (i !== 0) {
        sourceSheet.insertRowAfter(iPlus - 1);
        var notation = 'A' + iPlus + ':K' + iPlus;
        sourceSheet.getRange(notation).merge().setBackground(gray).setFontWeight('bold').setValue('TEXT CHRIS');
      }
      if (previousDriverEnd !== 0 && (previousDriverEnd - i) > 2) {
        var insertIndex = i + Math.floor((previousDriverEnd - i) / 2);
        sourceSheet.insertRowAfter(insertIndex + 2)
        var selectRangeIndex = insertIndex + 3;
        var fullRow = sourceSheet.getRange('A' + selectRangeIndex + ':K' + selectRangeIndex);
        fullRow.setBackground(gray).setFontWeight('bold');
        var firstHalf = sourceSheet.getRange('A' + selectRangeIndex + ':E' + selectRangeIndex);
        firstHalf.merge().setValue('Count if the number of remaining bags matches the following number and TEXT NUMBER to Chris:')
        var countCell = sourceSheet.getRange('F' + selectRangeIndex + ':F' + selectRangeIndex);
        // formula =COUNTIFS(Hx:Iy, ">0")
        countCell.setFormula('=COUNTIFS(H' + (selectRangeIndex) + ':I' + (previousDriverEnd + 3) + ', ">0")');
      }

      previousDriverEnd = i;
      previousDriver = driverName;
    }
  }
}
