package subscribe

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
	notificationType string
	conf             conf
	Files            []string
}

func NewRoute(notificationType string, rawConf json.RawMessage) *route {
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
		conf:             conf,
		notificationType: notificationType,
		Files:            notificationFiles,
	}
}

func (r *route) RunRoute(msg model.NotificationMessage) (string, []byte, error) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	num := r1.Intn(len(r.Files))
	name := r.Files[num]

	var subMsg subscribe
	json.Unmarshal(msg.Data, &subMsg)

	rBody := message{
		Title: fmt.Sprintf(r.conf.Title, subMsg.UserName, subMsg.Tier),
		Link:  r.getLink(name),
	}

	var network bytes.Buffer
	json.NewEncoder(&network).Encode(rBody)

	return r.notificationType, network.Bytes(), nil
}

func (r *route) getLink(name string) string {
	return filepath.Join(r.conf.Url, name)
}
