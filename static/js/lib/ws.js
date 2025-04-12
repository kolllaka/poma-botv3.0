export class ConnectWS {
	constructor(url) {
		this.url = url
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

	onOpen(event) {
		console.log(`ws conn open to  ${event.target.url}`);
	}

	onMessage(event) {
		const data = JSON.parse(event.data)
		console.log("data from ws", data);
	}

	onClose(event) {
		if (event.wasClean) {
			console.log(`[close] Соединение закрыто чисто, код=${event.code} причина=${event.reason}`);
			this.socket.close();
		} else {
			// например, сервер убил процесс или сеть недоступна
			// обычно в этом случае event.code 1006
			console.log('[close ws] Соединение прервано');
		}

		setTimeout(function () {
			this.socket = new WebSocket(this.url)
		}, 1000);
	}

	onError(error) {
		console.log(`[error] ${error}`);
		this.socket.close();
	}
}