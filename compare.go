package gogithub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CompareResponse struct {
	URL          string `json:"url"`
	HTMLURL      string `json:"html_url"`
	PermalinkURL string `json:"permalink_url"`
	DiffURL      string `json:"diff_url"`
	PatchURL     string `json:"patch_url"`
	Files        []File `json:"files"`
}

func (obj *Client) GetCompare(u string) (*CompareResponse, error) {
	req, err := http.NewRequest("GET", u, nil)
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
		return nil, fmt.Errorf("received %s response status for: %s", resp.Status, u)
	}

	// Read and print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return obj.parseCompare(body)
}

func (obj *Client) parseCompare(res []byte) (*CompareResponse, error) {
	var response CompareResponse
	err := json.Unmarshal(res, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
