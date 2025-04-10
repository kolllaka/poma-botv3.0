// song
// 	.name
// 	.duration
// 	.link
// 	.author

// info
// 	.count
// 	.duration

export function songServerToSong(data) {
	return {
		...data,
		name: data.data.title,
		duration: data.data.duration,
		link: data.data.link,
		author: data.data.author,
	}
}

export function playlistServerToSongArray(playlistData = {
	data: [],
}) {
	console.log("playlistFromServerToSongArray", playlistData);


	return playlistData.data.map((song) => {
		return {
			...song,
			name: song.title,
			duration: song.duration,
			link: song.link,
			author: song.author,
		}
	})

}

export function songToSongServer(reason, data) {
	return {
		reason: reason,
		data: data
	}
}