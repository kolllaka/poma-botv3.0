package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/adeithe/go-twitch"

	"github.com/kolllaka/poma-botv3.0/internal/model"
	"github.com/kolllaka/poma-botv3.0/internal/music/files"
	"github.com/kolllaka/poma-botv3.0/internal/music/youtube"
	"github.com/kolllaka/poma-botv3.0/internal/notifications"
	"github.com/kolllaka/poma-botv3.0/internal/rewards"
	"github.com/kolllaka/poma-botv3.0/internal/router"
	"github.com/kolllaka/poma-botv3.0/internal/services"
	"github.com/kolllaka/poma-botv3.0/internal/storage/sqlite"
	"github.com/kolllaka/poma-botv3.0/pkg/config"
	eventsub "github.com/kolllaka/poma-botv3.0/pkg/event-sub"
	"github.com/kolllaka/poma-botv3.0/pkg/logging"
	psqlite "github.com/kolllaka/poma-botv3.0/pkg/sqlite"
	"github.com/kolllaka/poma-botv3.0/pkg/youtubeapi"
)

var (
	logger       *logging.Logger
	twitchReader chan model.RewardMessage = make(chan model.RewardMessage)
)

const (
	LOGS_FOLDER    = "./logs"
	LOGS_LVL       = "debug"
	LOGS_IS_SOURCE = true
	LOGS_IS_JSON   = true

	REWARDS_CFG_PATH       = "./rewards.json"
	NOTIFICATIONS_CFG_PATH = "./notifications.json"
	DB_PATH                = "./storage/music.db"
)

func init() {
	// logger
	fileLogger := logging.NewFileWriter(
		fmt.Sprintf("%s/%s.log", LOGS_FOLDER, time.Now().Format("20060102T1504")),
	)
	multiWriter := io.MultiWriter(fileLogger, os.Stdout)
	logger = logging.NewLogger(
		multiWriter,
		logging.WithIsJSON(LOGS_IS_JSON),
		logging.WithLevel(LOGS_LVL),
		logging.WithAddSource(LOGS_IS_SOURCE),
	)

	// create folder log
	if err := os.MkdirAll(LOGS_FOLDER, os.FileMode(0644)); err != nil {
		logger.Error("Error create folder", logging.ErrAttr(err))

		os.Exit(1)
	}
}

