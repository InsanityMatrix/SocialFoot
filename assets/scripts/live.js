function signOut() {
  //ajax to signout for the server to interpret
  $.ajax({
    url: '/settings/user/signout'
  });
  window.location = "/assets/";
}
function loadFeed() {
  $.ajax({
    url: '/templates/post',
    success: serveFeed
  });
}
function serveFeed(template) {
  //First get Public Posts

  getPublicPosts();

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
        //get tags with regexp
        var tagpattern = /(#\w+)/g;
        var m;
        tagstext = ""
        do {
          m = tagpattern.exec(data[i].tags);
          if (m) {
            tagstext += "<p class='postTag'>" + m[1] + "</p>";
          }
        } while (m);

        var postData = {
          "userid":data[i].userid,
          "thisdate":thisdate,
          "imageLink":imageLink,
          "tags": tagstext,
          "caption":data[i].caption,
          "postid":data[i].postid
        };
        var parsedTemplate = executeHTMLTemplate(template, postData);
        $("#posts").html(stuff + parsedTemplate);
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
}
$(document).ready(function() {
  $("body").on("contextmenu",function(e){
    return false;
  });
});
