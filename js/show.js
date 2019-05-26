document.addEventListener("DOMContentLoaded", function() {
    let params = getJsonFromUrl(window.href);
    let user = params["user"].trim();
    let title = params["title"].trim();

    if(user.length == 0 || title.length == 0)
	return;

    postRequest("get_presentation_rest", "user=" + user + "&title=" + title, function(req) {
	if(req.readyState != 4) return;

	if(req.status != 200) {
	    document.write("404 not found");
	    return;
	}

	let presentation = JSON.parse(req.responseText);
	let cnv = new showdown.Converter();
	cnv.setOption("strikethrough", true);

	let left_btn = document.getElementById("left-btn");
	let right_btn= document.getElementById("right-btn");
	let viewer = document.getElementById("viewer-content");
	let slider = new Slider(viewer, left_btn, right_btn);

	viewer.classList.add(presentation.Style);
	
	let slides = presentation.Presentation.split("{{slide}}").filter(f => f.trim() != "").map(f => cnv.makeHtml(f));
	slides.forEach(html => {
	    let footer = document.createElement("div");
	    footer.classList.add("pres-footer");
	    footer.innerHTML = title;
	    
	    let div = document.createElement("div");
	    div.classList.add("slide");
	    div.innerHTML = html;
	    div.appendChild(footer);
	    viewer.appendChild(div);
	    slider.addSlide(div);
	});
    });
});

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
