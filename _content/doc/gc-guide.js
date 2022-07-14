function StackedAreaChart({
	xSeries,
	marginTop = 30, // top margin, in pixels
	marginRight = 100, // right margin, in pixels
	marginBottom = 60, // bottom margin, in pixels
	marginLeft = 50, // left margin, in pixels
} = {}) {
	const width = 756;
	const height = 189;
	const svg = d3.create("svg")
		.attr("preserveAspectRatio", "xMinYMin meet")
		.attr("viewBox", [0, 0, width, height]);

	const xRange = [marginLeft, width - marginRight]; // [left, right]
	const yRange = [height - marginBottom, marginTop]; // [bottom, top]

	// Add empty axes first.

	svg.append("g")
		.classed("x axis", true)
		.attr("transform", `translate(0,${height-marginBottom})`)

	svg.append("g")
		.classed("y axis", true)
		.attr("transform", `translate(${marginLeft},0)`)

	const update = function(data, mutTime) {
		const seriesKeys = Object.keys(data[0]).filter(s => s !== xSeries);

		let seriesColors = new Array();
		if (seriesKeys.length > 3) {
			const colorFn = d3.interpolateViridis;
			for (let i = 0; i < seriesKeys.length; i++) {
				seriesColors.push(colorFn((i / 10) - Math.floor(i/10)));
			}
		} else {
			seriesColors = ["#253443", "#007d9c", "#50b7e0"];
			if (seriesKeys.length < 3) {
				seriesColors = seriesColors.slice(seriesKeys.length-1);
			}
		}
		const seriesScale = d3.scaleOrdinal()
			.domain(seriesKeys)
			.range(seriesColors);

		const yStack = (d3.stack().keys(seriesKeys))(data);

		const xDomain = d3.extent(d3.map(data, p => p[xSeries]));
		const yDomain = d3.extent(d3.map(yStack[yStack.length-1], p => p[1]));
		yDomain[0] = 0;

		const xScale = d3.scaleLinear(xDomain, xRange);
		const yScale = d3.scaleLinear(yDomain, yRange);

		const xAxis = d3.axisBottom(xScale).tickFormat(x => `${x.toFixed(1)} s`);
		svg.selectAll("g.x.axis")
			.style("font-size", "11px")
			.call(xAxis);

		const yAxis = d3.axisLeft(yScale).ticks(5).tickFormat(x => `${x.toFixed(0)} MiB`);
		svg.selectAll("g.y.axis")
			.style("font-size", "11px")
			.call(yAxis);

		const area = d3.area()
			.curve(d3.curveLinear)
			.x(d => xScale(d.data[xSeries]))
			.y0(d => yScale(d[0]))
			.y1(d => yScale(d[1]));

		svg.selectAll("path.series")
			.data(yStack)
			.join("path")
				.classed("series", true)
				.attr("d", area)
				.style("fill", d => seriesScale(d.key));

		svg.selectAll("text.label")
			.data(seriesKeys)
			.join("text")
				.classed("label", true)
				.attr("text-anchor", "left")
				.attr("font-size", "12px")
				.attr("x", width-marginRight+20)
				.attr("y", d => (seriesKeys.length-1-seriesKeys.indexOf(d))*24+60)
				.attr("fill", "currentColor")
				.attr("display", (() => {
					if (seriesKeys.length <= 3) {
						return "inherit";
					}
					return "none";
				})())
				.text(d => d);

		svg.selectAll("rect.legend")
			.data(seriesKeys)
			.join("rect")
				.classed("legend", true)
				.attr("stroke", "none")
				.attr("x", width-marginRight+7)
				.attr("y", d => (seriesKeys.length-1-seriesKeys.indexOf(d))*24+51)
				.attr("width", 10)
				.attr("height", 10)
				.attr("display", (() => {
					if (seriesKeys.length <= 3) {
						return "inherit";
					}
					return "none";
				})())
				.attr("fill", d => seriesScale(d));

		svg.selectAll("text.duration")
			.data([xDomain[1]])
			.join("text")
				.classed("duration", true)
				.attr("text-anchor", "left")
				.attr("font-size", "10px")
				.attr("x", width-marginRight+5)
				.attr("y", height-marginBottom+10)
				.attr("fill", "currentColor")
				.attr("font-weight", "bold")
				.text(d => `Total: ${d.toFixed(2)} s`);

		svg.selectAll("text.results")
			.data([[(xDomain[1]-mutTime)/xDomain[1]*100, yDomain[1]]])
			.join("text")
				.classed("results", true)
				.attr("text-anchor", "middle")
				.attr("font-size", "12px")
				.attr("x", marginLeft + (width-marginLeft-marginRight)/2)
				.attr("y", height-marginBottom+37)
				.attr("fill", "currentColor")
				.attr("font-weight", "bold")
				.text(d => `GC CPU = ${d[0].toFixed(1)}%, Peak Mem = ${d[1].toFixed(1)} MiB`);

		const peakLive = d3.max(d3.map(data, p => p["Live Heap"]));
		const otherMem = d3.max(d3.map(data, p => p["Other Mem."]));

		svg.selectAll("text.subresults")
			.data([[peakLive, otherMem]])
			.join("text")
				.classed("subresults", true)
				.attr("text-anchor", "middle")
				.attr("font-size", "11px")
				.attr("x", marginLeft + (width-marginLeft-marginRight)/2)
				.attr("y", height-marginBottom+51)
				.attr("fill", "currentColor")
				.text(d => {
					let base = "";
					if (d[0]) {
						base += `Peak Live Mem = ${d[0].toFixed(1)} MiB`;
					}
					if (d[1]) {
						base += `, Other Mem = ${d[1].toFixed(1)} MiB`;
					}
					if (base !== "") {
						base = "(" + base + ")";
					}
					return base;
				});
	}
	return [svg.node(), update];
}

