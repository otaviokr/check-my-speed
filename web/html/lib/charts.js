function load() {
	google.charts.load('current', {packages: ['corechart', 'line', 'gauge']});
	google.charts.setOnLoadCallback(drawLine);
	google.charts.setOnLoadCallback(drawChart);
}
		
function drawChart() {
	var data = google.visualization.arrayToDataTable([
		['Label', 'Value'],
		['Min', minValue],
		['Max', maxValue],
		['Average', AvgValue],
		['Last', LastValue]
	]);

	var options = {
		width: 800, height: 240,
		animation: {
			duration: 5000,
			easing: 'in'},
		min: 0, max: nominalBand * 1.2,
		redFrom: nominalBand, redTo: nominalBand * 1.2,
		greenFrom:nominalBand * 0.8, greenTo: nominalBand,
		minorTicks: 5
	};

	var chart = new google.visualization.Gauge(document.getElementById('chart_div'));
	chart.draw(data, options);
}

function drawLine() {
	/*var data = new google.visualization.DataTable();
	data.addColumn('number', 'X');
	data.addColumn('number', 'Download');
	data.addColumn('number', 'Upload');

	data.addRows(points);*/
	var data = google.visualization.arrayToDataTable(points);

	var options = {
		title: 'Bandwidth (Mbit/s)',
		hAxis: {
			title: 'Time',
		},
		vAxis: {
			format:'#.##'
		},
		backgroundColor: '#ffffff'
	};

	var chart = new google.visualization.LineChart(document.getElementById('line_div'));
	chart.draw(data, options);
}

