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
	this.ajax = new XMLHttpRequest();
}

client.prototype.AddListenerOnReadyStateChange = function(callback) {
}

client.prototype.Post = function (url, post_content) {
}

client.prototype.Get = function (url) {
}


