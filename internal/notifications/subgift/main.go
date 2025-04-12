package subgift

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kolllaka/poma-botv3.0/internal/model"
	m "github.com/kolllaka/poma-botv3.0/internal/notifications/_misc"
)

type fillText func(msg subgift) string

var parcingMap map[string]fillText = map[string]fillText{
	"user": func(msg subgift) string {
		return msg.UserName
	},
	"gift_count": func(msg subgift) string {
		return fmt.Sprint(msg.Total)
	},
	"tier": func(msg subgift) string {
		return fmt.Sprint(msg.Tier / 1000)
	},
	"total_gift": func(msg subgift) string {
		return fmt.Sprint(msg.CumulativeTotal)
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
	var subgiftMsg subgift
	json.Unmarshal(msg.Data, &subgiftMsg)

	index, err := r.checks(subgiftMsg)
	if err != nil {
		return r.notificationType, nil, fmt.Errorf("%w: %d tier, %t isAnonym", err, subgiftMsg.Tier, subgiftMsg.IsAnonymous)
	}

	title := r.confs[index].Title
	words := m.GetArraySwitchingWordsFromTitle(title)

	for _, word := range words {
		newWord, ok := parcingMap[word]
		if !ok {
			continue
		}

		title = strings.Replace(title, fmt.Sprintf("${%s}", word), newWord(subgiftMsg), 1)
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

func (r *route) checks(subgiftMsg subgift) (int, error) {
	for i, conf := range r.confs {
		if (subgiftMsg.IsAnonymous == conf.conditions.IsAnonym) && (subgiftMsg.Tier >= conf.conditions.Tier) {
			return i, nil
		}
	}

	return -1, model.ErrorSubgiftNotAllowedConditions
}
