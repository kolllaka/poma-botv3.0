const INFO_SELECTOR = ".playlist__info"
const BODY_PLAYLIST_SELECTOR = ".playlist__body"

export class PlaylistUI {
	constructor(selector, templates = {}) {
		this.$el = document.getElementById(selector)
		this.templates = templates
	}

	drawPlaylist(datas = [], info = {}) {
		this.$el.querySelector(INFO_SELECTOR).innerHTML = this.templates.info(info)
		this.$el.querySelector(BODY_PLAYLIST_SELECTOR).innerHTML = datas.map(this.templates.song).join("")
	}

	addSong(song) {
		this.$el.querySelector(BODY_PLAYLIST_SELECTOR)
			.insertAdjacentHTML('beforeEnd', this.templates.song(song, this.$el.querySelector(".playlist__body").childElementCount))
	}

	deleteSong(index) {
		let isRemove = false
		Array.from(this.$el.querySelector(BODY_PLAYLIST_SELECTOR).children)
			.forEach((item, ind) => {
				if (!isRemove && ind === index) {
					item.remove()

					isRemove = true

					return
				}

				if (isRemove) {
					item.querySelector('.itemplaylist__cell').innerHTML = `${ind}`
					item.querySelector('[data-index]').dataset.index = ind - 1
				}

			})
	}

	updatePlaylist(datas) {
		this.$el.querySelector(BODY_PLAYLIST_SELECTOR).innerHTML = datas.map(this.templates.song).join("")
	}

	updateInfo(info) {
		this.$el.querySelector(INFO_SELECTOR).innerHTML = this.templates.info(info)
	}
}