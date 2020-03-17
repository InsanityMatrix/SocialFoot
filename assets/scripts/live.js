var allPosts;
var currPost;
var template;
function signOut() {
  //ajax to signout for the server to interpret
  $.ajax({
    url: "/settings/user/signout",
    success: function() {   window.location = "/assets/"; }
  });
}
function loadCustomFeed() {
  $.ajax({
    url: '/json/feed/custom',
    success: serveCustomFeed
  });
}
function serveCustomFeed(data) {
  //Data is already parsed as JSON
  allPosts = data;
  var length = data.length;
  if (length == 0) {
    currPost = 0;
    $("#posts").html("You are not Following anybody");
  }
  currPost = data.length;
  if(data.length > 10) {
    length = 10;
    currPost = 10;
  }
  $("#posts").html("");
  for(var i = 0; i < length; i++) {
    var imageLink;
    if (data[i].Type == "IMAGE") {
      imageLink = "/assets/uploads/imageposts/post" + data[i].Postid + data[i].Extension;
    } else {
      imageLink = "/assets/uploads/videoposts/post" + data[i].Postid + data[i].Extension;
    }
    var stuff = $("#posts").html();
    //Parse Date
    var d = new Date(data[i].Posted);
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
      m = tagpattern.exec(data[i].Tags);
      if (m) {
        tagstext += "<p class='postTag'>" + m[1] + "</p>";
      }
    } while (m);
    var content = "<img class='postimg' src='" + imageLink + "'>";
    if(data[i].Type == "VIDEO") {
      content = "<video class='postimg' controls loop><source src='" + imageLink + "' type='video/mp4'>Your browser doesnt support video</video>";
    }
    var postData = {
      "userid":data[i].Userid,
      "thisdate":thisdate,
      "content": content,
      "tags": tagstext,
      "likes": data[i].Likes,
      "caption":data[i].Caption,
      "postid":data[i].Postid
    };
    var parsedTemplate = executeHTMLTemplate(template, postData);
    $("#posts").html(stuff + parsedTemplate);
    $.ajax({
      url: "/json/user/id",
      method: "POST",
      data : {
        "userid": data[i].Userid
      },
      success: putPostUsernames
    });
  }

  function putPostUsernames(data) {
    var id = data[0].id;
    $('p[id="' + id + '"]').html(data[0].username);
  }
}
function loadFeed() {
  $.ajax({
    url: "/templates/post",
    success: serveFeed
  });
}
function serveFeed(templ) {
  //First get Public Posts
  template = templ;
  getPublicPosts();

  function getPublicPosts() {
    $.getJSON("/posts/public", function(data){
      allPosts = data;
      var length = data.length;
      currPost = data.length;
      if(data.length > 10) {
        length = 10;
        currPost = 10;
      }
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

function loadNextPosts() {
  let data = allPosts;
  if (currPost < data.length) {
    let length = data.length;
    if(currPost + 10 <= data.length) {
      length = currPost + 10;
    }

    for(var i = currPost; i < length; i++) {
      currPost++;
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
  } else {
    console.log("No more posts");

  }

  function putPostUsernames(dat) {
    var id = dat[0].id;
    $('p[id="' + id + '"]').html(dat[0].username);
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
