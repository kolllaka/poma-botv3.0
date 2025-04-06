package youtubeapi

import (
	"net/url"
)

const (
	getVideosMethod = "videos"
)

// TODO https://www.googleapis.com/youtube/v3/videos?id=%s&key=%s&part=snippet,contentDetails,statistics

func (yc *yClient) ReqSongInfoById(songId string) (SongResp, error) {
	q := url.Values{}
	q.Add("part", "snippet,contentDetails,statistics")
	q.Add("id", songId)

	body, err := yc.doRequest(getVideosMethod, q, nil)
	if err != nil {
		return SongResp{}, err
	}

	var resp SongResp
	if err := parseResponseBody(body, &resp); err != nil {
		return SongResp{}, err
	}

	return resp, nil
}
