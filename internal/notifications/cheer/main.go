package cheer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kolllaka/poma-botv3.0/internal/model"
	m "github.com/kolllaka/poma-botv3.0/internal/notifications/_misc"
)

type fillText func(msg cheer) string

var parcingMap map[string]fillText = map[string]fillText{
	"user": func(msg cheer) string {
		return msg.UserName
	},
	"bits": func(msg cheer) string {
		return fmt.Sprint(msg.Bits)
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
	var cheerMsg cheer
	json.Unmarshal(msg.Data, &cheerMsg)

	index, err := r.checks(cheerMsg)
	if err != nil {
		return r.notificationType, nil, fmt.Errorf("%w: %d bits, %t isAnonym", err, cheerMsg.Bits, cheerMsg.IsAnonymous)
	}

	title := r.confs[index].Title
	words := m.GetArraySwitchingWordsFromTitle(title)

	for _, word := range words {
		newWord, ok := parcingMap[word]
		if !ok {
			continue
		}

		title = strings.Replace(title, fmt.Sprintf("${%s}", word), newWord(cheerMsg), 1)
	}

	link, err := m.GetRandomFileLinkFromIndex(r.confs[index].Path)
	if err != nil {
		return r.notificationType, nil, err
	}

	rBody := message{
		Title:   title,
		Link:    link,
		Message: cheerMsg.Message,
	}

	var network bytes.Buffer
	json.NewEncoder(&network).Encode(rBody)

	return r.notificationType, network.Bytes(), nil
}

func (r *route) checks(cheerMsg cheer) (int, error) {
	for i, conf := range r.confs {
		if (cheerMsg.IsAnonymous == conf.conditions.IsAnonym) && (cheerMsg.Bits >= conf.conditions.Bits) {
			return i, nil
		}
	}

	return -1, model.ErrorCheerNotAllowedConditions
}
