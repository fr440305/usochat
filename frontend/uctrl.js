/* frontend/uctrl.js */
_CANV_PAINT_.width = window.innerWidth;
_CANV_PAINT_.height = window.innerHeight * 0.9;
_CANV_MASK_.ontouchmove = function(e) {
	e.preventDefault();
}
_CANV_PAINT_.onmousedown = function () {
	_CANV_PAINT_.onmousemove = function (e) {
		//FIXME - screenX does not works very well.
		var _x = e.screenX;
		var _y = e.screenY;
		console.log(_x, ", ", _y);
		e.preventDefault();
		_CANV_PAINT_.getContext("2d").fillStyle = "green";
		_CANV_PAINT_.getContext("2d").fillRect(_x, _y, 5, 5);
	};
};
_CANV_PAINT_.onmouseup = function () {
	_CANV_PAINT_.onmousemove = undefined;
};
_CANV_PAINT_.addEventListener("touchmove", function(e) {
	var _x = e.changedTouches[0].pageX - document.body.scrollLeft;
	var _y = e.changedTouches[0].pageY - document.body.scrollTop;
	console.log(_x, ", ", _y, document.body.scrollTop);
	e.preventDefault();
	_CANV_PAINT_.getContext("2d").fillStyle = "green";
	_CANV_PAINT_.getContext("2d").fillRect(_x, _y, 5, 5);
});
_DARK_.onclick = function () {
	document.body.style.backgroundColor = "#000";
	document.body.style.color = "#fff";
};
_LIGHT_.onclick = function () {
	document.body.style.backgroundColor = "#fff";
	document.body.style.color = "#000";
};
_SEND_OPTION_.onclick = function () {
	_SEND_OPTION_BOARD_.style.display = "block";
	_CANV_MASK_.style.display = "block";
	_CANV_PAINT_.style.display = "block";
};
_SEND_CANCEL_.onclick = function () {
	var w = _CANV_PAINT_.width;
	var h = _CANV_PAINT_.height;
	_CANV_PAINT_.getContext("2d").clearRect(0, 0, w, h);
	_CANV_MASK_.style.display = "none";
	_CANV_PAINT_.style.display = "none";
	_SEND_OPTION_BOARD_.style.display = "none";
}
if (window.WebSocket != undefined) {
	//change the abs-url to document.host:
	var ws_conn = new WebSocket("ws://" + document.location.host + "/ws");
	_TEST_.innerHTML += "<p> +++ client-WebSocket </p>";
	ws_conn.onopen = function () {
		var json_msg = {
			"source_node":"",
			"description":"user-login",
			"content":[]
		};
		_TEST_.innerHTML += "<p>+++ Server-websocket</p>";
		//ws_conn.send("_NEW_CLIENT_"); // will be a json. //should get rid of this line.
		ws_conn.send(JSON.stringify(json_msg)); // will be a json.

		_RECEIVED_.innerHTML += ("<p>Connection Opened.</p>");
	};
	ws_conn.onmessage = function(e) {
		//e.data will be jsonlified later on.
		//_RECEIVED_.innerHTML += ("<p>" + e.data + "</p>"); //test mode.
		var msg_content = JSON.parse(e.data);
		var msg_desp = "";
		if (msg_content["description"] !== undefined) {
			msg_desp = msg_content["description"];
		} else {
			//error: no description.
		}
		if (msg_desp === "user-login-0") {
			msg_content = msg_content["content"]; // now, msg_content is a string array.
			for (var i = 0; i < msg_content.length; i++) {
				_RECEIVED_.innerHTML += ("<p>" + msg_content[i] + "</p>");
			}
		} else if (msg_desp === "user-login-*") {
			msg_content = msg_content["content"][0]; // now, msg_content is a string array.
			_ONLINER_.innerHTML = msg_content;
		} else if (msg_desp === "user-msg-text-0") {
			// ... //
		} else if (msg_desp === "user-msg-text-*") {
			msg_content = msg_content["content"]; // now, msg_content is a string array.
			for (var i = 0; i < msg_content.length; i++) {
				_RECEIVED_.innerHTML += ("<p>" + msg_content[i] + "</p>");
			}
		} else if (msg_desp === "user-logout-*") {
			msg_content = msg_content["content"][0]; // now, msg_content is a string array.
			_ONLINER_.innerHTML = msg_content;
		} else {
			//error: invalid msg.
		}


	};
	ws_conn.onerror = function (e) {
		_TEST_.innerHTML += ("<p>" + e.data + "</p>");
	};
	ws_conn.onclose = function () {
		//The sever crash???
		//ws_conn has been closed in this time.
		_TEST_.innerHTML += "<p> ~~~ </p>";
		_RECEIVED_.innerHTML += ("<p>Connection Closed.</p>");
	};
	//To allow the user to send the messages by pressing enter.
	_SEND_CONTENT_.onkeyup = function (e) {
		e = e || windows.event;
		e = e.keyCode || e.which;
		if (e === 13) {
			_SEND_CONFIRM_.onclick();
			_SEND_CONTENT_.click();
		}
	};
	_SEND_PAINT_CONFIRM_.onclick = function () {
		console.log(_CANV_PAINT_.toDataURL());
		_IMG_TEST_.src = _CANV_PAINT_.toDataURL();
		_SEND_CANCEL_.onclick();
	};
	_SEND_CONFIRM_.onclick = function () {
		var json_msg = {
			"source_node":"",
			"description":"user-msg-text",
			"content":[]
		};
		if (_SEND_CONTENT_.value != '') {
			json_msg["content"].push(_SEND_CONTENT_.value);
			ws_conn.send(JSON.stringify(json_msg));
			_SEND_CONTENT_.value = '';
			_SEND_CONTENT_.click();
			_SEND_CONTENT_.focus();
		}
	};
} else {
	_TEST_.innerHTML += (
		"<p>No WebSocket Util Supported. </p>" 
		+"<p> 肥肠抱歉~ 您的浏览器是个战五渣呢。</p>"
		+"<p> 强烈建议您使用谷歌浏览器或者火狐浏览器。</p>"
	);
}
