// playlists
import { Controller } from '../../lib/playlists/controller.js';
import { PlaylistUI } from '../../lib/playlists/ui.js';
import { Playlist } from '../../lib/playlists/playlist.js';
// my Playlist
import * as myUi from './ui.js';

const myPlaylistUi = new PlaylistUI("myPlaylist", {
	song: myUi.getTemplateSong,
	info: myUi.getTemplateInfo
})

const myDatas = new Playlist()
export const myPlaylistController = new Controller("myPlaylist", myPlaylistUi, myDatas)
