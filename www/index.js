//  Name of file: index.js
//  Programmers: Colin Chan and Devansh Chopra
//  Team Name: Co.DEsign
//  Changes been made:
//          2018-11-17: created file
// Known Bugs:

function getTremors(since) {
	return fetch("/api/tremors?since=" + since).then(
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

function WeekFunction() {
	var oneWeekAgo = new Date();
	oneWeekAgo.setDate(oneWeekAgo.getDate() - 6);
	oneWeekAgo.setHours(0, 0, 0, 0);
	console.log(oneWeekAgo);

	//document.getElementsByTagName("BODY")[0].style.backgroundColor = "yellow";
	
	getTremors(oneWeekAgo.toISOString()).then(tremors => {
		getMedicines().then(medicines => {
			getExercises().then(exercises => {
				console.log("got exercises")
				makeGraph(tremors, medicines, exercises)
			})
		})
	})
}

function MonthFunction() {

	var oneMonthAgo = new Date();
	oneMonthAgo.setDate(oneMonthAgo.getDate() - 30);
	oneMonthAgo.setHours(0, 0, 0, 0);
	console.log(oneMonthAgo);

	getTremors(oneMonthAgo.toISOString()).then(tremors => {
		getMedicines().then(medicines => {
			getExercises().then(exercises => {
				makeGraph(tremors, medicines, exercises)
			})
		})
	})
	//document.getElementsByTagName("BODY")[0].style.backgroundColor = "yellow";
}

function YearFunction() {
	var oneYearAgo = new Date();
	oneYearAgo.setDate(oneYearAgo.getDate() - 365);
	oneYearAgo.setHours(0, 0, 0, 0);
	console.log(oneYearAgo);

	//var week = "2018-11-08T17:13:39Z"
	getTremors(oneYearAgo.toISOString()).then(tremors => {
		getMedicines().then(medicines => {
			getExercises().then(exercises => {
				makeGraph(tremors, medicines, exercises)
			})
		})
	})
	//document.getElementsByTagName("BODY")[0].style.backgroundColor = "yellow";
}

function f() {
	var date = new Date();
	var lastweek = new Date() - 1;
	var lastmonth = new Date(date.getFullYear(), date.getMonth(), 0);
	alert(firstDay + "===" + lastDay);
}

function getRandomColor() {
	var letters = '0123456789ABCDE';
	var color = '#';
	for (var i = 0; i < 6; i++) {
	  color += letters[Math.floor(Math.random() * 15)];
	}
	return color;
}

function loadCanvas() {
	WeekFunction()
}

function makeGraph(tremors, medicines, exercises) {
	var ctx = document.getElementById("myChart").getContext('2d');

	var scatterChartOptions = {
		type: 'scatter',

		data: {
			datasets: [{
				label: 'resting score',
				fill: false,
				showLine: false, // disable for a single dataset
				pointBorderColor: 'blue',
				backgroundColor: 'blue',
				pointBackgroundColor: 'blue',
			}, {
				label: 'postural score',
				showLine: false, // disable for a single dataset
				pointBorderColor: 'red',
				backgroundColor: 'red',
				pointBackgroundColor: 'red',
			}],
		},
		options: {
			scales: {
				xAxes: [{
					display: true,
					scaleLabel: {
						display: true,
						labelString: 'Time'
					},

					labelString: 'Time',
					type: 'time',
					position: 'bottom'
				}],
				yAxes: [{
					scaleLabel: {
						display: true,
						labelString: 'Severity Score'
					},
					ticks: {
						//suggestedMin: 0, // minimum will be 0, unless there is a lower value.
						// OR //
						beginAtZero: true, // minimum value will be 0.
						max: 10
					}

				}]
			}
		}
	};

	resting = tremors.map(tremor => {
		return {
			x: tremor.date,
			y: tremor.resting / 10
		}
	})
	postural = tremors.map(tremor => {
		return {
			x: tremor.date,
			y: tremor.postural / 10
		}
	})
	
	scatterChartOptions.data.datasets[0].data = resting
	scatterChartOptions.data.datasets[1].data = postural
	
	var y_value = 0.15;
	var offset = 0.15;

	medicines.forEach(medicine => {
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
		var color = getRandomColor();
		scatterChartOptions.data.datasets.push({
			label: medicine.name,
			fill: false,
			showLine: true, // disable for a single dataset
            borderColor: "" + color,
			data: medicineData,
			pointRadius: 0,
			borderWidth: 15,
		});
		y_value += offset;
	})
	exercises.forEach(exercise => {
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
		var color = getRandomColor();
		scatterChartOptions.data.datasets.push({
			label: exercise.name,
			fill: false,
			showLine: true, // disable for a single dataset
            borderColor: "" + color,
			data: exerciseData,
			pointRadius: 0,
			borderWidth: 15,
			steppedLine: true,
		});
		y_value += offset;
	})

	var scatterChart = new Chart(ctx, scatterChartOptions);

	// debug print
	for (var i = 0; i < scatterChart.data.datasets.length; i++) {
		for (var j = 0; j < scatterChart.data.datasets[i].data.length; j++) {
			console.log("data arr ", scatterChart.data.datasets[i].data[j])
		}
	}
}
