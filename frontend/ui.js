/* frontend/ui.js */

function Ui () {
	var ui = this;
	this.client = new Client(this);
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
	console.log(this);
	var client = this.client;
	_ROOM_.style.display = 'none';
	_EDEN_.style.diaplay = 'block';
	for (var i = 0; i < rlist.length; i++) {
		_EDEN_.innerHTML += ("<p id='_ROOM_" + i + "_'>" + rlist[i] + "</p>");
		var elem = document.getElementById("_ROOM_" + i + "_");
		elem.onclick = function () {
			client.Join(elem.innerHTML);
		};
	}
};

Ui.prototype.OnPlusRoom = function (newroom) {
};

Ui.prototype.OnMinusRoom = function (rlist) {
};

Ui.prototype.OnRoom = function (roomname, chist) {
	_EDEN_.innerHTML = "";
	_ROOM_.style.display = 'block';
	_EDEN_.style.diaplay = 'none';

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
