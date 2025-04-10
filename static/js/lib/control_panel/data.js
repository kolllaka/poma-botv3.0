export const DEFAULT_SONG_INFO = {
	name: "заказов нет",
	author: "-",
}


function mapping(data = {}) {
	return {
		song: {
			name: data.name || DEFAULT_SONG_INFO.name,
			author: data.author || DEFAULT_SONG_INFO.author,
		},
		volume: data.volume,
	}
}

export class Song {
	constructor(state = {
		song: {},
		volume: 0
	}) {
		this.state = mapping(state)
	}

	getSong() {
		return this.state.song
	}

	updateSong(song) {
		this.state.song = song
	}

	getVolume() {
		return this.state.volume
	}

	setVolume(volume) {
		this.state.volume = volume
	}

}