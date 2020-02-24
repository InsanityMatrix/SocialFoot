$("input[name='search']").keypress(function(e){
    var searchValue= $("input[name='search']").val();
    searchFor(searchValue);
});
var resultTemplate;
$(document).ready(function() {
  $.ajax({
    url: '/templates/result',
    success: function(data) { resultTemplate = data; }
  })
});
function search() {
  var searchValue = $("input[name='search']").val();
  searchFor(searchValue);
}
function searchFor(searchValue) {
    //DO AJAX HERE
    $.ajax({
		url: "/search",
		type: 'POST',
		data: {
			"term": searchValue
		},
		success: resultsSuccess
	});
}

function resultsSuccess(data) {
    var length = data.length;
    $("#search-results").html("");
    for(var i = 0; i < length; i++) {
      var results = $("#search-results").html();
      var gender;
      if(data[i].gender) {
        gender = "Male";
      } else {
        gender = "Female";
      }
      var userData = {
        "userid": data[i].id,
        "username": data[i].username,
        "gender": gender
      };
      var parsedTemplate = executeHTMLTemplate(resultTemplate, userData);
      $("#search-results").html(results + parsedTemplate);
    }
}
