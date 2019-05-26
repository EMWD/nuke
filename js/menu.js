class Menu {
    constructor(root, active_class, control_data, panel_data) {
	let controls = root.querySelectorAll("[data-" + control_data + "]");
	let panels = root.querySelectorAll("[data-" + panel_data + "]");
	this.active_class = active_class;

	for(let i = 0; i < controls.length; i++) {
	    let cls = this;
	    let control = controls[i];
	    let panel = panels[i];
	    control.onclick = function() {
		if(this.visible) {
		    cls.hidePanel(control, panel);
		} else {
		    cls.showPanel(control, panel);
		}
		this.visible = !this.visible;
	    }
	    control.onclick.visible = false;
	    this.hidePanel(control, panel);
	}
    }

    showPanel(control, panel) {
	control.classList.add(this.active_class);
	panel.style.display = "block";
	panel.style.visibility = "visible";
    }

    hidePanel(control, panel) {
	control.classList.remove(this.active_class);
	panel.style.display = "none";
	panel.style.visibility = "hidden";
    }
}
