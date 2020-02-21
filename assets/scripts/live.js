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
      //Parse Date
      var d = new Date(data[i].posted);
      var dd = d.getDate();
      var mm = d.getMonth() + 1;
      var yyyy = d.getFullYear();
      if (dd < 10) {
        dd = '0' + dd;
      }
      if (mm < 10) {
        mm = '0' + mm;
      }
      var thisdate = mm + '/' + dd + '/' + yyyy;
      var text = "<div class='col-xs-12 post'><div class='row userinfo'><p class='col-xs-6 userName' id='" + data[i].userid + "'></p><p class='col-xs-6 postDate'>" + thisdate + "</p></div>";
      text +=  "<div class='row'><div class='col-xs-12'><center><img class='postimg' src='" + imageLink + "' style='width:80%'></center></div></div>";
      text += "<div class='row'><div class='col-xs-4 postTags'></div><div class='col-xs-4 postCaption'>" + data[i].caption + "</div></div>";
      text += "</div>";
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
  var id = data[0].id;
  $('p[id="' + id + '"').html(data[0].username);
}
