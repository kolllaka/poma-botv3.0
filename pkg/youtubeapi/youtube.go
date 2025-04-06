package youtubeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

const (
	ytHost = "www.googleapis.com"
	ytPath = "youtube/v3"
)

type Api interface {
	ReqSongInfoById(song string) (SongResp, error)
}
type yClient struct {
	host       string
	basePath   string
	tokenQuery url.Values
	token      string
	client     http.Client
}

func New(
	token string,
) Api {
	q := url.Values{}
	q.Add("key", token)

	return &yClient{
		host:       ytHost,
		basePath:   ytPath,
		tokenQuery: q,
		token:      token,
		client:     http.Client{},
	}
}

func parseResponseBody(data Response, resp any) error {
	// TODO
	if data.StatusCode != http.StatusOK {
		return fmt.Errorf("error with code(%d)", data.StatusCode)
	}

	if err := json.Unmarshal(data.Data, &resp); err != nil {
		return err
	}

	return nil
}

// TODO Scheme https
// TODO Host www.googleapis.com
// TODO Base Path youtube/v3/
// TODO Method Path playlistItems
// TODO token Query key=%s
// TODO Query part=contentDetails&playlistId=%s&&maxResults=50

func (yc *yClient) doRequest(method string, query url.Values, data io.Reader) (Response, error) {
	const op = "clients.telegram.doRequest"

	u := url.URL{
		Scheme: "https",
		Host:   yc.host,
		Path:   path.Join(yc.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), data)
	if err != nil {
		return Response{}, fmt.Errorf("%s: can't do request: %w", op, err)
	}
	// req.URL.RawQuery = yc.tokenQuery.Encode()
	query.Add("key", yc.token)
	req.URL.RawQuery = query.Encode()
	if data != nil {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("accept", "application/json")
	}

	resp, err := yc.client.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("%s: can't send request: %w", op, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("%s: can't read response body: %w", op, err)
	}

	return Response{
		StatusCode: resp.StatusCode,
		Data:       body,
	}, nil
}
