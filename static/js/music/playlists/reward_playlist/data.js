import { Playlist } from '../playlist.js';

function mapping(data = []) {
	return {
		name: data.data.title,
		link: data.data.link,
		author: data.data.author,
		duration: data.data.duration
	}
}

export function data(datas = []) {
	return new Playlist(datas, mapping)
}