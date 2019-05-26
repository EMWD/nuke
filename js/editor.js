document.addEventListener("DOMContentLoaded", function() {
    let params = getJsonFromUrl(window.href);
    let update = params["update"];
    let old_title = params["title"];
    let ed_user = params["user"];

    if(update === "true" && old_title !== undefined && ed_user !== undefined) {
	let textarea = document.getElementById("editor-textarea");
	postRequest("/get_presentation_rest", "user=" + ed_user + "&title="+old_title, function(req) {
	    if(req.readyState != 4) return;
	    if(req.status == 200) {
		let presentation = JSON.parse(req.responseText);
		document.getElementById("pres-title").value = old_title;
		document.getElementById("send-btn").innerHTML = "Обновить";
		textarea.value = presentation.Presentation;
	    } else {
		alert("презентации " + old_title + " не существует");
		window.location.href = "/index";
	    }
	});
    }

    new Menu(document.getElementById("editor-menu"), "editor-active", "mcontrol", "menu");
    document.getElementById("view-btn").onclick = function() {
	let src = document.getElementById("editor-textarea").value;
	let cnv = new showdown.Converter();
	cnv.setOption("strikethrough", true);

	let left_btn = document.getElementById("left-btn");
	let right_btn= document.getElementById("right-btn");
	let viewer = document.getElementById("viewer-content");
	let slider = new Slider(viewer, left_btn, right_btn);

	let slides = src.split("{{slide}}").filter(f => f.trim() !== "").map(f => cnv.makeHtml(f));
	viewer.innerHTML = "";
	slides.forEach(html => {
	    let footer = document.createElement("div");
	    let title = document.getElementById("pres-title");
	    footer.classList.add("pres-footer");
	    footer.innerHTML = title.value;
	    
	    let div = document.createElement("div");
	    div.classList.add("slide");
	    div.innerHTML = html;
	    div.appendChild(footer);
	    viewer.appendChild(div);
	    slider.addSlide(div);
	});
    }

    document.getElementById("pres-style").onchange = function() {
	let value = this.options[this.selectedIndex].value;
	let viewer = document.getElementById("viewer-content");
	viewer.className = "";
	viewer.classList.add(value);
    }

    document.getElementById("pres-title").onchange = function() {
	let footers = document.querySelectorAll(".pres-footer");
	for(let i = 0; i < footers.length; i++) {
	    footers[i].innerHTML = this.value;
	}
    }

    document.getElementById("send-btn").onclick = function() {
	let styles = document.getElementById("pres-style");
	
	let title = document.getElementById("pres-title").value;
	let style = styles.options[styles.selectedIndex].value;
	let code = document.getElementById("editor-textarea").value;

	if(update === "true" && old_title !== undefined && ed_user !== undefined) {
	    postRequest("/update_presentation_rest", "title=" + old_title + "&new_title=" + title + "&new_style=" + style + "&new_code=" + code, function(req) {
		if(req.readyState != 4) return;
		
		let error = document.getElementById("error");
		if(req.status == 200) {
		    error.style.color = "green";
		    error.innerHTML = "Готово";
		} else {
		    error.innerHTML = req.responseText;
		    error.style.color = "red";
		}
	    });
	} else {
	    postRequest("/add_presentation_rest", "title=" + title + "&style=" + style + "&code=" + code, function(req) {
		if(req.readyState != 4) return;
		
		let error = document.getElementById("error");
		if(req.status == 200) {
		    error.style.color = "green";
		    error.innerHTML = "Готово";
		} else {
		    error.innerHTML = req.responseText;
		    error.style.color = "red";
		}
	    });
	}
    }
});

function add_header(size) {
    let textarea = document.getElementById("editor-textarea");
    textarea.value += "#".repeat(size) + " ";
    textarea.focus();
}

function add_bold() {
    let textarea = document.getElementById("editor-textarea");
    textarea.value += "****";
    textarea.selectionEnd-=2;
    textarea.focus();
}

function add_italic() {
    let textarea = document.getElementById("editor-textarea");
    textarea.value += "**";
    textarea.selectionEnd-=1;
    textarea.focus();
}

function add_underline() {
    let textarea = document.getElementById("editor-textarea");
    textarea.value += "<u></u>";
    textarea.selectionEnd-=4;
    textarea.focus();
}

function add_strike() {
    let textarea = document.getElementById("editor-textarea");
    textarea.value += "~~~~";
    textarea.selectionEnd-=2;
    textarea.focus();
}

function add_quote() {
    let textarea = document.getElementById("editor-textarea");
    textarea.value += "> ";
    textarea.focus();
}

function add_link() {
    let textarea = document.getElementById("editor-textarea");
    let link = prompt("Адрес", "https://");
    if(link === null) return;
    textarea.value += "[](" + link + ")";
    textarea.selectionEnd-=(link.length+3);
    textarea.focus();
}

function add_ul() {
    let textarea = document.getElementById("editor-textarea");
    textarea.value += "- "
    textarea.focus();
}

function add_ol() {
    let textarea = document.getElementById("editor-textarea");
    textarea.value += "1. "
    textarea.focus();
}

function add_image() {
    let textarea = document.getElementById("editor-textarea");
    let link = prompt("Адрес", "https://");
    if(link === null) return;
    textarea.value += "![image](" + link + ")";
    textarea.selectionEnd-=(link.length+3);
    textarea.focus();
}

function add_center() {
    let textarea = document.getElementById("editor-textarea");
    textarea.value += "<span></span>"
    textarea.selectionEnd-=7;
    textarea.focus();
}

function add_slide() {
    let textarea = document.getElementById("editor-textarea");
    textarea.value += "\n{{slide}}\n";
    textarea.focus();
}

function getJsonFromUrl(url) {
  if(!url) url = location.search;
  var query = url.substr(1);
  var result = {};
  query.split("&").forEach(function(part) {
    var item = part.split("=");
    result[item[0]] = decodeURIComponent(item[1]);
  });
  return result;
}
