/* app.js *
 * author - firerain *
 */

function Eventer () {
	this.conner = new Conner();
	this.bind();
}

Eventer.prototype.bind = function () {
	var Me = this;
	document.getElementById("sender").onclick = function () {
		Me.conner.Send(document.getElementById("say").value);
	}
}

function Conner () {
	this.ws = new WebSocket("ws://" + document.location.host + "/ws");
	this.wsInit();
}

Conner.prototype.wsInit = function () {
	this.ws.onopen = function () {
		document.title = "uso - open";
	}
	this.ws.onclose = function () {
		document.title = "uso - close";
	}
	this.ws.onmessage = function (evt) {
		document.getElementById("-uso-").innerHTML += ('<p>' + evt.data + '</p>');
	}
}

Conner.prototype.Send = function (msg) {
	this.ws.send(msg);
}

new Eventer();
