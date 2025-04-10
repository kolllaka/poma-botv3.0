export class Controller {
	constructor(selector, ui, data) {
		this.$el = document.getElementById(selector)

		this.controller = {
			data: data,
			ui: ui
		}

		this.#init()
	}

	#init() {
		const song = this.controller.data.getSong()
		this.controller.ui.updateInfo(song)

		const volume = this.controller.data.getVolume()
		this.controller.ui.setVolume(volume)

		this.clickHandler = this.clickHandler.bind(this)
		this.$el.addEventListener('click', this.clickHandler)

		this.inputHandler = this.inputHandler.bind(this)
		this.$el.addEventListener('input', this.inputHandler)
	}

	getSong() {
		return this.controller.data.getSong()
	}

	updateInfo(info) {
		this.controller.data.updateSong(info)

		this.controller.ui.updateInfo(info)
	}

	setVolume(volume) {
		this.controller.data.setVolume(volume)

		this.controller.ui.setVolume(volume)
	}

	getBtnPlay() {
		return this.controller.ui.getBtnPlay()
	}

	clickHandler($event) {
		switch ($event.target.dataset.btn) {
			case "skip":
				this.onSkipHandler($event.target)

				break

			case "shuffle":
				this.onShuffleHandler($event.target)

				break

			case "play":
			case "pause":
				this.onPlayHandler($event.target)

				break

		}
	}

	inputHandler($event) {
		this.onVolumeHandler($event.target.value)
	}

	onVolumeHandler(value) {
		console.log(value);
	}

	onSkipHandler(target) {
		console.log("not onSkipHandler implement");
	}

	onShuffleHandler(target) {
		console.log("not onShuffleHandler implement");
	}

	onPlayHandler(target) {
		console.log("not onPlayHandler implement");
	}
}