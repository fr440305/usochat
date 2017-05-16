/* frontend/ui.js */

function AlertPanel () {
}

function LoginPanel () {
}

function EdenPanel (client) {
	this.client = client;
	this.room_list = document.createElement("div");
	var input_roomname = document.createElement("input");
	input_roomname.type = "text";
	input_roomname.onkeydown = function () {
		//Handle the enter-key.
	};
	var confirm_roomname = document.createElement("button");
	confirm_roomname.innerHTML = "search(create)room";
	confirm_roomname.onclick = function () {
		var rname = input_roomname.value;
		input_roomname.value = "";
		client.Join(rname);
	};
	this.pan = document.createElement("div");
	this.pan.appendChild((function(){
		var eden_hint = document.createElement("p");
		eden_hint.innerHTML = "<h2>Pick A Room and Click It</h2>" + "<hr />";
		return eden_hint;
	})());
	this.pan.appendChild(this.room_list);
	this.pan.appendChild((function(){
		var eden_hint_nr = document.createElement("h2");
		eden_hint_nr.innerHTML = "...Or you can create a new Room:";
		return eden_hint_nr;
	})());
	this.pan.appendChild(input_roomname);
	this.pan.appendChild(confirm_roomname);
	document.getElementById("_APP_").appendChild(this.pan);
}

EdenPanel.prototype.Hide = function () {
	//this.room_list.innerHTML = "";
	this.pan.style.display = "none";
};

EdenPanel.prototype.Show = function () {
	this.pan.style.display = "block";
};

EdenPanel.prototype.AddRoomElem = function (room_name) {
	var relem = document.createElement("div");
	var client = this.client;
	relem.className = "_weak_button_";
	relem.innerHTML = room_name + "<hr />";
	relem.onclick = function(){client.Join(room_name)};
	this.room_list.appendChild(relem);
};

EdenPanel.prototype.SetRoomList = function (rlist) {
	this.room_list.innerHTML = "";
	for (var i = 0; i < rlist.length; i++) {
		this.AddRoomElem(rlist[i]);
	}
};

function RoomPanel (client) {
	this.client = client;
	this.pan = document.createElement("div");
	this.pan.style.display = "none";
	var exit_room = document.createElement("div");
	exit_room.className = "_weak_button_";
	exit_room.innerHTML = "&lt;-- Click here to exit this room." + "<hr />";
	exit_room.onclick = function(){
		client.Exitroom("rsv");
	};
	this.dialog_list = document.createElement("div");
	var say_input = document.createElement("input");
	say_input.onkeydown = function () {
		//handle enter-key.
	};
	var say_confirm = document.createElement("button");
	say_confirm.innerHTML = "Send";
	say_confirm.onclick = function() {
		var mydialog = say_input.value;
		say_input.value = "";
		client.Say(mydialog);
	};
	this.pan.appendChild(exit_room);
	this.pan.appendChild(this.dialog_list);
	this.pan.appendChild(say_input);
	this.pan.appendChild(say_confirm);
	document.getElementById("_APP_").appendChild(this.pan);
}

RoomPanel.prototype.Hide = function () {
	//this.dialog_list.innerHTML = "";
	this.pan.style.display = "none";
};

RoomPanel.prototype.Show = function () {
	this.pan.style.display = "block";
};

RoomPanel.prototype.AddDialogElem = function (usor, dialog) {
	var new_dialog_elem = document.createElement("p");
	new_dialog_elem.innerHTML = usor + "&nbsp;--&gt;&nbsp;" + dialog;
	this.dialog_list.appendChild(new_dialog_elem);
};

RoomPanel.prototype.AddSyshintElem = function (syshint) {
	this.AddDialogElem("System", syshint);
};

RoomPanel.prototype.SetChist = function (chist) {
	this.dialog_list.innerHTML = "";
	for (var i = 0; i < chist.length; i++) {
		this.AddDialogElem(chist[i][0], chist[i][1]);
	}
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
	.On('++dialog', this.OnPlusDialog);
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
	this.eden_panel.SetRoomList(rlist);
};

Ui.prototype.OnPlusRoom = function (newroom) {
	this.eden_panel.AddRoomElem(newroom);
};

Ui.prototype.OnMinusRoom = function (rlist) {
	this.eden_panel.SetRoomList(rlist);
};

Ui.prototype.OnRoom = function (roomname, chist) {
	this.eden_panel.Hide();
	this.room_panel.Show();
	this.room_panel.SetChist(chist);
	this.room_panel.AddSyshintElem("Now you are in room: " + roomname);
};

Ui.prototype.OnPlusUsor = function (new_usor, ulist) {
	this.room_panel.AddSyshintElem("New Usor called: " + new_usor + " Add in. Welcome!");
};

Ui.prototype.OnMinusUsor = function (gone_usor, ulist) {
	this.room_panel.AddSyshintElem("New Usor called: " + gone_usor + " gone. Bye!");
};

Ui.prototype.OnModiUsor = function (oldname, newname, ulist) {
};

Ui.prototype.OnPlusDialog = function (usor, dialog) {
	console.log(this);
	this.room_panel.AddDialogElem(usor, dialog);
};

new Ui();
