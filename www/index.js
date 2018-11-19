//  Name of file: index.js
//  Programmers: Colin Chan and Devansh Chopra and Nic Klaassen
//  Team Name: Co.DEsign
//  Changes been made:
//          2018-11-17: created file
// Known Bugs:

var MAX_COLOURS = 50;
var colourIndex = 0;
var global_chart = null

function loadCanvas() {
	fillRandomColor()
	loadData().then(data => {
		makeChart(data)
		WeekFunction()
	})
}

function WeekFunction() {
	let oneWeekAgo = new Date();
	oneWeekAgo.setDate(oneWeekAgo.getDate() - 6);
	oneWeekAgo.setHours(0, 0, 0, 0);

	global_chart.options.scales.xAxes[0].time.min = oneWeekAgo
	global_chart.update()
}

function MonthFunction() {
	let oneMonthAgo = new Date();
	oneMonthAgo.setDate(oneMonthAgo.getDate() - 30);
	oneMonthAgo.setHours(0, 0, 0, 0);

	global_chart.options.scales.xAxes[0].time.min = oneMonthAgo
	global_chart.update()
}

function YearFunction() {
	let oneYearAgo = new Date();
	oneYearAgo.setDate(oneYearAgo.getDate() - 365);
	oneYearAgo.setHours(0, 0, 0, 0);

	global_chart.options.scales.xAxes[0].time.min = oneYearAgo
	global_chart.update()
}

function getTremors() {
	return fetch("/api/tremors").then(
		response => response.json()
	)
}

function getMedicines() {
	return fetch("/api/meds").then(
		response => response.json()
	)
}

function getExercises() {
	return fetch("/api/exercises").then(
		response => response.json()
	)
}

function Data(tremors, medicines, exercises) {
	this.tremors = tremors
	this.medicines = medicines
	this.exercises = exercises
}

async function loadData() {
	// load all data in parallel
	tremorPromise = getTremors()
	medicinePromise = getMedicines()
	exercisePromise = getExercises()

	tremors = await tremorPromise.catch(err => console.log("failed to get tremors", err))
	medicines = await medicinePromise.catch(err => console.log("failed to get medicines", err))
	exercises = await exercisePromise.catch(err => console.log("failed to get exercises", err))

	return new Data(tremors, medicines, exercises)
}

function getRandomColor() {
	const letters = '0123456789ABCDE';
	let color = '#';
	for (let i = 0; i < 6; i++) {
		color += letters[Math.floor(Math.random() * 15)];
	}
	return color;
}

var colours = [];
function fillRandomColor() {
	let i = 0;
	while (i < 50) {
		colours.push("" + getRandomColor());
		i++;
	}
}

function makeChart(data) {
	let scatterChartOptions = {
		type: 'scatter',
		data: {
			datasets: [{
				label: 'resting score',
				fill: false,
				showLine: true,
				borderColor: '#00f5',
				pointBorderColor: 'blue',
				backgroundColor: 'blue',
				pointBackgroundColor: 'blue',
			}, {
				label: 'postural score',
				fill: false,
				showLine: true,
				borderColor: '#f005',
				pointBorderColor: 'red',
				backgroundColor: 'red',
				pointBackgroundColor: 'red',
			}],
		},
		options: {
			maintainAspectRatio: false,
			scales: {
				xAxes: [{
					type: 'time',
					scaleLabel: {
						display: true,
						labelString: 'Time'
					},
					time: {
						displayFormats: {
							hour: 'ddd, hA'
						}
					}
				}],
				yAxes: [{
					scaleLabel: {
						display: true,
						labelString: 'Severity Score'
					},
					ticks: {
						beginAtZero: true, // minimum value will be 0.
						max: 10
					}

				}]
			}
		}
	};

	resting = data.tremors.map(tremor => {
		return {
			x: tremor.date,
			y: tremor.resting / 10
		}
	})
	postural = data.tremors.map(tremor => {
		return {
			x: tremor.date,
			y: tremor.postural / 10
		}
	})
	scatterChartOptions.data.datasets[0].data = resting
	scatterChartOptions.data.datasets[1].data = postural
	
	var y_value = 0.15;
	var offset = 0.3;
	colourIndex = 0;
	data.medicines.forEach(medicine => {
		// if medicine.enddate is not set, treat it as today's date
		var enddate = new Date();
		enddate = enddate.toISOString()
		if (medicine.enddate != null) {
			enddate = medicine.enddate;
		}
		medicineData = [
		{
			x: medicine.startdate,
			y: y_value
		}, {
			x: enddate,
			y: y_value
		}]
		scatterChartOptions.data.datasets.push({
			label: medicine.name,
			fill: false,
			showLine: true, // disable for a single dataset
			borderColor: "" + colours[colourIndex],
			data: medicineData,
			pointRadius: 0,
			borderWidth: 15,
		});
		y_value += offset;
		if (colourIndex+1 <= MAX_COLOURS)
			colourIndex += 1;
	})
	data.exercises.forEach(exercise => {
		// if exercise.enddate is not set, treat it as today's date
		var enddate = new Date();
		enddate = enddate.toISOString()
		if (exercise.enddate != null) {
			enddate = exercise.enddate;
		}
		exerciseData = [
		{
			x: exercise.startdate,
			y: y_value
		}, {
			x: enddate,
			y: y_value
		}]
		scatterChartOptions.data.datasets.push({
			label: exercise.name,
			fill: false,
			showLine: true, // disable for a single dataset
			borderColor: "" + colours[colourIndex],
			data: exerciseData,
			pointRadius: 0,
			borderWidth: 15,
		});
		y_value += offset;
		if (colourIndex+1 <= MAX_COLOURS)
			colourIndex += 1
	})

	let ctx = document.getElementById("myChart").getContext('2d');
	global_chart = new Chart(ctx, scatterChartOptions);
}
