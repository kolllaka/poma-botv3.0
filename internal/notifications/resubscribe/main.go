package resubscribe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/kolllaka/poma-botv3.0/internal/model"
)

type route struct {
	rewardType        string
	conf              conf
	notificationFiles []string
}

func NewRoute(rewardType string, rawConf json.RawMessage) *route {
	var conf conf
	json.Unmarshal(rawConf, &conf)

	files, err := os.ReadDir(conf.Path)
	if err != nil {
		panic(err)
	}

	var notificationFiles []string

	for _, file := range files {
		notificationFiles = append(notificationFiles, file.Name())
	}

	return &route{
		conf:              conf,
		rewardType:        rewardType,
		notificationFiles: notificationFiles,
	}
}

func (r *route) RunRoute(msg model.RewardMessage) (string, []byte, error) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	num := r1.Intn(len(r.notificationFiles))
	name := r.notificationFiles[num]

	var reSubMsg reSubscribe
	json.Unmarshal(msg.Data, &reSubMsg)

	rBody := Message{
		Title: fmt.Sprintf(r.conf.Title, reSubMsg.UserName, reSubMsg.Tier, reSubMsg.CumulativeMonths),
		Link:  r.getLink(name),
		Msg:   reSubMsg.Message.Text,
	}

	var network bytes.Buffer
	json.NewEncoder(&network).Encode(rBody)

	return r.rewardType, network.Bytes(), nil
}

func (r *route) getLink(name string) string {
	return filepath.Join(r.conf.Url, name)
}
