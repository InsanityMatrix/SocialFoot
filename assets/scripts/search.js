$("input[name='search']").keypress(function(e){
    var searchValue= $("input[name='search']").val();
    searchFor(searchValue);
});

function searchFor(searchValue) {
    //DO AJAX HERE
    $.ajax({
		url: "/search",
		method: 'POST',
		data: {
			"term": searchValue
		},
		success: resultsSuccess
	});
}

function resultsSuccess(data) {
    //Parse JSON and display results
    console.log(data);
}
