/*api.js *
 * author - firerain *
 */

function Api() {
	this.Theard = new theard();
	this.Client = new client();
}

function theard() {
}

function client() {
	this.id = undefined;
	this.stat = undefined;
	this.ajax = new XMLHttpRequest();
	this.get_url = "get";
	this.post_url = "post";
}

client.prototype.AddEventListener = function(event_name, callback) {
	if (event_name === "onreadystatechange" ) {
		var Me = this;
		this.ajax.addEventListener("onreadystatechange", function(){callback(/*we can add some extra parameters here */)}, true);
	}
}

client.prototype.Say = function (dialog) {
	this.post("dialog=" + dialog);
}

client.prototype.SendGettingQuery = function () {
	this.get("conversation");
}

client.prototype.GetDialogJson = function () {
	return this.ajax.responseText;
}

client.prototype.post = function (post_content) {
	/* post_content : string = \" post_key \= post_value [{ \& post_key \= post_value }] \" */
	this.ajax.open("POST", this.post_url, true);
	this.ajax.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
	this.ajax.send(post_content);
}

client.prototype.get = function (get_req_content) {
	/* get_req_content : string = \" get_req_key \= get_req_value [{ \& get_req_key \= get_req_value }] \" */
	this.ajax.open("GET", this.get_url + '?' + get_req_content, true);
	this.ajax.send();
}


