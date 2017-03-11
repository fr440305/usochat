/* app.js *
 * author - firerain *
 */

function Eventer () {

}

function Conner () {
	this.ws = new WebSocket("ws://" + document.location.host + "/ws");
	this.wsInit();
}

Conner.prototype.wsInit = function () {
}