func main() {
	logger.Info("Starting application...")
	// env config
	envCfg := model.NewEnvConfig()
	if err := config.LoadEnv(envCfg); err != nil {
		logger.Error("Error loading env config", logging.ErrAttr(err))

		os.Exit(1)
	}

	logger.Debug("Env config load", logging.AnyAttr("cfg", envCfg))

	// rewards config
	rewardsCfg := model.NewRewardConfig()
	if err := config.LoadJsonByPath(REWARDS_CFG_PATH, rewardsCfg); err != nil {
		logger.Error("Error loading rewards config", logging.ErrAttr(err))

		os.Exit(1)
	}
	logger.Debug("Rewards config load", logging.AnyAttr("cfg", rewardsCfg))

	// notifications config
	notificationsCfg := model.NewNotificationConfig()
	if err := config.LoadJsonByPath(NOTIFICATIONS_CFG_PATH, notificationsCfg); err != nil {
		logger.Error("Error loading notifications config", logging.ErrAttr(err))

		os.Exit(1)
	}
	logger.Debug("Notifications config load", logging.AnyAttr("cfg", notificationsCfg))

	// youtube api
	yApi := youtubeapi.New(envCfg.YoutubeKey)

	// youtube music
	ymusic := youtube.New(logger, yApi)

	// file music
	fmusic := files.New(logger, envCfg.MyPlaylistPath)

	// db
	db, err := psqlite.NewDB(DB_PATH)
	if err != nil {
		logger.Error("Error to open db", logging.ErrAttr(err))

		os.Exit(1)
	}
	defer db.Close()

	// storage
	store := sqlite.New(db)

	// services
	services := services.New(logger, fmusic, ymusic, store)

	// init rewards handler
	rewards := rewards.New(logger, services, twitchReader)
	rewards.InitRewards(rewardsCfg)
	rewards.HandleReward()

	// event sub
	es := eventsub.New(eventsub.Secrets{
		SecretClientId: envCfg.ESBotClientId,
		SecretBearer:   envCfg.ESBotBearer,
	})
	es.SubscribeChannelFollow(eventsub.Condition{
		BroadcasterUserId: fmt.Sprint(envCfg.UserId),
		ModeratorUserId:   fmt.Sprint(envCfg.ESBotId),
	})
	es.SubscribeChannelRaid(eventsub.Condition{
		ToBroadcasterUserID: fmt.Sprint(envCfg.UserId),
	})
	es.SubscribeChannelSubscribe(eventsub.Condition{
		BroadcasterUserId: fmt.Sprint(envCfg.UserId),
	})
	es.SubscribeChannelSubscribeGift(eventsub.Condition{
		BroadcasterUserId: fmt.Sprint(envCfg.UserId),
	})
	es.SubscribeChannelReSubscribe(eventsub.Condition{
		BroadcasterUserId: fmt.Sprint(envCfg.UserId),
	})
	es.SubscribeChannelCheer(eventsub.Condition{
		BroadcasterUserId: fmt.Sprint(envCfg.UserId),
	})
	if err := es.Connect(); err != nil {
		logger.Error("Error to connect twitch event sub", logging.ErrAttr(err))

		os.Exit(1)
	}

	var notificationsReader chan model.NotificationMessage = make(chan model.NotificationMessage)

	go func() {
		for msg := range es.GetEventChan() {
			switch msg.EventType {
			case eventsub.TYPE_CHANNEL_FOLLOW:
				notificationsReader <- model.NotificationMessage{
					RouteType: model.NOTIFICATION_FOLLOW,
					Data:      msg.Event,
				}
			case eventsub.TYPE_CHANNEL_RAID:
				notificationsReader <- model.NotificationMessage{
					RouteType: model.NOTIFICATION_RAID,
					Data:      msg.Event,
				}
			case eventsub.TYPE_CHANNEL_SUBSCRIBE:
				notificationsReader <- model.NotificationMessage{
					RouteType: model.NOTIFICATION_SUBSCRIBE,
					Data:      msg.Event,
				}
			case eventsub.TYPE_CHANNEL_SUBSCRIBE_GIFT:
				notificationsReader <- model.NotificationMessage{
					RouteType: model.NOTIFICATION_SUBGIFT,
					Data:      msg.Event,
				}
			case eventsub.TYPE_CHANNEL_RE_SUBSCRIBE:
				notificationsReader <- model.NotificationMessage{
					RouteType: model.NOTIFICATION_RESUBSCRIBE,
					Data:      msg.Event,
				}
			case eventsub.TYPE_CHANNEL_CHEER:
				notificationsReader <- model.NotificationMessage{
					RouteType: model.NOTIFICATION_CHEER,
					Data:      msg.Event,
				}

			}
		}
	}()

	// init notifications handler
	notifications := notifications.New(logger, services, notificationsReader)
	notifications.InitNotifications(notificationsCfg)
	notifications.HandleNotification()

	// connect to PubSub
	ps := twitch.PubSub()
	ps.OnShardMessage(onMessage)
	ps.Listen("community-points-channel-v1", envCfg.UserId)
	defer ps.Close()

	// Start Server
	server := router.New(logger, envCfg, services, rewards, notifications)
	router := server.Start()

	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	router.Handle("/audio/", http.StripPrefix("/audio/", http.FileServer(http.Dir(envCfg.MyPlaylistPath))))

	for i, path := range rewards.GetPathToUrl() {
		server.RegFileServer(fmt.Sprintf("/%s%d/", model.REWARD_NAME, i), path)
	}

	go http.ListenAndServe(":"+strconv.Itoa(envCfg.Port), router)
	logger.Info("server start", logging.IntAttr("port", envCfg.Port))

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	<-sc
	logger.Info("Stopping...")

}

func onMessage(shard int, topic string, data []byte) {
	reward := model.Reward{}
	json.Unmarshal(data, &reward)

	type rewardMessageData struct {
		IsReward bool   `json:"is_reward,omitempty"`
		Username string `json:"username,omitempty"`
		Text     string `json:"text,omitempty"`
	}

	if reward.Type == "reward-redeemed" {
		rewardName := reward.Data.Redemption.Reward.Title

		rewardMessageData := rewardMessageData{
			IsReward: true,
			Username: reward.Data.Redemption.User.DisplayName,
			Text:     reward.Data.Redemption.UserInput,
		}

		data, err := json.Marshal(rewardMessageData)
		if err != nil {
			logger.Error("error on onMessage", logging.ErrAttr(err), logging.AnyAttr("data", rewardMessageData))

			return
		}

		rmsg := model.RewardMessage{
			RouteName: rewardName,

			Data: data,
		}

		twitchReader <- rmsg
	}
}
