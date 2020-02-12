function changePublicity() {
	var userID = document.getElementById("dataUserID").innerHTML;
	var publicityText = document.getElementById("publicityStatus").innerHTML;
	var Status;
	if(publicityText == "Public") {
		Status = "mPrivate";
	} else {
		Status = "mPublic";
	}
	$.ajax({
		url: "/settings/user/publicity",
		type: 'POST',
		data: {
			"userID": userID,
			"status": Status
		},
		success: publicityChangeSuccess
	});
}
function changeEmail(email) {
	var userID = document.getElementById("dataUserID").innerHTML;
	$.ajax({
		url: "/settings/user/email",
		type: 'POST',
		data: {
			"userID": userID,
			"email": email
		},
		success: emailChangeSuccess
	});
}
function changeLocation(location) {
	var userID = document.getElementById("dataUserID").innerHTML;
	$.ajax({
		url: "/settings/user/location",
		type: 'POST',
		data: {
			"userID": userID,
			"location": location
		},
		success: locationChangeSuccess
	});
}
function publicityChangeSuccess(data) {
	//Change all Elements here
	document.getElementById("publicityStatus").innerHTML = data;
	if(data == "Public") {
		document.getElementById("publicityStatus").classList.add("badge-success");
		document.getElementById("publicityStatus").classList.remove("badge-danger");
	} else {
		document.getElementById("publicityStatus").classList.add("badge-danger");
		document.getElementById("publicityStatus").classList.remove("badge-success");
	}
}
function emailChangeSuccess() {
	console.log("Changed email");
}
function locationChangeSuccess(data) {
	console.log("Change Location: " + data);
}
