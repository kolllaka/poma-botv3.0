export class Notification {
	constructor(selector, queue) {
		this.$el = document.getElementById(selector)
		this.queue = queue

		this.isPlay = false

		this.video = this.$el.querySelector("video")
		this.title = this.$el.querySelector("[data-field='title']")
		if (this.$el.querySelector("[data-field='message']")) {
			this.message = this.$el.querySelector("[data-field='message']")
		}
	}

	startNotification(data, timeout = 10000) {
		if (data) {
			this.isPlay = true
			this.video.src = data.link
			this.video.onloadedmetadata = () => {
				this.video.play()
				this.shuffleDisplayNotification("flex")
				// setTimeout(() => {
				// 	this.stopNotification();

				// 	this.startNotification(data, timeout = 10000)
				// }, timeout)
			}

			this.title.innerHTML = data.title
			if (data.message) {
				this.message.innerHTML = data.message
			}
		}
	}

	stopNotification() {
		this.video.pause()
		this.video.src = ""

		this.shuffleDisplayNotification("none")
	}

	shuffleDisplayNotification(display = "none") {
		this.$el.style.display = display
	}
}

