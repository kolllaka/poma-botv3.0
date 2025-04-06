

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
		const datas = this.controller.data.getSongs()
		const info = this.controller.data.getInfo()

		this.controller.ui.drawPlaylist(datas, info)

		this.clickHandler = this.clickHandler.bind(this)
		this.$el.addEventListener('click', this.clickHandler)
	}

	clickHandler($event) {
		switch ($event.target.dataset.btn) {
			case "del":
				this.onDeleteHandler($event.target)

				break
		}
	}

	shuffle() {
		const datas = this.controller.data.shuffle()
		this.controller.ui.updatePlaylist(datas)
	}

	addSong(data) {
		const song = this.controller.data.addSong(data)
		const info = this.controller.data.getInfo()

		this.controller.ui.addSong(song)
		this.controller.ui.updateInfo(info)
	}

	deleteSong(index) {
		const song = this.controller.data.deleteSong(index)
		const info = this.controller.data.getInfo()

		this.controller.ui.deleteSong(index)
		this.controller.ui.updateInfo(info)

		return song
	}

	onDeleteHandler(target) {
		const index = parseInt(target.dataset.index)

		this.deleteSong(index)
	}
}

