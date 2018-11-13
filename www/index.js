function loadCanvas() {
	var ctx = document.getElementById("scatterChart").getContext('2d');
	var restScore = []; // empty array
	restScore.push({ x: 0, y: 2});
	restScore.push({ x: 1, y: 2});
	restScore.push({ x: 2, y: 2});
	var scatterChart = new Chart(ctx, {
		type: 'scatter',
		data: {
			datasets: [{
				label: 'Scatter Dataset',
				data: restScore
			}]
    	},
		options: {
			scales: {
				xAxes: [{
					type: 'linear',
					position: 'bottom'
				}]
			}
		}
	}); 
}
