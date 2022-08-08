package status

import (
	"encoding/json"
	"net/http"
	"os"
)

type Resp struct {
	Players Players `json:"players"`
}

type Players struct {
	Now int `json:"now"`
}

var (
	apiURL   = "http://mcapi.us/server/status?ip="
	serverIP = os.Getenv("MC_SERVER_IP")
)

func GetPlayerCount() (int, error) {
	var resp Resp

	r, err := http.Get(apiURL + serverIP)
	if err != nil {
		return 0, err
	}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return 0, err
	}

	return resp.Players.Now, nil
}
