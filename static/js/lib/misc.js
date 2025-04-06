export function durationFormat(duration) {
	let hours = (duration / 3600) | 0

	let minutes = ((duration % 3600) / 60) | 0

	let seconds = (duration % 60) | 0
	if (seconds < 10) seconds = `0` + seconds

	if (hours == 0) {
		return `${minutes}:${seconds}`
	}

	if (minutes < 10) minutes = `0` + minutes
	return `${hours}:${minutes}:${seconds}`
}
export function shuffle(datas = []) {
	let currentIndex = datas.length, randomIndex;

	while (currentIndex > 0) {

		randomIndex = Math.floor(Math.random() * currentIndex);
		currentIndex--;

		[datas[currentIndex], datas[randomIndex]] = [
			datas[randomIndex], datas[currentIndex]];
	}
}

export function getLinkTemplate(song) {
	return `
		<a href="https://www.youtube.com/watch?v=${song.link}" target="_blank">${song.name}</a>
	`
}