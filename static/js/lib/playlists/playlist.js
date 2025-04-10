import { shuffle } from '/static/js/lib/misc.js';

export class Playlist {
	constructor(datas = []) {
		this.datas = datas
		this.info = this.#calcInfo()
	}

	addSong(song) {
		this.datas.push(song)
		this.info = {
			count: this.info.count + 1,
			duration: this.info.duration + song.duration
		}

		console.log("push", song);


		return song
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