package gogithub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	var allRepos []Repository
	iter := 0
	maxIter := 30
	for {
		iter += 1
		if iter >= maxIter {
			break
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return allRepos, err
		}
		token := obj.chooseToken()
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := obj.httpClient.Do(req)
		if err != nil {
			return allRepos, err
		}
		defer resp.Body.Close()

		// Check the response status
		if resp.StatusCode != http.StatusOK {
			return allRepos, fmt.Errorf("Error: received %s response status for: %s", resp.Status, url)
		}

		// Read and print the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return allRepos, err
		}
		var response []Repository
		err = json.Unmarshal(body, &response)
		if err != nil {
			return allRepos, err
		}
		allRepos = append(allRepos, response...)

		// Check for pagination (next page)
		linkHeader := resp.Header.Get("Link")
		if linkHeader == "" {
			break // No more pages, exit the loop
		}

		// Extract the URL of the next page from the "Link" header
		links := strings.Split(linkHeader, ",")
		var nextPageURL string
		for _, link := range links {
			if strings.Contains(link, "rel=\"next\"") {
				parts := strings.Split(link, ";")
				nextPageURL = strings.Trim(parts[0], "<> ")
				break
			}
		}

		// If there's no next page URL, stop the loop
		if nextPageURL == "" {
			break
		}

		// Set URL for the next request
		url = nextPageURL
	}

	return allRepos, nil
}
