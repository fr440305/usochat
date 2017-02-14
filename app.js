/* app.js *
 * author - firerain *
 */

/* import api.js in var::$api */
var $api = new API();

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
	this.dia_board.innerHTML = ''; /* not the best option. Need to be refactored */
}

ui.prototype.Var = function (key) {
	return {
		"dia-node-length": this.dia_borad.children.length
	}[key];
}

var $ui = new ui();

setInterval(function(){
	$ui.ClearDiaNodes();
	var dia_arr = $api.Client.Var("dialogs");
	for (var i = 0; i < dia_arr.length; i++) {
		$ui.AddDiaNode(dia_arr[i]);
	}
	//document.getElementById('-uso-').innerHTML = $api.Client.Var("dialogs");
	//console.log($api.Client.dialogs);
}, 1000, true);

