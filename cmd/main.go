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
	"github.com/kolllaka/poma-botv3.0/internal/rewards"
	"github.com/kolllaka/poma-botv3.0/internal/router"
	"github.com/kolllaka/poma-botv3.0/internal/services"
	"github.com/kolllaka/poma-botv3.0/internal/storage/sqlite"
	"github.com/kolllaka/poma-botv3.0/pkg/config"
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
	LOGS_LVL       = "info"
	LOGS_IS_SOURCE = false
	LOGS_IS_JSON   = true

	REWARDS_CFG_PATH = "./rewards.json"
	DB_PATH          = "./storage/music.db"
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
	reward := rewards.New(logger, services, twitchReader)
	reward.InitRewards(rewardsCfg)
	reward.HandleReward()

	// connect to PubSub
	ps := twitch.PubSub()
	ps.OnShardMessage(onMessage)
	ps.Listen("community-points-channel-v1", envCfg.UserId)
	defer ps.Close()

	// Start Server
	server := router.New(logger, envCfg, services, reward)
	router := server.Start()

	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	router.Handle("/audio/", http.StripPrefix("/audio/", http.FileServer(http.Dir(envCfg.MyPlaylistPath))))

	for i, path := range reward.GetPathToUrl() {
		server.RegFileServer(fmt.Sprintf("/%s%d/", model.REWARD_NAME, i), path)
	}

	go http.ListenAndServe(":"+strconv.Itoa(envCfg.Port), router)
	logger.Info("server start", logging.IntAttr("port", envCfg.Port))
	//
	//
	//
	//
	//

	//
	//
	//
	//
	// twitchReader <- model.RewardMessage{
	// 	IsReward: true,
	// 	Name:     "гадание",
	// 	Username: "Kollliaka",
	// 	Text:     "",
	// }

	// twitchReader <- model.RewardMessage{
	// 	IsReward: true,
	// 	Name:     "музик",
	// 	Username: "Kolliaka",
	// 	Text:     "https://www.youtube.com/watch?v=jO7UnKF-tEw",
	// }

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	<-sc
	logger.Info("Stopping...")

}

func onMessage(shard int, topic string, data []byte) {
	reward := model.Reward{}
	json.Unmarshal(data, &reward)

	if reward.Type == "reward-redeemed" {
		rewardName := reward.Data.Redemption.Reward.Title
		username := reward.Data.Redemption.User.DisplayName
		text := reward.Data.Redemption.UserInput

		rmsg := model.RewardMessage{
			IsReward: true,
			Name:     rewardName,
			Username: username,
			Text:     text,
		}

		twitchReader <- rmsg
	}
}
