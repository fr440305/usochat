/* app.js *
 * author - firerain *
 */
function ui () {
}

function eventer () {
}

function App () {
	this.stat = [this.Const('stat-list')[0], 0];
}

App.prototype.Const = function (const_name) {
	return {
		"stat-list": ["Input_Nickname", "Check_Nickname", "Isolator", "Be_Invited", "Go_Inviting", "Chatting"]
	}[const_name];
}

App.prototype.looper = function () {

}

