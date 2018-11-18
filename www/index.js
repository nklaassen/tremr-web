function getTremors(since, callback) {
	//function myFunction()
	fetch("/api/tremors?since=" + since)
		.then(response => response.json(),
			error => console.log(error))
		.then(
			result => {
				callback(result)
			},
			error => {
				console.log("error")
			}
		)
}

function WeekFunction() {
	var oneWeekAgo = new Date();
	oneWeekAgo.setDate(oneWeekAgo.getDate() - 6);
	oneWeekAgo.setHours(0, 0, 0, 0);
	console.log(oneWeekAgo);

	//var week = "2018-11-08T17:13:39Z"
	getTremors(oneWeekAgo.toISOString(), tremors => {
		makeGraph(tremors)
	})
	//document.getElementsByTagName("BODY")[0].style.backgroundColor = "yellow";
}

function MonthFunction() {

	var oneMonthAgo = new Date();
	oneMonthAgo.setDate(oneMonthAgo.getDate() - 30);
	oneMonthAgo.setHours(0, 0, 0, 0);
	console.log(oneMonthAgo);

	//var week = "2018-11-08T17:13:39Z"
	getTremors(oneMonthAgo.toISOString(), tremors => {
		makeGraph(tremors)
	})
	//document.getElementsByTagName("BODY")[0].style.backgroundColor = "yellow";
}

function YearFunction() {
	var oneYearAgo = new Date();
	oneYearAgo.setDate(oneYearAgo.getDate() - 365);
	oneYearAgo.setHours(0, 0, 0, 0);
	console.log(oneYearAgo);

	//var week = "2018-11-08T17:13:39Z"
	getTremors(oneYearAgo.toISOString(), tremors => {
		makeGraph(tremors)
	})
	//document.getElementsByTagName("BODY")[0].style.backgroundColor = "yellow";
}

function f() {
	var date = new Date();
	var lastweek = new Date() - 1;
	var lastmonth = new Date(date.getFullYear(), date.getMonth(), 0);
	alert(firstDay + "===" + lastDay);
}



function loadCanvas() {
	WeekFunction()
}

function makeGraph(tremors) {

	var ctx = document.getElementById("myChart").getContext('2d');

	var scatterChartOptions = {
		type: 'scatter',

		data: {
			datasets: [{
				label: 'resting score',
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

	var date = new Date();
	var lastweek = new Date() - 1;
	var firstDay = new Date(date.getFullYear(), date.getMonth() - 1, 0);

	//console.log(date);
	//console.log(firstDay);
	// var lastmonth = new Date(date.getFullYear(), date.getMonth(), 0);
	// alert(firstDay + "===" + lastDay);


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
	var scatterChart = new Chart(ctx, scatterChartOptions);

	// debug print
	for (var i = 0; i < scatterChart.data.datasets.length; i++) {
		for (var j = 0; j < scatterChart.data.datasets[i].data.length; j++) {
			console.log("data arr ", scatterChart.data.datasets[i].data[j])
		}
	}

}



something

var color = "" + something;