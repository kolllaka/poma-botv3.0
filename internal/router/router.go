package router

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/kolllaka/poma-botv3.0/internal/model"
	"github.com/kolllaka/poma-botv3.0/internal/notifications"
	"github.com/kolllaka/poma-botv3.0/internal/rewards"
	"github.com/kolllaka/poma-botv3.0/internal/services"
	"github.com/kolllaka/poma-botv3.0/pkg/logging"
)

const (
	AUG   = "aug"
	MUSIC = "music"

	TEMPLATE_PATH = "template/*.html"
)

var (
	tmpl     *template.Template
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Пропускаем любой запрос
		},
	}
)

func init() {
	tmpl = template.Must(template.ParseGlob(TEMPLATE_PATH))
}

type server struct {
	logger *logging.Logger
	conf   *model.EnvConfig

	services      services.Service
	rewards       rewards.Rewards
	notifications notifications.Notifications

	clients map[string]*websocket.Conn
	router  *http.ServeMux
}

func New(
	logger *logging.Logger,
	conf *model.EnvConfig,

	services services.Service,
	rewards rewards.Rewards,
	notifications notifications.Notifications,
) *server {
	return &server{
		logger:        logger,
		conf:          conf,
		services:      services,
		rewards:       rewards,
		notifications: notifications,
		clients:       make(map[string]*websocket.Conn),
	}
}

func (s *server) Start() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/"+AUG, s.aug)
	router.HandleFunc("/"+AUG+"/ws", s.augws)

	router.HandleFunc("/"+model.NOTIFICATION_RAID, s.raid)
	router.HandleFunc("/"+model.NOTIFICATION_RAID+"/ws", s.raidws)

	router.HandleFunc("/"+model.NOTIFICATION_SUBSCRIBE, s.subscribe)
	router.HandleFunc("/"+model.NOTIFICATION_SUBSCRIBE+"/ws", s.subscribews)

	router.HandleFunc("/"+model.NOTIFICATION_FOLLOW, s.follow)
	router.HandleFunc("/"+model.NOTIFICATION_FOLLOW+"/ws", s.followws)

	router.HandleFunc("/"+model.NOTIFICATION_SUBGIFT, s.subgift)
	router.HandleFunc("/"+model.NOTIFICATION_SUBGIFT+"/ws", s.subgiftws)

	router.HandleFunc("/"+model.NOTIFICATION_CHEER, s.cheer)
	router.HandleFunc("/"+model.NOTIFICATION_CHEER+"/ws", s.cheerws)

	router.HandleFunc("/"+MUSIC, s.music)
	router.HandleFunc("/"+MUSIC+"/ws", s.musicws)

	router.HandleFunc("/api/yplaylist", s.yplaylist)

	s.router = router

	return router
}

func (s *server) RegFileServer(path string, root string) {
	s.router.Handle(path, http.StripPrefix(path, http.FileServer(http.Dir(root))))
}

// augury
func (s *server) aug(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, AUG+".html", nil)
}

func (s *server) augws(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	s.clients[AUG] = conn
	defer delete(s.clients, AUG)

	go func() {
		for {
			mt, message, err := conn.ReadMessage()

			if err != nil || mt == websocket.CloseMessage {
				s.logger.Warn("error from socket", logging.ErrAttr(err))

				break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана
			}

			go s.handleMessage(message)
		}
	}()

	for {
		reward := <-s.rewards.GetRewardChannel(model.REWARD_AUGURY)

		s.writeByteMsg(AUG, reward)

		time.Sleep(10 * time.Second)
	}
}

// raid
func (s *server) raid(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, model.NOTIFICATION_RAID+".html", nil)
}

func (s *server) raidws(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	s.clients[model.NOTIFICATION_RAID] = conn
	defer delete(s.clients, model.NOTIFICATION_RAID)

	go func() {
		for {
			mt, message, err := conn.ReadMessage()

			if err != nil || mt == websocket.CloseMessage {
				s.logger.Warn("error from socket", logging.ErrAttr(err))

				break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана
			}

			go s.handleMessage(message)
		}
	}()

	for {
		notification := <-s.notifications.GetNotificationChannel(model.NOTIFICATION_RAID)

		s.logger.Debug("recieve notification",
			logging.AnyAttr("notification", notification),
		)

		s.writeByteMsg(model.NOTIFICATION_RAID, notification)

		time.Sleep(10 * time.Second)
	}
}

