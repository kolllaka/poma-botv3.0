import { shuffle } from '/static/js/lib/misc.js';

export class Playlist {
	constructor(datas = [], mapping) {
		this.mapping = mapping

		this.datas = datas.map(this.mapping)
		this.info = this.#calcInfo()
	}

	mapping(data = {}) {
		return {
			...data,
			duration: data.duration
		}
	}

	#calcInfo() {
		return {
			count: this.datas.length,
			duration: this.datas.reduce(
				(acc, data) => acc + data.duration,
				0
			)
		}
	}

	shuffle() {
		shuffle(this.datas)

		return this.datas
	}

	getSongs() {
		return this.datas
	}

	getSong(index) {
		return this.datas[index]
	}

	getInfo() {
		return this.info
	}

	addSong(dataSong) {
		const song = this.mapping(dataSong)

		this.datas.push(song)
		this.info = {
			count: this.info.count + 1,
			duration: this.info.duration + song.duration
		}

		return song
	}

	deleteSong(index) {
		const song = this.datas.splice(index, 1)[0]

		if (song) {
			this.info = {
				count: this.info.count - 1,
				duration: this.info.duration - song.duration
			}

			return song
		}

		return
	}

}