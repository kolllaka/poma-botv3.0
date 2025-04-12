// ws
import { ConnectWS } from '../../lib/ws.js';
// notification
import { Notification } from '../notification.js'

// init websocket
const url = 'ws://127.0.0.1:8080/subgift/ws'
const ws = new ConnectWS(url)
ws.onMessage = function (event) {
	const data = JSON.parse(event.data)
	console.log("data from ws", data, "isEmpty queue", queue.isEmpty);

	queue.push(data)
}
ws.setup()


class Queue {
	constructor() {
		this.queue = []
		this.isEmpty = true
	}

	push(data) {
		this.queue.push(data)
		this.isEmpty = false

		console.log("pushed:", this.queue);

		if (!raid.isPlay) {
			const data = this.pop()

			raid.startNotification(data)
		}
	}

	pop() {
		const data = this.queue.pop()

		console.log("poped:", this.queue);

		if (!data) {
			this.isEmpty = true
		}

		return data
	}
}

const queue = new Queue

const raid = new Notification("subgift", queue);
raid.video.addEventListener('ended', (e) => {
	raid.stopNotification()

	setTimeout(() => {
		raid.isPlay = false

		if (!queue.isEmpty) {
			const data = queue.pop()

			raid.startNotification(data)
		}
	}, 1000)
})


