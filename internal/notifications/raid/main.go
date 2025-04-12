package raid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/kolllaka/poma-botv3.0/internal/model"
)

type fillText func(msg raid) string

var parcingMap map[string]fillText = map[string]fillText{
	"channel": func(msg raid) string {
		return msg.FromBroadcasterUserName
	},
	"count": func(msg raid) string {
		return fmt.Sprint(msg.Viewers)
	},
}

type route struct {
	notificationType string
	confs            []conf
}

func NewRoute(notificationType string, rawConf json.RawMessage) *route {
	var confs []conf
	json.Unmarshal(rawConf, &confs)

	return &route{
		confs:            confs,
		notificationType: notificationType,
	}
}

func (r *route) RunRoute(msg model.NotificationMessage) (string, []byte, error) {
	var raidMsg raid
	json.Unmarshal(msg.Data, &raidMsg)

	index, err := r.checks(raidMsg)
	if err != nil {
		return r.notificationType, nil, fmt.Errorf("%w: %d viewers", err, raidMsg.Viewers)
	}

	title := r.getTitleFromIndex(index, raidMsg)
	link, err := r.getRandomFileLinkFromIndex(index)
	if err != nil {
		return r.notificationType, nil, err
	}

	rBody := message{
		Title: title,
		Link:  link,
	}

	var network bytes.Buffer
	json.NewEncoder(&network).Encode(rBody)

	return r.notificationType, network.Bytes(), nil
}

func (r *route) getRandomFileLinkFromIndex(index int) (string, error) {
	path := r.confs[index].Path

	isFile, err := isFile(path)
	if err != nil {
		return "", err
	}

	if isFile {
		return path, nil
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}
	if len(files) < 1 {
		return "", fmt.Errorf("%w: %s", ErrorEmptyDirectory, path)
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	num := r1.Intn(len(files))

	return fmt.Sprintf("%s/%s", path, files[num].Name()), nil
}

func (r *route) getTitleFromIndex(index int, msg raid) string {
	title := r.confs[index].Title

	find := 0
	var word string
	var words []string
	for _, s := range title {
		switch s {
		case '$':
			if find == 0 {
				find = 1
			}
		case '{':
			if find == 1 {
				find = 2
			}
		case '}':
			if find == 2 {
				words = append(words, word)
				find = 0
				word = ""
			}
		default:
			if find == 2 {
				word += string(s)
			}
		}
	}

	for _, word := range words {
		newWord, ok := parcingMap[word]
		if !ok {
			continue
		}

		title = strings.Replace(title, fmt.Sprintf("${%s}", word), newWord(msg), 1)
	}

	return title
}

func (r *route) checks(raidMsg raid) (int, error) {
	for i, conf := range r.confs {
		if raidMsg.Viewers >= conf.conditions.Viewers {
			return i, nil
		}
	}

	return -1, ErrorToLowViewers
}

func isFile(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return !f.IsDir(), nil
}
