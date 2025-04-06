export class Player {
	constructor(selector, events = {
		onReady: () => {
			console.log('not onPlayerReady implement');
		},
		onStateChange: () => {
			console.log('not onPlayerStateChange implement');
		}
	}) {
		this.player = this.#init(selector, events)
	}

	#init(selector, events) {
		return new YT.Player(selector, {
			height: '200',
			width: '300',
			playerVars: {
				'autoplay': 1,
				'start': 0
			},
			events: events
		});
	}

	playVideo() {
		this.player.playVideo()
	}

	pauseVideo() {
		this.player.pauseVideo()
	}

	stopVideo() {
		this.player.stopVideo()
	}

	loadVideoById(link) {
		this.player.loadVideoById(link)
	}

	setVolume(volume) {
		this.player.setVolume(volume)
	}

}