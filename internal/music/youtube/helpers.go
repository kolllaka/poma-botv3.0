package youtube

import (
	"fmt"
	"regexp"
	"strconv"
)

func parseSongIdByText(text string) (string, error) {
	urlReg := regexp.MustCompile(`https://www.youtube.com/watch\?v=([a-zA-Z0-9_-]*)|https://youtu.be/([a-zA-Z0-9_-]*)`)
	resUrl := urlReg.FindStringSubmatch(text)

	if len(resUrl) > 0 {
		if resUrl[1] == "" {
			return resUrl[2], nil
		}

		return resUrl[1], nil
	}

	return "", fmt.Errorf("unknown song: %s", text)
}

func formatTimeFromYoutube(time string) int {
	var fotmatTime int = 0
	reg := regexp.MustCompile(`PT(\d*[HMS])(\d*[HMS])(\d*[HMS])|PT(\d*[HMS])(\d*[HMS])|PT(\d*[HMS])`)
	res := reg.FindStringSubmatch(time)

	if res == nil {
		return 0
	}

	for _, val := range res[1:] {
		if val != "" {
			sim := string(val[len(val)-1])
			chislo, _ := strconv.Atoi(val[:len(val)-1])

			switch sim {
			case "H":
				fotmatTime += chislo * 3600
			case "M":
				fotmatTime += chislo * 60
			case "S":
				fotmatTime += chislo
			}
		}
	}

	return fotmatTime
}
