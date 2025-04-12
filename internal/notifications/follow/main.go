package follow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kolllaka/poma-botv3.0/internal/model"
	m "github.com/kolllaka/poma-botv3.0/internal/notifications/_misc"
)

type fillText func(msg follow) string

var parcingMap map[string]fillText = map[string]fillText{
	"user": func(msg follow) string {
		return msg.UserName
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
	var followMsg follow
	json.Unmarshal(msg.Data, &followMsg)

	if len(r.confs) < 1 {
		return r.notificationType, nil, model.ErrorEmptyFollowConf
	}

	title := r.confs[0].Title
	words := m.GetArraySwitchingWordsFromTitle(title)

	for _, word := range words {
		newWord, ok := parcingMap[word]
		if !ok {
			continue
		}

		title = strings.Replace(title, fmt.Sprintf("${%s}", word), newWord(followMsg), 1)
	}

	link, err := m.GetRandomFileLinkFromIndex(r.confs[0].Path)
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
