/* app.js *
 * author - firerain *
 */

/* import api.js in var::$api */
var $api = new API();

setInterval(function(){
	document.getElementById('-uso-').innerHTML = $api.Client.Var("dialogs");
	//console.log($api.Client.dialogs);
}, 1000, true);

