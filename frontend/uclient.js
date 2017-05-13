/*
 * uclient.js
 * __HUO_YU__
 *
 */

function Client (this_of_callbacks) {
	this.this_of_callbacks = this_of_callbacks;
	this.usor_name = name;
	this.evtlist = {
		'-)eden':function(rlist){
			console.log("IN EDEN", rlist);
		},
		'++room':function(new_room_name){
			console.log("ADD ROOM", new_room_name);
		},
		'--room':function(rlist){
			console.log("RM ROOM", rlist);
		},
		'-)room':function(room_name, chist){
			console.log("IN ROOM", room_name, chist);
		},
		'++usor':function(usor, ulist){
			console.log("ADD USOR", usor, ulist);
		},
		'--usor':function(usor, ulist){
			console.log("GONE USOR", usor, ulist);
		},
		'++dialog':function(usor, dialog){
			console.log("ADD DIALOG", usor, dialog);
		},
		'~~usor':function(ori_name, new_name, ulist){
			console.log("USOR NAME CHANGE", ori_name, new_name, ulist);
		},
		'~~name':function(new_name){
			console.log("SELF NAME CHANGE", new_name);
		},
		'error':function(hint){
			console.log("ERROR", hint);
		}
	}
};

Client.prototype.Connect = function () {
	if (window.WebSocket === undefined) {
		//error: unsupport.
		console.log("@err@ Browser does not support websocket.");
		return false;
	} else {
		this.ws_conn = new WebSocket(
			"ws://" + document.location.host + "/ws"
		);
		this.load_events();
		return true;
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
		switch (msg.Summary) {
		case '~~name':
			var name = msg.Content[0][0];
			client.evtlist['~~name'].call(client.this_of_callbacks, name);
			break;
		case 'error':
			var hint = msg.Content[0][0];
			client.evtlist['error'].call(client.this_of_callbacks, hint);
			break;
		case '-)eden':
			var rlist = msg.Content[0];
			client.evtlist['-)eden'].call(client.this_of_callbacks, rlist);
			break;
		case '++room':
			var new_room = msg.Content[0][0];
			client.evtlist['++room'].call(client.this_of_callbacks, new_room);
			break;
		case '--room':
			var rlist = msg.Content[0];
			client.evtlist['--room'].call(client.this_of_callbacks, rlist);
			break;
		case '-)room':
			var rname = msg.Content[0][0];
			msg.Content.shift();
			client.evtlist['-)room'].call(client.this_of_callbacks, rname, msg.Content);
			break;
		case '++usor':
			var new_usor = msg.Content[0][0];
			var ulist = msg.Content[1];
			client.evtlist['++usor'].call(client.this_of_callbacks, new_usor, ulist);
			break;
		case '--usor':
			var gone_usor = msg.Content[0][0];
			var ulist = msg.Content[1];
			client.evtlist['--usor'].call(client.this_of_callbacks, gone_usor, ulist);
			break;
		case '~~usor':
			var ori_usor = msg.Content[0][0];
			var now_usor = msg.Content[0][1];
			var ulist = msg.Content[1];
			client.evtlist['~~usor'].call(client.this_of_callbacks, ori_usor, now_usor, ulist);
			break;
		case '++dialog':
			client.evtlist['++dialog'].call(client.this_of_callbacks, msg.Content[0][0], msg.Content[0][1]);
			break;
		default:
			console.log("???");
		}
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

Client.prototype.On = function (event_name, callback) {
	//TODO provide that callback is a function
	this.evtlist[event_name] = callback;
	return this;
};
