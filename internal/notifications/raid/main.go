package raid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kolllaka/poma-botv3.0/internal/model"
	m "github.com/kolllaka/poma-botv3.0/internal/notifications/_misc"
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

	title := r.confs[index].Title
	words := m.GetArraySwitchingWordsFromTitle(title)

	for _, word := range words {
		newWord, ok := parcingMap[word]
		if !ok {
			continue
		}

		title = strings.Replace(title, fmt.Sprintf("${%s}", word), newWord(raidMsg), 1)
	}

	link, err := m.GetRandomFileLinkFromIndex(r.confs[index].Path)
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

func (r *route) checks(raidMsg raid) (int, error) {
	for i, conf := range r.confs {
		if raidMsg.Viewers >= conf.conditions.Viewers {
			return i, nil
		}
	}

	return -1, model.ErrorRaidToLowViewers
}
