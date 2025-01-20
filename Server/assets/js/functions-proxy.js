$.ajax({
    url: endpointAPI,
    method: 'GET',
    headers: {
        'Authorization': 'Bearer null' // Replace with your actual token
    },
    success: function(data) {

        // connection
        document.getElementById("active-miners").innerHTML = "Miners: " + data.miners.now + ' miners (max: ' + data.miners.max + ')';
        document.getElementById("uptime").innerHTML = "Uptime: " + (data.uptime / 3600).toFixed(2) + ' hours / ' + (data.uptime / 86400).toFixed(2) + ' days';

        document.getElementById("tot60").innerHTML = "Hashrate (60m) " + Number(data.hashrate.total[2] ).toLocaleString('en-GB') + ' Kh/s';

    },
    error: function(xhr, status, error) {
        console.error("Error fetching data:", error);
    }
});

// create meta tag for auto-refresh
if (timer > 0) {
	var meta = document.createElement('meta');
	meta.httpEquiv = "refresh";
	meta.content = timer;
	document.getElementsByTagName('head')[0].appendChild(meta);
}