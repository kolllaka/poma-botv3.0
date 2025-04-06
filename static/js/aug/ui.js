class Augury {
	constructor(selector) {
		this.$el = document.getElementById(selector)

		this.#render()
	}

	#render() {
		this.$el.classList.add("augury")
		this.$el.innerHTML = getTemplate()
	}

	change(options) {
		const video = this.$el.querySelector(".augury__source").querySelector("video")
		video.src = options.link
		video.onloadedmetadata = () => {
			video.play()
			video.addEventListener('ended', (e) => {
				aug.destroy();
			})
		}

		this.$el.querySelector(".augury__title").innerHTML = options.title
		this.$el.style.display = "block"
	}

	destroy() {
		const video = this.$el.querySelector(".augury__source").querySelector("video")
		video.pause()
		video.src = ""

		this.$el.style.display = "none"
	}
}

const getTemplate = () => {
	return `
	<div class="augury__body">
		<div class="augury__title"></div>
		<div class="augury__source">
			<video onloadstart="this.volume=0.5" src="">
			</video> 
		</div>
	</div>
	`
}