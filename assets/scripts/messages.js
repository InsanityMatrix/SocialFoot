var toMsgTemplate;
var fromMsgTemplate;
$(document).ready(function() {
  //Initialize the templates
  $.ajax({
    url: '/templates/tomsg',
    success: getToMsgSuccess
  });
  $.ajax({
    url: '/templates/frommsg',
    success: getFromMsgSuccess
  });
  function getToMsgSuccess(data) {
    toMsgTemplate = data;
  }
  function getFromMsgSuccess(data) {
    fromMsgTemplate = data;
  }
});

function getMessages(receiver, other) {

}
