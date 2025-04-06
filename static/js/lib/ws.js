let socket
function connectWS(WEBSOCKET, handler, onopenhandler) {

	socket = new WebSocket(WEBSOCKET);

	socket.onopen = function (e) {
		console.log(`ws conn open to  ${WEBSOCKET}`);

		if (!onopenhandler) {
			console.log("handler not allowed");

			return
		}

		onopenhandler()
	};

	socket.onmessage = function (e) {
		console.log(`[event] ${e}`);
	};

	socket.onclose = function (event) {
		if (event.wasClean) {
			console.log(`[close] Соединение закрыто чисто, код=${event.code} причина=${event.reason}`);
		} else {
			// например, сервер убил процесс или сеть недоступна
			// обычно в этом случае event.code 1006
			console.log('[close ws] Соединение прервано');
		}

		setTimeout(function () {
			connectWS();
		}, 1000);
	};

	socket.onerror = function (error) {
		console.log(`[error] ${error}`);
		socket.close();
	};
}


export class ConnectWS {
	constructor(url) {
		this.socket = this.connect(url)
	}

	setup() {
		this.socket.onopen = this.onOpen
		this.socket.onmessage = this.onMessage
		this.socket.onclose = this.onClose
		this.socket.onerror = this.onError
	}

	connect(url) {
		return new WebSocket(url);
	}

	onOpen(e) {
		console.log('not onOpen implement');
	}

	onMessage(e) {
		console.log('not onMessage implement');
	}

	onClose(e) {
		console.log('not onClose implement');
	}

	onError(e) {
		console.log('not onError implement');
	}
}