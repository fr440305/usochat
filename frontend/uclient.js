
function Client () {
	if (window.WebSocket === undefined) {
		//error: unsupport.
		console.log("@err@ Browser does not support websocket.");
	} else {
		this.ws_conn = new WebSocket(
			"ws://" + document.location.host + "/ws"
		);
		this.id = { 'r': undefined, 'u': undefined };
		this.load_events();
	}
};

Client.prototype.load_events = function () {
	this.ws_conn.onopen = function () {
	};
	this.ws_conn.onmessage = function (e) {
		console.log("@res@", e.data);

	};
	this.ws_conn.onclose = function () {
	};
	this.ws_conn.onerror = function (e) {
		console.log("@err@", e.data);
	};
};

Client.prototype.send_txt = function (txt) {
	if (this.ws_conn !== undefined) {
		this.ws_conn.send(txt);
	}
};

Client.prototype.send_login_msg = function () {
};


