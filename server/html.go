package server

const mainHTML = `
<html>

<head>
	<script>
function refresh(){
	document.getElementById("preview").src = "/render?w=800&h=600&cachebreak=" + Math.random();
}
window.setInterval(refresh, 1000)

	</script>
</head>

<body>
	<img id="preview" width=800 height=600 src=/render?w=800&h=600></img>
</body>

</html>
`
