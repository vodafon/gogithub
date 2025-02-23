package gogithub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type License struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SpdxID string `json:"spdx_id"`
	URL    string `json:"url"`
}

type Repository struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	FullName     string   `json:"full_name"`
	Owner        Owner    `json:"owner"`
	HTMLURL      string   `json:"html_url"`
	Description  string   `json:"description"`
	Fork         bool     `json:"fork"`
	LanguagesURL string   `json:"languages_url"`
	License      License  `json:"license"`
	Topics       []string `json:"topics"`
}

type Owner struct {
	Login string `json:"login"`
}

func (obj *Client) GetRepositories(username string) ([]Repository, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", username)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	token := obj.chooseToken()
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := obj.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: received %s response status for: %s", resp.Status, url)
	}

	// Read and print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response []Repository
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
