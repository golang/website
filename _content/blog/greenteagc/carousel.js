document.addEventListener('DOMContentLoaded', function() {
	setupFor("marksweep");
	setupFor("greentea");
});

function setupFor(id) {
	const next = document.getElementById(id+"-next");
	const prev = document.getElementById(id+"-prev");
	const caro = document.getElementById(id);
	next.addEventListener('click', scrollRight(next, prev, caro));
	prev.addEventListener('click', scrollLeft(next, prev, caro));
	prev.disabled = true;
	next.disabled = false;
	prev.hidden = false;
	next.hidden = false;
	caro.classList.add('hide-overflow');
}

function scrollRight(n, p, c) {
	return () => {
		c.scrollLeft += c.getBoundingClientRect().width;
		p.disabled = false;
		if (c.scrollLeft === c.scrollWidth-c.clientWidth) {
			n.disabled = true;
		}
	};
}

function scrollLeft(n, p, c) {
	return () => {
		c.scrollLeft -= c.getBoundingClientRect().width;
		n.disabled = false;
		if (c.scrollLeft === 0) {
			p.disabled = true;
		}
	};
}