function gcModel(workload, config) {
	let otherMem = config["otherMem"];
	if (typeof(otherMem) !== 'number') {
		otherMem = document.getElementById(config["otherMem"]).value;
	}
	let gogc = config["GOGC"];
	if (typeof(gogc) !== 'number') {
		gogc = document.getElementById(config["GOGC"]).value;
	}
	let memoryLimit = config["memoryLimit"];
	if (typeof(memoryLimit) !== 'number') {
		memoryLimit = document.getElementById(config["memoryLimit"]).value;
	}
	let initialLive = 0;
	if ("initialLive" in config) {
		initialLive = config["initialLive"];
	}
	let trackLive = false;
	if ("trackLive" in config) {
		trackLive = config["trackLive"];
		if (typeof(trackLive) !== 'boolean') {
			trackLive = document.getElementById(config["trackLive"]).checked;
		}
	}
	let fixedWindow = Infinity;
	if ("fixedWindow" in config) {
		fixedWindow = config["fixedWindow"];
	}
	const data = new Array();

	// State.
	const minHeapGoal = 4; // MiB
	let t = 0;
	let liveHeap = initialLive;
	let newHeap = 0;
	let liveFromCycle = new Array();
	liveFromCycle.push(initialLive);
	liveFromCycle.push(0);

	const computeHeapGoal = (liveHeap) => {
		let heapGoal = liveHeap*(1.0 + (gogc / 100));
		if (gogc === Infinity) {
			heapGoal = Infinity;
		}
		if (heapGoal+otherMem > memoryLimit) {
			heapGoal = memoryLimit - otherMem
		}
		if (gogc !== Infinity && heapGoal < minHeapGoal) {
			heapGoal = minHeapGoal
		}
		if (heapGoal < liveHeap + 0.0625) {
			heapGoal = liveHeap + 0.0625
		}
		return heapGoal
	}
	let heapGoal = computeHeapGoal(minHeapGoal / (1 + gogc/100)); // Fake a live heap for minHeapGoal.
	if (initialLive !== 0) {
		heapGoal = computeHeapGoal(initialLive);
	}

	let n = 0;
	const emit = function() {
		const datum = {"t": t};
		// The series will be automatically stacked, so for the best
		// possible presentation, we should make sure to put in
		// "other mem" first, then "live," then "new."
		// This is roughly in order of "least dynamic" series
		// to "most dynamic" which helps make the graph easier to
		// interpret.
		if (otherMem !== 0) {
			datum["Other Mem."] = otherMem;
		}
		if (trackLive) {
			for (let i = 0; i < liveFromCycle.length; i++) {
				datum[`Live Heap From GC ${i+1}`] = liveFromCycle[i];
			}
		} else {
			datum["Live Heap"] = liveHeap;
			datum["New Heap"] = newHeap;
		}
		data.push(datum)
	}

	// Emit points.
	emit();
	let nextLive = 0;
	let nextWillLive = 0;
	let nextWillDie = 0;
	let totalMutTime = 0;
	for (const work of workload) {
		let left = work.duration;
		let lastLive = liveHeap + nextLive;
		const willLive = work.duration * work.allocRate * work.newSurvivalRate;
		const willDie = lastLive * work.oldDeathRate;
		while (left > 0) {
			if (t >= fixedWindow) {
				break;
			} else if (t + left > fixedWindow) {
				left = fixedWindow - t;
			}
			let alloc = left * work.allocRate;
			let endCycle = false;
			if (liveHeap+newHeap+alloc > heapGoal) {
				alloc = heapGoal-liveHeap-newHeap;
				endCycle = true;
			}
			newHeap += alloc;

			// Calculate mutator time.
			const mutTime = alloc / work.allocRate;
			left -= mutTime;
			t += mutTime;
			totalMutTime += mutTime;
			nextLive += (willLive - willDie) * (mutTime / work.duration);

			// For tracking per-GC live memory.
			nextWillLive += willLive * (mutTime / work.duration);
			nextWillDie += willDie * (mutTime / work.duration);
			liveFromCycle[liveFromCycle.length-1] = newHeap;

			if (endCycle) {
				emit();

				liveHeap += nextLive;
				for (let i = 0; i < liveFromCycle.length; i++) {
					const live = liveFromCycle[i];
					if (live > 0) {
						if (live > nextWillDie) {
							liveFromCycle[i] -= nextWillDie;
							nextWillDie = 0;
							break;
						}
						nextWillDie -= live;
						liveFromCycle[i] = 0;
					}
				}
				liveFromCycle[liveFromCycle.length-1] = nextWillLive;

				nextLive = 0;
				nextWillLive = 0;
				nextWillDie = 0;
				newHeap = 0;
				const gcTime = liveHeap / work.scanRate + config.fixedCost;
				t += gcTime;

				emit();

				heapGoal = computeHeapGoal(liveHeap)

				liveFromCycle.push(newHeap);
			}
		}
		emit();
	}
	if (trackLive) {
		for (let i = 0; i < data.length; i++) {
			for (let j = 0; j < liveFromCycle.length; j++) {
				const key = `Live Heap From GC ${j+1}`;
				if (!(key in data[i])) {
					data[i][key] = 0;
				}
			}
		}
	}
	return [data, totalMutTime];
}

