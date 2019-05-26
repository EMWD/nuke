document.addEventListener("DOMContentLoaded", function() {
    document.getElementById("change-pass-btn").onclick = function() {
	let old_password = document.getElementById("old-password").value;
	let new_password = document.getElementById("new-password").value;
	let change_error = document.getElementById("change-error");
	let del_error = document.getElementById("del-error");

	postRequest("/change_password_rest", "old_password=" + old_password + "&new_password=" + new_password, function(req) {
	    if(req.readyState != 4) return;

	    if(req.status != 200) {
		change_error.innerHTML = req.responseText;
		change_error.style.color = "red";
		return;
	    }

	    change_error.innerHTML = "Готово";
	    change_error.style.color = "green";
	});
    }

    document.getElementById("del-btn").onclick = function() {
	let del_password = document.getElementById("del-password").value;
	postRequest("/delete_account_rest", "password=" + del_password, function(req) {
	    if(req.readyState != 4) return;

	    if(req.status != 200) {
		del_error.innerHTML = req.responseText;
		del_error.style.color = "red";
		return;
	    }
	    
	    window.location.href = "/index";
	});
    }
});

function showPresentation(login, name) {
    window.location.href = "/showpresentation?user=" + login + "&title=" + name;
}

function delPresentation(ind) {
    let presentation = document.querySelectorAll(".presentation")[ind];
    let title = presentation.dataset.title;
    postRequest("/del_presentation_rest", "title=" + title, function(req) {
	if(req.readyState != 4) return;

	if(req.status != 200) {
	    alert(req.status + ": " + req.statusText);
	    return;
	}

	presentation.remove();
    });
}

function editPresentation(login, ind) {
    let title = document.querySelectorAll(".presentation")[ind].dataset.title;
    window.location.href = "/editor?update=true&title=" + title + "&user=" + login;
}
