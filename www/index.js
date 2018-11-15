function getTremors(callback) {
	fetch("/api/tremors?since=2018-10-14T17:13:39Z")
		.then(response => response.json())
		.then(
			result => {
				callback(result)
			},
			error => {
				console.log("error")
			}
		)
}

function loadCanvas() {

	var ctx = document.getElementById("myChart").getContext('2d');

	var scatterChartOptions = {
		type: 'scatter',
		data: {
			datasets: [{
				label: 'resting score',
				showLine: false, // disable for a single dataset
				pointBorderColor: 'blue',
				backgroundColor: '#ADD8E6',
				pointBackgroundColor: 'black',
			}, {
				label: 'postural score',
				pointBackgroundColor: 'red',
				pointBorderColor: 'red'
			}]
		},
		options: {
			scales: {
				xAxes: [{
					display: true,
					labelString: 'screefgbnLeft',
					type: 'time',
					position: 'bottom'
				}]
			}
		}
	};

	getTremors(result => {
		resting = result.map(tremor => {
			return {
				x: tremor.date,
				y: tremor.resting / 10
			}
		})
		postural = result.map(tremor => {
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
	})
}