const graphs = document.querySelectorAll('.gc-guide-graph');

for (let i = 0; i < graphs.length; i++) {
	const workload = JSON.parse(graphs[i].getAttribute("data-workload"));
	const config = JSON.parse(graphs[i].getAttribute("data-config"));
	const [chart, update] = StackedAreaChart({xSeries: "t"});

	const setupSlider = function(parameter, f, fmt) {
		if (typeof(config[parameter]) !== 'number') {
			const id = config[parameter];
			const slider = document.getElementById(id);
			const display = document.getElementById(id+"-display");
			const value = f(slider.value);

			if (display) {
				display.innerHTML = fmt(value);
			}
			config[parameter] = value;

			slider.oninput = function() {
				const value = f(this.value);

				if (display) {
					display.innerHTML = fmt(value);
				}
				config[parameter] = value;

				const [data, mutTime] = gcModel(workload, config);
				update(data, mutTime);
			}
		}
	};
	const setupCheckbox = function(parameter) {
		if (parameter in config && typeof(config[parameter]) !== 'boolean') {
			const id = config[parameter];
			const checkbox = document.getElementById(id);

			config[parameter] = checkbox.checked;

			checkbox.oninput = function() {
				config[parameter] = checkbox.checked;

				const [data, mutTime] = gcModel(workload, config);
				update(data, mutTime);
			}
		}
	};
	setupSlider("otherMem", x => parseInt(x), x => `${x} MiB`);
	setupSlider("GOGC", x => {
		const v = Math.round(Math.pow(2, parseFloat(x)))
		if (v >= 1024) {
			return Infinity;
		}
		return v;
	}, x => {
		if (x === Infinity) {
			return "off";
		}
		return `${x}`;
	});
	setupSlider("memoryLimit", x => parseFloat(x), x => `${x.toFixed(1)} MiB`);
	setupCheckbox("trackLive");

	const [data, mutTime] = gcModel(workload, config);
	update(data, mutTime);
	graphs[i].appendChild(chart);
}
