class Slider {
    constructor(container, left, right) {
	this.container = container;
	this.curSlide = 0;
	this.slides = new Array();

	let cls = this;
	left.onclick = () => cls.prevSlide();
	right.onclick = () => cls.nextSlide();
    }

    addSlide(elem) {
	this.slides.push(elem);
	this.hideSlide(this.slides.length-1);
	if(this.slides.length === 1)
	    this.showSlide(0);
    }

    nextSlide() {
	if(this.slides.length === 0)
	    return;
	this.hideSlide(this.curSlide);
	this.curSlide++;
	if(this.curSlide === this.slides.length)
	    this.curSlide = 0;
	this.showSlide(this.curSlide);
    }

    prevSlide() {
	if(this.slides.length === 0)
	    return;
	this.hideSlide(this.curSlide);
	this.curSlide--;
	if(this.curSlide < 0)
	    this.curSlide = this.slides.length-1;
	this.showSlide(this.curSlide);
    }

    showSlide(ind) {
	let slide = this.slides[ind];
	slide.style.display = "block";
	slide.style.visibility = "visible";
    }
    
    hideSlide(ind) {
	let slide = this.slides[ind];
	slide.style.display = "none";
	slide.style.visibility = "hidden";
    }
}
