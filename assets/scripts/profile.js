function changePublicity() {
	var username = document.getElementById("dataUsername").innerHTML;
	var publicityText = document.getElementById("publicityStatus").innerHTML;
	var Status;
	if(publicityText == "Public") {
		Status = true;
	} else {
		Status = false;
	}
	$.ajax({
		url: "/settings/user/publicity",
		data: {
			"username": username,
			"status": Status;
		},
		success: publicityChangeSuccess
	});
}
function publicityChangeSuccess() {
	
}