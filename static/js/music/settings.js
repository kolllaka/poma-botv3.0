export const VOLUME_KEY = "volume"

const defaultSetting = {
	volume: 15
}

export function load() {
	const volume = (localStorage.getItem(VOLUME_KEY)) ? localStorage.getItem(VOLUME_KEY) : defaultSetting.volume
	localStorage.setItem(VOLUME_KEY, volume)

	return {
		volume: parseInt(volume),
	}
}