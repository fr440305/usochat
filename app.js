/* app.js *
 * author - firerain *
 */

/* import api.js in var::$api */
var $api = new API();

function app () {
	this.UI = new ui();
	this.req_interval = 250;
	var Me = this;
	setInterval(function(){			/* TODO - implement this code by using $api.Theard() */
		var dia_arr = $api.Client.Var("dialogs");
		if (dia_arr.length !== Me.UI.Var("dia-node-length") ) {
			Me.UI.ClearDiaNodes();
			for (var i = 0; i < dia_arr.length; i++) {
				Me.UI.AddDiaNode(dia_arr[i]);
			}
		}
	}, Me.req_interval, true);
}

function ui () {
	this.main = document.getElementById("ui");
	this.dia_board = document.getElementById("-uso-");
}

ui.prototype.AddDiaNode = function (dia_content) {
	var new_dia_node = document.createElement("p");
	new_dia_node.innerHTML = dia_content;
	this.dia_board.appendChild(new_dia_node);
}

ui.prototype.ClearDiaNodes = function () {
	this.dia_board.innerHTML = ''; /* TODO - not the best option. Need to be refactored */
}

ui.prototype.Var = function (key) {
	return {
		"dia-node-length": this.dia_board.childElementCount
	}[key];
}

var $app = new app();
