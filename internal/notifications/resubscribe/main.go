package resubscribe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kolllaka/poma-botv3.0/internal/model"
	m "github.com/kolllaka/poma-botv3.0/internal/notifications/_misc"
)

type fillText func(msg resubscribe) string

var parcingMap map[string]fillText = map[string]fillText{
	"tier": func(msg resubscribe) string {
		return fmt.Sprint(msg.Tier / 1000)
	},
	"user": func(msg resubscribe) string {
		return msg.UserName
	},
	"month": func(msg resubscribe) string {
		return fmt.Sprint(msg.CumulativeMonths)
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
	var resubscribeMsg resubscribe
	json.Unmarshal(msg.Data, &resubscribeMsg)

	index, err := r.checks(resubscribeMsg)
	if err != nil {
		return r.notificationType, nil, fmt.Errorf("%w: %d months, %d tier", err, resubscribeMsg.CumulativeMonths, resubscribeMsg.Tier)
	}

	title := r.confs[index].Title
	words := m.GetArraySwitchingWordsFromTitle(title)

	for _, word := range words {
		newWord, ok := parcingMap[word]
		if !ok {
			continue
		}

		title = strings.Replace(title, fmt.Sprintf("${%s}", word), newWord(resubscribeMsg), 1)
	}

	link, err := m.GetRandomFileLinkFromIndex(r.confs[index].Path)
	if err != nil {
		return r.notificationType, nil, err
	}

	rBody := message{
		Title:   title,
		Link:    link,
		Message: resubscribeMsg.Message.Text,
	}

	var network bytes.Buffer
	json.NewEncoder(&network).Encode(rBody)

	return r.notificationType, network.Bytes(), nil
}

func (r *route) checks(resubscribeMsg resubscribe) (int, error) {
	for i, conf := range r.confs {
		if (resubscribeMsg.Tier >= conf.conditions.Tier) && (resubscribeMsg.CumulativeMonths >= conf.conditions.Month) {
			return i, nil
		}
	}

	return -1, model.ErrorResubscribeNotAllowedConditions
}
