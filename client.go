package gogithub

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Client struct {
	httpClient *http.Client
	authTokens []string
	r          *rand.Rand
}

func NewClientWithTokensFile(filepath string) (*Client, error) {
	tokens, err := ghTokensFromFile(filepath)
	if err != nil {
		return nil, err
	}
	return NewClientWithTokens(tokens)
}

func NewClientWithToken(token string) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("github token is empty")
	}
	tokens := []string{token}
	return NewClientWithTokens(tokens)
}

func NewClientWithTokens(tokens []string) (*Client, error) {
	httpclient := &http.Client{
		Timeout: 5 * time.Second,
	}

	client := &Client{
		httpClient: httpclient,
		authTokens: tokens,
		r:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return client, nil
}

func IsTokenValid(token string) (bool, error) {
	httpclient := &http.Client{
		Timeout: 5 * time.Second,
	}

	// The /user endpoint is the standard way to check "Who am I?"
	// It returns 200 OK if the token is valid, and 401 Unauthorized if not.
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	// GitHub requires a User-Agent header; requests without it are often rejected.
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := httpclient.Do(req)
	if err != nil {
		return false, fmt.Errorf("network error during validation: %w", err)
	}
	defer resp.Body.Close()

	// 200 OK means the token is valid.
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	// 401 Unauthorized or 403 Forbidden usually means the token is invalid or expired.
	// We treat any non-200 status as an invalid token for this boolean check.
	return false, nil
}

func (obj *Client) chooseToken() string {
	randomIndex := rand.Intn(len(obj.authTokens))
	return obj.authTokens[randomIndex]
}

func ghTokensFromFile(filepath string) ([]string, error) {
	if filepath == "" {
		return nil, fmt.Errorf("path to file with tokens is empty")
	}
	// Open the file
	file, err := os.Open(filepath) // Replace with your file name
	if err != nil {
		return nil, fmt.Errorf("Error opening file %s: %v", filepath, err)
	}
	defer file.Close() // Ensure the file is closed after reading

	// Create a slice to hold the lines
	var lines []string

	// Create a new scanner
	scanner := bufio.NewScanner(file)

	// Read lines from the file
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading file %s: %v", filepath, err)
	}
	return lines, nil
}
