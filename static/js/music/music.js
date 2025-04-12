// ws
import { ConnectWS } from '../lib/ws.js';
// check duration
import { GetMetaData } from '../lib/get_meta_data.js';
// default settings
import * as settings from './settings.js';
// my Playlist
import { myPlaylistController } from './my_playlist/controller.js'
// reward Playlist
import { rewardPlaylistController } from './reward_playlist/controller.js'
// misc
import { getLinkTemplate } from '../lib/misc.js';
// control panel
import * as cpUi from '../lib/control_panel/ui.js';
import * as cpData from '../lib/control_panel/data.js';
import { Controller as CpController } from '../lib/control_panel/controller.js';
// my player
import { Player as MyPlayer } from './players/my_player/player.js';
// youtube player
import { Player as RewardPlayer } from './players/reward_player/player.js';
// players control
import * as pc from './players/controller.js';

// data mapping
import * as mapping from './mapping_data.js';


console.log(myPlaylist);

// ?check meta of song
const getDuration = new GetMetaData()

// console.log("myPlaylist", myPlaylistData);
// myPlaylistData = myPlaylistData.data


mapping.playlistServerToSongArray(myPlaylist).forEach((song) => {
	console.log(song);

	if (song.duration <= 0) {
		getDuration.addSong(song)

		return
	}

	myPlaylistController.addSong(song)
})
getDuration.onloadedmetadata = function () {
	const duration = Math.floor(getDuration.audio.duration)

	getDuration.audio.pause()
	getDuration.currentData.duration = duration


	console.log("send to socket", getDuration.currentData);

	sendSocket(mapping.songToSongServer("addDuration", getDuration.currentData))

	const data = getDuration.checkArray.pop()

	if (data) {
		getDuration.currentData = data
		getDuration.setSrc(getDuration.currentData)
	}
}


// init websocket
const url = 'ws://127.0.0.1:8080/music/ws'
const ws = new ConnectWS(url)
ws.onOpen = function (event) {
	console.log(`ws conn open to  ${event.target.url}`);
	getDuration.start()
}
ws.onMessage = function (event) {
	const data = JSON.parse(event.data)
	console.log("data from ws", data);

	proccessDataFromWS(data)
}
ws.setup()
const sendSocket = (data) => {
	ws.socket.send(JSON.stringify(data))
}

function proccessDataFromWS(data) {
	const song = mapping.songServerToSong(data)

	if (data.is_reward) {
		//! put on reward playlist
		rewardPlaylistController.addSong(song)

		return
	}

	//! put on my playlist
	myPlaylistController.addSong(song)

}




// load settings
const defaultsettings = settings.load()

myPlaylistController.onDeleteHandler = function (target) {
	const index = parseInt(target.dataset.index)

	const song = this.deleteSong(index)

	console.log(song);
}

// control panel
const newcpData = new cpData.Song({
	volume: defaultsettings.volume,
})
const newcpUi = new cpUi.SongUI("controlpanel")
const newcpController = new CpController("controlpanel", newcpUi, newcpData)
newcpController.onVolumeHandler = function (value) {
	myPlayer.setVolume(value)
	rewardPlayer.setVolume(value)
	localStorage.setItem(settings.VOLUME_KEY, value)
}
newcpController.onSkipHandler = function () {

	const song = newcpController.getSong()

	console.log("send to socket", song);

	if (song) {
		sendSocket(mapping.songToSongServer("skip", song))
	}

	nextSong()
}
newcpController.onShuffleHandler = function () {
	myPlaylistController.shuffle()
}
newcpController.onPlayHandler = function (target) {
	const status = target.dataset.btn
	target.dataset.btn = (status == cpUi.BTN_TEXT_PLAY) ? cpUi.BTN_TEXT_PAUSE : cpUi.BTN_TEXT_PLAY

	switch (status) {
		case cpUi.BTN_TEXT_PAUSE:
			pControl.changeStatus(pc.Status.PAUSE)

			myPlayer.pause()
			rewardPlayer.pauseVideo()

			break
		case cpUi.BTN_TEXT_PLAY:
			pControl.changeStatus(pc.Status.PLAY)


			switch (pControl.whoPlay) {
				case pc.Player.YOUTUBE:
					rewardPlayer.playVideo()

					break
				case pc.Player.AUDIO:
					myPlayer.play()

					break
				case pc.Player.NOTHING:
					nextSong()

					break
			}

			break
		default:

	}
}


// my playlist player
const myPlayer = new MyPlayer("myplayer", {
	volume: defaultsettings.volume,
})
myPlayer.onEndedHandler = function (e) {
	nextSong()
}

// reward playlist player
const rewardPlayer = new RewardPlayer("rewardplayer", {
	onReady: onPlayerReady,
	onStateChange: onPlayerStateChange,
	onError: onError,
})
function onPlayerReady() {
	rewardPlayer.loadVideoById("QxtKHo0iMa4")
	rewardPlayer.setVolume(defaultsettings.volume)

	rewardPlayer.playVideo()
}
function onPlayerStateChange(event) {

	if (event.data == YT.PlayerState.ENDED) {
		nextSong()
	}

	if (event.data == YT.PlayerState.PLAYING) {
		pControl.changePlayer(pc.Player.YOUTUBE)
	}

}
function onError(event) {
	sendSocket(mapping.songToSongServer("error", pControl.getCurrentSong(), event.data))
	nextSong()
}

// players controller
const pControl = new pc.Controller()


// next Song Handler
function nextSong() {
	const ysong = rewardPlaylistController.deleteSong(0)

	if (ysong) {
		console.log(ysong);
		rewardPlayer.loadVideoById(ysong.link)
		newcpController.updateInfo({
			author: ysong.author || "-",
			name: getLinkTemplate(ysong)
		})
		newcpController.getBtnPlay().dataset.btn = cpUi.BTN_TEXT_PAUSE
		pControl.setCurrentSong(ysong)

		return
	}

	rewardPlayer.stopVideo()
	pControl.changePlayer(pc.Player.AUDIO)
	const song = myPlaylistController.deleteSong(0)
	if (song) {
		myPlayer.setLink(song.link)
		myPlayer.play()
		newcpController.updateInfo({
			author: song.author || "-",
			name: song.name
		})
		newcpController.getBtnPlay().dataset.btn = cpUi.BTN_TEXT_PAUSE
		pControl.setCurrentSong(song)

		return
	}

	myPlayer.pause()
	myPlayer.setLink("")
	pControl.changePlayer(pc.Player.NOTHING)
	newcpController.updateInfo(cpData.DEFAULT_SONG_INFO)
	newcpController.getBtnPlay().dataset.btn = cpUi.BTN_TEXT_PLAY
	pControl.setCurrentSong({})

	return
}