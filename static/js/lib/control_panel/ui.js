const SONG_AUTHOR_SELECTOR = "[data-info='author']"
const SONG_NAME_SELECTOR = "[data-info='name']"
const VOLUME_SELECTOR = "[data-info='volume']"
const PLAY_BTN_SELECTOR = "[data-info='play']"

export const BTN_TEXT_PLAY = "play"
export const BTN_TEXT_PAUSE = "pause"

export class SongUI {
	constructor(selector) {
		this.$el = document.getElementById(selector)
	}

	updateInfo(info) {
		this.$el.querySelector(SONG_AUTHOR_SELECTOR).innerHTML = info.author
		this.$el.querySelector(SONG_NAME_SELECTOR).innerHTML = info.name
	}

	setVolume(volume) {
		this.$el.querySelector(VOLUME_SELECTOR).value = volume
	}

	getBtnPlay() {
		return this.$el.querySelector(PLAY_BTN_SELECTOR)
	}

}