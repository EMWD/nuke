document.addEventListener("DOMContentLoaded", function() {
    document.querySelector(".submit").onclick = function() {
	let login = document.querySelector("input[name='login']").value;
	let password = document.querySelector("input[name='password']").value;
	let error = document.querySelector("#error");
	
	postRequest("signup_rest", "login=" + login + "&password=" + password, function(req) {
	    if(req.readyState != 4) return;

	    if(req.status == 200) {
		window.location.href = "/index";
	    } else {
		error.innerHTML = req.responseText;
		error.style.color = "red";
	    }
	});
    }
});
