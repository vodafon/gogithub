package gogithub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Commit struct {
	Author struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Date  string `json:"date"`
	} `json:"author"`
	Committer struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Date  string `json:"date"`
	} `json:"committer"`
	Message string `json:"message"`
	Tree    struct {
		SHA string `json:"sha"`
		URL string `json:"url"`
	} `json:"tree"`
	URL string `json:"url"`
}

type Verification struct {
	Verified   bool    `json:"verified"`
	Reason     string  `json:"reason"`
	Signature  *string `json:"signature"`
	Payload    *string `json:"payload"`
	VerifiedAt *string `json:"verified_at"`
}

type Author struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	UserViewType      string `json:"user_view_type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Committer struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	UserViewType      string `json:"user_view_type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Stats struct {
	Total     int `json:"total"`
	Additions int `json:"additions"`
	Deletions int `json:"deletions"`
}

type File struct {
	SHA         string `json:"sha"`
	Filename    string `json:"filename"`
	Status      string `json:"status"`
	Additions   int    `json:"additions"`
	Deletions   int    `json:"deletions"`
	Changes     int    `json:"changes"`
	BlobURL     string `json:"blob_url"`
	RawURL      string `json:"raw_url"`
	ContentsURL string `json:"contents_url"`
	Patch       string `json:"patch"`
}

type Parent struct {
	SHA     string `json:"sha"`
	URL     string `json:"url"`
	HTMLURL string `json:"html_url"`
}

type CommitResponse struct {
	SHA          string       `json:"sha"`
	NodeID       string       `json:"node_id"`
	Commit       Commit       `json:"commit"`
	URL          string       `json:"url"`
	HTMLURL      string       `json:"html_url"`
	CommentsURL  string       `json:"comments_url"`
	Author       Author       `json:"author"`
	Committer    Committer    `json:"committer"`
	Parents      []Parent     `json:"parents"`
	Stats        Stats        `json:"stats"`
	Files        []File       `json:"files"`
	Verification Verification `json:"verification"`
}

func (obj *Client) GetCommitDiff(u string) (string, error) {
	commit, err := obj.GetCommit(u)
	if err != nil {
		return "", err
	}

	diff := ""
	for _, f := range commit.Files {
		diff += f.Patch + "\n"
	}

	return diff, nil
}

func (obj *Client) GetCommit(u string) (*CommitResponse, error) {
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
		return nil, fmt.Errorf("Error: received %s response status for: %s", resp.Status, u)
	}

	// Read and print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return obj.parseCommit(body)
}

func (obj *Client) parseCommit(res []byte) (*CommitResponse, error) {
	var response CommitResponse
	err := json.Unmarshal(res, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
