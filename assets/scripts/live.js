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
      var text = "<div class='post'><div class='row userinfo'><p class='userName'>" + data[i].userid + "</p></div>";
      text +=  "<div class='row'><div class='col-xs-12'><img src='" + imageLink + "' style='width:100%'></div></div></div>";
      $("#posts").html(stuff + text);
    }
  });
}
