function signOut() {
  //ajax to signout for the server to interpret
  $.ajax({
    url: '/settings/user/signout'
  });
  window.location = "/assets/";
}

function getPublicPosts() {
  $.getJSON("/posts/public", function(data){
    var length = data.length;
    $("#posts").html("");
    for(var i = 0; i < length; i++) {
      var imageLink = "/assets/uploads/imageposts/post" + data[i].postid + data[i].extension;
      var stuff = $("#posts").html();
      var text = "<div class='post'><div class='row userinfo'><p class='col-md-6 userName' id='" + data[i].userid + "'></p></div>";
      text +=  "<div class='row'><div class='col-xs-12'><center><img src='" + imageLink + "' style='width:80%'></center></div></div></div>";
      $("#posts").html(stuff + text);
      $.ajax({
        url: '/json/user/id',
        method: 'POST',
        data : {
          "userid": data[i].userid
        },
        success: putPostUsernames
      });
    }
  });
}

function putPostUsernames(data) {
  var obj = JSON.parse(data);
  var id = obj[0].id;
  $("#" + id).html(obj[0].username);
}
