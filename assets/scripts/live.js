function signOut() {
  //ajax to signout for the server to interpret
  $.ajax({
    url: '/settings/user/signout'
  });
  window.location = "/assets/";
}
