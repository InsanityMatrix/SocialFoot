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

function getMessages(receiver, convoid) {
  $.ajax({
    url: '/json/messages/convo',
    method: 'POST',
    data: { "convo": convoid },
    success: parseMessages
  });
  function parseMessages(data) {
    //DATA should already be json
    var response = "<center><h5>Conversation Created</h5></center>";
    let length = data.length;
    for(var i = 0; i < length; i++) {
      if (data[i].mfrom == receiver) {
        //this person sent that message
      } else {
        //this person received that message
      }
    }
  }
}

function sendMessage(from, to) {
  let message = $("input[name='messenger']").val();
  $.ajax({
    url: '',
    method: 'POST',
    data: {
      uidFrom: from,
      uidTo: to,
      message: message
    }
  });
}
