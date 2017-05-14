/* frontend/ui.js */

var Alert  = function (hint) {
	
}

function EdenPanel (client) {
	var eden_hint = document.createElement("span");
	eden_hint.innerHTML = "Pick A Room and Click It";
	this.client = client;
	this.pan = document.createElement("div");
	this.room_list = document.createElement("div");
	this.pan.appendChild(eden_hint);
	this.pan.appendChild(this.room_list);
	document.getElementById("_APP_").appendChild(this.pan);
}

EdenPanel.prototype.Hide = function () {
	this.room_list.innerHTML = "";
	this.pan.style.display = "none";
};

EdenPanel.prototype.Show = function () {
	this.pan.style.display = "block";
};

EdenPanel.prototype.AddRoomElem = function (room_name) {
	var relem = document.createElement("div");
	var client = this.client;
	relem.className = "_room_elem_";
	relem.innerHTML = room_name;
	relem.onclick = function(){client.Join(room_name)};
	this.room_list.appendChild(relem);
};

function RoomPanel (client) {
	this.client = client;
	this.pan = document.createElement("div");
	var exit_room = document.createElement("div");
	exit_room.innerHTML = "&lt;--";
	exit_room.onclick = function(){
		client.Exitroom("rsv");
	};
	this.pan.appendChild(exit_room);
	document.getElementById("_APP_").appendChild(this.pan);
}

RoomPanel.prototype.Hide = function () {
	this.pan.style.display = "none";
};

RoomPanel.prototype.Show = function () {
	this.pan.style.display = "block";
};

function Ui () {
	var ui = this;
	this.client = new Client(this);
	this.eden_panel = new EdenPanel(this.client);
	this.room_panel = new RoomPanel(this.client);
	this.load_events();
	this.client.Connect();
}

Ui.prototype.load_events = function () {
	this.client
	.On('error', this.OnError)
	.On('~~name', this.OnModiName)
	.On('-)eden', this.OnEden)
	.On('++room', this.OnPlusRoom)
	.On('--room', this.OnMinusRoom)
	.On('-)room', this.OnRoom)
	.On('++usor', this.OnPlusUsor)
	.On('--usor', this.OnMinusUsor)
	.On('~~usor', this.OnModiUsor)
	.On('++dialog', this.OnPlusUialog);
};

Ui.prototype.ParseUrl = function () {
};

Ui.prototype.OnError = function (hint) {
};

Ui.prototype.OnModiName = function (myname) {
	console.log("!!!");
};

Ui.prototype.OnEden = function (rlist) {
	this.eden_panel.Show();
	this.room_panel.Hide();
	for (var i = 0; i < rlist.length; i++) {
		this.eden_panel.AddRoomElem(rlist[i]);
	}
};

Ui.prototype.OnPlusRoom = function (newroom) {
};

Ui.prototype.OnMinusRoom = function (rlist) {
};

Ui.prototype.OnRoom = function (roomname, chist) {
	this.eden_panel.Hide();
	this.room_panel.Show();
};

Ui.prototype.OnPlusUsor = function (new_usor, ulist) {
};

Ui.prototype.OnMinusUsor = function (gone_usor, ulist) {
};

Ui.prototype.OnModiUsor = function (oldname, newname, ulist) {
};

Ui.prototype.OnPlusDialog = function (usor, dialog) {
};

var UI = new Ui();
