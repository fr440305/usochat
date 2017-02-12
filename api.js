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
	this.stat = undefined;
	this.ajax = new XMLHttpRequest();
	this.get_url = "get";
	this.post_url = "post";
}

client.prototype.Var = function(property) {
	/* accessor pattern */
	return {
		"http-status": this.ajax.status,
		"http-readystate": this.ajax.readyState,
		"dialog-json": this.ajax.responseText
	}[property];
}

client.prototype.Const = function (key) {
	/* only offer the constants which will be used by other object(s) */
	return ({
		"stat-initialize": 0,
		"stat-after-init": 1
	}[name]);
}

client.prototype.Say = function (dialog) {
	this.post({"dialog" : dialog});
}

client.prototype.FetchConversation = function () {
	this.get("conversation");
}

client.prototype.post = function (post_object) {
	var post_content;
	this.ajax.open("POST", this.post_url, true);
	this.ajax.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
	this.ajax.send(post_content);
}

client.fetchId = function () {
	this.get("initid");
}

client.prototype.get = function (get_req_content) {
	/* get_req_content : string = \" get_req_key \= get_req_value [{ \& get_req_key \= get_req_value }] \" */
	this.ajax.open("GET", this.get_url + '?' + get_req_content, true);
	this.ajax.send();
}

client.prototype.statWriter = function () {
}
