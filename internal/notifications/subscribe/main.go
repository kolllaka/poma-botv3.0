package subscribe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kolllaka/poma-botv3.0/internal/model"
	m "github.com/kolllaka/poma-botv3.0/internal/notifications/_misc"
)

type fillText func(msg subscribe) string

var parcingMap map[string]fillText = map[string]fillText{
	"tier": func(msg subscribe) string {
		return fmt.Sprint(msg.Tier / 1000)
	},
	"user": func(msg subscribe) string {
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
	var subscribeMsg subscribe
	json.Unmarshal(msg.Data, &subscribeMsg)

	index, err := r.checks(subscribeMsg)
	if err != nil {
		return r.notificationType, nil, fmt.Errorf("%w: %d tier, %t isGift", err, subscribeMsg.Tier, subscribeMsg.IsGift)
	}

	title := r.confs[index].Title
	words := m.GetArraySwitchingWordsFromTitle(title)

	for _, word := range words {
		newWord, ok := parcingMap[word]
		if !ok {
			continue
		}

		title = strings.Replace(title, fmt.Sprintf("${%s}", word), newWord(subscribeMsg), 1)
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

func (r *route) checks(subscribeMsg subscribe) (int, error) {
	for i, conf := range r.confs {
		if (subscribeMsg.IsGift == conf.conditions.IsGift) && (subscribeMsg.Tier >= conf.conditions.Tier) {
			return i, nil
		}
	}

	return -1, model.ErrorSubscribeNotAllowedConditions
}
