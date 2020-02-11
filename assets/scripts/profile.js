function changePublicity() {
	var username = document.getElementById("dataUsername").innerHTML;
	var publicityText = document.getElementById("publicityStatus").innerHTML;
	var status;
	if(publicityText == "Public") {
		status = true;
	} else {
		status = false;
	}
	$.ajax({
		url: "/settings/user/publicity",
		data: {
			
		},
	});
}
