package aug

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
	rewardType string
	conf       conf
	augFiles   []string
}

func NewRoute(rewardType string, rawConf json.RawMessage) *route {
	var conf conf
	json.Unmarshal(rawConf, &conf)

	files, err := os.ReadDir(conf.Path)
	if err != nil {
		panic(err)
	}

	var augFiles []string

	for _, file := range files {
		augFiles = append(augFiles, file.Name())
	}

	return &route{
		conf:       conf,
		rewardType: rewardType,
		augFiles:   augFiles,
	}
}

func (r *route) RunRoute(msg model.RewardMessage) (string, []byte, error) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	num := r1.Intn(len(r.augFiles))
	name := r.augFiles[num]

	rBody := Message{
		Title: fmt.Sprintf(r.conf.Title, msg.Username),
		Link:  r.getLink(name),
	}

	var network bytes.Buffer
	json.NewEncoder(&network).Encode(rBody)

	return r.rewardType, network.Bytes(), nil
}

func (r *route) getLink(name string) string {
	return filepath.Join(r.conf.Url, name)
}
