<html>
<head>
	<meta charset="utf-8">
<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.2.1/jquery.min.js"></script>
</head>

<body>
	<form action="calculate.php" method="GET">
		calls<input type="" name="calls" id="calls" value="1000">Кол-во за интервал<br>
		interval, sec<input type="" name="interval" id="interval" value="900">Интервал, c (900=15 мин)<br>
		maxq<input type="" name="maxq" id="maxq" value="10">Максимальная длина очереди<br>
		call time, sec<input type="" name="ctime" id="ctime" value="5">Время обработки<br>
		<br>
		answer time, sec<input type="" name="atime" id="atime" value="1">Время ожидания<br>
		procent <input type="" name="uo" id="uo" value="0.98">SLA 1=100%<br>
		<!-- <input type="submit" name="go"> -->
	</form>
	<button id="abutton">Calculate!</button>
	<h3><pre><div id='result'>result</div></pre></h3>
<script type="text/javascript">
	$('#abutton').click(function(){
		calls=$('#calls').val()
		interval=$('#interval').val()
		maxq=$('#maxq').val()
		ctime=$('#ctime').val()
		atime=$('#atime').val()
		uo=$('#uo').val()
		// alert(calls+' '+interval);
		var request
        $.ajax({
            url: "calculate.php",
            type: "POST",
            data: {calls:calls, interval:interval, ctime:ctime, atime:atime, uo:uo, maxq:maxq}
        })
        .done(function( data ) {
            if ( console && console.log ) {
                $("#result").html(data);
            }
        });
	});

</script>	
</body>

</html>
