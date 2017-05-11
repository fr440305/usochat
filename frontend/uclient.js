/*
 * uclient.js
 * __HUO_YU__
 *
 */

function Client (name) {
	this.usor_name = name;
	if (window.WebSocket === undefined) {
		//error: unsupport.
		console.log("@err@ Browser does not support websocket.");
	} else {
		this.ws_conn = new WebSocket(
			"ws://" + document.location.host + "/ws"
		);
		this.load_events();
	}
};

Client.prototype.load_events = function () {
	var client = this;
	this.ws_conn.onopen = function () {
		client.SetName(client.usor_name)
	};
	this.ws_conn.onmessage = function (e) {
		//console.log("Usor-->@res@", e.data);
		msg = JSON.parse(e.data);
		console.log(e.data);
	};
	this.ws_conn.onclose = function () {
		console.log("Usor-->@err@ Websocket Server Closed.")
	};
	this.ws_conn.onerror = function (e) {
		console.log("@err@", e.data);
	};
};

Client.prototype.send_txt = function (txt) {
	if (this.ws_conn !== undefined) {
		this.ws_conn.send(txt);
	} else {
		//error:
		console.log("@err@", "ws_conn === undefined !");
	}
};

Client.prototype.send_msg = function (summary, content) {
	var strjson = JSON.stringify({
		"Summary": summary.toString(),
		"Content": content
	});
	console.log("@std@ Client.send_msg", strjson)
	this.send_txt(strjson);
};

Client.prototype.SetName = function (usor_name) {
	this.send_msg("setname", [[usor_name]]);
};

Client.prototype.Say = function (txt) {
	this.send_msg("say", [[txt.toString()]])
};

Client.prototype.Join = function (room_name) {
	//If room_id == undefined, then it's a new room request.
	this.send_msg("join", [[room_name.toString()]]);
};

Client.prototype.Exitroom = function (if_rm_room) {
	//if_rm_room == {"rm"||"rsv"}
	this.send_msg("exitroom", [[if_rm_room]]);
};
