package bruteray
const mainHTML=`
<html>

<head>
	<title>bruteray</title>

	<script>

function refresh(){
	document.getElementById("render").src = "/render?cachebreak=" + Math.random();
}
window.setInterval(refresh, 1000)

	</script>
</head>

<body>
	<img id="render" width=600 height=400></img>
</body>

</html>
`
