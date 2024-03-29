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
function chooseConversation(elm) {
  let convoid = elm.firstElementChild.innerHTML;
  window.location = "/live/messages/" + convoid;
}
function getMessages(receiver, convoid) {
  $.ajax({
    url: '/json/messages/convo',
    method: 'POST',
    data: { "convo": convoid },
    success: parseMessages
  });
  function parseMessages(data) {
    //DATA should already be json
    let length = data.length;
    let lastMessage = data[length - 1];
    let cCont = $("#textList").html();
    if(cCont.includes("m" + lastMessage.MessageID)) {
      return;
    }
    $("#textList").html("<center><h5>Conversation Created</h5></center>")
    for(var i = 0; i < length; i++) {
      if (data[i].From == receiver) {
        //this person sent that message
        let currentContent = $("#textList").html();
        let mData = {
          "mid": data[i].MessageID,
          "content":data[i].Content
        };
        let newMessage = executeHTMLTemplate(toMsgTemplate, mData);
        $("#textList").html(currentContent + newMessage);
      } else {
        //this person received that message
        let currentContent = $("#textList").html();
        let mData = {
          "mid": data[i].MessageID,
          "content":data[i].Content
        };
        let newMessage = executeHTMLTemplate(fromMsgTemplate, mData);
        $("#textList").html(currentContent + newMessage);
      }
    }
    //Scroll to bottom
  }
}

function sendMessage(from, to) {
  let message = $("input[name='messenger']").val();
  message = escapeUserString(message);
  $.ajax({
    url: '/messages/send/text',
    method: 'POST',
    data: {
      uidFrom: from,
      uidTo: to,
      message: message
    }
  });
  $("input[name='messenger']").val("");
}
