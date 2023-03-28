package playerCount

import (
	"encoding/json"
	"net/http"
)

type Client interface {
	Get() (int, error)
}

type client struct {
	url string
}

func NewClient(mcServerIp string) Client {
	const apiURL = "http://mcapi.us/server/status?ip="

	return &client{
		url: apiURL + mcServerIp,
	}
}

func (pc *client) Get() (int, error) {
	var resp Resp

	r, err := http.Get(pc.url)
	if err != nil {
		return 0, err
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return 0, err
	}

	return resp.Players.Now, nil
}

// ------------------------------- API RESPONSE MODELS --------------------------------

type Resp struct {
	Players Players `json:"players"`
}

type Players struct {
	Now int `json:"now"`
}
