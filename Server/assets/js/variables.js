/**
  * Thank you for using the XMRig Monitor! If you have any
  * troubles, please open an issue on GitHub and I'll see what
  * I can do: https://github.com/curumimxara/xmrig-monitor
  *
  * This is the setup file. The only required information is
  * the API Endpoint address.
  *
 */

var endpointAPI = "http://127.0.0.1:16000/2/summary"; 	// (required) your API Endpoint address goes here. e.g. http://192.168.1.5:11000/ or http://myproxydomain.org:8080/
var timer = 60; 								// (optional) integrated timer for auto-refreshing the page (value in seconds). set to 0 if you don't want to auto-refresh
var alarmThreshold = 0; 						// (optional) if the 10-seconds average goes below this number, your tab will play a sound and an alert (for proxy it looks at the 60-minutes average)
var aggregateThreads = 0; 						// (optional) if you're using AMD Radeon Vega and you want to see one card per row (each card has two threads) change this to '1'