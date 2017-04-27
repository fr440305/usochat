//Websocket Client
function Client (ws_addr) {
	this.ws_conn = new WebSocket(ws_addr);
	this.usor_id = undefined;
	this.room_id = undefined;
	this.if_dbg = true;
}

Client.prototype.SendMsg = function (usor_id, room_id, summary, content) {

};

Client.prototype.Send = function (raw) {
};
