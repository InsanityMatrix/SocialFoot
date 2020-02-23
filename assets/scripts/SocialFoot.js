function executeHTMLTemplate(template, data) {
  var dataLength = Object.keys(data).length;
  for (var i = 0; i < dataLength; i++) {
    var key = Object.keys(data)[i];
    if (template.includes("{{." + key +"}}")) {
      template = template.replaceAll("{{." + key + "}}", data[key]);
    }
  }
  return template;
}
