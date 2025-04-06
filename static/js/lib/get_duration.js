export class GetDuration {
	constructor() {
		this.audio = new Audio()
		this.checkArray = []
		this.currentData = {}

		this.#setup()
	}

	#setup() {
		this.handler = this.handler.bind(this)
		this.audio.addEventListener('loadedmetadata', this.handler)
	}

	handler() {
		this.onloadedmetadata()
	}

	onloadedmetadata = function (e) {
		console.log("not onloadedmetadata implement");
	}

	setSrc(data) {
		this.audio.src = encodeURI(data.data.link)
	}

	start() {
		this.currentData = this.checkArray.pop()

		this.setSrc(this.currentData)
	}

	addSong(song) {
		this.checkArray.push(song)
	}

}