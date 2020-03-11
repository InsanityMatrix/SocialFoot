function signOut() {
  //ajax to signout for the server to interpret
  $.ajax({
    url: "/settings/user/signout",
    success: function() {   window.location = "/assets/"; }
  });
}
function loadFeed() {
  $.ajax({
    url: "/templates/post",
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
        var imageLink;
        if (data[i].type == "IMAGE") {
          imageLink = "/assets/uploads/imageposts/post" + data[i].postid + data[i].extension;
        } else {
          imageLink = "/assets/uploads/videoposts/post" + data[i].postid + data[i].extension;
        }
        var stuff = $("#posts").html();
        //Parse Date
        var d = new Date(data[i].posted);
        var dd = d.getDate();
        var mm = d.getMonth() + 1;
        var yyyy = d.getFullYear();
        if (dd < 10) {
          dd = "0" + dd;
        }
        if (mm < 10) {
          mm = "0" + mm;
        }
        var thisdate = mm + "/" + dd + "/" + yyyy;
        //get tags with regexp
        var tagpattern = /(#\w+)/g;
        var m;
        tagstext = "";
        do {
          m = tagpattern.exec(data[i].tags);
          if (m) {
            tagstext += "<p class='postTag'>" + m[1] + "</p>";
          }
        } while (m);
        var content = "<img class='postimg' src='" + imageLink + "'>";
        if(data[i].type == "VIDEO") {
          content = "<video class='postimg' controls loop><source src='" + imageLink + "' type='video/mp4'>Your browser doesnt support video</video>";
        }
        var postData = {
          "userid":data[i].userid,
          "thisdate":thisdate,
          "content": content,
          "tags": tagstext,
          "likes": data[i].likes,
          "caption":data[i].caption,
          "postid":data[i].postid
        };
        var parsedTemplate = executeHTMLTemplate(template, postData);
        $("#posts").html(stuff + parsedTemplate);
        $.ajax({
          url: "/json/user/id",
          method: "POST",
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
    $('p[id="' + id + '"]').html(data[0].username);
  }
}
var isFlwd = false;
var postsRec = false;
function follow(userid, profileid) {
  $.ajax({
    url: "/user/follow",
    data: {
      "userid":userid,
      "profileid":profileid
    },
    success: followedUser
  });
  function followedUser(data) {
    if(data == "Successfully followed this user!") {
      followingAlert();
      isFlwd = true;
    } else {
      isFlwd = false;
    }
  }
}
var template;
$(document).ready(function() {
  $("body").on("contextmenu",function(e){
    return false;
  });
  $.ajax({
    url: "/templates/post",
    success: function(data) {
      template = data;
    }
  });
});
function getUserPosts() {
  if(postsRec) {
    return;
  }
  var uid = parseInt(document.getElementById("profileUserId").innerHTML);
  $.ajax({
    url: '/live/user/posts',
    method: 'POST',
    data: {
      "uid": uid
    },
    success: function(data) {
      postsRec = true;
      var length = data.length;
      $("#posts").html("");
      for(var i = 0; i < length; i++) {
        var imageLink;
        if (data[i].type == "IMAGE") {
          imageLink = "/assets/uploads/imageposts/post" + data[i].postid + data[i].extension;
        } else {
          imageLink = "/assets/uploads/videoposts/post" + data[i].postid + data[i].extension;
        }
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
        var content = "<img class='postimg' src='" + imageLink + "'>";
        if(data[i].type == "VIDEO") {
          content = "<video class='postimg' preload='metadata' controls loop><source src='" + imageLink + "' type='video/mp4'>Your browser doesnt support video</video>";
        }
        var postData = {
          "userid":data[i].userid,
          "thisdate":thisdate,
          "content":content,
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
  }
});
  function putPostUsernames(data) {
    var id = data[0].id;
    $('p[id="' + id + '"]').html(data[0].username);
  }
}
