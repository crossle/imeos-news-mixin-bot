package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var apiEndpoint = "https://api.bitzhidao.com/newsflashes"

type HoldData struct {
	Items []NewsFlash `json:"data"`
}

type NewsFlash struct {
	ID        string `json:"_id"`
	Content   string `json:"content"`
	IssueTime int64  `json:"issueTime"`
}

func GetStories() ([]NewsFlash, error) {
	content, err := getJSON(apiEndpoint)
	if err != nil {
		return nil, err
	}
	var f HoldData
	if err := json.Unmarshal(content, &f); err != nil {
		return nil, err
	}
	return f.Items, nil
}

func getJSON(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Close = true // connection reset

	client := new(http.Client)
	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	return content, err
}
