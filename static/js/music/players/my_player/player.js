export class Player {
	constructor(selector, state = {
		volume: 0,
	}) {
		this.src = ""
		this.player = this.#init(selector)
		this.setVolume(state.volume)

		this.#setup()
	}

	#init(selector) {
		document.getElementById(selector).innerHTML = this.#getAudioPlayerTemplate()

		return document.getElementById(selector).querySelector('audio')
	}

	#setup() {
		this.handler = this.handler.bind(this)
		this.player.addEventListener('ended', this.handler)
	}

	play() {
		this.player.play()
	}

	pause() {
		this.player.pause()
	}

	setVolume(volume) {
		this.player.volume = volume / 100
	}

	setLink(link) {
		this.player.src = link
		this.src = link
	}

	handler(e) {
		this.onEndedHandler(e)
	}

	onEndedHandler(e) {
		console.log('not onEndedHandler implement');
	}

	#getAudioPlayerTemplate() {
		return `<audio src="" controls type="audio/mpeg"></audio>`
	}
}