// subscribe
func (s *server) subscribe(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, model.NOTIFICATION_SUBSCRIBE+".html", nil)
}
func (s *server) subscribews(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	s.clients[model.NOTIFICATION_SUBSCRIBE] = conn
	defer delete(s.clients, model.NOTIFICATION_SUBSCRIBE)

	go func() {
		for {
			mt, message, err := conn.ReadMessage()

			if err != nil || mt == websocket.CloseMessage {
				s.logger.Warn("error from socket", logging.ErrAttr(err))

				break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана
			}

			go s.handleMessage(message)
		}
	}()

	for {
		notification := <-s.notifications.GetNotificationChannel(model.NOTIFICATION_SUBSCRIBE)

		s.logger.Debug("recieve notification",
			logging.AnyAttr("notification", notification),
		)

		s.writeByteMsg(model.NOTIFICATION_SUBSCRIBE, notification)

		time.Sleep(10 * time.Second)
	}
}

// follow
func (s *server) follow(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, model.NOTIFICATION_FOLLOW+".html", nil)
}
func (s *server) followws(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	s.clients[model.NOTIFICATION_FOLLOW] = conn
	defer delete(s.clients, model.NOTIFICATION_FOLLOW)

	go func() {
		for {
			mt, message, err := conn.ReadMessage()

			if err != nil || mt == websocket.CloseMessage {
				s.logger.Warn("error from socket", logging.ErrAttr(err))

				break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана
			}

			go s.handleMessage(message)
		}
	}()

	for {
		notification := <-s.notifications.GetNotificationChannel(model.NOTIFICATION_FOLLOW)

		s.logger.Debug("recieve notification",
			logging.AnyAttr("notification", notification),
		)

		s.writeByteMsg(model.NOTIFICATION_FOLLOW, notification)

		time.Sleep(10 * time.Second)
	}
}

// subgift
func (s *server) subgift(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, model.NOTIFICATION_SUBGIFT+".html", nil)
}
func (s *server) subgiftws(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	s.clients[model.NOTIFICATION_SUBGIFT] = conn
	defer delete(s.clients, model.NOTIFICATION_SUBGIFT)

	go func() {
		for {
			mt, message, err := conn.ReadMessage()

			if err != nil || mt == websocket.CloseMessage {
				s.logger.Warn("error from socket", logging.ErrAttr(err))

				break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана
			}

			go s.handleMessage(message)
		}
	}()

	for {
		notification := <-s.notifications.GetNotificationChannel(model.NOTIFICATION_SUBGIFT)

		s.logger.Debug("recieve notification",
			logging.AnyAttr("notification", notification),
		)

		s.writeByteMsg(model.NOTIFICATION_SUBGIFT, notification)

		time.Sleep(10 * time.Second)
	}
}

// cheer
func (s *server) cheer(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, model.NOTIFICATION_CHEER+".html", nil)
}
func (s *server) cheerws(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	s.clients[model.NOTIFICATION_CHEER] = conn
	defer delete(s.clients, model.NOTIFICATION_CHEER)

	go func() {
		for {
			mt, message, err := conn.ReadMessage()

			if err != nil || mt == websocket.CloseMessage {
				s.logger.Warn("error from socket", logging.ErrAttr(err))

				break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана
			}

			go s.handleMessage(message)
		}
	}()

	for {
		notification := <-s.notifications.GetNotificationChannel(model.NOTIFICATION_CHEER)

		s.logger.Debug("recieve notification",
			logging.AnyAttr("notification", notification),
		)

		s.writeByteMsg(model.NOTIFICATION_CHEER, notification)

		time.Sleep(10 * time.Second)
	}
}

