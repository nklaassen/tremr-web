function loadCanvas() {
	var ctx = document.getElementById("scatterChart").getContext('2d');
	var restScore = []; // empty array
	restScore.push({ x: 0, y: 2});
	restScore.push({ x: 1, y: 2});
	restScore.push({ x: 2, y: 2});
	restScore.push({ x: 3, y: 2});
	restScore.push({ x: 4, y: 2});
	restScore.push({ x: 5, y: 2});
	restScore.push({ x: 6, y: 2});
	var scatterChart = new Chart(ctx, {
		type: 'scatter',
		data: {
			datasets: [{
				label: 'Scatter Dataset',
				for (var i = 0; i < restScore.length; i++) {
					data: [{
						x: 0,
						y: 1
					}, {
						x: 0,
						y: 10
					}, {
						x: 10,
						y: 5
					}]
				}
				}
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


