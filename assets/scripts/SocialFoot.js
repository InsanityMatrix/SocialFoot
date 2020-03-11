function executeHTMLTemplate(template, data) {
  var dataLength = Object.keys(data).length;
  for (var i = 0; i < dataLength; i++) {
    var key = Object.keys(data)[i];
    while (template.includes("{{." + key +"}}")) {
      template = template.replace("{{." + key + "}}", data[key]);
    }
  }
  return template;
}
function escapeUserString(data) {
  return data.replace(/</g, "&lt;").replace(/>/g, "&gt;");
}
function fallbackCopyTextToClipboard(text) {
  var textArea = document.createElement("textarea");
  textArea.value = text;

  // Avoid scrolling to bottom
  textArea.style.top = "0";
  textArea.style.left = "0";
  textArea.style.position = "fixed";

  document.body.appendChild(textArea);
  textArea.focus();
  textArea.select();

  try {
    var successful = document.execCommand('copy');
    var msg = successful ? 'successful' : 'unsuccessful';
    console.log('Fallback: Copying text command was ' + msg);
  } catch (err) {
    console.error('Fallback: Oops, unable to copy', err);
  }

  document.body.removeChild(textArea);
}
function copyTextToClipboard(text) {
  if (!navigator.clipboard) {
    fallbackCopyTextToClipboard(text);
    return;
  }
  navigator.clipboard.writeText(text).then(function() {
    console.log('Async: Copying to clipboard was successful!');
  }, function(err) {
    console.error('Async: Could not copy text: ', err);
  });
}