// music
func (s *server) music(w http.ResponseWriter, r *http.Request) {
	playlist, err := s.services.GetMyPlaylist(false)
	if err != nil {
		s.logger.Error("error on GetMyPlaylist", logging.AnyAttr("playlist", playlist), logging.ErrAttr(err))
	}

	s.logger.Debug("get playlist", logging.AnyAttr("playlist", playlist))

	var network bytes.Buffer
	json.NewEncoder(&network).Encode(playlist)

	tmpl.ExecuteTemplate(w, MUSIC+".html", network.String())
}
func (s *server) musicws(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	s.clients[MUSIC] = conn
	defer delete(s.clients, MUSIC)

	go func() {
		for {
			mt, message, err := conn.ReadMessage()

			if err != nil || mt == websocket.CloseMessage {
				s.logger.Warn("error from socket", logging.ErrAttr(err))

				break // Выходим из цикла, если клиент пытается закрыть соединение или связь прервана
			}

			data := MsgFromSocket{}
			if err := json.Unmarshal(message, &data); err != nil {
				s.logger.Error("Unmarshal error", logging.ErrAttr(err))
			}

			switch data.Reason {
			case "addDuration":
				var dataMusic SongFromSocket
				if err := json.Unmarshal(data.Data, &dataMusic); err != nil {
					s.logger.Error("Unmarshal error", logging.AnyAttr("data", data), logging.ErrAttr(err))

					continue
				}

				s.logger.Debug("get data", logging.AnyAttr("song", dataMusic))

				if err := s.services.StoreDuration(&model.Music{
					Link:     dataMusic.Link,
					Duration: dataMusic.Duration,
				}); err != nil {
					s.logger.Error("error from StoreDuration", logging.AnyAttr("song", dataMusic), logging.ErrAttr(err))

					continue
				}

				s.logger.Debug("store duration", logging.AnyAttr("song", dataMusic))
			}

			go s.handleMessage(message)
		}
	}()

	// ! DELETE
	// go func() {

	// 	data := services.Responce{
	// 		Source:   1,
	// 		IsReward: true,
	// 		Data: model.Music{
	// 			Name:     "Upbeat Battle Music　テンションが上がるバトル音楽　「 Devine Fencer 」",
	// 			Link:     "jO7UnKF-tEw",
	// 			Author:   "Kolliaka",
	// 			Duration: 407,
	// 		},
	// 	}

	// 	var network bytes.Buffer
	// 	json.NewEncoder(&network).Encode(data)

	// 	s.writeByteMsg(MUSIC, network.Bytes())
	// }()

	for {
		mBody := <-s.rewards.GetRewardChannel(model.REWARD_MUSIC)

		s.writeByteMsg(MUSIC, mBody)
	}
}

// TODO
func (s *server) yplaylist(w http.ResponseWriter, r *http.Request) {
	// type Playlist struct {
	// 	Link string `json:"link,omitempty"`
	// }

	// w.Header().Set("Content-Type", "application/json")
	// resp := Playlist{}
	// err := json.NewDecoder(r.Body).Decode(&resp)
	// if err != nil {
	// 	s.logger.Error("Decode error", logging.ErrAttr(err))
	// }

	// s.logger.Info("request from server", logging.StringAttr("link", resp.Link))

	// playlist := s.services.GetPlaylistBy(resp.Link)

	// songs := playlist.ListOfSongsFromPlaylist(resp.Link, s.conf.YoutubeKey, "")
	// var songsLink []string
	// listOfSongs := []*model.Playlist{}

	// for i, song := range songs {
	// 	songsLink = append(songsLink, song.Link)

	// 	if (i+1)%10 == 0 {
	// 		listSongsInfo := playlist.ReqSongInfo(strings.Join(songsLink, ","), s.conf.YoutubeKey)
	// 		listOfSongs = append(listOfSongs, listSongsInfo...)

	// 		songsLink = []string{}
	// 	}

	// 	if i+1 == len(songs) {
	// 		listSongsInfo := playlist.ReqSongInfo(strings.Join(songsLink, ","), s.conf.YoutubeKey)
	// 		listOfSongs = append(listOfSongs, listSongsInfo...)
	// 	}
	// }

	// json.NewEncoder(w).Encode(listOfSongs)
}

// sss
func (s *server) writeByteMsg(typeMsg string, message []byte) {
	conn := s.clients[typeMsg]
	conn.WriteMessage(websocket.TextMessage, message)
}

func (s *server) handleMessage(message []byte) {
	s.logger.Info("message from socket", logging.AnyAttr("message", message))
}
