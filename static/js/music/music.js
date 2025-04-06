// ws
import { ConnectWS } from '../lib/ws.js';
// check duration
import { GetDuration } from '../lib/get_duration.js';
// default settings
import * as settings from './settings.js';
// playlists
import { Controller } from './playlists/controller.js';
import { PlaylistUI } from './playlists/ui.js';
// my Playlist
import * as myUi from './playlists/my_playlist/ui.js';
import * as myData from './playlists/my_playlist/data.js';
// reward Playlist
import * as rewardUi from './playlists/reward_playlist/ui.js';
import * as rewardData from './playlists/reward_playlist/data.js';
// misc
import { getLinkTemplate } from '../lib/misc.js';
// control panel
import * as cpUi from './control_panel/ui.js';
import * as cpData from './control_panel/data.js';
import { Controller as CpController } from './control_panel/controller.js';
// my player
import { Player as MyPlayer } from './players/my_player/player.js';
// youtube player
import { Player as RewardPlayer } from './players/reward_player/player.js';
// players control
import * as pc from './players/controller.js';






// init websocket
const url = 'ws://127.0.0.1:8080/music/ws'
const ws = new ConnectWS(url)
ws.onOpen = function (event) {
	console.log(`ws conn open to  ${url}`);
	getDuration.start()
}
ws.onMessage = function (event) {
	const data = JSON.parse(event.data)
	console.log("data from ws", data);

	proccessDataFromWS(data)
}
ws.onClose = function (event) {
	if (event.wasClean) {
		console.log(`[close] Соединение закрыто чисто, код=${event.code} причина=${event.reason}`);
		ws.socket.close();
	} else {
		// например, сервер убил процесс или сеть недоступна
		// обычно в этом случае event.code 1006
		console.log('[close ws] Соединение прервано');
	}

	setTimeout(function () {
		ws.connect(url)
	}, 1000);
}
ws.onError = function (error) {
	console.log(`[error] ${error}`);
	ws.socket.close();
}
ws.setup()
const sendSocket = (data) => {
	ws.socket.send(JSON.stringify(data))
}

function proccessDataFromWS(data) {
	if (data.is_reward) {
		//! put on reward playlist

		rewardPlaylistController.addSong(data)

		return
	}

	//! put on my playlist
	myPlaylistController.addSong(data)

}

// ?check meta of song
const getDuration = new GetDuration()


myPlaylist.forEach((data) => {
	getDuration.addSong(data)
})
getDuration.onloadedmetadata = function () {
	const duration = Math.floor(getDuration.audio.duration)

	getDuration.audio.pause()
	getDuration.currentData.data.duration = duration


	console.log("send to socket", getDuration.currentData);

	sendSocket({ data: getDuration.currentData, reason: "addDuration" })

	const data = getDuration.checkArray.pop()

	if (data) {
		getDuration.currentData = data
		getDuration.setSrc(getDuration.currentData)
	}
}


// load settings
const defaultsettings = settings.load()

// my Playlist
const myPlaylistUi = new PlaylistUI("myPlaylist", {
	song: myUi.getTemplateSong,
	info: myUi.getTemplateInfo
})

const myDatas = myData.data(myPlaylist)
const myPlaylistController = new Controller("myPlaylist", myPlaylistUi, myDatas)

// reward Playlist
const rewardPlaylistUi = new PlaylistUI("rewardPlaylist", {
	song: rewardUi.getTemplateSong,
	info: rewardUi.getTemplateInfo
})
const rewardDatas = rewardData.data([])
const rewardPlaylistController = new Controller("rewardPlaylist", rewardPlaylistUi, rewardDatas)
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
		sendSocket({
			data: song,
			reason: "skip"
		})
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
	onStateChange: onPlayerStateChange
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

		return
	}

	myPlayer.pause()
	myPlayer.setLink("")
	pControl.changePlayer(pc.Player.NOTHING)
	newcpController.updateInfo(cpData.DEFAULT_SONG_INFO)
	newcpController.getBtnPlay().dataset.btn = cpUi.BTN_TEXT_PLAY

	return
}