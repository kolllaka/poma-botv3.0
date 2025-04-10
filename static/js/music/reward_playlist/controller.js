// playlists
import { Controller } from '../../lib/playlists/controller.js';
import { PlaylistUI } from '../../lib/playlists/ui.js';
import { Playlist } from '../../lib/playlists/playlist.js';
// my Playlist
import * as rewardUi from './ui.js';

const rewardPlaylistUi = new PlaylistUI("rewardPlaylist", {
	song: rewardUi.getTemplateSong,
	info: rewardUi.getTemplateInfo
})
const rewardDatas = new Playlist()
export const rewardPlaylistController = new Controller("rewardPlaylist", rewardPlaylistUi, rewardDatas)

