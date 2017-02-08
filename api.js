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

client.prototype.AddEventListener = function(event_name, callback) {
	if (event_name === "onreadystatechange" ) {
		var Me = this;
		this.ajax.addEventListener("onreadystatechanged", function(){callback(/*we can add some extra parameters here */)}, true);
	}
}

client.prototype.Post = function (url, post_content) {
}

client.prototype.Get = function (url) {
}


