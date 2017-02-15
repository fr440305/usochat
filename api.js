/*api.js *
 * author - firerain *
 */

function API() {
	this.Theard = new theard();
	this.Client = new client();
}

function theard() {
}

function client() {
	this.id = undefined;
	this.dialogs = undefined;
	this.stat = undefined;
	this.ajax = new XMLHttpRequest();
	this.post_queue = undefined;		/* a string array for storing the dialogs that needed to be sent */
	this.get_url = "get";
	this.post_url = "post";
	var Me = this;
	this.ajax.onreadystatechange = function() {
		//console.log('qweqwe');
		if (Me.ajax.readyState === 4) {
			Me.dialogs = Me.ajax.responseText;
			//console.log(Me.dialogs);
		}
	}
	setInterval(function(){Me.fetchConversation()}, 500, true);
}

client.prototype.Var = function(v_name) {
	/* accessor pattern */
	var dia = this.dialogs;
	if (v_name === "dialogs") { 
		return (dia === "null" || dia === undefined || dia === "undefined") ? [] : JSON.parse(this.dialogs);
	}
	/*
	return {
		"http-status": this.ajax.status,
		"http-readystate": this.ajax.readyState,
		"dialog-json": this.ajax.responseText
	}[property];
	*/
}

client.prototype.Say = function (dialog) {
	/* TODO - push this dialog to this.push_queue */
	this.post({"dialog" : dialog});		
}

client.prototype.fetchConversation = function () {
	this.get({"conversation":"-"});
}

client.prototype.post = function (post_object) {
	var post_string_array = [];
	var keys = Object.keys(post_object);
	for (var key_index = 0; key_index < keys.length; key_index++) {
		var k = keys[key_index];
		var v = post_object[k].toString();
		post_string_array.push(k.toString() + '=' + v + '&'); /* if v is empty then the equal sign should not be written. */
	}
	var post_string = post_string_array.join(''); /* .slice(0, -1) */
	post_string = post_string.substring(0, post_string.length - 1);
	this.ajax.open("POST", this.post_url, true);
	this.ajax.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
	this.ajax.send(post_string);
}

client.fetchId = function () {
	this.get({"initid":"-"});
}

client.prototype.get = function (get_object) {
	var get_string_array = [];
	var keys = Object.keys(get_object);
	for (var key_index = 0; key_index < keys.length; key_index++) {
		var k = keys[key_index];
		var v = get_object[k].toString();
		get_string_array.push(k.toString() + '=' + v + '&');
	}
	var get_string = get_string_array.join(''); /* .slice(0, -1) */
	get_string = get_string.substring(0, get_string.length - 1);
	/* get_req_content : string = \" get_req_key \= get_req_value [{ \& get_req_key \= get_req_value }] \" */
	this.ajax.open("GET", this.get_url + '?' + get_string, true);
	this.ajax.send();
}

client.prototype.looper = function () {
	/* TODO - analyse and update the status, including post the dialogs in this.post_queue to server , *
	 * fetch the conversations from server, *
	 * and update properties of this object */
}
