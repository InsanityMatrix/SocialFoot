function signOut() {
  //ajax to signout for the server to interpret
  $.ajax({
    url: '/settings/user/signout'
  });
  window.location = "/assets/";
}

function getPublicPosts() {
  $.getJSON("/posts/public", function(data){
    console.log(data);
  });
}
