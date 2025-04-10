export const Player = Object.freeze({
	NOTHING: Symbol(0),
	YOUTUBE: Symbol(1),
	AUDIO: Symbol(2),
});

export const Status = Object.freeze({
	PLAYING: Symbol(0),
	PAUSE: Symbol(1),
	STOP: Symbol(2),
});

export class Controller {
	constructor() {
		this.song = {}
		this.whoPlay = Player.YOUTUBE
		this.status = Status.PAUSE
	}

	changePlayer(player) {
		this.whoPlay = player
	}

	changeStatus(status) {
		this.status = status
	}

	setCurrentSong(song) {
		this.song = song
	}

	getCurrentSong() {
		return this.song
	}

}