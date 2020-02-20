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
    document.getElementById("posts").innerHTML = "";
    for(var i = 0; i < length; i++) {
      var imageLink = "/assets/uploads/imageposts/post" + data[i].postid + data[i].extension;
      document.getElementById("posts").innerHTML += "<div class='post'><div class='row userinfo'><p class='userName'>" + data[i].userid + "</p></div>";
      document.getElementById("posts").innerHTML += "<div class='row'><div class='col-xs-12'><img src='" + imageLink + "' style='width:100%'></div></div></div>";
    }
  });
}
